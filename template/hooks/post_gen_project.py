import subprocess
import os
import sys

# Get the day and year parameters from Cookiecutter context
day = "{{ cookiecutter.dayNumber }}"
year = "{{ cookiecutter.year }}"

# Run the command and capture the output
result = subprocess.run(["aocd", day, year], capture_output=True, text=True)

# Write the output to a file named 'input.txt'
output_file = os.path.join(os.getcwd(), "input.txt")
with open(output_file, "w") as file:
    file.write(result.stdout)

# truncate the output path a bit
dirname = os.path.dirname(output_file)
basename = os.path.basename(output_file)
truncated_dirname = os.path.join("...", *dirname.split(os.sep)[-2:])
truncated_path = os.path.join(truncated_dirname, basename)

# Verify if the file was written successfully
if os.path.exists(output_file):
    print(f"Output written to {truncated_path}")
else:
    print(f"Failed to write output to {truncated_path}")
    sys.exit(1)
