package knwldge

type Knowledge struct {
	AvailableBehaviours               map[string]func()
	LastBehaviour                     string
	CurrentSecurityLevelOfEnvironment string
	CurrentSecurityLevelOfApplication string
}

var KnowledgeDatabase = NewKnowledge()

func NewKnowledge() *Knowledge {
	return &Knowledge{LastBehaviour: "DefaultBehaviour", CurrentSecurityLevelOfApplication: "", CurrentSecurityLevelOfEnvironment: ""}
}
