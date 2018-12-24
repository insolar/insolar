#!/bin/bash

#for i in `seq 1 9`;
#do
#	./bin/benchmark -c=3 -r=100 -k=scripts/insolard/configs/root_member_keys.json &
#done

./bin/benchmark -c=3 -r=100 -k=scripts/insolard/configs/root_member_keys.json -u http://localhost:19191/api &
./bin/benchmark -c=3 -r=100 -k=scripts/insolard/configs/root_member_keys.json -u http://localhost:19192/api &
./bin/benchmark -c=3 -r=100 -k=scripts/insolard/configs/root_member_keys.json -u http://localhost:19193/api &
