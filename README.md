# Pomodoro CLI

This is a command-line Pomodoro timer designed to help you stay focused and productive. It's built with Go and uses audio cues to notify you of round completions, rest periods, and upcoming rounds.

## Features

- **Customizable Timers:** Set your own durations for rounds, rests, and notifications.
- **Audio Notifications:** Get audio alerts for different phases of the Pomodoro cycle.
- **Cross-Platform:** Build and run on Linux, macOS, and Windows.
- **Easy to Use:** Simple command-line interface.

## Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/shrtyk/pomodoro-cli.git
    cd pomodoro-cli
    ```

2.  **Build the application:**
    - **Linux:**
      ```bash
      # binary will be placed into ./bin/linux
      make build/linux
      ```
    - **macOS:**
      ```bash
      # binary will be placed into ./bin/mac
      make build/mac
      ```
    - **Windows:**
      ```bash
      # binary will be placed into ./bin/win
      make build/win
      ```

## Usage

To run the application use builded binary or use the following command:

```bash
go run ./cmd/app [flags]
```

### Flags

- `--rounds`: Number of Pomodoro rounds. (default: `1`)
- `--round_duration`: Duration of a single round. (default: `25m`)
- `--rest_duration`: Duration of the rest period. (default: `5m`)
- `--notify_before_rest_end`: Time to notify before the end of the rest period. (default: `30s`)
- `--notification_file`: Path to the audio file for rest end notifications. (default: `notify.mp3`)
- `--done_file`: Path to the audio file for the end of a rest period. (default: `done.mp3`)
- `--round_file`: Path to the audio file for the start of a new round. (default: `round.mp3`)

## Configuration

You can configure the Pomodoro timer by providing command-line flags. For example, to run 2 rounds of 30 minutes each, with a 10-minute rest, you would use the following command:

```bash
go run ./cmd/app --rounds 2 --round_duration 30m --rest_duration 5m
```

You can also change the notification sounds by providing the path to your own audio files (`.mp3` or `.wav`).
