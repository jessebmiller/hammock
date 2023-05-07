use clap::Parser;

#[derive(Parser, Debug)]
#[clap(name = "hmm", version = "0.1.0", about = "Manages workspaces, including kanban boards")]
pub struct Args {
    #[clap(subcommand)]
    pub command: Option<Command>,
}

#[derive(Parser, Debug)]
pub enum Command {
    #[command(about = "Manage the current workspace kanban board")]
    Kanban,
    #[command(about = "Manage the current workspace notes")]
    Notes,
    #[command(about = "Add a new note to the workspace")]
    Note { text: String },
    #[command(about = "Build and serve documentation for the current workspace")]
    Docs,
}

