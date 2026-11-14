# Flow Free Bridges Solver
 
Solver berbasis *backtracking* untuk varian Flow Free: Bridges.

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

## Antarmuka Web

```bash
./flowweb
# buka http://localhost:8080
```

Aplikasi web ini menyediakan:

- **Editor peta** - grid gelap bergaya game; letakkan titik berwarna, jembatan, dan dinding dengan mengklik sel
- **Koordinat teks** - tempel atau ubah format file puzzle secara manual
- **Salin koordinat** - ekspor isi editor ke clipboard (`rows cols`, baris warna, baris `X`/`W`)
- **Solve** - menjalankan solver *backtracking* dan menampilkan pratinjau input + solusinya

## Format puzzle

```
rows cols
R r1 c1 r2 c2
G r1 c1 r2 c2
X row col    # jembatan (bisa dilintasi silang)
W row col    # dinding
```

Koordinat dimulai dari 0 (*0-indexed*). Baris yang diawali `#` dianggap komentar.
