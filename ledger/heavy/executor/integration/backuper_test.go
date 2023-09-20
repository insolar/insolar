// +build slowtest
// +build !coverage

package intergration

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/stretchr/testify/require"
)

type testKey struct {
	id uint64
}

func (t *testKey) ID() []byte {
	bs := make([]byte, 8)
	binary.PutUvarint(bs, t.id)
	return bs
}

func (t *testKey) Scope() store.Scope {
	return store.ScopeJetDrop
}

func makeBackuperConfig(t *testing.T, prefix string, badgerDir string, recoverDBDir string) (configuration.Ledger, string) {

	cwd, err := os.Getwd()
	if err != nil {
		require.NoError(t, err)
	}

	tmpDir, err := ioutil.TempDir("", "bdb-backup-test-integr-")
	require.NoError(t, err)

	cfg := configuration.Backup{
		ConfirmFile:          "BACKUPED",
		MetaInfoFile:         "META.json",
		TargetDirectory:      tmpDir + "/TARGET/" + prefix,
		TmpDirectory:         tmpDir + "/TMP",
		DirNameTemplate:      "pulse-%d",
		BackupWaitPeriod:     10,
		BackupFile:           "incr.bkp",
		Enabled:              true,
		PostProcessBackupCmd: []string{"bash", "-c", cwd + "/post_process_backup.sh" + " " + badgerDir + " " + recoverDBDir},
	}

	err = os.MkdirAll(cfg.TargetDirectory, 0777)
	require.NoError(t, err)
	err = os.MkdirAll(cfg.TmpDirectory, 0777)
	require.NoError(t, err)

	return configuration.Ledger{
		Backup: cfg,
		Storage: configuration.Storage{
			DataDirectory: badgerDir,
		},
	}, tmpDir
}

func clearData(t *testing.T, tmpDir string) {
	err := os.RemoveAll(tmpDir)
	require.NoError(t, err)
}

func TestBackuper(t *testing.T) {

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	recovTmpDir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(recovTmpDir)

	cfg, tmpDir := makeBackuperConfig(t, t.Name(), tmpdir, recovTmpDir)
	defer clearData(t, tmpDir)

	db, err := store.NewBadgerDB(badger.DefaultOptions(tmpdir))
	require.NoError(t, err)
	defer db.Stop(context.Background())

	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, insolar.GenesisPulse.PulseNumber, db)
	require.NoError(t, err)

	savedKeys := make(map[store.Key]insolar.PulseNumber, 0)

	var stopWriting uint32
	sgWriteStopped := sync.WaitGroup{}
	sgWriteStopped.Add(1)

	testPulse := insolar.GenesisPulse.PulseNumber + insolar.PulseNumber(rand.Int()%20000+1)
	// writing data to db
	go func() {
		for i := 0; ; i++ {
			if atomic.LoadUint32(&stopWriting) != 0 {
				break
			}
			key := &testKey{id: uint64(i)}
			value := testPulse + insolar.PulseNumber(i)
			err := db.Set(key, value.Bytes())
			require.NoError(t, err)
			savedKeys[key] = value
			time.Sleep(time.Duration(rand.Int()%10) * time.Millisecond)
		}
		sgWriteStopped.Done()
	}()

	wgBackup := sync.WaitGroup{}
	numIterations := 15

	wgBackup.Add(numIterations)
	// doing backups
	go func() {
		for i := 0; i < numIterations; i++ {
			err := bm.MakeBackup(context.Background(), testPulse+insolar.PulseNumber(i))
			require.NoError(t, err)
			wgBackup.Done()
			time.Sleep(time.Duration(rand.Int()%1000) * time.Millisecond)
		}
	}()

	// wait for all backups done
	wgBackup.Wait()
	// stop writing to db
	atomic.StoreUint32(&stopWriting, 1)
	// wait for stopping
	sgWriteStopped.Wait()

	require.NotEqual(t, 0, len(savedKeys))

	// final backup to collect all rest records
	err = bm.MakeBackup(context.Background(), testPulse+insolar.PulseNumber(numIterations))
	require.NoError(t, err)

	// load all backups and check all records
	{
		recoveredDB, err := store.NewBadgerDB(badger.DefaultOptions(recovTmpDir))
		require.NoError(t, err)
		defer recoveredDB.Stop(context.Background())

		for k, v := range savedKeys {
			gotRawValue, err := recoveredDB.Get(k)
			require.NoError(t, err)
			gotPulseNumber := insolar.NewPulseNumber(gotRawValue)
			require.Equal(t, v, gotPulseNumber)
		}
	}
}

var binaryPath string

func init() {
	var ok bool

	binaryPath, ok = os.LookupEnv("BIN_DIR")
	if !ok {
		wd, err := os.Getwd()
		binaryPath = filepath.Join(wd, "..", "..", "..", "..", "bin")

		if err != nil {
			panic(err.Error())
		}
	}
}

