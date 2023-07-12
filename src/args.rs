use clap::{ValueEnum, Parser};

#[derive(Parser, Debug)]
#[clap(
    name = "hmm",
    version = "0.1.0",
    about = "Manages workspaces, including kanban boards"
)]
pub struct Args {
    #[clap(subcommand)]
    pub command: Option<Command>,
}

#[derive(Parser, Debug)]
pub enum Command {
    #[command(about = "Manage the current workspace kanban board")]
    Kanban,

    #[command(about = "Add a card to the leftmost column of the kanban board")]
    Card { headline: Option<String> },

    #[command(about = "Move a card to the left or right")]
    Move { headline: String, direction: Option<Direction> },

    #[command(about = "Manage the current workspace notes")]
    Notes,

    #[command(about = "Add a new note to the workspace")]
    Note { text: String },

    #[command(about = "Build and serve documentation for the current workspace")]
    Docs,

    #[command(about = "Show various objects in the current workspace")]
    Show {
        #[clap(subcommand)]
        object: ShowObject,
    },
}

#[derive(Parser, Debug, Clone, ValueEnum)]
pub enum Direction {
    Left,
    Right,
}

pub fn default_direction() -> Direction {
    Direction::Right
}

#[derive(Parser, Debug)]
pub enum ShowObject {
    #[command(about = "Show the current workspace kanban board")]
    Kanban,

    #[command(about = "Show the current workspace notes")]
    Notes,
}
