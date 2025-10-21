#!/usr/bin/env bash
set -euo pipefail

# Minimal health check for 4-node compose cluster
# Wait then probe /health on ports 4600..4603

echo "Waiting for cluster to stabilize..."
sleep 60

NODES=(aequa-node-0 aequa-node-1 aequa-node-2 aequa-node-3)
PORTS=(4600 4601 4602 4603)

echo "Starting health checks..."
for i in ${!NODES[@]}; do
  NODE_NAME=${NODES[$i]}
  PORT=${PORTS[$i]}
  URL="http://127.0.0.1:${PORT}/health"
  echo "Checking ${NODE_NAME} at ${URL}..."
  if ! curl --fail --silent --retry 5 --retry-delay 10 "$URL"; then
    echo "::error::Health check failed for ${NODE_NAME}"
    exit 1
  fi
  echo "${NODE_NAME} is healthy."
done

echo "All nodes are healthy. Cluster is up!"
exit 0