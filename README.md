# Hammock

CLI tools for hammock driven development, or whatever drives my
development.

# Install

I haven't done any work on distrubution and don't expect to for a while.
However you can of course clone this repo and do a 

`$ cargo install --path .`

from the root of the project. Of couse you'll need to
[install cargo](https://doc.rust-lang.org/cargo/getting-started/installation.html)
fisrt if you haven't already.

# Workspaces

Hammock recognizes workspaces under ~/work or `HAMMOCK_WORKSPACE_ROOT`

Every folder in the workspace root is a workspace.

Each workspace can be "worked on" with a workon command.

If there is a workspace at `~/work/my-project` you can go `workon
my-project` and it'll change directory to the project and run
`./workon.sh`

## Installing the workon command

source the shell_functions.sh file in your main rc file. For instance
add `source ~/work/hammock/src/shell_functions.sh` to `~/.zshrc`

# Kanban

A workspace may have work to manage. Maybe it's a project, maybe just
fun stuff you are working on but you might want a Kanban board and
some cards in the workspace.

a special `<workspace>/.kanban` folder gives your workspace a kanban
board.

see [kanban docs](./docs/kanban.md) for more info.

# Planning

One often wants some space to do some planning when working on
something. While you can plan in kanban cards, it's often nice to have
another place to do it.

Eventually maybe we'll build tools for managing the planning
space. But for now let's just use markdown files in a docs folder.

# Notes

It'll be nice to have a flexible way to deal with notes. These could
be //TODOs or //NOTEs in code or notes in a notes file, or notes thrown in on the
command line for later.

# Config

```
HAMMOCK_WORKSPACE_ROOT=~/work
```

# CLI

| command                            | effect                                  |
| `workspace > hmm`                  | Show the board and latest notes         |
| `workspace > hmm kanban`           | Launch the kanban TUI                   |
| `workspace > hmm notes`            | Launch the noted TUI                    |
| `workspace > hmm note "note text"` | Create a new note                       |
| `workspace > hmm docs`             | Build docs and open them in the browser |






