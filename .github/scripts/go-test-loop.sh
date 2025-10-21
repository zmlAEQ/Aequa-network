#!/usr/bin/env bash
set -euo pipefail

# go-test-loop.sh PKG MINUTES LOGFILE [EXTRA...]
# Repeatedly runs `go test` on a package in ~1m windows until wall-clock MINUTES elapse.

PKG="${1:?package path required}"
TOTAL_MIN="${2:?minutes required}"
LOGFILE="${3:-gotest.out}"
shift 3 || true
EXTRA=("$@")

ONE_MIN=1
END=$(( $(date +%s) + TOTAL_MIN*60 ))

echo "[go-test-loop] package=${PKG} minutes=${TOTAL_MIN}" | tee -a "$LOGFILE"

fail=0
iter=0
while [ "$(date +%s)" -lt "$END" ]; do
  iter=$((iter+1))
  echo "[go-test-loop] iter=$iter start=$(date -Iseconds)" | tee -a "$LOGFILE"
  set +e
  go test "$PKG" -count=1 -timeout 3m "${EXTRA[@]}" >>"$LOGFILE" 2>&1
  rc=$?
  set -e
  if [ $rc -ne 0 ]; then
    echo "[go-test-loop] iter=$iter rc=$rc (continuing)" | tee -a "$LOGFILE"
    fail=1
  fi
done

echo "[go-test-loop] completed at $(date -Iseconds) fail=$fail" | tee -a "$LOGFILE"
exit $fail

