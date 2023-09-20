// +build slowtest

package pulse

import (
	"context"
	"math/rand"
	"os"
	"sort"
	"sync"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/assert"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/ledger/heavy/migration"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/tests/common"
	"github.com/stretchr/testify/require"
)

var db *PostgresDB

var (
	poolLock     sync.Mutex
	globalPgPool *pgxpool.Pool
)

func setPool(pool *pgxpool.Pool) {
	poolLock.Lock()
	defer poolLock.Unlock()
	globalPgPool = pool
}

func getPool() *pgxpool.Pool {
	poolLock.Lock()
	defer poolLock.Unlock()
	return globalPgPool
}

// TestMain does the before and after setup
func TestMain(m *testing.M) {
	ctx := context.Background()
	log.Info("[TestMain] About to start PostgreSQL...")
	pgURL, stopPostgreSQL := common.StartPostgreSQL()
	log.Info("[TestMain] PostgreSQL started!")

	pool, err := pgxpool.Connect(ctx, pgURL)
	if err != nil {
		stopPostgreSQL()
		log.Panicf("[TestMain] pgxpool.Connect() failed: %v", err)
	}

	migrationPath := "../../insolar-scripts/migration"
	cwd, err := os.Getwd()
	if err != nil {
		stopPostgreSQL()
		panic(errors.Wrap(err, "[TestMain] os.Getwd failed"))
	}
	log.Infof("[TestMain] About to run PostgreSQL migration, cwd = %s, migration migrationPath = %s", cwd, migrationPath)
	ver, err := migration.MigrateDatabase(ctx, pool, migrationPath)
	if err != nil {
		stopPostgreSQL()
		panic(errors.Wrap(err, "Unable to migrate database"))
	}
	log.Infof("[TestMain] PostgreSQL database migration done, current schema version: %d", ver)

	setPool(pool)

	// Run all tests
	code := m.Run()

	log.Info("[TestMain] Cleaning up...")
	stopPostgreSQL()
	os.Exit(code)
}

func cleanupDatabase() {
	ctx := context.Background()
	conn, err := getPool().Acquire(ctx)
	if err != nil {
		panic("Unable to acquire a database connection")
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "DELETE FROM pulses CASCADE")
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(ctx, "DELETE FROM key_value")
	if err != nil {
		panic(err)
	}
}

func generatePulse(pn insolar.PulseNumber, prev insolar.PulseNumber, next insolar.PulseNumber) *insolar.Pulse {
	conf1 := insolar.PulseSenderConfirmation{
		PulseNumber:     pn,
		ChosenPublicKey: "ololo",
		Entropy:         [insolar.EntropySize]byte{3, 3, 2, 2, 1, 1},
		Signature:       []byte{1, 1, 2, 2, 3, 3},
	}
	conf2 := insolar.PulseSenderConfirmation{
		PulseNumber:     pn,
		ChosenPublicKey: "trololo",
		Entropy:         [insolar.EntropySize]byte{3, 3, 2, 2, 1, 1},
		Signature:       []byte{1, 1, 2, 2, 3, 3},
	}
	signs := make(map[string]insolar.PulseSenderConfirmation, 1)
	signs[conf1.ChosenPublicKey] = conf1
	signs[conf2.ChosenPublicKey] = conf2
	return &insolar.Pulse{
		PulseNumber:      pn,
		PrevPulseNumber:  prev,
		NextPulseNumber:  next,
		PulseTimestamp:   123456789,
		EpochPulseNumber: pulse.Epoch(1234),
		OriginID:         [insolar.OriginIDSize]byte{3, 2, 1},
		Entropy:          [insolar.EntropySize]byte{1, 2, 3},
		Signs:            signs,
	}
}

func TestPostgresWriteReadAndLatest(t *testing.T) {
	defer cleanupDatabase()

	ctx := context.Background()
	pn := gen.PulseNumber()
	db := NewPostgresDB(getPool())

	// Make sure there is no such pulse in PostgresDB yet
	_, err := db.ForPulseNumber(ctx, pn)
	require.Error(t, err)

	writePulse := generatePulse(pn, gen.PulseNumber(), gen.PulseNumber())

	// Write the pulse to the database
	err = db.Append(ctx, *writePulse)
	require.NoError(t, err)

	// Read the pulse from the database
	readPulse, err := db.ForPulseNumber(ctx, pn)
	require.NoError(t, err)
	require.Equal(t, *writePulse, readPulse)

	// Make sure .Latest returns something now when we know there is data in the database
	_, err = db.Latest(ctx)
	require.NoError(t, err)
}

func TestPostgresForwardsBackwards(t *testing.T) {
	defer cleanupDatabase()

	ctx := context.Background()
	db := NewPostgresDB(getPool())
	pulsesNum := 10
	pulseNumbers := make([]insolar.PulseNumber, pulsesNum+2)
	pulses := make([]*insolar.Pulse, pulsesNum+2)
	for i := 0; i < len(pulseNumbers); i++ {
		pulseNumbers[i] = gen.PulseNumber()
	}
	sort.Slice(pulseNumbers, func(i, j int) bool {
		return pulseNumbers[i] < pulseNumbers[j]
	})

	startPulseIdx := 1
	endPulseIdx := len(pulseNumbers) - 1
	for i := startPulseIdx; i < endPulseIdx; i++ {
		pn := pulseNumbers[i]
		prev := pulseNumbers[i-1]
		next := pulseNumbers[i+1]
		pulses[i] = generatePulse(pn, prev, next)
		err := db.Append(ctx, *pulses[i])
		require.NoError(t, err)
	}

	// Make sure Forwards/Backwards happy path
	foundPulse, err := db.Forwards(ctx, pulseNumbers[startPulseIdx], 5)
	require.NoError(t, err)
	require.Equal(t, *pulses[startPulseIdx+5], foundPulse)

	foundPulse, err = db.Backwards(ctx, pulseNumbers[startPulseIdx+9], 9)
	require.NoError(t, err)
	require.Equal(t, *pulses[startPulseIdx], foundPulse)

	// Also check `not found` path
	_, err = db.Forwards(ctx, pulseNumbers[endPulseIdx-4], 5)
	require.Error(t, err)

	_, err = db.Backwards(ctx, pulseNumbers[startPulseIdx+6], 10)
	require.Error(t, err)
}

