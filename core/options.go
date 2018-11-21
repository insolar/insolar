package core

type Options struct {
	Receiver *RecordRef
	Token    DelegationToken
}

type Option func(*Options)

func Receiver(receiver *RecordRef) Option {
	return func(args *Options) {
		args.Receiver = receiver
	}
}

func Token(token DelegationToken) Option {
	return func(args *Options) {
		args.Token = token
	}
}