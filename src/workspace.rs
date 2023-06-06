use crate::kanban::board::{load_board_from_file, Board};
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

        let board = load_board_from_file(&root.join(".kanban/board.toml")).ok();
        Self {
            root: root.clone(),
            name,
            updated_at,
            board,
            notes: load_notes(root),
        }
    }

    pub fn summary(&self) -> String {
        let mut summary = String::new();
        summary.push_str(&format!("{}: ", self.name));

        // show card headlines in To do and In progress
        if let Some(board) = &self.board {
            if let Some(c) = board.columns.iter().find(|c| c.name == "To do") {
                summary.push_str(&format!("\n  To do:"));
                for c in &c.cards {
                    if let Some(headline) = &c.headline {
                        summary.push_str(&format!("\n    {}", headline));
                    }
                }
            }
            if let Some(c) = board.columns.iter().find(|c| c.name == "In progress") {
                summary.push_str(&format!("\n  In progress:"));
                for c in &c.cards {
                    if let Some(headline) = &c.headline {
                        summary.push_str(&format!("\n    {}", headline));
                    }
                }
            }
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
