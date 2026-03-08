# broom

A fast, no-nonsense Windows system cleanup tool. Run it, let it sweep, move on.

![Windows](https://img.shields.io/badge/platform-Windows-blue)
![Go](https://img.shields.io/badge/built%20with-Go-00ADD8)

## What it does

Broom clears out the junk that Windows accumulates over time:

- **Temp files** — `C:\Windows\Temp` and user temp directories
- **Prefetch cache**
- **Windows logs** and error reports (WER)
- **Thumbnail cache** — kills Explorer briefly to unlock files, then restarts it
- **Delivery Optimization cache**
- **Recycle Bin**
- **DNS cache** — runs `ipconfig /flushdns`
- **Windows Update cache** — stops `wuauserv`/`bits`, clears downloads, restarts both
- **Windows.old** — takes ownership and removes the entire folder if present
- **Microsoft Store cache**
- **DISM component cleanup**

All tasks run concurrently for speed. Everything gets logged to `cleanup_log.txt`.

## Getting started

### Option 1: Download from Releases

Grab the latest `.exe` from the [Releases](../../releases) page. Two builds are available:

| File | Description |
|------|-------------|
| `broom.exe` | GUI mode — runs silently without a console window |
| `broom-console.exe` | Console mode — shows progress output in terminal |

Double-click it. Broom will ask for admin privileges (UAC prompt), and on first run it installs itself to `%LOCALAPPDATA%\broom` and adds that to your PATH. After that, just open any terminal and type:

```
broom
```

### Option 2: Build from source

Requires Go 1.24+.

```bash
git clone https://github.com/your-username/broom.git
cd broom

# GUI build (no console window)
make build

# Console build (shows output)
make build-console
```

The binary lands in `bin/broom.exe`.

### Updating

Broom can update itself to the latest version published in GitHub releases. Simply run:

```bash
broom update
```

It will fetch the latest executable, replace the current one, and keep the older version backed up as `broom.old.exe`.

## How it works

1. Broom checks if it's running as Administrator. If not, it re-launches itself elevated via UAC.
2. On first run, it copies itself to `%LOCALAPPDATA%\broom` and adds that folder to your user PATH (so you can type `broom` from anywhere going forward).
3. Cleanup tasks kick off in parallel — temp files, caches, logs, etc.
4. Results are printed to the console (in console mode) and always appended to `cleanup_log.txt`.

## Project structure

```
cmd/broom/          Entry point
internal/
  cleanup/          All the cleanup routines
  elevate/          UAC elevation via ShellExecuteW
  install/          Self-install to PATH on first run
  logger/           File + console logging
  update/           Self-updating from GitHub releases
```

## Requirements

- Windows 10/11
- Administrator privileges (Broom will request them automatically)

## License

MIT
