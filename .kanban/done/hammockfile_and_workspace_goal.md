---
last_moved_at: "2024-01-06T23:14:37.734601789+00:00"

[events]
  [events.created]
  time = "2024-01-06T23:14:37.734607122+00:00"
---
# Hammockfile and workspace goal
---

`Hammock` file keeps a priority project and a goal

```
priority_workspace = "name"
goal = "My goal right now is to ..."
```

Workspace `.conf.toml` specifies a project goal and shows it in the summary

Right now the priority summary is shown on every command. Maybe just the main
goal and then the name of the priority project and in progress and top of todo

Show the whole summary for the default command when not in a workspace

only show goal and priority project stuff on the default command
