## Profiling insolar network

Start insolar network:

    ./scripts/insolard/launchnet.sh -g

Check that all nodes are in the complete network state:

    ./scripts/insolard/check_status.sh

Start benchmark:

    ./bin/benchmark -c=10 -r=1000 -k=scripts/insolard/configs/root_member_keys.json

Start profiler:

    ./scripts/insolard/profile.sh

As soon as profiler collects statistics (default 30s), web pages with profile info will be opened for each node.