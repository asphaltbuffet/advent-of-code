# Suggested Commands

All `elf` commands run from repo root with exercise paths relative to repo root (not from inside exercise dirs).

```
elf download <puzzle-url> --lang <go|rs|py>      # scaffold (re-run with new --lang to add language)
elf solve <path> --lang=<lang> --plain            # run example tests then solve real input
elf solve <path> --lang=<lang> --plain -X         # skip example tests, check recorded answer only
elf benchmark <path> --plain                      # timing per part per language
elf analyze <path>                                # day path → day run-times.png; year dir → year graph
elf visualize <path> --lang=go --outdir=<abs-dir> # invoke Exercise.Vis
```

VCS:
```
jj commit -m "solve <year> day <N> in <lang>"
jj commit -m "add <lang> to <year> day <N>"
```
