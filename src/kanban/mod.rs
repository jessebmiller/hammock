pub mod board;
pub mod card;
pub mod tui;

use std::path::PathBuf;

pub fn find_current_workspace_root() -> anyhow::Result<PathBuf> {
    let current_dir = std::env::current_dir()?;
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

pub fn find_kanban() -> anyhow::Result<PathBuf> {
    let kanban_path = find_current_workspace_root()?.join(".kanban");
    if !kanban_path.exists() {
        return Err(anyhow::anyhow!("No .kanban directory found in current workspace"));
    }
    Ok(kanban_path)
}

