use crate::utils::{find_kanban, get_editor};
use chrono::{DateTime, Utc};
use gray_matter::engine::TOML;
use gray_matter::Matter;
use serde::de::{Deserialize, Deserializer};
use serde::ser::{Serialize, Serializer};
use std::fs::File;
use std::io::Read;
use std::io::Write;
use std::path::PathBuf;
use std::process::{exit, Command};
use tempfile::NamedTempFile;

use super::board::{load_board_from_file, save_board_to_file};

#[derive(Debug, Clone)]
pub struct Card {
    pub file_path: String,
    pub headline: Option<String>,
    pub last_moved_at: Option<DateTime<Utc>>,
}

impl Card {
    pub fn new(headline: Option<String>) -> anyhow::Result<Self> {
        let editor = get_editor().expect("No editor set set $VISUAL or $EDITOR");
        let temp_file_path = NamedTempFile::new()?.into_temp_path();
        let mut temp_file = File::create(&temp_file_path)?;
        let boilerplate = format!(
            r#"---
last_moved_at: "{}"

[events]
  [events.created]
  time = "{}"
---
# {}
---

Card content here.
"#,
            Utc::now().to_rfc3339(),
            Utc::now().to_rfc3339(),
            headline.unwrap_or_default()
        );
        temp_file.write_all(boilerplate.as_bytes())?;
        let editor_arg = format!("{}", temp_file_path.display());
        let status = Command::new(editor).arg(editor_arg).status()?;
        if !status.success() {
            println!("Editor exited with non-zero status code");
            exit(1);
        }
        let mut edited_temp_file = File::open(&temp_file_path)?;
        let mut contents = String::new();
        edited_temp_file.read_to_string(&mut contents)?;
        let matter: Matter<TOML> = Matter::new();
        let parsed_card = matter.parse(&contents);
        let file_name = match parsed_card.excerpt.clone() {
            Some(excerpt) => {
                excerpt
                    .replace("#", "")
                    .trim()
                    .replace(" ", "_")
                    .to_lowercase()
                    .to_owned()
                    + ".md"
            }
            None => {
                println!("No headline found in card");
                exit(1);
            }
        };

        let kanban_path = find_kanban()?;
        let file_path = kanban_path.join("cards").join(file_name);

        let mut card_file = File::create(&file_path)?;
        card_file.write_all(contents.as_bytes())?;

        Ok(Card {
            file_path: file_path.to_str().unwrap().to_string(),
            headline: parsed_card.excerpt,
            last_moved_at: None,
        })
    }

    pub fn add_to_board(self) -> anyhow::Result<Card> {
        let board_path = find_kanban()?.join("board.toml");
        let mut board = load_board_from_file(board_path.clone())?;
        board.columns[0].cards.insert(0, self.clone());
        save_board_to_file(&board, board_path)?;
        Ok(self)
    }
}

impl Serialize for Card {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        serializer.serialize_str(&self.file_path)
    }
}

impl<'de> Deserialize<'de> for Card {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        let file_path: PathBuf = PathBuf::deserialize(deserializer)?;

        // Open the card file listed in the board column card list
        let mut file = File::open(&file_path)
            .expect(format!("file ({}) not found", file_path.display()).as_str());
        let mut contents = String::new();
        file.read_to_string(&mut contents).unwrap();

        // Parse the front matter and excerpt
        let matter: Matter<TOML> = Matter::new();
        let parsed_card = matter.parse(&contents);
        let data = parsed_card.data.as_ref().unwrap();

        // If there is a last moved at date, parse it
        let last_moved_at_value = data["last_moved_at"].clone().as_string().ok();
        let last_moved_at = last_moved_at_value
            .map(|value| DateTime::parse_from_rfc3339(&value).ok())
            .flatten();

        // If there is a date, convert it to a DateTime<Utc>
        let last_moved_at_utc = last_moved_at.map(|dt| dt.with_timezone(&Utc));

        Ok(Card {
            file_path: file_path.to_str().unwrap().to_string(),
            headline: parsed_card.excerpt,
            last_moved_at: last_moved_at_utc,
        })
    }
}
