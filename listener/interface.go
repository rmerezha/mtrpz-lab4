package listener

type Listener interface {
	Listen(stopCh <-chan struct{})
}
