#!/usr/bin/env bash
# Show which languages were used to solve each AoC day, grouped by year.
set -euo pipefail

EXERCISES="$(cd "$(dirname "$0")" && pwd)/exercises"

# Known language directory names (everything else inside a day dir is data).
LANGS="go|py|python|bash|lua|rs|f77|fortran"

declare -A solved  # solved[year/day]="lang1 lang2 ..."

while IFS= read -r langdir; do
  lang=$(basename "$langdir")
  daydir=$(dirname "$langdir")
  day=$(basename "$daydir")
  year=$(basename "$(dirname "$daydir")")
  key="$year/$day"
  solved[$key]="${solved[$key]:+${solved[$key]} }$lang"
done < <(fd --type d --min-depth 3 --max-depth 3 -E '__pycache__' \
           --regex "^($LANGS)$" "$EXERCISES" | sort)

# Collect sorted years
mapfile -t years < <(printf '%s\n' "${!solved[@]}" | cut -d/ -f1 | sort -u)

for year in "${years[@]}"; do
  echo "── $year ──────────────────────────"
  # Collect days for this year, sorted by day number prefix
  mapfile -t days < <(printf '%s\n' "${!solved[@]}" | grep "^$year/" | \
                      sed "s|^$year/||" | sort -t- -k1,1n)
  for day in "${days[@]}"; do
    langs="${solved[$year/$day]}"
    printf "  %-35s %s\n" "$day" "$langs"
  done
done
