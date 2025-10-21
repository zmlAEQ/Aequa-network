#!/usr/bin/env bash
set -euo pipefail

# go-fuzz-loop.sh PKG MINUTES LOGFILE
# Repeatedly runs `go test -fuzz=Fuzz` on a package in CHUNK windows until
# the given wall-clock minutes elapse. Aggregates logs into LOGFILE.

PKG="${1:?package path required}"
TOTAL_MIN="${2:?minutes required}"
LOGFILE="${3:-fuzz.out}"

CHUNK_MIN=5
CHUNK="${CHUNK_MIN}m"
END=$(( $(date +%s) + TOTAL_MIN*60 ))

echo "[fuzz-loop] package=${PKG} minutes=${TOTAL_MIN} chunk=${CHUNK}" | tee -a "$LOGFILE"

fail=0
iter=0
while [ "$(date +%s)" -lt "$END" ]; do
  iter=$((iter+1))
  echo "[fuzz-loop] iter=$iter start=$(date -Iseconds)" | tee -a "$LOGFILE"
  set +e
  go test "$PKG" -run '^$' -fuzz=Fuzz -fuzztime="$CHUNK" -timeout $((CHUNK_MIN+1))m >>"$LOGFILE" 2>&1
  rc=$?
  set -e
  if [ $rc -ne 0 ]; then
    echo "[fuzz-loop] iter=$iter rc=$rc (continuing)" | tee -a "$LOGFILE"
    fail=1
  fi
done

echo "[fuzz-loop] completed at $(date -Iseconds) fail=$fail" | tee -a "$LOGFILE"
exit $fail

