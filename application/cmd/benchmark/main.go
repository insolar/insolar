//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/insolar/insolar/application/api/sdk"
	"github.com/insolar/insolar/insolar/defaults"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/testutils"
)

const (
	defaultStdoutPath   = "-"
	createMemberRetries = 5
	balanceCheckRetries = 10
	balanceCheckDelay   = 5 * time.Second
)

var (
	defaultMemberFile         = filepath.Join(defaults.ArtifactsDir(), "bench-members", "members.txt")
	defaultDiscoveryNodesLogs = defaults.LaunchnetDiscoveryNodesLogsDir()

	memberFile          string
	output              string
	concurrent          int
	repetitions         int
	memberKeys          string
	adminAPIURLs        []string
	publicAPIURLs       []string
	logLevel            string
	logLevelServer      string
	saveMembersToFile   bool
	useMembersFromFile  bool
	noCheckBalance      bool
	checkMembersBalance bool
	checkAllBalance     bool
	checkTotalBalance   bool
	scenarioName        string
	discoveryNodesLogs  string
	maxRetries          int
	retryPeriod         time.Duration
)

func parseInputParams() {
	pflag.StringVarP(&output, "output", "o", defaultStdoutPath, "output file (use - for STDOUT)")
	pflag.IntVarP(&concurrent, "concurrent", "c", 1, "concurrent users")
	pflag.IntVarP(&repetitions, "repetitions", "r", 1, "repetitions for one user")
	pflag.StringVarP(&memberKeys, "memberkeys", "k", "", "path to dir with members keys")
	pflag.StringArrayVarP(&adminAPIURLs, "adminurls", "a", []string{"http://localhost:19001/admin-api/rpc"}, "url to admin api")
	pflag.StringArrayVarP(&publicAPIURLs, "publicurls", "p", []string{"http://localhost:19101/api/rpc"}, "url to public api")
	pflag.StringVarP(&logLevel, "loglevel", "l", "info", "log level for benchmark")
	pflag.StringVarP(&logLevelServer, "loglevelserver", "L", "", "server log level")
	pflag.BoolVarP(&saveMembersToFile, "savemembers", "s", false, "save members to file")
	pflag.BoolVarP(&useMembersFromFile, "usemembers", "m", false, "use members from file")
	pflag.StringVarP(&memberFile, "members-file", "", defaultMemberFile, "dir for saving members data")
	pflag.BoolVarP(&noCheckBalance, "nocheckbalance", "b", false, "don't check balance at the end")
	pflag.BoolVarP(&checkMembersBalance, "check-members-balance", "", false, "check balance of every ordinary member from file, don't run any scenario")
	pflag.BoolVarP(&checkAllBalance, "check-all-balance", "", false, "check balance of every object from file, and don't run any scenario")
	pflag.BoolVarP(&checkTotalBalance, "check-total-balance", "", false, "check total balance of members from file, don't run any scenario")
	pflag.StringVarP(&scenarioName, "scenarioname", "t", "", "name of scenario")
	pflag.StringVarP(&discoveryNodesLogs, "discovery-nodes-logs-dir", "", defaultDiscoveryNodesLogs, "launchnet logs dir for checking errors")
	pflag.IntVarP(&maxRetries, "retries", "R", 0, "number of request attempts after getting -31429 error. -1 retries infinitely")
	pflag.DurationVarP(&retryPeriod, "retry-period", "P", 0, "delay between retries")
	pflag.Parse()
}

func chooseOutput(path string) (io.Writer, error) {
	var res io.Writer
	if path == defaultStdoutPath {
		res = os.Stdout
	} else {
		var err error
		res, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't open file for writing")
		}
	}
	return res, nil
}

func writeToOutput(out io.Writer, data string) {
	_, err := out.Write([]byte(data))
	check("Can't write data to output", err)
}