// prepareBackup uses backupmanager utility to prepare backup for usage
func prepareBackup(t *testing.T, dbDir string) {
	println("=====> Start preparing backup")
	cmd := exec.Command(binaryPath+"/backupmanager", "prepare_backup", "-d", dbDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	require.NoError(t, err)
	err = cmd.Wait()
	require.NoError(t, err)
	println("<===== Finish preparing backup")
}

// createDirForBackup uses backupmanager utility to create empty badger
func createDirForBackup(t *testing.T, dbDir string) {
	println("=====> Start creating db for backup")
	cmd := exec.Command(binaryPath+"/backupmanager", "create", "-d", dbDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	require.NoError(t, err)
	err = cmd.Wait()
	require.NoError(t, err)
	println("<===== Finish creating db for backup")
}

// loadIncrementalBackup uses backupmanager utility to roll backups
func loadIncrementalBackup(t *testing.T, dbDir string, backupFile string) {
	println("=====> Start loading backup")
	cmd := exec.Command(binaryPath+"/backupmanager", "merge", "-t", dbDir, "-n", backupFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	require.NoError(t, err)
	err = cmd.Wait()
	require.NoError(t, err)
	println("<===== Finish loading backup")
}

func makeCurrentBkpDir(cfg configuration.Backup, pulse insolar.PulseNumber) string {
	return filepath.Join(cfg.TargetDirectory, fmt.Sprintf(cfg.DirNameTemplate, pulse))
}

func calculateFileHash(t *testing.T, fileName string) string {
	f, err := os.Open(fileName)
	require.NoError(t, err)
	defer f.Close()
	hasher := sha256.New()
	_, err = io.Copy(hasher, f)
	require.NoError(t, err)

	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func TestBackupSendDeleteRecords(t *testing.T) {

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	recovTmpDir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(recovTmpDir)

	cfg, tmpDir := makeBackuperConfig(t, t.Name(), tmpdir, recovTmpDir)
	defer clearData(t, tmpDir)

	db, err := store.NewBadgerDB(badger.DefaultOptions(tmpdir))
	require.NoError(t, err)
	defer db.Stop(context.Background())

	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, insolar.GenesisPulse.PulseNumber, db)
	require.NoError(t, err)

	key := &testKey{id: uint64(3)}
	deletedKey := &testKey{id: uint64(4)}

	err = db.Set(key, []byte{})
	require.NoError(t, err)

	err = db.Set(deletedKey, []byte{})
	require.NoError(t, err)

	err = db.Delete(deletedKey)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), 100000)
	require.NoError(t, err)
	err = bm.MakeBackup(context.Background(), 100001)
	require.NoError(t, err)

	recoveredDB, err := store.NewBadgerDB(badger.DefaultOptions(recovTmpDir))
	require.NoError(t, err)
	_, err = recoveredDB.Get(key)
	require.NoError(t, err)
	recoveredDB.Stop(context.Background())

	err = db.Delete(key)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), 100002)
	require.NoError(t, err)

	recoveredDB, err = store.NewBadgerDB(badger.DefaultOptions(recovTmpDir))
	require.NoError(t, err)
	defer recoveredDB.Stop(context.Background())

	_, err = recoveredDB.Get(key)
	require.EqualError(t, err, store.ErrNotFound.Error())

	_, err = recoveredDB.Get(deletedKey)
	require.EqualError(t, err, store.ErrNotFound.Error())
}

func TestBackup_FullCycle(t *testing.T) {
	ctx := inslogger.TestContext(t)
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	recovTmpDir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(recovTmpDir)

	cfg, tmpDir := makeBackuperConfig(t, t.Name(), tmpdir, recovTmpDir)
	defer clearData(t, tmpDir)

	db, err := store.NewBadgerDB(badger.DefaultOptions(tmpdir))
	require.NoError(t, err)
	defer db.Stop(context.Background())

	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, insolar.GenesisPulse.PulseNumber, db)
	require.NoError(t, err)

	testPulse := insolar.GenesisPulse.PulseNumber + 10
	testJet := insolar.ZeroJetID

	pulsesDB := pulse.NewBadgerDB(db)
	err = pulsesDB.Append(ctx, insolar.Pulse{PulseNumber: insolar.GenesisPulse.PulseNumber})
	require.NoError(t, err)
	err = pulsesDB.Append(ctx, insolar.Pulse{PulseNumber: testPulse})
	require.NoError(t, err)

	jetsDB := jet.NewBadgerDBStore(db)
	err = jetsDB.Update(ctx, testPulse, true, testJet)
	require.NoError(t, err)

	jetKeeper := executor.NewBadgerJetKeeper(jetsDB, db, pulsesDB)

	err = jetKeeper.AddHotConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)
	err = jetKeeper.AddDropConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)

	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())

	err = bm.MakeBackup(ctx, testPulse)
	require.NoError(t, err)

	prepareBackup(t, recovTmpDir)
	recoveredDB, err := store.NewBadgerDB(badger.DefaultOptions(recovTmpDir))
	require.NoError(t, err)
	defer recoveredDB.Stop(context.Background())

	recoveredJetKeeper := executor.NewBadgerJetKeeper(jet.NewBadgerDBStore(recoveredDB), recoveredDB, pulse.NewBadgerDB(recoveredDB))

	// pulse must be finalized when prepare_backup complete without error
	require.Equal(t, testPulse, recoveredJetKeeper.TopSyncPulse())
}

