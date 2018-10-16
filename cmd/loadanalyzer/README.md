LoadAnalyzer
===============

Usage
----------
#### Build

    make loadanalizer
   
#### Start insolard

    ./scripts/insolard/launch.sh
   
#### Start loadanalyzer

    ./bin/loadanalyzer -c=3 -r=1 --with_init

### Options

        -c concurrency
                Number of concurrent users. Default is one.

        -r repetitions
                Number of repetitions for one user. Default is one.

        -i input
                Path to file with initial data - references of members.

        -o output
                Path to output file (use - for STDOUT).

        --with_init
                Do initialization before run. In that case you don't have to provide input file, new members will be generated automatically. Default is false.
