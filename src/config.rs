use std::collections::HashMap;

use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Config {
    pub general: General,
    pub rules: Vec<Rule>,
}

#[derive(Debug, Deserialize)]
pub struct General {
    pub i3statuscmd: String,
}

#[derive(Debug, Deserialize)]
pub struct Rule {
    /// Unique rule ID
    pub id: String,

    /// Name of the block, if set
    pub name: Option<String>,

    /// Instance of the block, if set
    pub instance: Option<String>,

    /// Map buttons to executable commands.
    // key: X11 button ID (for example 1 to 3 for left/middle/right mouse button)
    // value: a command
    pub actions: HashMap<String, String>,
}
