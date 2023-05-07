use super::Card;
use serde::de::{self, Deserialize, Deserializer};
use std::path::PathBuf;

pub fn deserialize<'de, D>(deserializer: D) -> Result<Card, D::Error>
where
    D: Deserializer<'de>,
{
    let path = PathBuf::deserialize(deserializer)?;
    let card = Card::from_path(&path).map_err(de::Error::custom)?;
    Ok(card)
}

