use chrono::{DateTime, Utc};
use serde::de::{Deserialize, Deserializer};
use std::path::PathBuf;
use gray_matter::Matter;
use gray_matter::engine::TOML;
use std::fs::File;
use std::io::Read;

#[derive(Debug)]
pub struct Card {
    pub file_path: String,
    pub headline: Option<String>,
    pub last_moved_at: Option<DateTime<Utc>>,
}

impl<'de> Deserialize<'de> for Card {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        let file_path: PathBuf = PathBuf::deserialize(deserializer)?;
        let mut file = File::open(&file_path).expect(
            format!("file ({}) not found", file_path.display()).as_str(),
        );
        let mut contents = String::new();
        file.read_to_string(&mut contents).unwrap();
        let matter: Matter<TOML> = Matter::new();
        let parsed_card = matter.parse(&contents);
        let data = parsed_card.data.as_ref().unwrap();
        let last_moved_at_value = data["last_moved_at"].clone().as_string().ok();
        let last_moved_at = last_moved_at_value.map(
            |value| DateTime::parse_from_rfc3339(&value).ok()
        ).flatten();
        let last_moved_at_utc = last_moved_at.map(|dt| dt.with_timezone(&Utc));
        Ok(Card {
            file_path: file_path.to_str().unwrap().to_string(),
            headline: parsed_card.excerpt,
            last_moved_at: last_moved_at_utc,
        })
    }
}
