---
last_moved_at: "2023-06-10T15:57:39.620893250+00:00"

[events]
  [events.created]
  time = "2023-06-10T15:57:39.620907296+00:00"
---
# Manual edits of files should never break the system
---

for instance if I rename the workspace folder from jesse.works to
clutter hammock shouldn't care at all.

Currently the board keeps explicit paths (which break in this case) to
cards in each column. This should be changed to keep each card in a
folder named the same as the column.

That then means that the column name should only be specified by the
folder name, and the board.toml file probably goes away, or is broken
up into column specific config files within each folder.
