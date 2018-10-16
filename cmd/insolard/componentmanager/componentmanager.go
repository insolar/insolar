package componentmanager

import (
	"reflect"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/bootstrap"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/messagebus"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/networkcoordinator"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/updater"
)

type ComponentManager struct {
	Components core.Components
}

var instance *ComponentManager

func GetComponentManager() *ComponentManager {
	if instance == nil {
		instance = &ComponentManager{} // <--- NOT THREAD SAFE
	}
	return instance
}

// linkAll - link dependency for all components
func (cm *ComponentManager) linkAll() {
	v := reflect.ValueOf(cm.Components)
	for i := 0; i < v.NumField(); i++ {
		componentName := v.Field(i).String()
		log.Infof("Starting component `%s` ...", componentName)
		err := v.Field(i).Interface().(core.Component).Start(cm.Components)
		if err != nil {
			log.Fatalf("failed to start component %s : %s", componentName, err.Error())
		}

		log.Infof("Component `%s` successfully started", componentName)
	}
}

// stopAll - reverse order stop all components
func (cm *ComponentManager) StopAll() {
	v := reflect.ValueOf(cm.Components)
	for i := v.NumField() - 1; i >= 0; i-- {
		err := v.Field(i).Interface().(core.Component).Stop()
		log.Infoln("Stop component: ", v.String())
		if err != nil {
			log.Errorf("failed to stop component %s : %s", v.String(), err.Error())
		}
	}
}

func New(cfgHolder *configuration.Holder) (*ComponentManager, *servicenetwork.ServiceNetwork) {
	cm := GetComponentManager()
	nw, err := servicenetwork.NewServiceNetwork(cfgHolder.Configuration)
	if err != nil {
		log.Fatalln("failed to start Network: ", err.Error())
	}
	cm.Components.Network = nw

	cm.Components.Ledger, err = ledger.NewLedger(cfgHolder.Configuration.Ledger)
	if err != nil {
		log.Fatalln("failed to start Ledger: ", err.Error())
	}

	cm.Components.LogicRunner, err = logicrunner.NewLogicRunner(&cfgHolder.Configuration.LogicRunner)
	if err != nil {
		log.Fatalln("failed to start LogicRunner: ", err.Error())
	}

	cm.Components.MessageBus, err = messagebus.NewMessageBus(cfgHolder.Configuration)
	if err != nil {
		log.Fatalln("failed to start LogicRunner: ", err.Error())
	}

	cm.Components.Bootstrapper, err = bootstrap.NewBootstrapper(cfgHolder.Configuration)
	if err != nil {
		log.Fatalln("failed to start Bootstrapper: ", err.Error())
	}

	cm.Components.APIRunner, err = api.NewRunner(&cfgHolder.Configuration.APIRunner)
	if err != nil {
		log.Fatalln("failed to start ApiRunner: ", err.Error())
	}

	cm.Components.Metrics, err = metrics.NewMetrics(cfgHolder.Configuration.Metrics)
	if err != nil {
		log.Fatalln("failed to start Metrics: ", err.Error())
	}

	cm.Components.NetworkCoordinator, err = networkcoordinator.New()
	if err != nil {
		log.Fatalln("failed to start NetworkCoordinator: ", err.Error())
	}

	cm.Components.Updater, err = updater.NewUpdater(&cfgHolder.Configuration.Updater)
	if err != nil {
		log.Fatalln("failed to start Update Service: ", err.Error())
	}

	cm.linkAll()
	err = cm.Components.LogicRunner.OnPulse(*pulsar.NewPulse(cfgHolder.Configuration.Pulsar.NumberDelta, 0, &entropygenerator.StandardEntropyGenerator{}))
	if err != nil {
		log.Fatalln("failed init pulse for LogicRunner: ", err.Error())
	}

	defer func() {
		cm.StopAll()
	}()
	return cm, nw
}
