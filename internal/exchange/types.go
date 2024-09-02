package exchange

import "sync"

type exchange struct {
	mu  sync.RWMutex
	Exc map[interface{}]*chan StatusChannel
}

type StatusChannel struct {
	Module string
	Start  bool
	Stop   bool
	Update bool
	Error  error
	Data   interface{}
}

func Init() *exchange {
	exc := exchange{
		Exc: make(map[interface{}]*chan StatusChannel),
	}
	return &exc
}

func (ex *exchange) NewChannel(idChannel int) {
	ex.mu.Lock()
	defer ex.mu.Unlock()
	ex.Exc[idChannel] = createChannel()
}

func (ex *exchange) ReadChannel(idChannel int) <-chan StatusChannel {
	var outCh = make(chan StatusChannel, 10)
	go func() {
		for s := range *ex.Exc[idChannel] {
			outCh <- s
		}
	}()
	return outCh
}

func (ex *exchange) WriteChannel(idChannel int, status StatusChannel) {
	ex.mu.Lock()
	defer ex.mu.Unlock()
	*ex.Exc[idChannel] <- status
}

func createChannel() *chan StatusChannel {
	ch := make(chan StatusChannel, 10)
	return &ch
}
