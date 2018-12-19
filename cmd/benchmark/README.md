Benchmark
===============

Usage
----------
#### Build

    make benchmark
   
#### Start insolard

    ./scripts/insolard/launchnet.sh -g
   
#### Start benchmark

    ./bin/benchmark -c=3 -r=1 -k=scripts/insolard/configs/root_member_keys.json

### Options

        -c concurrency
                Number of concurrent users. Default is one. 

        -r repetitions
                Number of repetitions for one user. Default is one.

        -i input
                Path to file with initial data - references of members.
                If you don't provide input file, new members will be generated automatically.

        -o output
                Path to output file (use - for STDOUT).

        -k rootmemberkeys
                Path to file with RootMember keys.

        -u apiurl (may be specified multiple times for roundrobin requests)
                API url for requests (default - http://localhost:19191/api).

        -l loglevel
                Log level (default - info).
