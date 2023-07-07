use super::card::Card;
use super::find_kanban;
use crate::kanban::card::load_column_cards;
use serde::{Deserialize, Serialize};
use std::fs;
use std::path::Path;

#[derive(Debug, Deserialize, Serialize)]
pub struct Board {
    pub columns: Vec<Column>,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct Column {
    pub name: String,
    pub dir_name: String,
    pub limit: Option<u32>,
    pub display: Option<u32>,
    pub sort_by: Option<String>,
    pub sort_order: Option<String>,
    pub show_headlines_in_summary: Option<bool>,
}

impl Column {
    pub fn get_cards(&self) -> anyhow::Result<Vec<Card>> {
        load_column_cards(&self.dir_name)
    }
}

pub fn load_board() -> anyhow::Result<Board> {
    load_board_at(find_kanban()?)
}

pub fn load_board_at<P: AsRef<Path>>(path: P) -> anyhow::Result<Board> {
    let contents = fs::read_to_string(&path.as_ref().join("board").join(".conf.toml"))?;
    let board: Board = toml::from_str(&contents)?;
    Ok(board)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_load_board() {
        let board = load_board().unwrap();
        assert_eq!(board.columns.len(), 3);
        assert_eq!(board.columns[0].name, "To Do");
        assert_eq!(board.columns[1].name, "In Progress");
        assert_eq!(board.columns[2].name, "Done");
    }

    #[test]
    fn test_get_cards() {
        let board = load_board().unwrap();
        let cards = board.columns[0].get_cards();
        assert!(cards.is_ok());
    }
}
