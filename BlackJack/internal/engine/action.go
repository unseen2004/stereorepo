package engine

type Action string

const (
	ActionHit       Action = "hit"
	ActionStand     Action = "stand"
	ActionDouble    Action = "double"
	ActionSplit     Action = "split"
	ActionSurrender Action = "surrender"
)

func AllActions() []Action {
	return []Action{ActionHit, ActionStand, ActionDouble, ActionSplit, ActionSurrender}
}
