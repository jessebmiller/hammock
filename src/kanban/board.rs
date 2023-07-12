use super::card::Card;
use crate::find_current_workspace;
use crate::kanban::card::load_column_cards;
use crate::args::{Direction, default_direction};
use serde::{Deserialize, Serialize};
use std::fs;
use std::path::{PathBuf, Path};

#[derive(Debug, Deserialize, Serialize)]
pub struct Board {
    pub columns: Vec<Column>,
    pub path: Option<PathBuf>,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct Column {
    pub name: String,
    pub dir_name: String,
    pub path: Option<PathBuf>,
    pub limit: Option<u32>,
    pub display: Option<u32>,
    pub sort_by: Option<String>,
    pub sort_order: Option<String>,
    pub show_headlines_in_summary: Option<bool>,
}

impl Board {
    pub fn move_card(self, headline: String, direction: Option<Direction>) {
        let mut found: Option<Card> = None;
        let mut found_in_column_index: Option<usize> = None;
        for (i, column) in self.columns.iter().enumerate() {
            for card in column.get_cards().unwrap() {
                if let Some(card_headline) = &card.headline {
                    if card_headline.starts_with(&headline) {
                        if found.is_some() {
                            println!("Found multiple cards with headline starting with {}", headline);
                            return;
                        }
                        found_in_column_index = Some(i);
                        found = Some(card);
                    }
                }
            }
        }

        let mut card = match found {
            Some(card) => card,
            None => {
                println!("No card found with headline starting with {}", headline);
                return;
            }
        };

        let column_index = match found_in_column_index {
            Some(i) => i,
            None => {
                println!("No column found with card with headline starting with {}", headline);
                return;
            }
        };

        match direction.unwrap_or(default_direction()) {
            Direction::Left => {
                if column_index == 0 {
                    println!("Card is already in leftmost column");
                    return;
                }
                match card.move_to_column(&self.columns[column_index - 1]) {
                    Ok(_) => {
                        println!(
                            "Moved \"{}\" into column {}",
                            card.headline.as_ref().unwrap(),
                            self.columns[column_index - 1].name,
                        );
                    }
                    Err(e) => {
                        println!("Failed to move card: {}", e);
                    }
                }
            }
            Direction::Right => {
                if column_index == self.columns.len() - 1 {
                    println!("Card is already in rightmost column");
                    return;
                }
                match card.move_to_column(&self.columns[column_index + 1]) {
                    Ok(_) => {
                        println!(
                            "Moved \"{}\" into column {}",
                            card.headline.as_ref().unwrap(),
                            self.columns[column_index + 1].name,
                        );
                    }
                    Err(e) => {
                        println!("Failed to move card: {}", e);
                    }
                }
            }
        }
    }
}

impl Column {
    pub fn get_cards(&self) -> anyhow::Result<Vec<Card>> {
        if let Some(path) = &self.path {
            load_column_cards(path)
        } else {
            Err(anyhow::anyhow!("No path for column"))
        }
    }
}

pub fn load_board() -> anyhow::Result<Board> {
    load_board_at(find_current_workspace()?)
}

pub fn load_board_at<P: AsRef<Path>>(path: P) -> anyhow::Result<Board> {
    let contents = fs::read_to_string(&path.as_ref().join(".kanban").join(".conf.toml"))?;
    let mut board: Board = toml::from_str(&contents)?;
    board.path = Some(path.as_ref().to_path_buf());
    for column in &mut board.columns {
        column.path = Some(path.as_ref().join(".kanban").join(&column.dir_name));
    }
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
