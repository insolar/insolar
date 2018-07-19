/*
 *    Copyright 2018 INS Ecosystem
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

/*
Package message provides insolar messaging protocol and serialization layer.

Message can be created with shortcut:

	senderAddr, _ := node.NewAddress("127.0.0.1:1337")
	receiverAddr, _ := node.NewAddress("127.0.0.1:1338")

	sender := node.NewNode(senderAddr)
	receiver := node.NewNode(receiverAddr)

	msg := message.NewPingMessage(sender, receiver)

	// do something with message


Or with builder:

	builder := message.NewBuilder()

	senderAddr, _ := node.NewAddress("127.0.0.1:1337")
	receiverAddr, _ := node.NewAddress("127.0.0.1:1338")

	sender := node.NewNode(senderAddr)
	receiver := node.NewNode(receiverAddr)

	msg := builder.
		Sender(sender).
		Receiver(receiver).
		Type(message.TypeFindNode).
		Request(&message.RequestDataFindNode{}).
		Build()

	// do something with message


Message may be serialized:

	msg := &message.Message{}
	serialized, err := message.SerializeMessage(msg)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(serialized)


And deserialized therefore:

	var buffer bytes.Buffer

	// Fill buffer somewhere

	msg, err := message.DeserializeMessage(buffer)

	if err != nil {
		panic(err.Error())
	}

	// do something with message

*/
package message
