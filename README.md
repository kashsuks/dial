# Dial

Dial is a lightweight time tracker with two faces on the same tool: a fast CLI for starting/stopping/logging work from the terminal, and a native desktop GUI (built with [Wails](https://wails.io) + Svelte) for a visual tracker and dashboard. Both share the same local SQLite database, so you can start a session from the CLI and see it live in the GUI, or vice versa.

> **Status:** Dial is an early-stage personal project. It has not been broadly tested across environments and the database schema/CLI flags may change without notice between releases. Use it for your own time tracking, but don't build critical workflows on top of it yet, and keep in mind your tracked data lives in a local SQLite file with no built-in backup or sync.

## How it works

- Running `dial` with no arguments launches the GUI.
- Running `dial` with a subcommand (`start`, `stop`, `status`, `log`, `today`) runs the CLI.
- Both read/write the same database at `~/.dial/data.db`, created automatically on first run.

---

## Installing from a GitHub Release

1. Go to the project's [Releases](../../releases) page and download the archive for your OS/architecture (e.g. `dial-darwin-arm64.tar.gz`, `dial-linux-amd64.tar.gz`, `dial-windows-amd64.zip`).
2. Extract the archive to get the `dial` binary (`dial.exe` on Windows).
3. Move it somewhere on your `PATH`, for example:

   ```bash
   # macOS / Linux
   tar -xzf dial-<platform>.tar.gz
   chmod +x dial
   sudo mv dial /usr/local/bin/dial
   ```

   ```powershell
   # Windows (PowerShell) - extract the zip, then add the folder to PATH
   # or move dial.exe into a directory already on PATH
   ```

4. Verify it's installed:

   ```bash
   dial status
   ```

### macOS Gatekeeper note

Since the binary/app isn't notarized by an Apple-registered developer, macOS may block it on first launch ("cannot be opened because the developer cannot be verified"). If that happens, either:

```bash
xattr -d com.apple.quarantine /usr/local/bin/dial
# or, for the .app bundle:
xattr -d com.apple.quarantine /Applications/Dial.app
```

or right-click the app in Finder and choose "Open" to bypass the check once.

### Windows SmartScreen note

Windows may show an "unrecognized app" SmartScreen warning for the same reason (no code-signing certificate). Click **More info → Run anyway** if you trust the source.

---

## Local Setup (build from source)

### Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- [Node.js](https://nodejs.org/) + npm (for the frontend)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2 (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`) if you want to use `wails dev`/`wails build` directly
- A C compiler toolchain (needed by the SQLite driver on some platforms) — Xcode Command Line Tools on macOS, `build-essential` on Linux, or MSVC/TDM-GCC on Windows

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

---

## GUI usage

Launch the GUI by running `dial` with no arguments (or by opening the installed app).

The window has two parts:

- **Tracker panel** — enter a task (and optionally a project/tags), then **Start**. While a session is running you can **Pause**, **Resume**, or **Stop** it. The elapsed time updates live.
- **Dashboard panel** — a pie chart of time by tag, a bar chart of daily totals, and summary stats (total time, session count, top tag, streak), selectable over **Today / This Week / This Month / This Year**.

Starting or stopping a session in the GUI refreshes the dashboard automatically, and since the GUI and CLI share the same database, a session started from the CLI will show up next time the dashboard refreshes.

---

## Data & storage

Dial stores all data locally in a SQLite database at `~/.dial/data.db` (created on first run). There is no cloud sync, backup, or telemetry — the data never leaves your machine. Back up that file yourself if you want to preserve your history.

## Disclaimers

- This is an unsigned/unnotarized personal-project binary — see the Gatekeeper/SmartScreen notes above before running it on macOS or Windows.
- No warranty is provided; use at your own risk. See [LICENSE](LICENSE) (MIT).
- The CLI and database schema are still evolving and may introduce breaking changes between releases.
