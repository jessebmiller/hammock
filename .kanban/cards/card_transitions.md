---
last_moved_at: "2023-06-06T01:32:49.393960014+00:00"

[events]
  [events.created]
  time = "2023-06-06T01:32:49.394011154+00:00"
---
# Card transitions
---

## Start

Start working on the top card of the column to the right of the column
with transition = "start" if there is no card in that column,
interactively ask if they'd like to move the top card from To do (or
column configured with transition = "start" into the next column. In
addition, if it doesn't exist create and switch to a new branch for
the work (probably pulling main first etc.). If git is dirty tell the
user to commit changes etc.

Keep track of what card is being worked on and keep branches open for
them

## Start --card "<unique card prefix>"

Switch to working on the matched card. If the match is not unique
offer an interactive list of cards to work on. Move it to in progress
no matter where it was and do all the git stuff as above
