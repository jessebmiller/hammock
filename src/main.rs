mod args;
mod kanban;
mod notes;
mod utils;
mod workspace;

//use std::path::PathBuf;
//use std::env;

use gray_matter::Matter;
use gray_matter::engine::TOML;

use args::{Args, Command};
use clap::Parser;
use kanban::card::Card;
use kanban::board::load_board;
use workspace::find_current_workspace;

use crate::utils::get_editor;

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
        Some(Command::Edit { headline }) => {
            let found = load_board().unwrap().find_cards_by_headline_prefix(&headline);
            if found.len() != 1 {
                println!("Found {} cards matching headline prefix '{}'", found.len(), headline);
                return;
            }
            let card = &found[0].card;
            let editor = get_editor().expect("No editor set set $VISUAL or $EDITOR");
            std::process::Command::new(editor.clone())
                .arg(card.file_path.clone())
                .status()
                .expect(format!("Failed to open editor {}", editor).as_str());
        }
        Some(Command::Move { headline, direction }) => {
            load_board().unwrap().move_card(headline, direction);
        }
        None => {
            match find_current_workspace() {
                Some(workspace) => {
                    println!("{}", workspace.summary());
                    return;
                }
                None => {
                    println!("not in a workspace, summarizing all workspaces");
                }
            }
            for w in &workspace::workspaces() {
                println!("{}", w.summary());
            }
        }
    }
}

