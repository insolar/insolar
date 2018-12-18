package logicrunner

/*

func TestExecute_DontStartQueueProcessorWhen(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	am := testutils.NewArtifactManagerMock(t)
	//pm := testutils.NewPulseManagerMock(t)
	nn := network.NewNodeNetworkMock(t)
	jc := testutils.NewJetCoordinatorMock(t)
	lr, _ := NewLogicRunner(&configuration.LogicRunner{})
	lr.ArtifactManager = am
	///lr.PulseManager = pm
	lr.NodeNetwork = nn
	lr.JetCoordinator = jc

	nn.GetOriginMock.Return(network.NewNodeMock(t).IDMock.Return(Ref{}))

	jc.IsAuthorizedMock.Return(true, nil)

	///pm.CurrentMock.Return(&core.Pulse{}, nil)

	requestID := testutils.RandomID()

	am.RegisterRequestMock.Return(&requestID, nil)
	parcel := testutils.NewParcelMock(t)
	parcel.MessageMock.Return(&message.CallMethod{})

	reply, err := lr.Execute(ctx, parcel)
	require.NoError(t, err)
	require.NotNil(t, reply)
}

*/
