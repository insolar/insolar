## Profiling insolar network

Start insolar network:

    ./insolar-scripts/insolard/launchnet.sh -g

Check that all nodes are in the complete network state:

    ./insolar-scripts/insolard/check_status.sh

Start profiler:

    ./insolar-scripts/insolard/profile.sh

As soon as profiler collects statistics (default 30s), web pages with profile info will be opened for each node.
