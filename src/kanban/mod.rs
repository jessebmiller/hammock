pub mod board;
pub mod card;
pub mod tui;

use std::path::PathBuf;
use std::env;

pub fn find_kanban() -> anyhow::Result<PathBuf> {
    let current_dir = env::current_dir()?;
    let mut current_path = current_dir.as_path();

    if current_path.join(".kanban").is_dir() {
        return Ok(current_path.join(".kanban").to_path_buf());
    }

    while let Some(parent) = current_path.parent() {
        let kanban_folder = parent.join(".kanban");
        if kanban_folder.is_dir() {
            return Ok(kanban_folder.to_path_buf());
        }
        current_path = parent;
    }

    anyhow::bail!("Could not find .kanban folder")
}

