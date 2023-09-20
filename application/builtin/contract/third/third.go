package third

import (
	"github.com/insolar/insolar/application/builtin/proxy/third"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type Third struct {
	foundation.BaseContract
	SagaCallsNum int
}

func New() (*Third, error) {
	return &Third{SagaCallsNum: 0}, nil
}

func (c *Third) GetName() (string, error) {
	return "YOU ARE ROBBED!", nil
}

func (c *Third) Transfer(delta int) error {
	proxy := third.GetObject(c.GetReference())
	err := proxy.Accept(delta)
	if err != nil {
		return err
	}
	return nil
}

func (c *Third) GetSagaCallsNum() (int, error) {
	return c.SagaCallsNum, nil
}

//ins:saga(Rollback)
func (c *Third) Accept(delta int) error {
	c.SagaCallsNum += delta
	return nil
}
func (c *Third) Rollback(delta int) error {
	c.SagaCallsNum -= delta
	return nil
}

func (c *Third) DoNothing() error {
	return nil
}
