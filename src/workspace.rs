use crate::kanban::board::{load_board_at, Board};
use crate::notes::{load_notes, Note};
use std::path::PathBuf;
use std::time::UNIX_EPOCH;

#[derive(Debug)]
pub struct Workspace {
    pub root: PathBuf,
    pub name: String,
    pub updated_at: u64,
    pub board: Option<Board>,
    pub notes: Vec<Note>,
}

impl Workspace {
    pub fn new(root: PathBuf) -> Self {
        let name = root.file_name().unwrap().to_str().unwrap().to_string();
        let updated_at = std::fs::metadata(&root)
            .unwrap()
            .modified()
            .unwrap()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs();

        let board = load_board_at(&root);

        Self {
            root: root.clone(),
            name,
            updated_at,
            board: board.ok(),
            notes: load_notes(root),
        }
    }

    pub fn summary(&self) -> String {
        let mut summary = String::new();
        summary.push_str(&format!("{}:\n", self.name));

        // show card headlines for columns set to show headlines in summary
        if let Some(board) = &self.board {
            board.columns.iter().filter(|c| c.show_headlines_in_summary.unwrap_or(false)).for_each(|c| {
                summary.push_str(&format!("  {}:\n", c.name));
                match c.get_cards() {
                    Ok(cards) => {
                        cards.iter().for_each(|c| {
                            summary.push_str(&format!(
                                "    {}\n",
                                c.clone().headline.unwrap_or("<missing headline>".to_string())
                            ));
                        });
                    }
                    Err(e) => {
                        println!("Error loading cards for column {}: {}", c.name, e);
                    }
                }
            });
        }

        // show up to 3 notes
        let show_notes_count = 3;
        if !self.notes.is_empty() {
            summary.push_str(&format!("\n  Notes:"));
            for note in self.notes.iter().take(show_notes_count) {
                let indented_note = note.text.replace("\n", "\n    ");
                summary.push_str(&format!("\n    {}", indented_note));
            }
        }

        summary
    }
}

pub fn workspaces() -> Vec<Workspace> {
    let mut workspaces = Vec::new();
    let home = std::env::var("HOME").unwrap();
    let root = PathBuf::from(home).join("work");
    for entry in root.read_dir().unwrap() {
        let entry = entry.unwrap();
        let path = entry.path();
        if path.is_dir() {
            workspaces.push(Workspace::new(path));
        }
    }
    workspaces.sort_by(|a, b| b.updated_at.cmp(&a.updated_at));
    workspaces
}
