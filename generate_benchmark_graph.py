#!/usr/bin/env python3

import json
import matplotlib.pyplot as plt
import matplotlib.patches as patches
import os
import sys
import re

OUTPUT_FILE = sys.argv[1]
YEAR = sys.argv[2]

COLORS = {"Golang": "#00ADD8", "Python": "#3776AB"}

MAX_Y_VALUE = 1

exercise_dir_regex = re.compile("""(?m)^(\d{2})-([a-zA-Z]+)$""")

directories = []
path = os.path.join("exercises", YEAR)
for filename in os.listdir(path):
    if os.path.isdir(os.path.join(path, filename)) and exercise_dir_regex.match(
        filename
    ):
        directories.append(filename)

files = [os.path.join(x, "benchmark.json") for x in directories]

benchmark_data = {
    "Golang": {},
    "Python": {},
}  # adding dicts here sets the order of points being plotted

for filename in files:
    fpath = os.path.join(path, filename)
    try:
        f = open(fpath)
    except FileNotFoundError:
        print(f"Warning: missing file {fpath}")
        continue

    data = json.load(f)
    f.close()

    for language in data["implementations"]:
        x = benchmark_data.get(language, {})
        x[str(data["day"]) + ".1"] = data["implementations"][language]["part.1.avg"]
        x[str(data["day"]) + ".2"] = data["implementations"][language]["part.2.avg"]

        for key in x:
            if x[key] == 0:
                x[key] = None
        benchmark_data[language] = x

all_days = set()

for lang in benchmark_data:
    for key in benchmark_data[lang]:
        day = int(key.split(".", 1)[0])
        all_days.add(day)

figure = plt.figure(figsize=(25 / 2, 5))
axp1 = figure.add_subplot(1, 2, 1)
axp2 = figure.add_subplot(1, 2, 2, sharey=axp1)

axp1.axhline(y=15, color="#fc8080", linestyle="dotted")
axp2.axhline(y=15, color="#fc8080", linestyle="dotted")

for i, language in enumerate(benchmark_data):
    data = benchmark_data[language]
    part_one_times = []
    part_two_times = []
    days = []

    for key in data:
        day = int(key.split(".", 1)[0])
        if day not in days:
            days.append(day)

        if key.endswith(".1"):
            part_one_times.append(data[key])
        if key.endswith(".2"):
            part_two_times.append(data[key])

    color = COLORS.get(language)

    p1 = axp1.scatter(days, part_one_times, color=color)
    p2 = axp2.scatter(days, part_two_times, color=color)

    for i, day in enumerate(days):
        if i + 1 >= len(days):
            continue
        if days[i + 1] == day + 1:
            axp1.plot(
                (day, days[i + 1]),
                (part_one_times[i], part_one_times[i + 1]),
                "-",
                color=color,
            )
            axp2.plot(
                (day, days[i + 1]),
                (part_two_times[i], part_two_times[i + 1]),
                "-",
                color=color,
            )

figure.suptitle(f"Average {YEAR} exercise running time")
axp1.set_title("Part one")
axp2.set_title("Part two")


def do_auxillary_parts(axis):
    plt.sca(axis)
    plt.xticks(list(all_days), [str(y) for y in all_days])
    plt.ylabel("Running time (seconds)")
    plt.yscale("log")
    plt.xlabel("Day")
    plt.legend(
        handles=[patches.Patch(color=COLORS[label], label=label) for label in COLORS]
    )
    # plt.ylim([0, MAX_Y_VALUE])
    # plt.legend(legends)


do_auxillary_parts(axp1)
do_auxillary_parts(axp2)

plt.tight_layout()
plt.savefig(OUTPUT_FILE)
