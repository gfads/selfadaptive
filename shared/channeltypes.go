package shared

type ToManagedChan struct {
	Behaviours        map[string]func()
	SelectedBehaviour string
}

type ToPlannerChan struct {
	ChangeRequest     string
	SelectedBehaviour string
}
