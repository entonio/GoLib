package lang

var shutdownListeners = []ShutdownListener{}

type ShutdownListener interface {
	OnError(err error, message string, stack string)
}

func AddShutdownListener(listener ShutdownListener) {
	shutdownListeners = append(shutdownListeners, listener)
}
