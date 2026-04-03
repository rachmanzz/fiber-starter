package cores

type AppContracts struct {
}

func CreateContract() *AppContracts {
	return &AppContracts{}
}

func (c *AppContracts) Initialize() {
	once.Do(func() {
		NewLogger()
	})
}
