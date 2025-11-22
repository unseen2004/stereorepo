# 15 Puzzle Bot

A C++ implementation of a solver for the classic 15 puzzle using the A* search algorithm.

---

## Overview

Solves the 15 puzzle (or 8 puzzle in 3x3 mode) using A* search with two heuristics:
- Hamming distance (misplaced tiles count)
- Manhattan distance (sum of tile distances from goal positions)

Can demonstrate individual solutions or compare heuristic performance.

---

## Building

Requires C++20 compiler and CMake 3.30+:
```bash
mkdir build && cd build
cmake ..
make
```

---

## Running

Execute the solver:
```bash
./15_puzzle_bot
```

Interactive prompts will ask you to:
1. Select puzzle size (3x3 or 4x4)
2. Choose demonstration mode or performance tests
3. Select heuristic (if in demo mode)

---

## Features

- A* search algorithm implementation
- Multiple heuristic options
- Performance comparison between heuristics
- Solution path visualization
- Statistics: moves, states visited, time taken
