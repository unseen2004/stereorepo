# Tic-Tac-Toe 5x5 Server Bot

Networked implementation of 5x5 Tic-Tac-Toe with special rules, featuring AI bot with Minimax algorithm.

---

## Game Rules

- 5x5 board with X and O players
- Win by connecting 4 marks in a row (any direction)
- Lose by creating 3-in-a-row that cannot extend to 4
- Draw after 25 moves (full board)

---

## Structure

```
.
├── minimax_player.cpp     # AI bot with human mode option
└── server/
    ├── board.hpp          # Game board implementation
    ├── game_client.cpp    # Human player client
    └── game_random_bot.cpp # Random move bot
```

---

## Building

Server and basic clients:
```bash
cd server
mkdir build && cd build
cmake ..
make
```

Minimax AI clients:
```bash
g++ -o minimax_player minimax_player.cpp -O3
```

---

## Running

Start server:
```bash
./game_server <IP> <PORT>
```

Connect human player:
```bash
./game_client <IP> <PORT> <PLAYER_TYPE> <NAME>
```

Connect random bot:
```bash
./game_random_bot <IP> <PORT> <PLAYER_TYPE> <NAME>
```

Connect Minimax AI:
```bash
./minimax_player <IP> <PORT> <PLAYER_TYPE> <NAME> [DEPTH] [AI]
```

Parameters:
- PLAYER_TYPE: 1=X, 2=O
- DEPTH: AI search depth (1-10, default 5)
- AI: 0=human mode, 1=AI mode (default 0)

---

## Note

Two players with different types must connect to start a game.