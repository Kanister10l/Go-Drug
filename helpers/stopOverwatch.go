package helpers

import (
	"os"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
)

//Overwatch is supposed to spread signal between every registered entity.
//His main purpose is to catch interrupt signal and notify every registered app component.
type Overwatch struct {
	Sugar     *zap.SugaredLogger
	SigInt    chan bool
	Children  []chan bool
	Final     chan bool
	retries   int
	desiredGR int
	Mutex     sync.Mutex
}

//NewOverwatch creates, starts and returns new Overwatch
func NewOverwatch(sugar *zap.SugaredLogger, retries int) *Overwatch {
	ov := &Overwatch{
		Sugar:     sugar,
		SigInt:    make(chan bool),
		Children:  []chan bool{},
		Final:     make(chan bool),
		retries:   retries,
		desiredGR: runtime.NumGoroutine() + 1,
		Mutex:     sync.Mutex{},
	}

	go ov.catch()

	return ov
}

func (ov *Overwatch) catch() {
	<-ov.SigInt
	ov.notify()
}

func (ov *Overwatch) notify() {
	ov.Mutex.Lock()
	defer ov.Mutex.Unlock()
	for _, child := range ov.Children {
		child <- true
	}

	try := 0
	for {
		if runtime.NumGoroutine() <= ov.desiredGR {
			break
		}
		time.Sleep(100 * time.Millisecond)
		ov.Sugar.Info("Waiting additional 100ms for goroutines to finish tasks.")

		try++
		if try > ov.retries {
			ov.Sugar.Error("Some goroutines are not responding. Performing hard shutdown.")
			os.Exit(127)
		}
	}
	ov.Final <- true
}

//Register is used to register component for signal notification
func (ov *Overwatch) Register() chan bool {
	child := make(chan bool)

	ov.Mutex.Lock()
	ov.Children = append(ov.Children, child)
	ov.Mutex.Unlock()

	return child
}
