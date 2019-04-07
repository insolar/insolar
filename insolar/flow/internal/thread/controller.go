package thread

type Controller struct {
	cancel chan struct{}
}

func NewController() *Controller {
	return &Controller{cancel: make(chan struct{})}
}

func (c *Controller) Cancel() <-chan struct{} {
	return c.cancel
}

func (c *Controller) Pulse() {
	toClose := c.cancel
	c.cancel = make(chan struct{})
	close(toClose)
}
