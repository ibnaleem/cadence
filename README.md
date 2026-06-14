# cadence

Terminal-based habit tracker. Tracks daily and weekly habits, streaks, and completion history — all local, no account required.

```
cadence · Show up. Every day.

  Sunday, June 14

  ████████████████░░░░░░░░░░░░  4/5 · 80%

  ✓  Drink water     [daily]   🔥3
  ✓  Read            [daily]   🔥3
  ✓  Meditate        [daily]   🔥1
  ✓  Exercise        [weekly]  🔥1
  ○  Journal         [daily]

  This week  on a roll 🚀

  ●●●●●●●  Drink water    7x  [daily]
  ●●●●●●●  Read           7x  [daily]
  ●●●●●░░  Meditate       5x  [daily]
  ●●░░░░░  Exercise       2x  [weekly]
  ●●●●░░░  Journal        4x  [daily]
```

## Install

**Prerequisites:** Go 1.25+, [Ollama](https://ollama.com/download) (for name-based habit lookup).

```bash
git clone https://github.com/ibnaleem/cadence
cd cadence
go install .
```

Pull the embedding model (used by `done <name>` and `add` similarity check):

```bash
ollama pull embeddinggemma
```

## Usage

```
cadence            Show today's dashboard
cadence setup      Initialise the database (auto-runs on first use)
cadence add        Add a new habit
cadence list       List all habits
cadence done       Log a habit completion for today
cadence edit       Update a habit's name or description
cadence delete     Remove a habit and its history
cadence streak     Show current streak for each habit
```

### add

```bash
cadence add "Drink water"
cadence add "Read" --description "30 mins before bed" --frequency weekly
```

Flags: `--description / -d`, `--frequency / -f` (default: `daily`).

If Ollama is running, `add` checks for similar existing habits before inserting and prompts you if one is found above the similarity threshold.

### done

```bash
cadence done 1              # by ID (fast, no Ollama needed)
cadence done "drink water"  # by name — uses cosine similarity to find the best match
```

### edit

```bash
cadence edit 1 --name "Drink 8 glasses"
cadence edit 1 --description "Before every meal"
cadence edit 1 --name "Drink 8 glasses" --description "Before every meal"
```

### delete

```bash
cadence delete 1
```

Removes the habit and all its completion history.

### streak

```bash
cadence streak
```

```
  Drink water     🔥 7-day streak
  Read            🔥 3-day streak
  Meditate        no streak
```

A streak counts consecutive days ending today or yesterday. Missing a day resets it to zero.

## Data

Stored in `~/.cadence/cadence.db` (SQLite). No cloud sync, no telemetry.

## Build

| Task  | Command          |
|-------|------------------|
| Build | `go build ./...` |
| Run   | `go run .`       |
| Test  | `go test ./...`  |
| Vet   | `go vet ./...`   |
