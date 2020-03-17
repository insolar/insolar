Benchmark
===============

Usage
----------
#### Build

    make benchmark
   
#### Start insolard

    ./insolar-scripts/insolard/launchnet.sh -g
   
#### Start benchmark

    ./bin/benchmark -c=4 -r=25 -k=.artifacts/launchnet/configs/

or you can run benchmark with

    ./scripts/bench.sh

### Options
        --check-all-balance
                If true, just check balance of every object from file and don't run any scenario. Default is false.

        --check-members-balance
                If true, just check balance of every ordinary member from file, (without general entities), and don't run any scenario. Default is false.

        --check-total-balance
                If true, check total balance of members from file and don't run any scenario. Default is false.

        -c concurrency
                Number of concurrent users. Default is one. 

        -r repetitions
                Number of repetitions for one user. Default is one.

        -o output
                Path to output file (use - for STDOUT).

        -k rootmemberkeys
                Path to file with RootMember keys.

        -a adminapiurl (may be specified multiple times for roundrobin requests)
                API url for requests (default - http://localhost:19001/admin-api/rpc).
                
        -p publicapiurl (may be specified multiple times for roundrobin requests)
                API url for requests (default - http://localhost:19101/api/rpc).

        -l loglevel
                Log level (default - info).

        -s savemembers
                Saves members to file .artifacts/bench-members/members.txt (file can be change with members-file option).
                If false, file wont be updated. Default is false.
                If nocheckbalance set to false, and run was successful, balances in file will be updated after scenario.

        -m usemembers
                Use members from file .artifacts/bench-members/members.txt (file can be change with members-file option).
                If false, wright info about created members in this file. Default is false. 
                If nocheckbalance set to false, and run was successful, balances in file will be updated after scenario.

        --members-file
                Path to file for saving members data

        -b nocheckbalance
                If true, don't check balance at the start/end of transfers. Default is false.
                If false, and savemembers or usemembers provided, and run was successful, balances in file will be updated after scenario.

        -t scenarioname
                Name of scenario. Default scenario is "transfer" scenario.
                You can choose "createMember" for create member scenario.
                You can choose "migration" for migration scenario.
                You can choose "transferTwoSides" for two sides transfer scenario.
                You can choose "depositTransfer" for transfer money from deposit to account.

        --discovery-nodes-logs-dir
                Launchnet logs dir for checking errors

        -R --retries
                Number of request retries if ServiceUnavailable error received

        -P --retry-period
                Time to wait between retries (accepts go duration format: 0.5s, 500ms, 1m30s etc.)
