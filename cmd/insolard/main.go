/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/insolar/insolar/core/utils"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/version"
)

// ComponentManager is deprecated and will be removed after completly switching to component.Manager
type ComponentManager struct {
	components core.Components
}

// linkAll - link dependency for all components
func (cm *ComponentManager) linkAll(ctx context.Context) {
	inslog := inslogger.FromContext(ctx)
	v := reflect.ValueOf(cm.components)
	for i := 0; i < v.NumField(); i++ {

		if component, ok := v.Field(i).Interface().(core.Component); ok {
			componentName := v.Field(i).String()
			inslog.Infof("==== Old ComponentManager: Starting component `%s` ...", componentName)
			err := component.Start(ctx, cm.components)
			if err != nil {
				inslog.Fatalf("==== Old ComponentManager: failed to start component %s : %s", componentName, err.Error())
			}
			inslog.Infof("==== Old ComponentManager: Component `%s` successfully started", componentName)
		}
	}
}

type inputParams struct {
	configPath               string
	isBootstrap              bool
	bootstrapCertificatePath string
	traceEnabled             bool
	nodeKeysPath             string
}

func parseInputParams() inputParams {
	var rootCmd = &cobra.Command{Use: "insolard"}
	var result inputParams
	rootCmd.Flags().StringVarP(&result.configPath, "config", "c", "", "path to config file")
	rootCmd.Flags().BoolVarP(&result.isBootstrap, "bootstrap", "b", false, "is bootstrap mode")
	rootCmd.Flags().StringVarP(&result.bootstrapCertificatePath, "cert_out", "r", "", "path to write bootstrap certificate")
	rootCmd.Flags().BoolVarP(&result.traceEnabled, "trace", "t", false, "enable tracing")
	rootCmd.Flags().StringVarP(&result.nodeKeysPath, "nodeskeys", "k", "", "path to dir with node keys files for bootstrap")
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Wrong input params:", err)
	}

	if result.isBootstrap && len(result.bootstrapCertificatePath) == 0 {
		log.Fatal("flag '--cert_out|-r' must not be empty, if '--bootstrap|-b' exists")
	}

	if result.isBootstrap && len(result.nodeKeysPath) == 0 {
		log.Fatal("flag '--nodeskeys|-k' must not be empty, if '--bootstrap|-b' exists")
	}
	return result
}

func keysFiles(ctx context.Context, nodeKeysPath string) []string {
	var keysFiles []string
	files, err := ioutil.ReadDir(nodeKeysPath)
	checkError(ctx, err, "failed to read dir for nodes keys")
	for _, f := range files {
		if f.IsDir() {
			keysFiles = append(keysFiles, f.Name())
		}
	}
	return keysFiles
}

var roles = []string{"virtual", "virtual", "heavy_material", "light_material"}

func bootstrapNodesInfo(ctx context.Context, nodeKeysPath string) []map[string]string {
	files := keysFiles(ctx, nodeKeysPath)
	if len(roles) < len(files) {
		checkError(ctx, errors.New("not enough roles"), "there more files with nodes keys than roles for them")
	}

	info := []map[string]string{}
	for i, f := range files {
		cert, err := certificate.NewCertificatesWithKeys(f)
		checkError(ctx, err, "can't create certificate from file: "+f)
		nodeInfo := map[string]string{
			"public_key": cert.PublicKey,
			"role":       roles[i],
		}
		info = append(info, nodeInfo)
	}

	return info
}

func registerCurrentNode(ctx context.Context, host string, bootstrapCertificatePath string, cert core.Certificate, nc core.NetworkCoordinator) {
	role := "virtual"
	publicKey, err := cert.GetPublicKey()
	checkError(ctx, err, "failed to get public key")

	rawCertificate, err := nc.RegisterNode(ctx, publicKey, 0, 0, role, host)
	checkError(ctx, err, "can't register node")

	err = ioutil.WriteFile(bootstrapCertificatePath, rawCertificate, 0644)
	checkError(ctx, err, "can't write certificate")
}

