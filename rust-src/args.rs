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

    #[command(about = "Edit a card")]
    Edit { headline: String },

    #[command(about = "Show various objects in the current workspace")]
    Show {
        #[clap(subcommand)]
        object: ShowObject,
    },

    #[command(about = "Initialize a workspace")]
    Init { name: Option<String> },

    #[command(about = "Go to the priority workspace")]
    Go,
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
    #[command(about = "Show the matching card")]
    Card { headline: String },
}
