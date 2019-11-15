package actor

// Pid is an unique ID of the Actor.
type Pid uint64

// Message represents a message sent between actors.
type Message interface {}



// Actor is an entity that processes messages, sends messages to other
// actors and stores some state.
type Actor interface {
	// Receive is called when Actor receives a new Message.
	// Returns a new Actor state and/or error.
	Receive(message Message) (Actor, error)
}

// Constructor is a procedure that creates a new Actor.
// It is called when Actor is created, before receiving any Messages.
// `state` sets the initial actor state. If `limit` is > 0 it sets the
// maximum size of the mailbox (only for regular, not stashed, not
// priority messages).
type Constructor func(system System, pid Pid) (state Actor, limit int)

// System is a class responsible for creating, scheduling and otherwise
// controlling actors.
type System interface {
	// Spawn creates a new Actor and returns it's Pid.
	Spawn(constructor Constructor) Pid

	// Send sends a Message to the Actor with a given Pid. InvalidPid is
	// returned if actor with a given Pid doesn't exist or was terminated.
	// MailboxFull is returned if acrtor's mailbox is full.
	Send(pid Pid, message Message) error

	// SendPriority sends a priority Message to the Actor with a given Pid.
	// Priority messages are processed before any other messages. InvalidPid is
	// returned if actor with a given Pid doesn't exist or was terminated.
	SendPriority(pid Pid, message Message) error

	// AwaitTermination returns when all spawned Actors terminate.
	AwaitTermination()
}
