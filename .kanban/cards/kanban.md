---
last_moved_at = "2023-05-05T19:50:00Z"

[events]

  [events.created]
  time = "2023-05-05T19:50:00Z"

  [[events.moved_into_column]]
  name = "In progress"
  time = "2023-05-05T19:50:00Z"
---
# Build the kanban system
---

CLI and/or TUI kanban board stored in text files in git

A workspace may have work to manage. Maybe it's a project, maybe just
fun stuff you are working on but you might want a Kanban board and
some cards in the workspace.

a special `<workspace>/.kanban` folder gives your workspace a kanban
board.

## Design

[docs](../../docs/kanban.md)

# Tasks

- [ ] View the current board state
- [ ] Add a card
- [ ] View a card
- [ ] Edit a card
- [ ] Move a card into a column
- [ ] Change the order of cards in a column
