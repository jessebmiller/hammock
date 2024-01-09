use super::find_kanban;
use super::board::Column;
use crate::utils::get_editor;
use chrono::{DateTime, Utc};
use gray_matter::engine::TOML;
use gray_matter::Matter;
use std::fs::{read_dir, File};
use std::io::{Read, Write};
use std::path::{Path, PathBuf};
use std::process::{exit, Command};
use tempfile::NamedTempFile;
use std::ffi::OsStr;

#[derive(Debug, Clone, Ord, PartialOrd, Eq, PartialEq)]
pub struct Card {
    pub file_path: PathBuf,
    pub headline: Option<String>,
    pub last_moved_at: Option<DateTime<Utc>>,
}

pub fn load_column_cards<P: AsRef<Path>>(column_path: P) -> anyhow::Result<Vec<Card>> {
    let cards = read_dir(&column_path)?.filter_map(|entry| {
        let entry = entry.ok()?;
        let file_path = entry.path();
        if file_path.extension() == Some(OsStr::new("md")) {
            load_card(&file_path).ok()
        } else {
            None
        }
    }).collect();

    Ok(cards)
}

fn load_card<P: AsRef<Path>>(path: P) -> anyhow::Result<Card> {
    let mut file = match File::open(&path) {
        Ok(file) => file,
        Err(_) => {
            eprintln!("Failed to open card file: {:?}", path.as_ref());
            exit(1);
        }
    };

    let mut contents = String::new();
    file.read_to_string(&mut contents).unwrap();

    let matter: Matter<TOML> = Matter::new();
    let parsed_card = matter.parse(&contents);
    let data = parsed_card.data.as_ref().unwrap();

    let last_moved_at_value = data["last_moved_at"].clone().as_string().ok();
    let last_moved_at = last_moved_at_value
        .map(|value| DateTime::parse_from_rfc3339(&value).ok())
        .flatten();

    let last_moved_at_utc = last_moved_at.map(|dt| dt.with_timezone(&Utc));

    Ok(Card {
        file_path: path.as_ref().to_path_buf(),
        headline: parsed_card.excerpt.map(|s| s.replace("# ", "").trim().to_string()),
        last_moved_at: last_moved_at_utc,
    })
}

impl Card {
    pub fn new(headline: Option<String>, column: &Column) -> anyhow::Result<Self> {
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

        let file_path = find_kanban()?.join(&column.dir_name).join(&file_name);
        println!("Saving card to {:?}", file_path);

        let mut card_file = File::create(&file_path)?;
        card_file.write_all(contents.as_bytes())?;

        Ok(Card {
            file_path,
            headline: parsed_card.excerpt.map(|excerpt| excerpt.replace("#", "").trim().to_owned()),
            last_moved_at: None,
        })
    }

    pub fn move_to_column(&mut self, column: &Column) -> anyhow::Result<()> {
        let mut file = File::open(&self.file_path)?;
        let mut contents = String::new();
        file.read_to_string(&mut contents)?;
        let matter: Matter<TOML> = Matter::new();
        let mut parsed_card = matter.parse(&contents);
        if let Some(data) = parsed_card.data.as_mut() {
            data["last_moved_at"] = Utc::now().to_rfc3339().into();
        }
        let mut card_file = File::create(&self.file_path)?;
        card_file.write_all(contents.as_bytes())?;
        let new_file_path = find_kanban()?.join(&column.dir_name).join(&self.file_path.file_name().unwrap());
        std::fs::rename(&self.file_path, &new_file_path)?;
        self.file_path = new_file_path;
        Ok(())
    }
}
