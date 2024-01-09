# Implement card ranking

Cards should be arbitrarily rankable in the Kanban bard

They'll need to be rankable against all other card sources that end up in the
board.

keep the rank in metamatter (anything at the bottom of the file between lines of +++
see below

The tool should make ranking and changing ranks easy. It may rewrite all ranks
(though remembering where in the file to put it back could be hard so maybe we
just support metamatter at the bottom where it's most convinient for now)

so commands like

List the backlog in rank order with

[x] `hammock list`

Update card with rank x to rank y by shuffling around all the ranks to make room

[ ] `hammock move <x> <y>`

By default add new cards to the top of the list, in the future accept an
optional rank or explicit `bottom`

[ ] `hammock new`

Future work could also include adding functionality for other new thigns like
projects, or workspaces or other newly invented things. Or a quick card mode
where you can input the content inline (like in `git commit -m "commit
message"`)

+++
created_at=2024-01-07T12:34:00
hammock_type = "Card"
priority_rank = 1
+++
