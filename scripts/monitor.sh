#!/bin/bash
cd scripts
docker-compose down
docker-compose up -d

docker-compose ps
echo "# Grafana: http://localhost:3000 admin:pass"
echo "# Jaeger: http://localhost:16686"
echo "# Prometheus: http://localhost:9090/targets"
