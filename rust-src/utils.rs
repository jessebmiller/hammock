use std::env;

pub fn get_editor() -> Option<String> {
    if let Ok(editor) = env::var("VISUAL") {
        return Some(editor);
    }
    if let Ok(editor) = env::var("EDITOR") {
        return Some(editor);
    }
    None
}
