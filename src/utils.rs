use std::env;
use std::path::{Path, PathBuf};

pub fn get_editor() -> Option<String> {
    if let Ok(editor) = env::var("VISUAL") {
        return Some(editor);
    }
    if let Ok(editor) = env::var("EDITOR") {
        return Some(editor);
    }
    None
}

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