func TestPostgresTruncateHead(t *testing.T) {
	defer cleanupDatabase()

	ctx := context.Background()
	db := NewPostgresDB(getPool())
	pulsesNum := 10
	pulseNumbers := make([]insolar.PulseNumber, pulsesNum+2)
	pulses := make([]*insolar.Pulse, pulsesNum+2)
	for i := 0; i < len(pulseNumbers); i++ {
		pulseNumbers[i] = gen.PulseNumber()
	}
	sort.Slice(pulseNumbers, func(i, j int) bool {
		return pulseNumbers[i] < pulseNumbers[j]
	})

	startPulseIdx := 1
	endPulseIdx := len(pulseNumbers) - 1
	for i := startPulseIdx; i < endPulseIdx; i++ {
		pn := pulseNumbers[i]
		prev := pulseNumbers[i-1]
		next := pulseNumbers[i+1]
		pulses[i] = generatePulse(pn, prev, next)
		err := db.Append(ctx, *pulses[i])
		require.NoError(t, err)
	}

	// Call TruncateHead
	err := db.TruncateHead(ctx, pulseNumbers[startPulseIdx+pulsesNum/2])
	require.NoError(t, err)

	// Make sure half of the pulses are still in the database...
	for i := startPulseIdx; i < startPulseIdx+pulsesNum/2; i++ {
		readPulse, err := db.ForPulseNumber(ctx, pulseNumbers[i])
		require.NoError(t, err)
		require.Equal(t, *pulses[i], readPulse)
	}

	// ...and another half is gone
	for i := startPulseIdx + pulsesNum/2; i < endPulseIdx; i++ {
		_, err := db.ForPulseNumber(ctx, pulseNumbers[i])
		require.Error(t, err)
	}
}

func TestPostgresPulse_Components(t *testing.T) {
	defer cleanupDatabase()

	ctx := inslogger.TestContext(t)

	memStorage := NewStorageMem()
	dbStorage := NewPostgresDB(getPool())

	var pulses []insolar.Pulse
	f := fuzz.New().Funcs(func(p *insolar.Pulse, c fuzz.Continue) {
		p.PulseNumber = gen.PulseNumber()
		p.Signs = make(map[string]insolar.PulseSenderConfirmation, 1)
		_, err := rand.Read(p.Entropy[:])
		require.NoError(t, err)
	})
	f.NilChance(0).NumElements(10, 20)
	f.Fuzz(&pulses)

	var appended []insolar.Pulse
	latest := pulses[0]
	for i, p := range pulses {
		// Append appends if Pulse is greater.
		memErr := memStorage.Append(ctx, p)
		dbErr := dbStorage.Append(ctx, p)
		if p.PulseNumber <= latest.PulseNumber && i > 0 {
			assert.Equal(t, ErrBadPulse, memErr)
			assert.Equal(t, ErrBadPulse, dbErr)
			continue
		}
		latest = p
		appended = append(appended, p)

		// Latest returns correct Pulse.
		memLatest, memErr := memStorage.Latest(ctx)
		dbLatest, dbErr := dbStorage.Latest(ctx)
		assert.NoError(t, memErr)
		assert.NoError(t, dbErr)
		assert.Equal(t, p, memLatest)
		assert.Equal(t, p, dbLatest)

		// ForPulse returns correct value
		memForPulse, memErr := memStorage.ForPulseNumber(ctx, p.PulseNumber)
		dbForPulse, dbErr := dbStorage.ForPulseNumber(ctx, p.PulseNumber)
		assert.NoError(t, memErr)
		assert.NoError(t, dbErr)
		assert.Equal(t, p, memForPulse)
		assert.Equal(t, p, dbForPulse)
	}

	// Forwards returns correct value.
	{
		steps := rand.Intn(len(appended))
		memPulse, memErr := memStorage.Forwards(ctx, appended[0].PulseNumber, steps)
		dbPulse, dbErr := dbStorage.Forwards(ctx, appended[0].PulseNumber, steps)
		assert.NoError(t, memErr)
		assert.NoError(t, dbErr)
		assert.Equal(t, appended[steps], memPulse)
		assert.Equal(t, appended[steps], dbPulse)
	}
	// Backwards returns correct value.
	{
		steps := rand.Intn(len(appended))
		memPulse, memErr := memStorage.Backwards(ctx, appended[len(appended)-1].PulseNumber, steps)
		dbPulse, dbErr := dbStorage.Backwards(ctx, appended[len(appended)-1].PulseNumber, steps)
		assert.NoError(t, memErr)
		assert.NoError(t, dbErr)
		assert.Equal(t, appended[len(appended)-steps-1], memPulse)
		assert.Equal(t, appended[len(appended)-steps-1], dbPulse)
	}
}
