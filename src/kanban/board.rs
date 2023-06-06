use super::card::Card;
use serde::{Deserialize, Serialize};
use std::{fs, path::Path};

#[derive(Debug, Deserialize, Serialize)]
pub struct Board {
    pub columns: Vec<Column>,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct Column {
    pub name: String,
    pub limit: Option<u32>,
    pub display: Option<u32>,
    pub sort_by: Option<String>,
    pub sort_order: Option<String>,
    pub cards: Vec<Card>,
}

pub fn load_board_from_file<P: AsRef<Path>>(file_path: P) -> anyhow::Result<Board> {
    let contents = fs::read_to_string(file_path)?;
    let board: Board = toml::from_str(&contents)?;
    Ok(board)
}

pub fn save_board_to_file<P: AsRef<Path>>(board: &Board, file_path: P) -> anyhow::Result<()> {
    let contents = toml::to_string(board)?;
    fs::write(file_path, contents)?;
    Ok(())
}
