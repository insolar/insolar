package core

type SendOptions struct {
	Receiver *RecordRef
	Token    DelegationToken
}

type SendOption func(*SendOptions)

func Receiver(receiver *RecordRef) SendOption {
	return func(args *SendOptions) {
		args.Receiver = receiver
	}
}

func Token(token DelegationToken) SendOption {
	return func(args *SendOptions) {
		args.Token = token
	}
}