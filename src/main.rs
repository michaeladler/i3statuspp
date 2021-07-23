mod config;
mod i3;

use anyhow::{anyhow, Result};
use log::{debug, error, trace};
use smol::{
    fs::File,
    future,
    io::{AsyncBufReadExt, AsyncReadExt, AsyncWriteExt, BufReader},
    prelude::*,
    process::{Command, Stdio},
    Unblock,
};

use config::Config;
use i3::ClickEvent;

static IPC_HEADER: &[u8] = b"{\"version\":1,\"click_events\":true}\n";

#[derive(Debug)]
enum Source {
    I3Status,
    Stdin,
}

#[derive(Debug)]
struct Container<T> {
    source: Source,
    data: T,
}

fn main() -> Result<()> {
    env_logger::init();
    smol::block_on(async {
        let cfg = parse_config().await?;

        let mut child = Command::new("i3status").stdout(Stdio::piped()).spawn()?;

        let stdout = child
            .stdout
            .take()
            .ok_or_else(|| anyhow!("Unable to access the child's stdout"))?;
        let mut status_reader = BufReader::new(stdout).lines().map(|data| Container {
            source: Source::I3Status,
            data,
        });

        let mut stdout = Unblock::new(std::io::stdout());

        // skip header line which is missing the click_events flag
        status_reader.next().await;
        // instead, send our own header
        stdout.write_all(IPC_HEADER).await?;

        let stdin = Unblock::new(std::io::stdin());
        let mut click_reader = BufReader::new(stdin).lines().map(|data| Container {
            source: Source::Stdin,
            data,
        });

        loop {
            // we prefer click events over i3status output
            let next = future::or(click_reader.next(), status_reader.next()).await;
            if let Some(next) = next {
                trace!("Received from {:?}: {:?}", next.source, next.data);
                match (next.source, next.data) {
                    (Source::I3Status, Ok(data)) => {
                        // just forward
                        stdout.write_all(data.as_bytes()).await?;
                    }
                    (Source::Stdin, Ok(data)) => {
                        // handle click event
                        let bytes = data.as_bytes();
                        let mut idx_left = None;
                        for (i, &item) in bytes.iter().enumerate() {
                            if item == b'{' {
                                idx_left = Some(i);
                                break;
                            }
                        }
                        if idx_left.is_none() {
                            continue;
                        }
                        let idx_left = idx_left.unwrap();

                        let mut idx_right = bytes.len();
                        while idx_right > idx_left {
                            idx_right -= 1;
                            if bytes[idx_right] == b'}' {
                                break;
                            }
                        }

                        let filtered_data = &bytes[idx_left..=idx_right];
                        if filtered_data.is_empty() {
                            continue;
                        }
                        let parse_result: Result<ClickEvent, serde_json::Error> =
                            serde_json::from_slice(filtered_data);
                        match parse_result {
                            Ok(ce) => {
                                debug!("Parsed click event: {:?}", ce);
                                // find matching rule
                                for rule in &cfg.rules {
                                    if rule.name.is_some() && rule.name != ce.name {
                                        continue;
                                    }
                                    if rule.instance.is_some() && rule.instance != ce.instance {
                                        continue;
                                    }
                                    debug!("Rule {} matches", rule.id);
                                    if let Some(action) = rule.actions.get(&ce.button.to_string()) {
                                        if let Some(action) = shlex::split(action) {
                                            debug!("Launching command: {:?}", action);
                                            Command::new(&action[0])
                                                .args(&action[1..])
                                                .stdout(Stdio::null())
                                                .spawn()?;
                                        }
                                    }
                                }
                            }
                            Err(e) => {
                                error!(
                                    "Failed to parse click event: data: {}, reason: {}",
                                    data, e
                                );
                            }
                        };
                    }
                    (_, Err(e)) => {
                        error!("No data available: {}", e);
                    }
                }
            }
        }
    })
}

async fn parse_config() -> Result<Config> {
    let dir = dirs_next::config_dir().ok_or_else(|| anyhow!("Unable to get config dir"))?;
    let mut file = File::open(dir.join("i3statuspp").join("config.json")).await?;
    let mut content = String::with_capacity(4096);
    file.read_to_string(&mut content).await?;

    let cfg: config::Config = serde_json::from_str(&content)?;
    Ok(cfg)
}
