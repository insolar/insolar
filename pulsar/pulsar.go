package pulsar

import (
	"bytes"
	"fmt"
	"github.com/insolar/insolar/configuration"
	"net"
	"os"
	"strconv"
)

type Pulsar struct {
	Sock net.Listener
}

func NewPulsar(configuration configuration.Pulsar) *Pulsar {
	// Listen for incoming connections.
	l, err := net.Listen(configuration.Type, configuration.ListenAddress)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	return &Pulsar{Sock: l}
}

func (pulsar *Pulsar) Listen() {
	for {
		// Listen for an incoming connection.
		conn, err := pulsar.Sock.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Builds the message.
	message := "Hi, I received your message! It was "
	message += strconv.Itoa(reqLen)
	message += " bytes long and that's what it said: \""
	n := bytes.Index(buf, []byte{0})
	message += string(buf[:n-1])
	message += "\" ! Honestly I have no clue about what to do with your messages, so Bye Bye!\n"

	// Write the message in the connection channel.
	conn.Write([]byte(message))
	// Close the connection when you're done with it.
	conn.Close()
}
