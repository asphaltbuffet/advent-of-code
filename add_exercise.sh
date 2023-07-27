#!/bin/bash

# Check if both year and day arguments are provided
if [ $# -ne 2 ]; then
    echo "Usage: $0 <year> <day>"
    exit 1
fi

year=$1
day=$2

# Interval in seconds
interval=30

# File to store the timestamp of the last execution
timestamp_file="/tmp/aoc_timestamp.txt"  # Specify the desired path and filename

# Get the current timestamp
current_timestamp=$(date +%s)

# Read the last execution timestamp from the file (if it exists)
if [ -f "$timestamp_file" ]; then
    last_timestamp=$(cat "$timestamp_file")
else
    last_timestamp=0
fi

# Calculate the time elapsed since the last execution
elapsed_time=$((current_timestamp - last_timestamp))

# Check if the required interval has passed
if [ "$elapsed_time" -lt "$interval" ]; then
    echo "Script can only be executed once every 30s. Please wait for the next interval."
    exit 1
fi

# Fetch the puzzle page
url="https://adventofcode.com/$year/day/$day"
puzzle_page=$(curl -b session=$(cat ~/.config/aocd/token) -s "$url")

# Extract the puzzle name
puzzle_name=$(echo "$puzzle_page" | grep -oP '(?<=<h2>--- Day '$day': ).*?(?= ---</h2>)')

if [ -z "$puzzle_name" ]; then
    echo "Failed to retrieve puzzle name for year $year, day $day."
    exit 1
fi

#echo "Puzzle name for year $year, day $day: $puzzle_name"

echo $( cookiecutter template --no-input year=$year dayNumber=$day exerciseTitle="$puzzle_name" -o exercises/$year )

# Store the current timestamp as the last execution timestamp
echo "$current_timestamp" > "$timestamp_file"
