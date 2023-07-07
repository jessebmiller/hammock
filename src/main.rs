mod args;
mod kanban;
mod notes;
mod utils;
mod workspace;

use args::{Args, Command};
use clap::Parser;
use kanban::card::Card;
use kanban::board::load_board;

fn main() {
    let args = Args::parse();
    match args.command {
        Some(Command::Kanban) => {
            kanban::tui::run().expect("Failed to run kanban TUI");
        }
        Some(Command::Card { headline }) => {
            let column = load_board().unwrap().columns[0].name.clone();
            match Card::new(headline, column) {
                Ok(card) => {
                    println!(
                        "Created card: {}",
                        card.headline.unwrap_or("<empty headline>".to_string())
                    );
                }
                Err(e) => {
                    println!("Failed to create card: {}", e);
                }
            }
        }
        Some(Command::Notes) => {
            println!("Running Notes TUI (not implemented)");
        }
        Some(Command::Note { text }) => {
            println!("Adding note: {} (not implemented)", text);
        }
        Some(Command::Docs) => {
            println!("Building and serving Docs (not implemented)");
        }
        Some(Command::Show { object }) => {
            println!("Showing some object (not implemented) {:?}", object);
        }
        None => {
            default_command();
        }
    }
}

fn default_command() {
    for w in &workspace::workspaces() {
        println!("{}", w.summary());
    }
}