func check(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func newTransferDifferentMemberScenarios(out io.Writer, insSDK *sdk.SDK, concurrent int, repetitions int) benchmark {
	return benchmark{
		scenario: &walletToWalletTransferScenario{
			insSDK: insSDK,
		},
		concurrent:  concurrent,
		repetitions: repetitions,
		name:        "TransferDifferentMembers",
		out:         out,
	}
}

func newTransferTwoSidesScenario(out io.Writer, insSDK *sdk.SDK, concurrent int, repetitions int) benchmark {
	return benchmark{
		scenario: &walletToWalletTwoSidesScenario{
			insSDK: insSDK,
		},
		concurrent:  concurrent,
		repetitions: repetitions,
		name:        "TransferTwoSides",
		out:         out,
	}
}

func newCreateMemberScenarios(out io.Writer, insSDK *sdk.SDK, concurrent int, repetitions int) benchmark {
	return benchmark{
		scenario: &createMemberScenario{
			insSDK: insSDK,
		},
		concurrent:  concurrent,
		repetitions: repetitions,
		name:        "CreateMember",
		out:         out,
	}
}

func newMigrationScenarios(out io.Writer, insSDK *sdk.SDK, concurrent int, repetitions int) benchmark {
	return benchmark{
		scenario: &migrationScenario{
			insSDK: insSDK,
		},
		concurrent:  concurrent,
		repetitions: repetitions,
		name:        "Migration",
		out:         out,
	}
}

func newDepositTransferScenarios(out io.Writer, insSDK *sdk.SDK, concurrent int, repetitions int) benchmark {
	return benchmark{
		scenario: &depositTransferScenario{
			insSDK: insSDK,
		},
		concurrent:  concurrent,
		repetitions: repetitions,
		name:        "DepositTransfer",
		out:         out,
	}
}

func startScenario(ctx context.Context, b benchmark) {
	err := b.scenario.canBeStarted()
	check(fmt.Sprintf("Scenario %s can not be started:", b.getName()), err)

	writeToOutput(b.getOut(), fmt.Sprintf("Scenario %s started: \n", b.getName()))

	start := time.Now()
	logReaderCloseChan := testutils.NodesErrorLogReader(discoveryNodesLogs, b.getOut())

	b.start(ctx)
	elapsed := time.Since(start)
	writeToOutput(b.getOut(), fmt.Sprintf("Scenario %s took: %s \n", b.getName(), elapsed))

	close(logReaderCloseChan)
	printResults(b)
}

func printResults(b benchmark) {
	speed := b.getOperationPerSecond()
	writeToOutput(b.getOut(), fmt.Sprintf("Scenario %s: Speed - %f resp/s \n", b.getName(), speed))
	writeToOutput(
		b.getOut(),
		fmt.Sprintf(
			"Scenario %s: Average Request Duration - %s\n",
			b.getName(), b.getAverageOperationDuration(),
		),
	)
	b.printResult()
}

func createMembers(insSDK *sdk.SDK, count int, migration bool) []sdk.Member {
	var (
		members []sdk.Member
		member  sdk.Member
		traceID string
		err     error
	)

	for i := 0; i < count; i++ {
		retries := createMemberRetries
		for retries > 0 {
			member, traceID, err = createMember(insSDK, migration)
			if err != nil {
				fmt.Printf("Retry to create member. TraceID: %s Error is: %s\n", traceID, err.Error())
				retries--
				continue
			}
			members = append(members, member)
			break
		}
		check(fmt.Sprintf("Couldn't create member after retries: %d", createMemberRetries), err)
	}
	return members
}

func createMember(insSDK *sdk.SDK, migration bool) (sdk.Member, string, error) {
	var (
		member  sdk.Member
		traceID string
		err     error
	)

	if migration {
		member, traceID, err = insSDK.MigrationCreateMember()
	} else {
		member, traceID, err = insSDK.CreateMember()
	}

	if err != nil {
		return nil, traceID, errors.Wrap(err, "Failed to create member")
	}

	traceID, err = insSDK.Transfer("100000000000000", insSDK.GetRootMember(), member)

	return member, traceID, errors.Wrap(err, "Failed to transfer initial amount")
}

func getTotalBalance(insSDK *sdk.SDK, members []sdk.Member) (*big.Int, map[string]*big.Int) {
	type Result struct {
		num     int
		balance *big.Int
		err     error
	}
	nmembers := len(members)

	membersWithBalanceMap := make(map[string]*big.Int, nmembers)
	membersWithBalanceMapLock := sync.Mutex{}

	var wg sync.WaitGroup
	wg.Add(nmembers)
	results := make(chan Result, nmembers)

	// execute all queries in parallel
	for i := 0; i < nmembers; i++ {
		go func(m sdk.Member, num int) {
			res := Result{num: num}
			balance, deposits, err := insSDK.GetBalance(m)
			if err == nil {
				for _, d := range deposits {
					depositBalanceStr, ok := d.(map[string]interface{})["balance"].(string)
					if !ok {
						err = errors.New("failed to get balance from deposit")
					}
					depositBalance, ok := new(big.Int).SetString(depositBalanceStr, 10)
					if !ok {
						err = errors.New("failed to parse balance to big.Int")
					}

					balance = balance.Add(balance, depositBalance)
				}
			}
			res.balance, res.err = balance, err
			results <- res
			membersWithBalanceMapLock.Lock()
			membersWithBalanceMap[m.GetReference()] = res.balance
			membersWithBalanceMapLock.Unlock()
			wg.Done()
		}(members[i], i)
	}

	wg.Wait()
	totalBalance := big.NewInt(0)
	for i := 0; i < nmembers; i++ {
		res := <-results
		if res.err != nil {
			fmt.Printf("Can't get balance for %v-th member: %v\n", res.num, res.err)
			continue
		}
		b := totalBalance
		totalBalance.Add(b, res.balance)
	}

	return totalBalance, membersWithBalanceMap
}

func getMembers(insSDK *sdk.SDK, number int, migration bool) ([]sdk.Member, error) {
	var members []sdk.Member
	var err error

	if useMembersFromFile {
		// from file we load not just number of members, but also migration admin or fee member
		for i := 0; i < number+2; i++ {
			if migration {
				members = append(members, &sdk.MigrationMember{})
			} else {
				members = append(members, &sdk.CommonMember{})
			}
		}
		err = loadMembers(&members)
		if err != nil {
			return nil, errors.Wrap(err, "error while loading members: ")
		}
	} else {
		start := time.Now()
		members = createMembers(insSDK, number, migration)
		creationTime := time.Since(start)
		fmt.Printf("Members were created in %s\n", creationTime)
		fmt.Printf("Average creation of member time - %s\n", time.Duration(int64(creationTime)/int64(concurrent*2)))
	}

	return members, nil
}

func saveMembers(members []sdk.Member) error {
	dir, _ := path.Split(memberFile)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return errors.Wrap(err, "couldn't create dir for file")
	}
	file, err := os.Create(memberFile)
	if err != nil {
		return errors.Wrap(err, "couldn't create file")
	}
	defer file.Close() // nolint: errcheck

	result, err := json.MarshalIndent(members, "", "    ")
	if err != nil {
		return errors.Wrap(err, "couldn't marshal members in json")
	}
	_, err = file.Write(result)
	return errors.Wrap(err, "couldn't save members in file")
}

