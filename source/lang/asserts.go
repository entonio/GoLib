package lang

func Assert(condition bool, message string) string {
	if !condition {
		panic(message)
	}
	return "[ASSERT OK]"
}

func AssertNil(err error) {
	if err != nil {
		if len(shutdownListeners) > 0 {
			message := err.Error()
			stack := errorStack(err, 2)
			for _, listener := range shutdownListeners {
				listener.OnError(err, message, stack)
			}
		}
		panic(err)
	}
}
