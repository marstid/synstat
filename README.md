# Synstat

A terminal dashboard for monitoring your Synthetic API usage.

## Installation

```bash
go install github.com/marstid/synstat@latest
```

Or run directly:

```bash
go run github.com/marstid/synstat@latest
```

## Usage

Set your Synthetic API key as an environment variable:

```bash
export SYNTHETIC_API_KEY="your-api-key"
synstat
```

### Configuration

You can set a default theme using the `SYNTHETIC_THEME` environment variable:

```bash
export SYNTHETIC_THEME="Hacker Terminal"
synstat
```

Or use partial names:
```bash
export SYNTHETIC_THEME="hacker"
export SYNTHETIC_THEME="ocean"
export SYNTHETIC_THEME="purple"
```

## Controls

- `←/→` or `h/l` - Change theme
- `r` - Refresh data
- `q` or `Ctrl+C` - Quit

## Themes

The dashboard includes 11 different color themes:
- Cyan (default)
- Blue
- Green
- Purple
- Candy
- Soda-pop
- Hacker Terminal
- Sunset
- Ocean
- Forest
- Monochrome

## Requirements

- Go 1.23 or later
- A Synthetic API key
