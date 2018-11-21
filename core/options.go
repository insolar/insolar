package core

type SendOptions struct {
	Receiver *RecordRef
	Token    DelegationToken
}

type SendOption func(*SendOptions)

func SendOptionDestination(receiver *RecordRef) SendOption {
	return func(args *SendOptions) {
		args.Receiver = receiver
	}
}

func SendOptionToken(token DelegationToken) SendOption {
	return func(args *SendOptions) {
		args.Token = token
	}
}