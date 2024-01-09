use std::path::PathBuf;

#[derive(Debug)]
pub struct Note {
    pub text: String,
    pub path: PathBuf,
    pub line: usize,
}

pub fn load_notes(project_dir: PathBuf) -> Vec<Note> {
    let path = project_dir.join("notes.md");
    let notes_string = match std::fs::read_to_string(path.clone()) {
        Ok(s) => s,
        Err(_) => return vec![],
    };
    // load all notes from notes.md. They are separated by a blank line
    let mut notes = Vec::new();
    // path is pwd/notes.md
    let mut note_text = String::new();
    for (i, line) in notes_string.lines().enumerate() {
        if line.is_empty() && !note_text.is_empty() {
            // we have reached the end of a note
            notes.push(Note {
                text: note_text,
                path: path.clone(),
                line: i,
            });
            note_text = String::new();
            continue;
        }
        if line.is_empty() {
            // we are still looking for the first line of a note
            continue;
        }
        // we are in the middle of a note
        note_text.push_str(line);
        note_text.push('\n');
    }
    if !note_text.is_empty() {
        // we have one last note to add
        notes.push(Note {
            text: note_text,
            path,
            line: notes_string.lines().count(),
        });
    }
    // TODO look for notes in other files, code commetns, etc.
    notes
}
