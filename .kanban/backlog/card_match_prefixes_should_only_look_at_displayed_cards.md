---
last_moved_at: "2024-01-07T00:00:28.670591829+00:00"

[events]
  [events.created]
  time = "2024-01-07T00:00:28.670609600+00:00"
---
# Card match prefixes should only look at displayed cards
---

For instance when I do `hmm` in this repo I get

```
hammock:
  Backlog:
    Design HERD style commands
    Automatic git branch when moving something to In Progress
    don't save blank cards
    Archive projects
    Display board columns side by side
    Distraction towards progress
    Do we want branching kanban boards
    ls that is Hammock aware
    Metrics
    Tab completion
    Track move events
    Card transitions
    Hammock is the TUI, uh is the CLI
    New cards go in to do
  To Do:
    Schedule work on calendar with notifications
    Show naive schedule
    Write post on my org system
  In Progress:
    Manual ordering of cards
```

But when I go `hmm move Manual` I get a conflict.
There is a card in Done that starts with `Manual`

I expect it know that I only mean to move things I
can see in the summary