func mergeConfigAndCertificate(ctx context.Context, cfg *configuration.Configuration) {
	inslog := inslogger.FromContext(ctx)
	if len(cfg.CertificatePath) == 0 {
		inslog.Info("No certificate path - No merge")
		return
	}
	cert, err := certificate.NewCertificate(cfg.KeysPath, cfg.CertificatePath)
	checkError(ctx, err, "can't create certificate")

	cfg.Host.BootstrapHosts = []string{}
	for _, bn := range cert.BootstrapNodes {
		cfg.Host.BootstrapHosts = append(cfg.Host.BootstrapHosts, bn.Host)
	}
	cfg.Node.Node.ID = cert.Reference
	cfg.Host.MajorityRule = cert.MajorityRule

	inslog.Infof("Add %d bootstrap nodes. Set node id to %s. Set majority rule to %d",
		len(cfg.Host.BootstrapHosts), cfg.Node.Node.ID, cfg.Host.MajorityRule)
}

func main() {
	params := parseInputParams()

	jww.SetStdoutThreshold(jww.LevelDebug)
	cfgHolder := configuration.NewHolder()
	var err error
	if len(params.configPath) != 0 {
		err = cfgHolder.LoadFromFile(params.configPath)
	} else {
		err = cfgHolder.Load()
	}
	if err != nil {
		log.Warnln("failed to load configuration from file: ", err.Error())
	}

	err = cfgHolder.LoadEnv()
	if err != nil {
		log.Warnln("failed to load configuration from env:", err.Error())
	}

	cfg := &cfgHolder.Configuration
	traceid := utils.RandTraceID()
	ctx := inslogger.ContextWithTrace(context.Background(), traceid)
	ctx, inslog := initLogger(ctx, cfg.Log)

	if !params.isBootstrap {
		mergeConfigAndCertificate(ctx, cfg)
	}
	cfg.Metrics.Namespace = "insolard"

	fmt.Print("Starts with configuration:\n", configuration.ToString(cfgHolder.Configuration))

	jaegerflush := func() {}
	if params.traceEnabled {
		jconf := cfg.Tracer.Jaeger
		jaegerflush = instracer.ShouldRegisterJaeger(ctx, "insolard", jconf.AgentEndpoint, jconf.CollectorEndpoint)
		ctx = instracer.SetBaggage(ctx, instracer.Entry{Key: "traceid", Value: traceid})
	}
	defer jaegerflush()

	cm, cmOld, repl, err := InitComponents(ctx, *cfg, params.isBootstrap, params.nodeKeysPath)
	checkError(ctx, err, "failed to init components")

	cmOld.linkAll(ctx)

	err = cm.Start(ctx)
	checkError(ctx, err, "failed to start components")

	defer func() {
		inslog.Warn("DEFER STOP APP")
		err = cm.Stop(ctx)
		checkError(ctx, err, "failed to stop components")
	}()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		inslog.Debugln("caught sig: ", sig)

		inslog.Warn("GRACEFULL STOP APP")
		err = cm.Stop(ctx)
		jaegerflush()
		checkError(ctx, err, "failed to graceful stop components")
		os.Exit(0)
	}()

	// move to bootstrap component
	if params.isBootstrap {
		registerCurrentNode(
			ctx,
			cfg.Host.Transport.Address,
			params.bootstrapCertificatePath,
			cmOld.components.Certificate,
			cmOld.components.NetworkCoordinator,
		)
		inslog.Info("It's bootstrap mode, that is why gracefully stop daemon by sending SIGINT")
		gracefulStop <- syscall.SIGINT
	}

	fmt.Println("Version: ", version.GetFullVersion())
	fmt.Println("Running interactive mode:")
	repl.Start(ctx)
}

func initLogger(ctx context.Context, cfg configuration.Log) (context.Context, core.Logger) {
	inslog, err := log.NewLog(cfg)
	if err != nil {
		panic(err)
	}
	err = inslog.SetLevel(cfg.Level)
	if err != nil {
		inslog.Errorln(err.Error())
	}
	ctx = inslogger.SetLogger(ctx, inslog)
	return ctx, inslog
}

func checkError(ctx context.Context, err error, message string) {
	if err == nil {
		return
	}
	inslog := inslogger.FromContext(ctx)
	log.WithSkipDelta(inslog, +1).Fatalf("%v: %v", message, err.Error())
}
