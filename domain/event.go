package domain

type Event struct {
	Action string `bson:"action" json:"action"`
	Target string `bson:"target" json:"target"`
	Message string `bson:"message" json:"message"`
}
