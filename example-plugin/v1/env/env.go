package env

type Environment struct{}

func NewEnvironment() *Environment {
	return &Environment{}
}

func (Environment) Start() {
	forever := make(chan int)
	for {
		<-forever
	}
}
