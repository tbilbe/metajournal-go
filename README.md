# Metajournal-GO

Golang TUI application to track my wrok exploits at cinch.
This data will saved into the brag repo and used alongside the commits history to track my work progress

Built with Charm’s BubbleTea.

## ✨ Features

    •	Daily or Weekly journal mode
    •	Dynamic multi-item inputs
    •	Styled markdown preview
    •	Writes structured markdown with frontmatter
    •	Terminal-based, fast and distraction-free

Installation

1. Download the latest binary

(Or build yourself from source — see below)

```bash
# Clone the repo
git clone https://github.com/YOUR_GITHUB_USERNAME/metajournal-go.git
cd metajournal-go

# Build the binary
go build -o metajournal

# Move to a system-wide binary folder
sudo mv metajournal /usr/local/bin/
```

2. Set your save path

This tool uses an environment variable to know where to save your journal files.

Choose a location where your journals should live, e.g.:

```bash
export METAJOURNAL_SAVE_PATH="$HOME/dev/work/professional-brag/data/journal"
```

To make it permanent, add it to your shell configuration:
(if using zsh)

```bash
echo 'export METAJOURNAL_SAVE_PATH="$HOME/dev/work/professional-brag/data/journal"' >> ~/.zshrc
source ~/.zshrc
```

(If you use bash, replace .zshrc with .bashrc.)

Usage

Just type:

```bash
metajournal
```

You’ll be prompted to select Daily or Weekly.
• Fill out your journal interactively.
• A Markdown file will be saved at:

```
$METAJOURNAL_SAVE_PATH/week-beginning/<week-start-date>/<date>_<daily|weekly>_entry.md
```

    •	Example:

```
~/dev/work/professional-brag/data/journal/week-beginning/2025-04-28/2025-04-28_daily_entry.md
```

Development

If you want to build from source manually:

```bash
go build -o metajournal
./metajournal
```

Notes
• If METAJOURNAL_SAVE_PATH is not set, the tool will fallback to ./data/journal/.
• Markdown is structured for easy AI parsing, mentoring, and self-review.
• Designed to be minimal, fast, and extensible.
