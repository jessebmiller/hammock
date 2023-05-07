mod args;
mod kanban;

use args::{Args, Command};
use clap::Parser;

fn main() {
    let args = Args::parse();
    match args.command {
        Some(Command::Kanban) => {
            println!("Running Kanban TUI (not finished)");
            kanban::tui::run().expect("Kanban TUI Failed");
        }
        Some(Command::Notes) => {
            println!("Running Notes TUI (not implemented)");
        }
        Some(Command::Note{ text }) => {
            println!("Adding note: {} (not implemented)", text);
        }
        Some(Command::Docs) => {
            println!("Building and serving Docs (not implemented)");
        }
        None => {
            println!("Running default command (not implemented)");
        }
    }
}
