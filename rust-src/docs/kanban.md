# Kanban

Workspaces can add a kanban board by creating a `board.toml` file in a
`.kanban` folder in the root of the workspace.

The `board.toml` file should contain a list of column objects.

Columns have the following fields

```
name: string               // Required. The column display name
limit: int                 // Optional. Number of cards to mark as "over the limit"
display: int               // Optional. Number of cards to display
sort_by: string            // Optional. Card property to sort by
sort_order: "asc" | "desc" // Optional.
```

For example:

```
[[columns]]
name = "To do"
limit = 7

[[columns]]
name = "In progress"
limit = 1

[[columns]]
name = "Review"
limit = 3

[[columns]]
name = "Done"
display = 10
sort_by = "last_moved_into_current_column_at"
sort_order = "desc"
```

# Basic Features

* View the current board state
* View a card
* Edit a card
* Move a card into a column
* Change the order of cards in a column

# Other features

* Archive a card
* View archived cards
* Search cards

# Reporting

* View a cumulative flow diagram across a given date range
  * Track the history of a card. when it went into each column. etc.
* View rolling lead time and cycle time metrics
  * 50, 90, 99th percentile
  * 7, 30, 90 days
 rolling average

# Architecture

## View the current board state

`workspace/ $ hmm kanban`

Should print the current board to the terminal

1. Read .kanban/board.toml
1. scan all cards, putting the first heading of each into the correct column
  1. identify which column based on most recent "moved_into_column" event.
1. Display it in a nice way. Show the name so we can use the view command below

## View a card

`workspace/ $ hmm show card my-card`

Display the contents of the card at `workspace/.kanban/cards/my-card.md`

Just the contents, not the frontmatter

## Edit a card

`workspace/ $ hmm edit card my-card`

run `emacsclient -n workspace/.kanban/cards/my-card.md` to open it in
the running emacs window or open a new emacs window if that doesn't
work

## Move a card into a column

`workspace/ $ hmm move my-card "My column"`

Add an event to the frontmatter of `workspace/.kanban/cards/my-card.md`

```
moved_into_column
name: "My column"
time: [datetime]
```

## Change the order of cards in a column

