package knwldge

type Knowledge struct {
	AvailableBehaviours               []func()
	LastBehaviour                     int
	CurrentSecurityLevelOfEnvironment int
	CurrentSecurityLevelOfApplication int
}

var KnowledgeDatabase = NewKnowledge()

func NewKnowledge() *Knowledge {
	return &Knowledge{LastBehaviour: 0}
}
