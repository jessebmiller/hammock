---
last_moved_at = "2023-06-06T00:29:39.731738086+00:00"

[events]
  [events.created]
  time = "2023-06-06T00:29:39.731738086+00:00"
---
# Design HERD style commands
---

Since the tool is aware of the board, it can be used to create
branches, run tests, "integrate", make a branch to work on a
card. Manage environment variables etc.

Outline some options there and add cards for them

rely on the workspace configuration for the commands,

Expect and support 3 commands
* develop: for setting up and watching a development environment
* integrate: for integrating code with other developers
* release: for releasing to audiences

develop should do everything you're supposed to do. Like pulling
others changes, watching tests, setting up local databases with good
test data etc.

integrate should do everything you're supposed to do for the
continuous integration cycle. pull from mainline, test locally, push
to CI server, pass CI, then and only then merge into
mainline. Deploying to environments for inspection and debugging all
along the way.

Release exposes a given deployment (created by an integration) to an
audience (URL, or sets permissions, or however you want to do that)
