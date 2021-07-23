use std::collections::HashSet;

use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct ClickEvent {
    /// Name of the block, if set
    pub name: Option<String>,
    /// Instance of the block, if set
    pub instance: Option<String>,
    /// X11 button ID (for example 1 to 3 for left/middle/right mouse button)
    pub button: u8,
    // An array of the modifiers active when the click occurred. The order in which modifiers are listed is not guaranteed.
    pub modifiers: HashSet<String>,

    /// X11 root window coordinates where the click occurred
    pub x: u32,
    pub y: u32,

    // Coordinates where the click occurred, with respect to the top left corner of the block
    pub relative_x: u32,
    pub relative_y: u32,

    // Coordinates relative to the current output where the click occurred
    pub output_x: u32,
    pub output_y: u32,

    //Width and height (in px) of the block
    pub width: u32,
    pub height: u32,
}
