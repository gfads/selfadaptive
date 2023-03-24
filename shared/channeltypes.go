package shared

type ToManagedChan struct {
	Behaviours        []func()
	SelectedBehaviour int
}

type ToPlannerChan struct {
	ChangeRequest     string
	SelectedBehaviour int
}
