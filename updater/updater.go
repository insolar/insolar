package updater

import (
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/version"
	"github.com/pkg/errors"
	"time"
)

type Updater struct {
	ServersList       []string
	BinariesList      []string
	LastSuccessServer string
	CurrentVer        string
	Delay             int64
	started           bool
}

func NewUpdater(cfg *configuration.Updater) (*Updater, error) {
	if cfg == nil {
		return nil, errors.New("[ NewUpdater ] config is nil")
	}
	delay := cfg.Delay
	if delay < 0 {
		return nil, errors.New("[ NewUpdater ] Delay value is out of bounds")
	} else if delay == 0 {
		log.Warn("[ NewUpdater ] The update service is DISABLED, to ENABLE the update service, set the DELAY value not equal to zero")
	}
	if len(cfg.BinariesList) == 0 {
		log.Warn("[ NewUpdater ] The list of binaries is clean, the update service will be disabled")
		delay = 0
	}
	if len(cfg.ServersList) == 0 {
		log.Warn("[ NewUpdater ] The list of update servers is clean, the update service will be disabled")
		delay = 0
	}
	updater := Updater{
		[]string{"http://localhost:2345"},
		[]string{"insgocc", "insgorund", "insolar", "insolard", "pulsard", "insupdater"},
		"",
		version.Version,
		delay,
		false,
	}
	return &updater, nil
}

// Start is implementation of core.Component interface.
func (up *Updater) Start(components core.Components) error {
	log.Info("Update service starting...")
	if up.Delay == 0 {
		log.Warn("The update service is DISABLED")
	}
	delay := time.Duration(up.Delay)
	go func() {
		ticker := time.NewTicker(time.Minute * delay)
		defer func() {
			log.Info("Stopping update service")
			ticker.Stop()
		}()

		for range ticker.C {
			err := up.verifyAndUpdate()
			if err != nil {
				log.Warn(err)
			}
		}
	}()
	return nil
}

// Stop is implementation of core.Component interface.
func (up *Updater) Stop() error {
	const timeOut = 5
	log.Infoln("Shutting down update service")
	count := 0
	for {
		if !up.started || count > 36 {
			break
		}
		count++
		time.Sleep(time.Second * time.Duration(timeOut))
	}
	return nil
}
