function workon {
    if [ -z "$1" ]; then
        echo "Usage: workon <workspace>"
        return 1
    fi

    # workspaces are in ~/work by default
    # that's configurable with $HAMMOCK_WORKSPACE_ROOT
    local workspace_root=${HAMMOCK_WORKSPACE_ROOT:-~/work}

    if [ ! -d "$workspace_root/$1" ]; then
        echo "Workspace $1 does not exist"
        return 1
    fi

    # change directory to the workspace
    cd "$workspace_root/$1"

    # if there's a .workon file, source it
    if [ -f .workon ]; then
        source .workon
    fi
}
