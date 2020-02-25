package errors

import "fmt"

// InvalidPid is returned by System.Send when actor with given a Pid doesn't
// exist (e.g. it was terminated).
var InvalidPid = fmt.Errorf("Invalid actor Pid")

// MailboxFull is returned by System.Send if actor's Mailbox is full.
var MailboxFull = fmt.Errorf("Mailbox is full")

// UnexpectedMessage is returned by Actor.Receive when the actor didn't
// expect a given message.
var UnexpectedMessage = fmt.Errorf("Unexpected message")

// Stash is returned by Actor.Receive when actor wants to delay the processing
// of the message until the next actor state change.
var Stash = fmt.Errorf("Stash")

// Terminate is returned by Actor.Receive when actor wants to terminate.
var Terminate = fmt.Errorf("Terminate")
