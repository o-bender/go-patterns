package singleton

type Instance struct {
	connection string
}

var instance *Instance = nil

func GetSingleton() *Instance {
	return instance
}

func NewSingleton(connection string) *Instance {
	if instance == nil {
		instance = &Instance{
			connection: connection,
		}
	}
	return instance
}

func ResetSingleton() {
	instance = nil
}
