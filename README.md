# Flow Free Bridges Solver

Backtracking solver for the Flow Free: Bridges variant.

## Build

```bash
go build -o flowcli ./cmd/flowcli
go build -o flowweb ./cmd/flowweb
```

## CLI

```bash
./flowcli solve puzzles/classic_4x4.txt
./flowcli solve -trace puzzles/bridge_5x5.txt
./flowcli bench puzzles/
./flowcli verify puzzles/classic_4x4.txt
```

## Web UI

```bash
./flowweb
# open http://localhost:8080
```

The web app includes:

- **Map editor** - game-style dark grid; place colored dots, bridges, and walls by clicking cells
- **Text coordinates** - paste or edit the puzzle file format manually
- **Copy coordinates** - export editor state to clipboard (`rows cols`, color lines, `X`/`W` lines)
- **Solve** - runs the backtracking solver and shows input preview + solution

## Puzzle format

```
rows cols
R r1 c1 r2 c2
G r1 c1 r2 c2
X row col    # bridge (allows crossing)
W row col    # wall
```

Coordinates are 0-indexed. Lines starting with `#` are comments.
