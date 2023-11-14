package shared

type ToManagedChan struct {
	Behaviours        map[string]func()
	SelectedBehaviour string
}

type ToPlannerChan struct {
	ChangeRequest     string
	SelectedBehaviour string
}

type SubscriberToAdapter struct {
	QueueSize        int
	ReceivedMessages int
	D                float64
}
