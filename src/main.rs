mod args;
mod kanban;
mod workspace;
mod notes;
mod utils;

use args::{Args, Command};
use clap::Parser;
use kanban::card::Card;

fn main() {
    let args = Args::parse();
    match args.command {
        Some(Command::Kanban) => {
            println!("Running Kanban TUI (not finished)");
            kanban::tui::run().expect("Kanban TUI Failed");
        }
        Some(Command::Card{ headline }) => {
            match Card::new(headline) {
                Ok(card) => {
                    println!("Created card: {}", card.headline
                             .unwrap_or("<empty headline>".to_string()));
                }
                Err(e) => {
                    println!("Failed to create card: {}", e);
                }
            }
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
            default_command();
        }
    }
}

fn default_command() {
    for w in &workspace::workspaces() {
        println!("{}\n", w.summary());
    }
}
