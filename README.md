# Dial

Dial is a lightweight time tracker with two faces on the same tool: a fast CLI for starting/stopping/logging work from the terminal, and a native desktop GUI (built with [Wails](https://wails.io) + Svelte) for a visual tracker and dashboard. Both share the same local SQLite database, so you can start a session from the CLI and see it live in the GUI, or vice versa.

## How it works

- Running `dial` with no arguments launches the GUI.
- Running `dial` with a subcommand (`start`, `stop`, `status`, `log`, `today`) runs the CLI.
- Both read/write the same database at `~/.dial/data.db`, created automatically on first run.


## Installing from a GitHub Release

Go to the project's [Releases](../../releases) page and download the archive for your OS/architecture:

| Platform | Archive | Contains |
|----------|---------|----------|
| macOS (Apple Silicon) | `dial-darwin-arm64.zip` | `Dial.app` |
| macOS (Intel) | `dial-darwin-amd64.zip` | `Dial.app` |
| Windows (x64) | `dial-windows-amd64.zip` | `dial.exe` |
| Windows (ARM64) | `dial-windows-arm64.zip` | `dial.exe` |
| Linux (x64) | `dial-linux-amd64.tar.gz` | `dial` |

### macOS

1. Unzip the archive and drag `Dial.app` to `/Applications`.
2. Double-clicking `Dial.app` launches the GUI directly.
3. `Dial.app` is also the CLI binary — the executable inside the bundle is a normal `dial` binary that behaves like any other CLI tool when given arguments. To use it from the terminal, symlink it onto your `PATH`:

   ```bash
   sudo ln -s /Applications/Dial.app/Contents/MacOS/dial /usr/local/bin/dial
   dial status
   ```

#### Gatekeeper note

The app isn't notarized/signed by an Apple-registered developer, so macOS may block it on first launch ("cannot be opened because the developer cannot be verified"). Either right-click `Dial.app` in Finder and choose **Open** to bypass the check once, or remove the quarantine flag:

```bash
xattr -d com.apple.quarantine /Applications/Dial.app
```

### Windows

1. Extract `dial-windows-<arch>.zip` to get `dial.exe`.
2. Double-click it to launch the GUI, or move it onto your `PATH` to use it from a terminal:

   ```powershell
   dial status
   ```

#### SmartScreen note

Since the binary isn't code-signed, Windows may show an "unrecognized app" SmartScreen warning. Click **More info → Run anyway** if you trust the source.

### Linux

```bash
tar -xzf dial-linux-amd64.tar.gz
chmod +x dial
sudo mv dial /usr/local/bin/dial
dial status       # CLI
dial               # GUI (requires GTK3 + WebKit2GTK installed, see below)
```

The Linux GUI is built against `webkit2gtk-4.1`/GTK3, which need to be installed on your system (most desktop Linux distros already have these as part of GNOME/GTK apps). On Debian/Ubuntu: `sudo apt install libgtk-3-0 libwebkit2gtk-4.1-0`.

## Local Setup (build from source)

### Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- [Node.js](https://nodejs.org/) + npm (for the frontend)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2 (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`) if you want to use `wails dev`/`wails build` directly
- A native GUI toolchain for your OS — Xcode Command Line Tools on macOS, GTK3 + WebKit2GTK dev packages on Linux (`sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev`), or a working `go build` toolchain on Windows. (The SQLite driver itself is pure Go and needs no C compiler; the GUI's native webview bindings are what require these.)

### Clone and install dependencies

```bash
git clone <this-repo-url>
cd dial
cd frontend && npm install && cd ..
```

### Build using the dev script

The repo includes `scripts/dial.sh`, a helper for common dev tasks:

```bash
./scripts/dial.sh build   # go build -> bin/dial
./scripts/dial.sh run     # build + run (args after `run` are passed through)
./scripts/dial.sh test    # go test with race detector + coverage
./scripts/dial.sh lint    # gofmt check + go vet (+ golangci-lint if installed)
./scripts/dial.sh fmt     # gofmt -w .
./scripts/dial.sh tidy    # go mod tidy
./scripts/dial.sh check   # lint + test (what CI should run)
./scripts/dial.sh clean   # remove bin/ and coverage.out
./scripts/dial.sh all     # tidy + fmt + lint + test + build
```

`scripts/dial.sh build` runs a plain `go build`, which produces a working CLI binary but embeds whatever is already in `frontend/dist` for the GUI. If you're changing frontend code, build the frontend first:

```bash
cd frontend && npm run build && cd ..
./scripts/dial.sh build
```

> **Note:** a plain `go build` (without the `production` build tag) links in a stub Wails backend that does nothing when you run `dial` with no arguments — the CLI subcommands still work, but the GUI window won't appear. This is fine for CLI-only development. If you need a build where the GUI actually launches, use `go build -tags production` or, better, `wails build` (see below).

### Build/run with Wails directly (GUI development)

For live-reloading GUI development:

```bash
wails dev
```

To produce a full desktop app bundle (e.g. a `.app` on macOS):

```bash
wails build
```

The output binary/bundle will be under `build/bin/`.

### Cross-platform builds

Dial's GUI links against a different native toolkit per OS (Cocoa on macOS, WebView2 on Windows, GTK3/WebKit2GTK on Linux), so a single machine can't produce all release artifacts — each platform has to be built on that platform. `.github/workflows/release.yml` handles this: pushing a `v*.*.*` tag spins up a macOS runner (arm64 + amd64), a Windows runner (amd64 + arm64), and a Linux runner (amd64) in parallel, and publishes all five archives to a GitHub Release.

---

## CLI usage

```
dial <command> [flags]
```

### `dial start <task>`

Start tracking a task. Automatically stops any currently running session first.

| Flag        | Description                    |
|-------------|---------------------------------|
| `--project` | Project name                    |
| `--tag`     | Comma-separated tags            |

```bash
dial start "Write proposal"
dial start "Fix login bug" --project dial --tag bug,urgent
```

### `dial stop`

Stop the currently running session, if any.

```bash
dial stop
```

### `dial status`

Show the currently running session (task, elapsed time, project, tags), or say nothing is running.

```bash
dial status
```

### `dial log <task>`

Log a completed session retroactively (e.g. work you forgot to track live).

| Flag         | Description                                          | Required |
|--------------|-------------------------------------------------------|----------|
| `--duration` | Duration, e.g. `30m`, `1h15m`                          | Yes      |
| `--at`       | End time in `HH:MM` (defaults to now)                  | No       |
| `--project`  | Project name                                           | No       |
| `--tag`      | Comma-separated tags                                   | No       |

```bash
# Log 45 minutes ending now
dial log "Code review" --duration 45m

# Log a session that ended at 14:30
dial log "Team standup" --duration 15m --at 14:30 --project dial --tag meetings
```

### `dial today`

Show all sessions tracked today and the total time.

```bash
dial today
```

## GUI usage

Launch the GUI by running `dial` with no arguments (or by opening the installed app).

The window has two parts:

- **Tracker panel** — enter a task (and optionally a project/tags), then **Start**. While a session is running you can **Pause**, **Resume**, or **Stop** it. The elapsed time updates live.
- **Dashboard panel** — a pie chart of time by tag, a bar chart of daily totals, and summary stats (total time, session count, top tag, streak), selectable over **Today / This Week / This Month / This Year**.

Starting or stopping a session in the GUI refreshes the dashboard automatically, and since the GUI and CLI share the same database, a session started from the CLI will show up next time the dashboard refreshes.

## Data & storage

Dial stores all data locally in a SQLite database at `~/.dial/data.db` (created on first run). There is no cloud sync, backup, or telemetry — the data never leaves your machine. Back up that file yourself if you want to preserve your history.

## Disclaimers

- This is an unsigned/unnotarized personal-project binary — see the Gatekeeper/SmartScreen notes above before running it on macOS or Windows.

## AI Declaration

I used AI for debugging the dashboard wheel along with some styling issues for different screen sizes.
