package listener

type GlobalListener struct {
	Listeners []Listener
}

func (pl *GlobalListener) Listen(stopCh <-chan struct{}) {
	for _, listener := range pl.Listeners {
		go func(l Listener) {
			l.Listen(stopCh)
		}(listener)
	}
}
