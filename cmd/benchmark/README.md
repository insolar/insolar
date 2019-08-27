Benchmark
===============

Usage
----------
#### Build

    make benchmark
   
#### Start insolard

    ./scripts/insolard/launchnet.sh -g
   
#### Start benchmark

    ./bin/benchmark -c=4 -r=25 -k=.artifacts/launchnet/configs/

or you can run benchmark with

    ./scripts/bench.sh

### Options

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
                Saves members to file .artifacts/bench-members/members.txt.
                If false, file wont be updated. Default is false.

        -m usemembers
                Use members from file .artifacts/bench-members/members.txt.
                If false, wright info about created members in this file. Default is false. 
        -members-dir
                Dir for saving memebers data
        -b nocheckbalance
                If true, don't check balance at the start/end of transfers. Default is false. 
        -discovery-nodes-logs-dir
                Launchnet logs dir for checking errors
