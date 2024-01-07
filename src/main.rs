mod args;
mod kanban;
mod notes;
mod utils;
mod workspace;

//use std::path::PathBuf;
//use std::env;

use gray_matter::Matter;
use gray_matter::engine::TOML;
use serde::Deserialize;
use std::fs;
use std::path::PathBuf;

use args::{Args, Command};
use clap::Parser;
use kanban::card::Card;
use kanban::board::load_board;
use workspace::{get_workspace_root, find_current_workspace, Workspace};

use crate::utils::get_editor;

#[derive(Debug, Deserialize)]
struct Hammock {
    priority_workspace: String,
    goals: String,
}

impl Hammock {
    fn new() -> Self {
        let hammock_file = get_workspace_root().join("Hammock");
        let config = fs::read_to_string(&hammock_file).unwrap();
        return toml::from_str(&config).unwrap();
    }
}

fn main() {
    let hammock = Hammock::new();
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
        Some(Command::Init { name }) => {
            workspace::init_workspace(name).expect("Failed to initialize workspace");
        }
        None => {
            println!("\n*** GOAL\n{}", hammock.goals);
            match find_current_workspace() {
                Some(workspace) => {
                    println!("\n{}", workspace.summary());
                    return;
                }
                None => {
                    let priority_workspace = Workspace::new(get_workspace_root().join(hammock.priority_workspace));
                    println!("\n*** Priority\n{}", priority_workspace.summary());
                }
            }
        }
    }
}

