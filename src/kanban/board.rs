use serde::Deserialize;
use std::{fs, path::Path};
use super::card::Card;

#[derive(Debug, Deserialize)]
pub struct Board {
    pub columns: Vec<Column>,
}

#[derive(Debug, Deserialize)]
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

