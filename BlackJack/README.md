# Blackjack

CLI blackjack game with European and American rules, local play, network play, bots, and ML-friendly mode.

## Build

```sh
cd /home/maks/repos/BlackJackStrategy/BlackJack
go build ./cmd/blackjack
```

## Run (local)

```sh
./blackjack -mode local -variant european -humans 1 -bots 1
```

## Run (host)

```sh
./blackjack -mode host -addr :9000 -humans 1 -remotes 1
```

## Run (join)

```sh
./blackjack -mode join -addr 127.0.0.1:9000
```

## Run (ml)

```sh
./blackjack -mode ml -bots 1
```

## Flags

All rules and configuration are exposed via flags:

- `-variant` european|american
- `-decks` number of decks
- `-blackjack_payout` payout multiplier for blackjack
- `-dealer_hits_soft17`
- `-double_after_split`
- `-resplit`
- `-max_splits`
- `-surrender`
- `-insurance`
- `-split_aces`
- `-one_card_split_aces`
- `-min_bet`
- `-max_bet`
- `-penetration` (percent of shoe dealt before shuffle, e.g. 0.75)
- `-humans`
- `-remotes`
- `-bots`
- `-bot_type` basic|strategy|hilo
- `-european_obo`
- `-bankroll`
- `-seed`
- `-rounds`
- `-mode` local|host|join|ml
- `-addr` network address
- `-wait_timeout` host wait timeout (0 for no timeout)

