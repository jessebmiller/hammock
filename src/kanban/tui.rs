use crossterm::{
    event::{poll, read, Event, KeyCode},
    execute, terminal,
};
use std::{io, time::Duration};
use tui::{
    backend::CrosstermBackend,
    layout::{Constraint, Direction, Layout},
    style::{Color, Style},
    text::Span,
    widgets::{Block, Borders, List, ListItem},
    Terminal,
};

use super::board::{load_board, Board};

pub fn run() -> Result<(), anyhow::Error> {
    let stdout = io::stdout();
    let backend = CrosstermBackend::new(stdout);
    let mut terminal = Terminal::new(backend)?;

    terminal.hide_cursor()?;
    terminal::enable_raw_mode()?;
    execute!(terminal.backend_mut(), terminal::EnterAlternateScreen)?;

    let board: Board = load_board()?;

    loop {
        terminal.draw(|rect| {
            let size = rect.size();

            let layout = Layout::default()
                .direction(Direction::Horizontal)
                .margin(1)
                .constraints([Constraint::Percentage(30), Constraint::Percentage(70)].as_ref())
                .split(size);

            let block = Block::default().borders(Borders::ALL);
            rect.render_widget(block, layout[0]);

            let list_items: Vec<ListItem> = board.columns[0]
                .get_cards()
                .expect("Failed to get cards")
                .iter()
                .map(|card| {
                    ListItem::new(Span::styled(
                        card.headline
                            .clone()
                            .unwrap_or("No headline (no excerpt)".to_string()),
                        Style::default().fg(Color::Yellow),
                    ))
                })
                .collect::<Vec<ListItem>>();
            let list = List::new(list_items).block(Block::default().borders(Borders::ALL));
            rect.render_widget(list, layout[1]);
        })?;

        if poll(Duration::from_millis(500))? {
            let event = read()?;
            match event {
                Event::Key(key_event) => {
                    if key_event.code == KeyCode::Esc {
                        break;
                    }
                }
                _ => {}
            }
        }
    }

    execute!(terminal.backend_mut(), terminal::LeaveAlternateScreen)?;
    terminal::disable_raw_mode()?;
    terminal.show_cursor()?;
    terminal::disable_raw_mode()?;
    Ok(())
}
