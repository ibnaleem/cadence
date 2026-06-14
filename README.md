<img src="https://github.com/ibnaleem/cadence/actions/workflows/go.yml/badge.svg?event=push" alt="GitHub Actions Badge"> 

# cadence

Terminal habit tracker. Tracks daily and weekly habits with streaks and history, all stored locally in SQLite.

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

Requires Go 1.25+ and [Ollama](https://ollama.com/download) (only needed for name-based habit lookup).

#### Quick one-liner
```bash
go install https://github.com/ibnaleem/cadence@latest
```
#### Manual
```bash
git clone https://github.com/ibnaleem/cadence
cd cadence
go install .
```

Pull the embedding model used by `done <name>` and the duplicate check in `add`:

```bash
ollama pull embeddinggemma:latest
```

## Commands

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

If Ollama is running, `add` checks whether a similar habit already exists before inserting. If the similarity score crosses the threshold, it prompts you before proceeding.

### done

```bash
cadence done 1              # by ID, no Ollama needed
cadence done "drink water"  # by name, uses cosine similarity to find the best match
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

Deletes the habit and all its completions.

### streak

```bash
cadence streak
```

```
  Drink water     🔥 7-day streak
  Read            🔥 3-day streak
  Meditate        no streak
```

A streak counts consecutive days ending today or yesterday. Miss a day and it resets.

## Data

Everything lives in `~/.cadence/cadence.db`. No cloud sync, no telemetry, nothing leaving your machine.

## Build

```bash
go build ./...   # build
go run .         # run
go test ./...    # test
go vet ./...     # vet
```

## LICENSE
This project is licensed under the GNU General Public License - see the [LICENSE](https://github.com/ibnaleem/gosearch/blob/main/LICENSE) file for details.

## Support
[![BuyMeACoffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-ffdd00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black)](https://buymeacoffee.com/ibnaleem)
[![Thanks.dev](https://img.shields.io/badge/thanks.dev-0a0a0a?style=for-the-badge&logo=tv-time&logoColor=white)](https://thanks.dev/u/gh/ibnaleem)

                        
## Stargazers over time
[![Stargazers over time](https://starchart.cc/ibnaleem/cadence.svg?variant=adaptive)](https://starchart.cc/ibnaleem/cadence)
