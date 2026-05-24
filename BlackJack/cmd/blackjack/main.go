package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"BlackJack/internal/bot"
	"BlackJack/internal/cli"
	"BlackJack/internal/engine"
	"BlackJack/internal/ml"
	"BlackJack/internal/netplay"
	"BlackJack/internal/rules"
)

func main() {
	mode := flag.String("mode", "local", "local|host|join|ml")
	addr := flag.String("addr", ":9000", "network address")
	variant := flag.String("variant", "european", "european|american")
	decks := flag.Int("decks", 6, "number of decks")
	bjPayout := flag.Float64("blackjack_payout", 1.5, "blackjack payout multiplier")
	hitSoft17 := flag.Bool("dealer_hits_soft17", false, "dealer hits on soft 17")
	das := flag.Bool("double_after_split", true, "allow double after split")
	resplit := flag.Bool("resplit", true, "allow resplit")
	maxSplits := flag.Int("max_splits", 3, "max split count")
	surrender := flag.Bool("surrender", false, "allow surrender")
	insurance := flag.Bool("insurance", true, "allow insurance")
	splitAces := flag.Bool("split_aces", true, "allow split aces")
	oneCardSplitAces := flag.Bool("one_card_split_aces", true, "one card after split aces")
	euroOBO := flag.Bool("european_obo", true, "european no-hole-card original bet only")
	minBet := flag.Int("min_bet", 1, "minimum bet")
	maxBet := flag.Int("max_bet", 500, "maximum bet")
	penetration := flag.Float64("penetration", 0.25, "shuffle penetration (0-1)")
	humans := flag.Int("humans", 1, "local human players")
	remotes := flag.Int("remotes", 0, "remote players (host mode)")
	bots := flag.Int("bots", 0, "bot players")
	botType := flag.String("bot_type", "basic", "basic|strategy|hilo")
	waitTimeout := flag.Duration("wait_timeout", 0, "host wait timeout (0 for no timeout)")
	bankroll := flag.Int("bankroll", 1000, "starting bankroll")
	seed := flag.Int64("seed", 0, "rng seed (0 for random)")
	rounds := flag.Int("rounds", 1, "number of rounds (0 for infinite)")
	flag.Parse()

	variantValue, err := parseVariant(*variant)
	if err != nil {
		log.Fatal(err)
	}
	r := rules.DefaultRules(variantValue)
	r.Decks = *decks
	r.BlackjackPayout = *bjPayout
	r.DealerHitsSoft17 = *hitSoft17
	r.AllowDoubleAfterSplit = *das
	r.AllowResplit = *resplit
	r.MaxSplits = *maxSplits
	r.AllowSurrender = *surrender
	r.AllowInsurance = *insurance
	r.AllowSplitAces = *splitAces
	r.OneCardOnSplitAces = *oneCardSplitAces
	r.EuropeanOBO = *euroOBO
	r.MinBet = *minBet
	r.MaxBet = *maxBet
	r.Penetration = *penetration

	shoe := engine.NewShoe(r.Decks, r.Penetration, *seed)

	switch *mode {
	case "local":
		players := buildLocalPlayers(*humans, *bots, *bankroll, *botType)
		runGame(r, shoe, players, *rounds)
	case "host":
		players, server := buildHostPlayers(*humans, *remotes, *bots, *bankroll, *addr, *botType, *waitTimeout)
		defer server.Close()
		runGame(r, shoe, players, *rounds)
	case "join":
		client, err := netplay.NewClient(*addr)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		if err := client.Run(); err != nil {
			log.Fatal(err)
		}
	case "ml":
		players := buildMLPlayers(*bots, *bankroll, *botType)
		runGame(r, shoe, players, *rounds)
	default:
		log.Fatalf("unknown mode: %s", *mode)
	}
}

func parseVariant(value string) (rules.Variant, error) {
	switch strings.ToLower(value) {
	case string(rules.VariantEuropean):
		return rules.VariantEuropean, nil
	case string(rules.VariantAmerican):
		return rules.VariantAmerican, nil
	default:
		return "", fmt.Errorf("unknown variant: %s", value)
	}
}

