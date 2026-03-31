#!/usr/bin/env bash
set -euo pipefail

RUNS="${1:-100}"
shift || true
EXTRA_ARGS=("$@")

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

go build .

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

calc_stats() {
  local label="$1"
  local file_path="$2"

  awk -v label="$label" '
    {
      x[NR] = $1
      sum += $1
    }
    END {
      if (NR == 0) {
        printf("%s,%d,NaN,NaN\n", label, 0)
        exit
      }

      mean = sum / NR
      var = 0
      for (i = 1; i <= NR; i++) {
        d = x[i] - mean
        var += d * d
      }
      var /= NR
      std = sqrt(var)

      printf("%s,%d,%.12f,%.12f\n", label, NR, mean, std)
    }
  ' "$file_path"
}

run_continuous() {
  local objective="$1"
  local out_file="$TMP_DIR/continuous_${objective}.txt"

  go run . -mode continuous -objective "$objective" -runs "$RUNS" "${EXTRA_ARGS[@]}" \
    | awk '/^Run [0-9]+: Continuous EMAS \(/ { print $NF }' > "$out_file"

  calc_stats "continuous:${objective}" "$out_file"
}

run_tsp() {
  local out_file="$TMP_DIR/tsp.txt"

  go run . -mode tsp -runs "$RUNS" "${EXTRA_ARGS[@]}" \
    | awk '/^Run [0-9]+: EMAS result / { print $NF }' > "$out_file"

  calc_stats "tsp" "$out_file"
}

echo "problem,runs,mean,std"
run_continuous "sphere"
run_continuous "rastrigin"
run_continuous "rosenbrock"
run_tsp
