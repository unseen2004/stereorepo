package rules

type Variant string

const (
	VariantEuropean Variant = "european"
	VariantAmerican Variant = "american"
)

type Rules struct {
	Variant                Variant
	Decks                  int
	BlackjackPayout        float64
	DealerHitsSoft17       bool
	AllowDoubleAfterSplit  bool
	AllowResplit           bool
	MaxSplits              int
	AllowSurrender         bool
	AllowInsurance         bool
	AllowSplitAces         bool
	OneCardOnSplitAces     bool
	EuropeanOBO            bool
	MinBet                 int
	MaxBet                 int
	Penetration            float64
}

func DefaultRules(variant Variant) Rules {
	r := Rules{
		Variant:               variant,
		Decks:                 6,
		BlackjackPayout:       1.5,
		DealerHitsSoft17:      false,
		AllowDoubleAfterSplit: true,
		AllowResplit:          true,
		MaxSplits:             3,
		AllowSurrender:        false,
		AllowInsurance:        true,
		AllowSplitAces:        true,
		OneCardOnSplitAces:    true,
		MinBet:                1,
		MaxBet:                500,
		Penetration:           0.25,
	}

	if variant == VariantEuropean {
		r.AllowSurrender = false
		r.AllowInsurance = true
		r.DealerHitsSoft17 = false
		r.EuropeanOBO = true
	}

	if variant == VariantAmerican {
		r.AllowSurrender = true
		r.AllowInsurance = true
		r.DealerHitsSoft17 = true
	}

	return r
}