func buildLocalPlayers(humans, botsCount, bankroll int, botType string) []*engine.Player {
	players := make([]*engine.Player, 0, humans+botsCount)
	for i := 0; i < humans; i++ {
		name := fmt.Sprintf("Player%d", i+1)
		players = append(players, &engine.Player{
			ID:       name,
			Name:     name,
			Bankroll: bankroll,
			Provider: cli.NewLocalPlayer(name),
		})
	}
	for i := 0; i < botsCount; i++ {
		name := fmt.Sprintf("Bot%d", i+1)
		players = append(players, &engine.Player{
			ID:       name,
			Name:     name,
			Bankroll: bankroll,
			Provider: selectBotProvider(botType, name),
		})
	}
	return players
}

func buildHostPlayers(humans, remotes, botsCount, bankroll int, addr string, botType string, waitTimeout time.Duration) ([]*engine.Player, *netplay.Server) {
	server, err := netplay.NewServer(addr)
	if err != nil {
		log.Fatal(err)
	}

	remoteProviders, err := server.WaitForPlayers(remotes, waitTimeout)
	if err != nil {
		log.Fatal(err)
	}

	players := make([]*engine.Player, 0, humans+remotes+botsCount)
	for i := 0; i < humans; i++ {
		name := fmt.Sprintf("Player%d", i+1)
		players = append(players, &engine.Player{
			ID:       name,
			Name:     name,
			Bankroll: bankroll,
			Provider: cli.NewLocalPlayer(name),
		})
	}
	for i, provider := range remoteProviders {
		name := fmt.Sprintf("Remote%d", i+1)
		provider.Name = name
		players = append(players, &engine.Player{
			ID:       name,
			Name:     name,
			Bankroll: bankroll,
			Provider: provider,
		})
	}
	for i := 0; i < botsCount; i++ {
		name := fmt.Sprintf("Bot%d", i+1)
		players = append(players, &engine.Player{
			ID:       name,
			Name:     name,
			Bankroll: bankroll,
			Provider: selectBotProvider(botType, name),
		})
	}
	return players, server
}

func buildMLPlayers(botsCount, bankroll int, botType string) []*engine.Player {
	players := make([]*engine.Player, 0, botsCount+1)
	name := "Agent1"
	players = append(players, &engine.Player{
		ID:       name,
		Name:     name,
		Bankroll: bankroll,
		Provider: ml.NewAgent(name),
	})
	for i := 0; i < botsCount; i++ {
		botName := fmt.Sprintf("Bot%d", i+1)
		players = append(players, &engine.Player{
			ID:       botName,
			Name:     botName,
			Bankroll: bankroll,
			Provider: selectBotProvider(botType, botName),
		})
	}
	return players
}

func selectBotProvider(botType, name string) engine.DecisionProvider {
	switch strings.ToLower(botType) {
	case "strategy":
		return &bot.BasicStrategyBot{Name: name}
	case "hilo":
		return &bot.HiLoBot{Name: name}
	default:
		return &bot.BasicBot{Name: name}
	}
}

func runGame(r rules.Rules, shoe *engine.Shoe, players []*engine.Player, rounds int) {
	for i := 0; rounds == 0 || i < rounds; i++ {
		round := engine.NewRound(r, shoe, players)
		results, err := round.Play()
		if err != nil {
			if errors.Is(err, engine.ErrNoActiveHands) {
				fmt.Println("No active hands; ending session")
				return
			}
			log.Fatal(err)
		}
		printRound(round, results)
	}
}

func printRound(round *engine.Round, results []engine.RoundResult) {
	fmt.Printf("Dealer: %s (%d)\n", cardsToString(round.Dealer.Cards), round.Dealer.BestTotal())
	for _, p := range round.Players {
		for i, h := range p.Hands {
			fmt.Printf("%s hand %d: %s (%d) bet=%d\n", p.Name, i+1, cardsToString(h.Cards), h.BestTotal(), h.Bet)
		}
	}
	for _, res := range results {
		fmt.Printf("Result %s hand %d: %s payout=%d\n", res.PlayerID, res.Hand+1, res.Outcome, res.Payout)
	}
}

func cardsToString(cards []engine.Card) string {
	parts := make([]string, 0, len(cards))
	for _, c := range cards {
		parts = append(parts, c.String())
	}
	return strings.Join(parts, ",")
}
