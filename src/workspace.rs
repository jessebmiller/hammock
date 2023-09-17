use crate::kanban::board::{load_board_at, Board};
use crate::notes::{load_notes, Note};
use std::path::PathBuf;
use std::time::UNIX_EPOCH;
use std::fs::{create_dir_all, File};
use std::io::prelude::*;

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

pub fn find_current_workspace() -> Option<Workspace> {
    let home = std::env::var("HOME").unwrap();
    let root = PathBuf::from(home).join("work").canonicalize().ok()?;
    let mut current_dir = std::env::current_dir().unwrap().canonicalize().unwrap();
    if current_dir.starts_with(&root) && current_dir != root {
        while current_dir.parent() != Some(&root) {
            current_dir.pop();
        }
        Some(Workspace::new(current_dir))
    } else {
        None
    }
}

fn get_workspace_root() -> PathBuf {
    let home = std::env::var("HOME").unwrap();
    PathBuf::from(home).join("work").canonicalize().unwrap()
}

pub fn workspaces() -> Vec<Workspace> {
    let mut workspaces = Vec::new();
    let root = get_workspace_root();
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

const DEFAULT_KANBAN_CONF: &str = r#"
[[columns]]
name = "To Do"
dir_name = "to_do"
limit = 7
show_headlines_in_summary = true

[[columns]]
name = "In Progress"
dir_name = "in_progress"
limit = 1
show_headlines_in_summary = true

[[columns]]
name = "Done"
dir_name = "done"
display = 10
sort_by = "last_moved_at"
sort_order = "desc"
[[column]]
"#;

pub fn init_workspace(name: Option<String>) -> anyhow::Result<()> {
    let root = get_workspace_root();
    let path = match name {
        Some(name) => root.join(name),
        None => {
            let pwd = std::env::current_dir()?;
            if pwd.parent() == Some(&root) {
                pwd.clone()
            } else {
                println!("Cannot initialize workspace in {} must be in {}", pwd.display(), root.display());
                return Err(anyhow::anyhow!("Failed to initialize workspace"));
            }
        }
    };

    if !path.exists() {
        create_dir_all(&path)?;
    }
    if path.join(".kanban").exists() {
        println!("Workspace already initialized at {}", path.display());
        return Ok(());
    }
    create_dir_all(path.join(".kanban"))?;
    let mut conf_file = File::create(path.join(".kanban").join(".conf.toml"))?;
    conf_file.write_all(DEFAULT_KANBAN_CONF.as_bytes())?;
    create_dir_all(path.join(".kanban").join("to_do"))?;
    create_dir_all(path.join(".kanban").join("in_progress"))?;
    create_dir_all(path.join(".kanban").join("done"))?;
    println!("Initialized workspace at {}", path.display());
    Ok(())
}
