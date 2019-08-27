#!/usr/bin/env python3
import csv
import sys
seen = {}

with open(sys.argv[1], 'r') if len(sys.argv) > 1 else sys.stdin as f:
    csv_in = csv.DictReader(f)
    csv_out = csv.DictWriter(sys.stdout, fieldnames=csv_in.fieldnames)
    csv_out.writeheader()
    for row in csv_in:
        if seen.get(row['name']):
            continue
        # print(row)
        seen[row['name']] = True
        csv_out.writerow(row)
