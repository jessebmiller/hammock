use chrono::{DateTime, Utc};
use serde::de::{Deserialize, Deserializer};
use std::path::PathBuf;
use gray_matter::Matter;
use gray_matter::engine::TOML;
use std::fs::File;
use std::io::Read;
use crate::utils::{find_kanban, get_editor};
use std::io::Write;
use std::process::{Command, exit};
use tempfile::NamedTempFile;

#[derive(Debug)]
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
        let boilerplate = format!(r#"---
[events]
  [events.created]
  time = "{}"
---
# {}
---

Card content here.
"#,
            Utc::now().to_rfc3339(),
            headline.unwrap_or_default());
        temp_file.write_all(boilerplate.as_bytes())?;
        let editor_arg = format!("{}", temp_file_path.display());
        let status = Command::new(editor)
            .arg(editor_arg)
            .status()?;
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
                    .to_owned() + ".md"
            },
            None => {
                println!("No headline found in card");
                exit(1);
            }
        };

        let kanban_path = find_kanban()?;
        let file_path = kanban_path.join(file_name);

        let mut card_file = File::create(&file_path)?;
        card_file.write_all(contents.as_bytes())?;

        println!("adding card to board NOT IMPLEMENTED!");

        Ok(Card {
            file_path: file_path.to_str().unwrap().to_string(),
            headline: parsed_card.excerpt,
            last_moved_at: None,
        })
    }
}

impl<'de> Deserialize<'de> for Card {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        let file_path: PathBuf = PathBuf::deserialize(deserializer)?;

        // Open the card file listed in the board column card list
        let mut file = File::open(&file_path).expect(
            format!("file ({}) not found", file_path.display()).as_str(),
        );
        let mut contents = String::new();
        file.read_to_string(&mut contents).unwrap();

        // Parse the front matter and excerpt
        let matter: Matter<TOML> = Matter::new();
        let parsed_card = matter.parse(&contents);
        let data = parsed_card.data.as_ref().unwrap();

        // If there is a last moved at date, parse it
        let last_moved_at_value = data["last_moved_at"].clone().as_string().ok();
        let last_moved_at = last_moved_at_value.map(
            |value| DateTime::parse_from_rfc3339(&value).ok()
        ).flatten();

        // If there is a date, convert it to a DateTime<Utc>
        let last_moved_at_utc = last_moved_at.map(|dt| dt.with_timezone(&Utc));

        Ok(Card {
            file_path: file_path.to_str().unwrap().to_string(),
            headline: parsed_card.excerpt,
            last_moved_at: last_moved_at_utc,
        })
    }
}

