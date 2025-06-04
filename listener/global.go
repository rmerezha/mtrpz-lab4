package listener

import "sync"

type GlobalListener struct {
	Listeners []Listener
}

func (pl *GlobalListener) Listen(stopCh <-chan struct{}) {
	var wg sync.WaitGroup
	wg.Add(len(pl.Listeners))

	for _, listener := range pl.Listeners {
		go func(l Listener) {
			defer wg.Done()
			l.Listen(stopCh)
		}(listener)
	}

	wg.Wait()
}