func loadMembers(members *[]sdk.Member) error {
	rawMembers, err := ioutil.ReadFile(memberFile)
	if err != nil {
		return errors.Wrap(err, "can't read members from file")
	}

	err = json.Unmarshal(rawMembers, members)
	if err != nil {
		return errors.Wrap(err, "can't unmarshal members from file")
	}

	return nil
}

func main() {
	parseInputParams()

	// Start benchmark time
	t := time.Now()
	fmt.Printf("Start: %s\n\n", t.String())

	err := log.SetLevel(logLevel)
	check(fmt.Sprintf("Can't set '%s' level on logger:", logLevel), err)

	out, err := chooseOutput(output)
	check("Problems with output file:", err)

	insSDK, err := sdk.NewSDK(adminAPIURLs, publicAPIURLs, memberKeys, sdk.Options{
		RetryPeriod: retryPeriod,
		MaxRetries:  maxRetries,
	})
	check("SDK is not initialized: ", err)

	err = insSDK.SetLogLevel(logLevelServer)
	check("Failed to parse log level: ", err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP)

	b := switchScenario(out, insSDK)

	go func() {
		stopGracefully := true
		for {
			sig := <-sigChan

			switch sig {
			case syscall.SIGHUP:
				printResults(b)
			case syscall.SIGINT:
				if !stopGracefully {
					log.Fatal("Force quiting.")
				} else {
					log.Info("Gracefully finishing benchmark. Press Ctrl+C again to force quit.")
				}

				stopGracefully = false
				cancel()
			}
		}
	}()

	if checkMembersBalance || checkTotalBalance || checkAllBalance {
		var commonMembers []*sdk.CommonMember
		rawMembers, err := ioutil.ReadFile(memberFile)
		check("Can't read members from file: ", err)

		err = json.Unmarshal(rawMembers, &commonMembers)
		check("Error while loading members for checking balances: ", err)
		var members []sdk.Member

		feeMemberRef := insSDK.GetFeeMember().GetReference()
		migrationAdminRef := insSDK.GetMigrationAdminMember().GetReference()
		for _, m := range commonMembers {
			if checkMembersBalance {
				if m.GetReference() == feeMemberRef {
					continue
				}
				if m.GetReference() == migrationAdminRef {
					continue
				}
			}
			members = append(members, m)
		}

		totalBalance, membersWithBalanceMap := getTotalBalance(insSDK, members)
		totalFileBalance := checkBalanceAtFile(members, membersWithBalanceMap)
		if totalFileBalance.Cmp(totalBalance) != 0 {
			log.Fatalf("Total balance mismatch: all members balance at file - %s, all members balance at system - %s \n", totalFileBalance, totalBalance)
		}
		log.Info("Balances for members from file was successfully checked.")
		return
	}

	b.scenario.prepare(repetitions)

	var totalBalanceBefore *big.Int
	if !noCheckBalance {
		totalBalanceBefore, _ = getTotalBalance(insSDK, b.scenario.getBalanceCheckMembers())
	}

	if saveMembersToFile {
		err = saveMembers(b.scenario.getBalanceCheckMembers())
		check("Error while saving members before scenario: ", err)
	}

	startScenario(ctx, b)

	// Finish benchmark time
	t = time.Now()
	fmt.Printf("\nFinish: %s\n\n", t.String())

	if !noCheckBalance {
		membersWithBalanceMap := checkBalance(insSDK, totalBalanceBefore, b.scenario.getBalanceCheckMembers())
		// update balances in file
		for _, m := range b.scenario.getBalanceCheckMembers() {
			b := membersWithBalanceMap[m.GetReference()]
			m.SetBalance(b)
		}
		if saveMembersToFile || useMembersFromFile {
			err := saveMembers(b.scenario.getBalanceCheckMembers())
			check("Error while saving members after scenario: ", err)
		}
	}
}

