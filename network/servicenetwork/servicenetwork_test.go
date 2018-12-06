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

package servicenetwork

/*
func newTestNodeKeeper(nodeID core.RecordRef, address string, isBootstrap bool) (network.NodeKeeper, core.Node) {
	origin := nodenetwork.NewNode(nodeID, nil, nil, 0, address, "")
	keeper := nodenetwork.NewNodeKeeper(origin)
	if isBootstrap {
		keeper.AddActiveNodes([]core.Node{origin})
	}
	return keeper, origin
}

func mockCryptographyService(t *testing.T) core.CryptographyService {
	mock := testutils.NewCryptographyServiceMock(t)
	mock.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		signature := core.SignatureFromBytes(nil)
		return &signature, nil
	}
	mock.VerifyFunc = func(p crypto.PublicKey, p1 core.Signature, p2 []byte) (r bool) {
		return true
	}
	return mock
}

func mockParcelFactory(t *testing.T) message.ParcelFactory {
	mock := mockCryptographyService(t)
	delegationTokenFactory := delegationtoken.NewDelegationTokenFactory()
	parcelFactory := messagebus.NewParcelFactory()
	cm := &component.Manager{}
	cm.Register(platformpolicy.NewPlatformCryptographyScheme())
	cm.Inject(delegationTokenFactory, parcelFactory, mock)
	return parcelFactory
}

func initComponents(t *testing.T, nodeID core.RecordRef, address string, isBootstrap bool) (core.CryptographyService, network.NodeKeeper, core.Node) {
	key, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	require.NotNil(t, key)
	cs := cryptography.NewKeyBoundCryptographyService(key)
	kp := platformpolicy.NewKeyProcessor()
	pk, _ := cs.GetPublicKey()
	_, err := certificate.NewCertificatesWithKeys(pk, kp)

	require.NoError(t, err)
	keeper, origin := newTestNodeKeeper(nodeID, address, isBootstrap)

	mock := mockCryptographyService(t)
	return mock, keeper, origin
}

func TestNewServiceNetwork(t *testing.T) {
	cfg := configuration.NewConfiguration()
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	sn, err := NewServiceNetwork(cfg, scheme)
	require.NoError(t, err)
	require.NotNil(t, sn)
}

func TestServiceNetwork_GetAddress(t *testing.T) {
	cfg := configuration.NewConfiguration()
	scheme := platformpolicy.NewPlatformCryptographyScheme()
	network, err := NewServiceNetwork(cfg, scheme)
	require.NoError(t, err)
	err = network.Init(context.Background())
	require.NoError(t, err)
	require.True(t, strings.Contains(network.GetAddress(), strings.Split(cfg.Host.Transport.Address, ":")[0]))
}

func TestServiceNetwork_SendMessage(t *testing.T) {
	cfg := configuration.NewConfiguration()
	scheme := platformpolicy.NewPlatformCryptographyScheme()
	serviceNetwork, err := NewServiceNetwork(cfg, scheme)

	key, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	require.NotNil(t, key)
	cs := cryptography.NewKeyBoundCryptographyService(key)
	kp := platformpolicy.NewKeyProcessor()
	pk, _ := cs.GetPublicKey()
	serviceNetwork.Certificate, _ = certificate.NewCertificatesWithKeys(pk, kp)
	require.NoError(t, err)

	ctx := context.TODO()
	serviceNetwork.CryptographyService, serviceNetwork.NodeKeeper, _ = initComponents(t, testutils.RandomRef(), "", true)

	err = serviceNetwork.Init(context.Background())
	require.NoError(t, err)
	err = serviceNetwork.Start(ctx)
	require.NoError(t, err)

	e := &message.CallMethod{
		ObjectRef: core.NewRefFromBase58("test"),
		Method:    "test",
		Arguments: []byte("test"),
	}

	pf := mockParcelFactory(t)
	parcel, err := pf.Create(ctx, e, serviceNetwork.GetNodeID(), nil)
	require.NoError(t, err)

	ref := testutils.RandomRef()
	serviceNetwork.SendMessage(ref, "test", parcel)
}

func mockServiceConfiguration(host string, bootstrapHosts []string, nodeID string) configuration.Configuration {
	cfg := configuration.NewConfiguration()
	transport := configuration.Transport{Protocol: "UTP", Address: host, BehindNAT: false}
	h := configuration.HostNetwork{
		Transport:      transport,
		IsRelay:        false,
	}

	n := configuration.NodeNetwork{Node: &configuration.Node{ID: nodeID}}
	cfg.Host = h
	cfg.Node = n

	return cfg
}

func TestServiceNetwork_SendMessage2(t *testing.T) {
	ctx := context.TODO()
	firstNodeId := "4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj"
	secondNodeId := "53jNWvey7Nzyh4ZaLdJDf3SRgoD4GpWuwHgrgvVVGLbDkk3A7cwStSmBU2X7s4fm6cZtemEyJbce9dM9SwNxbsxf"
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	firstNode, err := NewServiceNetwork(mockServiceConfiguration(
		"127.0.0.1:10000",
		[]string{"127.0.0.1:10001"},
		firstNodeId), scheme)
	require.NoError(t, err)
	secondNode, err := NewServiceNetwork(mockServiceConfiguration(
		"127.0.0.1:10001",
		nil,
		secondNodeId), scheme)
	require.NoError(t, err)

	var fNode, sNode core.Node
	secondNode.CryptographyService, secondNode.NodeKeeper, fNode = initComponents(t, core.NewRefFromBase58(secondNodeId), "127.0.0.1:10001", true)
	err = secondNode.Init(context.Background())
	require.NoError(t, err)
	err = secondNode.Start(ctx)
	require.NoError(t, err)
	firstNode.CryptographyService, firstNode.NodeKeeper, sNode = initComponents(t, core.NewRefFromBase58(firstNodeId), "127.0.0.1:10000", false)
	firstNode.NodeKeeper.AddActiveNodes([]core.Node{fNode, sNode})
	secondNode.NodeKeeper.AddActiveNodes([]core.Node{fNode, sNode})

	err = firstNode.Init(context.Background())
	require.NoError(t, err)
	err = firstNode.Start(ctx)
	require.NoError(t, err)

	defer func() {
		firstNode.Stop(ctx)
		secondNode.Stop(ctx)
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	secondNode.RemoteProcedureRegister("test", func(args [][]byte) ([]byte, error) {
		wg.Done()
		return nil, nil
	})

	e := &message.CallMethod{
		ObjectRef: core.NewRefFromBase58("test"),
		Method:    "test",
		Arguments: []byte("test"),
	}

	pf := mockParcelFactory(t)
	parcel, err := pf.Create(ctx, e, firstNode.GetNodeID(), nil)
	require.NoError(t, err)

	firstNode.SendMessage(core.NewRefFromBase58(secondNodeId), "test", parcel)
	success := utils.WaitTimeout(&wg, 100*time.Millisecond)

	require.True(t, success)
}

func TestServiceNetwork_SendCascadeMessage(t *testing.T) {
	ctx := context.TODO()
	firstNodeId := "4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj"
	secondNodeId := "53jNWvey7Nzyh4ZaLdJDf3SRgoD4GpWuwHgrgvVVGLbDkk3A7cwStSmBU2X7s4fm6cZtemEyJbce9dM9SwNxbsxf"
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	firstNode, err := NewServiceNetwork(mockServiceConfiguration(
		"127.0.0.1:10100",
		[]string{"127.0.0.1:10101"},
		firstNodeId), scheme)
	require.NoError(t, err)
	secondNode, err := NewServiceNetwork(mockServiceConfiguration(
		"127.0.0.1:10101",
		nil,
		secondNodeId), scheme)
	require.NoError(t, err)

	var fNode, sNode core.Node
	secondNode.CryptographyService, secondNode.NodeKeeper, fNode = initComponents(t, core.NewRefFromBase58(secondNodeId), "127.0.0.1:10101", true)
	err = secondNode.Init(context.Background())
	require.NoError(t, err)
	err = secondNode.Start(ctx)
	require.NoError(t, err)
	firstNode.CryptographyService, firstNode.NodeKeeper, sNode = initComponents(t, core.NewRefFromBase58(firstNodeId), "127.0.0.1:10100", false)
	firstNode.NodeKeeper.AddActiveNodes([]core.Node{fNode, sNode})
	secondNode.NodeKeeper.AddActiveNodes([]core.Node{fNode, sNode})
	err = firstNode.Init(context.Background())
	require.NoError(t, err)
	err = firstNode.Start(ctx)
	require.NoError(t, err)

	defer func() {
		firstNode.Stop(ctx)
		secondNode.Stop(ctx)
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	secondNode.RemoteProcedureRegister("test", func(args [][]byte) ([]byte, error) {
		wg.Done()
		return nil, nil
	})

	e := &message.CallMethod{
		ObjectRef: core.NewRefFromBase58("test"),
		Method:    "test",
		Arguments: []byte("test"),
	}

	c := core.Cascade{
		NodeIds:           []core.RecordRef{core.NewRefFromBase58(secondNodeId)},
		ReplicationFactor: 2,
		Entropy:           core.Entropy{0},
	}

	pf := mockParcelFactory(t)
	parcel, err := pf.Create(ctx, e, firstNode.GetNodeID(), nil)
	require.NoError(t, err)

	err = firstNode.SendCascadeMessage(c, "test", parcel)
	success := utils.WaitTimeout(&wg, 100*time.Millisecond)

	require.NoError(t, err)
	require.True(t, success)

	err = firstNode.SendCascadeMessage(c, "test", nil)
	require.Error(t, err)
	c.ReplicationFactor = 0
	err = firstNode.SendCascadeMessage(c, "test", parcel)
	require.Error(t, err)
	c.ReplicationFactor = 2
	c.NodeIds = nil
	err = firstNode.SendCascadeMessage(c, "test", parcel)
	require.Error(t, err)
}

func TestServiceNetwork_SendCascadeMessage2(t *testing.T) {
	ctx := context.TODO()
	t.Skip("fix data race INS-534")
	nodeIds := []core.RecordRef{
		core.NewRefFromBase58("4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj"),
		core.NewRefFromBase58("53jNWvey7Nzyh4ZaLdJDf3SRgoD4GpWuwHgrgvVVGLbDkk3A7cwStSmBU2X7s4fm6cZtemEyJbce9dM9SwNxbsxf"),
		core.NewRefFromBase58("9uE5MEWQB2yfKY8kTgTNovWii88D4anmf7GAiovgcxx6Uc6EBmZ212mpyMa1L22u9TcUUd94i8LvULdsdBoG8ed"),
		core.NewRefFromBase58("4qXdYkfL9U4tL3qRPthdbdajtafR4KArcXjpyQSEgEMtpuin3t8aZYmMzKGRnXHBauytaPQ6bfwZyKZzRPpR6gyX"),
		core.NewRefFromBase58("5q5rnvayXyKszoWofxp4YyK7FnLDwhsqAXKxj6H7B5sdEsNn4HKNFoByph4Aj8rGptdWL54ucwMQrySMJgKavxX1"),
		core.NewRefFromBase58("5tsFDwNLMW4GRHxSbBjjxvKpR99G4CSBLRqZAcpqdSk5SaeVcDL3hCiyjjidCRJ7Lu4VZoANWQJN2AgPvSRgCghn"),
		core.NewRefFromBase58("48UWM6w7YKYCHoP7GHhogLvbravvJ6bs4FGETqXfgdhF9aPxiuwDWwHipeiuNBQvx7zyCN9wFxbuRrDYRoAiw5Fj"),
		core.NewRefFromBase58("5owQeqWyHcobFaJqS2BZU2o2ZRQ33GojXkQK6f8vNLgvNx6xeWRwenJMc53eEsS7MCxrpXvAhtpTaNMPr3rjMHA"),
		core.NewRefFromBase58("xF12WfbkcWrjrPXvauSYpEGhkZT2Zha53xpYh5KQdmGHMywJNNgnemfDN2JfPV45aNQobkdma4dsx1N7Xf5wCJ9"),
		core.NewRefFromBase58("4VgDz9o23wmYXN9mEiLnnsGqCEEARGByx1oys2MXtC6M94K85ZpB9sEJwiGDER61gHkBxkwfJqtg9mAFR7PQcssq"),
		core.NewRefFromBase58("48g7C8QnH2CGMa62sNaL1gVVyygkto8EbMRHv168psCBuFR2FXkpTfwk4ZwpY8awFFXKSnWspYWWQ7sMMk5W7s3T"),
		core.NewRefFromBase58("Lvssptdwq7tatd567LUfx2AgsrWZfo4u9q6FJgJ9BgZK8cVooZv2A7F7rrs1FS5VpnTmXhr6XihXuKWVZ8i5YX9"),
	}

	prefix := "127.0.0.1:"
	port := 10000
	bootstrapNodes := nodeIds[len(nodeIds)-2:]
	bootstrapHosts := make([]string, 0)
	var wg sync.WaitGroup
	wg.Add(len(nodeIds) - 1)
	services := make([]*ServiceNetwork, 0)

	defer func() {
		for _, service := range services {
			service.Stop(ctx)
		}
	}()

	scheme := platformpolicy.NewPlatformCryptographyScheme()
	// init node and register test function
	initService := func(node string, bHosts []string) (service *ServiceNetwork, host string) {
		host = prefix + strconv.Itoa(port)
		service, _ = NewServiceNetwork(mockServiceConfiguration(host, bHosts, node), scheme)
		service.Start(ctx)
		service.RemoteProcedureRegister("test", func(args [][]byte) ([]byte, error) {
			wg.Done()
			return nil, nil
		})
		port++
		services = append(services, service)
		return
	}

	for _, node := range bootstrapNodes {
		_, host := initService(node.String(), nil)
		bootstrapHosts = append(bootstrapHosts, host)
	}
	nodes := nodeIds[:len(nodeIds)-2]
	// first node that will send cascade message to all other nodes
	var firstService *ServiceNetwork
	for i, node := range nodes {
		service, _ := initService(node.String(), bootstrapHosts)
		if i == 0 {
			firstService = service
		}
	}

	e := &message.CallMethod{
		ObjectRef: core.NewRefFromBase58("test"),
		Method:    "test",
		Arguments: []byte("test"),
	}
	// send cascade message to all nodes except the first
	c := core.Cascade{
		NodeIds:           nodeIds[1:],
		ReplicationFactor: 2,
		Entropy:           core.Entropy{0},
	}

	pf := mockParcelFactory(t)
	parcel, err := pf.Create(ctx, e, firstService.GetNodeID(), nil)
	require.NoError(t, err)

	firstService.SendCascadeMessage(c, "test", parcel)
	success := utils.WaitTimeout(&wg, 100*time.Millisecond)

	require.True(t, success)
}
*/
