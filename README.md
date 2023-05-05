# Hammock

CLI tools for hammock driven development, or whatever drives my
development.

# Workspaces

Hammock recognizes workspaces under (at ~/work) or under
`HAMMOCK_WORKSPACE_ROOT`

Every folder in the workspace root is a workspace.

Each workspace can be "worked on" with a workon command.

If there is a workspace at `~/work/my-project` you can go `workon
my-project` and it'll change directory to the project and run
`./workon.sh`

# Config

```
HAMMOCK_WORKSPACE_ROOT=~/work
```
