#!/usr/bin/env bash
# Print a TOML table of answers for an AoC year, e.g.:
#   1.1 = "1234"
#   1.2 = ""
#
# Usage: answers-toml.sh <path/to/year>
set -euo pipefail

year_dir="${1:?usage: answers-toml.sh <path/to/year>}"

for day_dir in "$year_dir"/*/; do
    info="${day_dir}info.json"
    [[ -f "$info" ]] || continue

    day=$(jq -r '.day' "$info")
    a=$(jq -c '.data.answers.a // "" | tostring' "$info")
    b=$(jq -c '.data.answers.b // "" | tostring' "$info")

    printf '%s.1 = %s\n' "$day" "$a"
    printf '%s.2 = %s\n' "$day" "$b"
done | sort -t. -k1,1n -k2,2n
