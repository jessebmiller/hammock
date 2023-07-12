pub mod board;
pub mod card;
pub mod tui;

use std::path::PathBuf;
use crate::find_current_workspace;

pub fn find_kanban() -> anyhow::Result<PathBuf> {
    let kanban_path = find_current_workspace()?.join(".kanban");
    if !kanban_path.exists() {
        return Err(anyhow::anyhow!("No .kanban directory found in current workspace"));
    }
    Ok(kanban_path)
}

