#!/usr/bin/env bash

grep "BFT consensus passed" ./scripts/insolard/discoverynodes/*/output.log | python3 ./scripts/insolard/analyze_consensus/analyze_logs_consensus.py