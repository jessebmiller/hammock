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
        let mut file = File::open(&file_path).unwrap();
        let mut contents = String::new();
        file.read_to_string(&mut contents).unwrap();
        let matter: Matter<TOML> = Matter::new();
        let parsed_card = matter.parse(&contents);
        Ok(Card {
            file_path: file_path.to_str().unwrap().to_string(),
            headline: parsed_card.excerpt,
            last_moved_at: parsed_card.data.as_ref()
                .and_then(|data| Some(data["last_moved_at"].clone()))
                .and_then(|last_moved_at| last_moved_at.as_string().ok())
                .and_then(
                    |last_moved_at|
                    DateTime::parse_from_rfc3339(&last_moved_at)
                        .and_then(|dt| Ok(dt.with_timezone(&Utc)))
                        .ok()
                )
        })
    }
}
