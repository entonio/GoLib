package lang

import (
	"time"

	"golib/log"
)

func WaitUntil(durationGoalFunction func() (time.Duration, bool)) bool {
	return WaitUntilTimeout(durationGoalFunction, 0)
}

func WaitUntilTimeout(durationGoalFunction func() (time.Duration, bool), timeout time.Duration) bool {
	started := time.Now()
	var elapsed time.Duration
	success := false
	hadToWait := false

	for {
		duration, achieved := durationGoalFunction()

		elapsed = time.Since(started)

		if achieved {
			success = true
			log.Trace("Condition met, will proceed")
			break
		}

		if timeout > 0 && elapsed > timeout {
			log.Debug("Condition not met after %s", timeout)
			break
		}

		log.Debug("Condition not yet met, will wait for %s", duration)
		time.Sleep(duration)
		hadToWait = true
	}

	if hadToWait {
		log.Debug("Total wait: %s", elapsed)
	}

	return success
}
