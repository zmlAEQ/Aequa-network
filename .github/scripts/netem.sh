#!/usr/bin/env bash
set -euo pipefail
# netem.sh cmd
# cmds:
#   init           : ensure iproute2 in netem-* containers
#   apply MODE     : MODE in {delay,loss,reorder,partition}
#   clear          : clear qdisc rules
#   random_tick    : randomly apply or clear one mode

nodes=(0 1 2 3)

_install_once() {
  local n=$1
  docker exec netem-$n sh -lc 'apk add --no-cache iproute2 >/dev/null 2>&1 || true'
}

apply_mode() {
  local mode=$1
  for n in "${nodes[@]}"; do
    case "$mode" in
      delay)     docker exec netem-$n sh -lc 'tc qdisc replace dev eth0 root netem delay 100ms 20ms' ;;
      loss)      docker exec netem-$n sh -lc 'tc qdisc replace dev eth0 root netem loss 2% 25%' ;;
      reorder)   docker exec netem-$n sh -lc 'tc qdisc replace dev eth0 root netem reorder 25% 50% delay 20ms' ;;
      partition) docker exec netem-$n sh -lc 'tc qdisc replace dev eth0 root netem loss 100%' ;;
      *) echo "unknown mode: $mode"; exit 1 ;;
    esac
  done
}

clear_mode() {
  for n in "${nodes[@]}"; do docker exec netem-$n sh -lc 'tc qdisc del dev eth0 root || true'; done
}

case "${1:-}" in
  init)
    for n in "${nodes[@]}"; do _install_once $n; done ;;
  apply)
    shift; apply_mode "${1:?mode required}" ;;
  clear)
    clear_mode ;;
  random_tick)
    # 0:clear 1:delay 2:loss 3:reorder 4:partition
    m=$(( RANDOM % 5 ))
    case "$m" in
      0) clear_mode ;;
      1) apply_mode delay ;;
      2) apply_mode loss ;;
      3) apply_mode reorder ;;
      4) apply_mode partition ;;
    esac ;;
  *) echo "usage: $0 {init|apply MODE|clear|random_tick}"; exit 1 ;;
 esac