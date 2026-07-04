# Task Completion

After implementing a solution or adding a language:

1. `elf solve <exercise-path> --lang=<lang> --plain` — must show PASS for all tests and recorded answers.
2. `elf benchmark <exercise-path> --plain`
3. `elf analyze <exercise-path>` — writes day `run-times.png`
4. `elf analyze exercises/<year>` — updates year aggregate graph
5. Update day README (approach prose + `elf solve` output block + run-times.png embed).
6. Update year README row (title link, stars, language links).
7. `jj commit -m "<message>"`

For visualizations: run `elf visualize`, convert with `magick <img> -colorspace Gray <out>` to verify grayscale, then add `## Visualization` section to day README.