func switchScenario(out io.Writer, insSDK *sdk.SDK) benchmark {
	var b benchmark

	switch scenarioName {
	case "transferTwoSides":
		b = newTransferTwoSidesScenario(out, insSDK, concurrent, repetitions)
	case "createMember":
		b = newCreateMemberScenarios(out, insSDK, concurrent, repetitions)
	case "migration":
		b = newMigrationScenarios(out, insSDK, concurrent, repetitions)
	case "depositTransfer":
		b = newDepositTransferScenarios(out, insSDK, concurrent, repetitions)
	default:
		b = newTransferDifferentMemberScenarios(out, insSDK, concurrent, repetitions)
	}

	return b
}

func checkBalance(insSDK *sdk.SDK, totalBalanceBefore *big.Int, balanceCheckMembers []sdk.Member) map[string]*big.Int {
	totalBalanceAfter := big.NewInt(0)
	var membersWithBalanceMap map[string]*big.Int

	for nretries := 0; nretries < balanceCheckRetries; nretries++ {
		totalBalanceAfter, membersWithBalanceMap = getTotalBalance(insSDK, balanceCheckMembers)
		if totalBalanceAfter.Cmp(totalBalanceBefore) == 0 {
			break
		}
		fmt.Printf("Total balance before and after don't match: %v vs %v - retrying in %s ...\n",
			totalBalanceBefore, totalBalanceAfter, balanceCheckDelay)
		time.Sleep(balanceCheckDelay)

	}

	fmt.Printf("Total balance before: %v and after: %v\n", totalBalanceBefore, totalBalanceAfter)
	if totalBalanceAfter.Cmp(totalBalanceBefore) != 0 {
		log.Fatal("Total balance mismatch!\n")
	}

	for n := 0; n < 2; n++ {
		totalBalanceAfter, membersWithBalanceMap = getTotalBalance(insSDK, balanceCheckMembers)
		if totalBalanceAfter.Cmp(totalBalanceBefore) != 0 {
			log.Fatal("Total balance mismatch!\n")
		}

		fmt.Println("Wait if balance changes after matching: ", n)
		time.Sleep(balanceCheckDelay)
	}

	fmt.Printf("Total balance successfully matched\n")
	return membersWithBalanceMap
}

func checkBalanceAtFile(members []sdk.Member, membersWithBalanceMap map[string]*big.Int) *big.Int {
	totalFileBalance := big.NewInt(0)

	for _, m := range members {
		b := m.GetBalance()
		totalFileBalance = totalFileBalance.Add(totalFileBalance, b)

		if checkMembersBalance || checkAllBalance {
			if membersWithBalanceMap[m.GetReference()] == nil {
				log.Fatalf("Balance mismatch: member with ref %s exists in file, but we didn't get its system balance. Balance at file - %s. \n", m.GetReference(), m.GetBalance())
			}
			if b.Cmp(membersWithBalanceMap[m.GetReference()]) != 0 {
				log.Fatalf("Balance mismatch: member with ref %s, balance at file - %s, balance at system - %s \n", m.GetReference(), m.GetBalance(), membersWithBalanceMap[m.GetReference()])
			}
		}
	}
	return totalFileBalance
}