func copyDir(src, dst string) error {
	cmd := exec.Command("cp", "-vR", src+"/", dst)
	output, err := cmd.CombinedOutput()
	println("copyDir: ", string(output))

	return err
}

// 1. Create db
// 2. Add not all confirmations
// 3. Copy db to different place - backup place
// 4. Finalize current pulse and add confirmations for next one
// 5. Make backup and merge it to backup db
// 6. Launch on backup db and check top sync pulse
func TestBackup_UseMainDBAsBackup(t *testing.T) {
	ctx := inslogger.TestContext(t)
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	backupTmpDir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(backupTmpDir)

	cfg, tmpDir := makeBackuperConfig(t, t.Name(), tmpdir, backupTmpDir)
	defer clearData(t, tmpDir)

	db, err := store.NewBadgerDB(badger.DefaultOptions(tmpdir))
	require.NoError(t, err)

	testPulse := insolar.GenesisPulse.PulseNumber + 10
	testJet := insolar.ZeroJetID

	pulsesDB := pulse.NewBadgerDB(db)
	err = pulsesDB.Append(ctx, insolar.Pulse{PulseNumber: insolar.GenesisPulse.PulseNumber})
	require.NoError(t, err)
	err = pulsesDB.Append(ctx, insolar.Pulse{PulseNumber: testPulse})
	require.NoError(t, err)

	jetsDB := jet.NewBadgerDBStore(db)
	err = jetsDB.Update(ctx, testPulse, true, testJet)
	require.NoError(t, err)

	jetKeeper := executor.NewBadgerJetKeeper(jetsDB, db, pulsesDB)

	err = jetKeeper.AddHotConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)
	err = jetKeeper.AddDropConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)

	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())

	// Stop db and copy it to backup place
	err = db.Stop(ctx)
	require.NoError(t, err)

	{
		// -------------------- Copy db to backup db

		err = copyDir(tmpdir, backupTmpDir)
		require.NoError(t, err)
	}

	{
		// -------------------- run on db again
		db, err = store.NewBadgerDB(badger.DefaultOptions(tmpdir))
		require.NoError(t, err)
		defer db.Stop(ctx)

		// -------------------- finalize pulse
		pulsesDB = pulse.NewBadgerDB(db)
		jetsDB = jet.NewBadgerDBStore(db)
		jetKeeper = executor.NewBadgerJetKeeper(jetsDB, db, pulsesDB)
		jetKeeper.AddBackupConfirmation(ctx, testPulse)
		require.Equal(t, testPulse, jetKeeper.TopSyncPulse())

		// -------------------- and prepare next
		nextPulse := testPulse + 10
		err = jetKeeper.AddHotConfirmation(ctx, nextPulse, testJet, false)
		require.NoError(t, err)
		err = jetKeeper.AddDropConfirmation(ctx, nextPulse, testJet, false)
		require.NoError(t, err)
		err = jetsDB.Update(ctx, nextPulse, true, testJet)
		require.NoError(t, err)
		err = pulsesDB.Append(ctx, insolar.Pulse{PulseNumber: nextPulse})
		require.NoError(t, err)

		// -------------------- make backup
		bm, err := executor.NewBackupMaker(context.Background(), db, cfg, insolar.GenesisPulse.PulseNumber, db)
		require.NoError(t, err)
		err = bm.MakeBackup(ctx, nextPulse)
		require.NoError(t, err)

		// -------------------- merge backup

		prepareBackup(t, backupTmpDir)
		recoveredDB, err := store.NewBadgerDB(badger.DefaultOptions(backupTmpDir))
		require.NoError(t, err)
		defer recoveredDB.Stop(context.Background())

		// check that db is ok
		recoveredJetKeeper := executor.NewBadgerJetKeeper(
			jet.NewBadgerDBStore(recoveredDB),
			recoveredDB,
			pulse.NewBadgerDB(recoveredDB))

		// pulse must be finalized when prepare_backup complete without error
		require.Equal(t, nextPulse, recoveredJetKeeper.TopSyncPulse())
	}
}
