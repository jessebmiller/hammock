mod args;
mod kanban;
mod notes;
mod utils;
mod workspace;

use std::path::PathBuf;
use std::env;

use gray_matter::Matter;
use gray_matter::engine::TOML;

use args::{Args, Command};
use clap::Parser;
use kanban::card::Card;
use kanban::board::load_board;

pub fn find_current_workspace() -> anyhow::Result<PathBuf> {
    let current_dir = env::current_dir()?;
    let mut current_path = current_dir.as_path();

    if current_path.join(".kanban").is_dir() {
        return Ok(current_path.to_path_buf());
    }

    while let Some(parent) = current_path.parent() {
        let kanban_folder = parent.join(".kanban");
        if kanban_folder.is_dir() {
            return Ok(parent.to_path_buf());
        }
        current_path = parent;
    }

    Err(anyhow::anyhow!("No workspace found"))
}

fn main() {
    let args = Args::parse();
    match args.command {
        Some(Command::Kanban) => {
            kanban::tui::run().expect("Failed to run kanban TUI");
        }
        Some(Command::Card { headline }) => {
            let column = &load_board().unwrap().columns[0];
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
            match object {
                args::ShowObject::Card { headline } => {
                    load_board().unwrap().find_cards_by_headline_prefix(&headline).iter().for_each(|found| {
                        let contents = std::fs::read_to_string(found.card.file_path.clone()).unwrap();
                        let matter: Matter<TOML> = Matter::new();
                        let parsed_card = matter.parse(&contents);
                        println!("{}", termimad::inline(&parsed_card.content));
                    });
                }
                args::ShowObject::Kanban => {
                    println!("Showing kanban board (not implemented)");
                }
                args::ShowObject::Notes => {
                    println!("Showing notes (not implemented)");
                }
            }
        }
        Some(Command::Move { headline, direction }) => {
            load_board().unwrap().move_card(headline, direction);
        }
        None => {
            for w in &workspace::workspaces() {
                println!("{}", w.summary());
            }
        }
    }
}

