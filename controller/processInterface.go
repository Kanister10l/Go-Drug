package controller

type Process interface {
	NewProcess(interface{}) error
	Stop() error
	Halt() error // Forcibly shut down process (in case of process being stuck)
	GetState() (StateSchema, interface{}, error)
}

type StateSchema int

// Process state schema descriptor
const (
	Internal StateSchema = iota
	Local
	External
)
