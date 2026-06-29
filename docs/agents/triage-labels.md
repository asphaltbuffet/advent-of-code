# Triage Labels

Triage is **not used** in this repo. Issues are not run through the triage
state machine, and the `triage` skill is out of scope here.

If you later decide to adopt triage, map the five canonical roles below to the
actual GitHub label strings you configure, and update the `## Agent skills`
block in `CLAUDE.md` to point skills at this file.

| Label in skills   | Meaning                                  |
| ----------------- | ---------------------------------------- |
| `needs-triage`    | Maintainer needs to evaluate this issue  |
| `needs-info`      | Waiting on reporter for more information |
| `ready-for-agent` | Fully specified, ready for an AFK agent  |
| `ready-for-human` | Requires human implementation            |
| `wontfix`         | Will not be actioned                     |
