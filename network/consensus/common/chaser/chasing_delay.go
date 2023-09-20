package chaser

import "time"

func NewChasingTimer(chasingDelay time.Duration) ChasingTimer {
	return ChasingTimer{chasingDelay: chasingDelay}
}

type ChasingTimer struct {
	chasingDelay time.Duration
	timer        *time.Timer
	wasCleared   bool
}

func (c *ChasingTimer) IsEnabled() bool {
	return c.chasingDelay > 0
}

func (c *ChasingTimer) WasStarted() bool {
	return c.timer != nil
}

func (c *ChasingTimer) RestartChase() {

	if c.chasingDelay <= 0 {
		return
	}

	if c.timer == nil {
		c.timer = time.NewTimer(c.chasingDelay)
		return
	}

	// Restart chasing timer from this moment
	if !c.wasCleared && !c.timer.Stop() {
		<-c.timer.C
	}
	c.wasCleared = false
	c.timer.Reset(c.chasingDelay)
}

func (c *ChasingTimer) Channel() <-chan time.Time {
	if c.timer == nil {
		return nil // receiver will wait indefinitely
	}
	return c.timer.C
}

func (c *ChasingTimer) ClearExpired() {
	c.wasCleared = true
}
