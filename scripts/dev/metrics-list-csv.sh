#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
PROM_SERVER=${PROM_SERVER:-"http://localhost:9090"}

# metadata API (experimental)

FETCH_CSV='.artifacts/fetch_series.csv'
CLEANED_CSV='scripts/insolar_metrics.csv'

echo "name, type, description" > ${FETCH_CSV}
curl -G -sS -g "$PROM_SERVER/api/v1/targets/metadata" \
    --data-urlencode 'match_target={job=~"light_material|heavy_material|virtual"}' \
    |  jq -r '.data[] | "\(.metric),\(.type),\(.help)"' \
    | sort | grep -v '^go_' | grep -v '^insolar_badger_' >> ${FETCH_CSV}
./scripts/dev/metrics-filter.py ${FETCH_CSV} > ${CLEANED_CSV}
echo ""
echo "open ${CLEANED_CSV}"