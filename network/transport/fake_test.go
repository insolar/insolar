// +build networktest

package transport

import (
	"context"
	"io"
	"log"
	"net"
)

type fakeTcp struct {
}

func (fakeTcp) Dial(ctx context.Context, address string) (io.ReadWriteCloser, error) {
	_, conn2 := net.Pipe()
	log.Println(conn2.LocalAddr().String())
	log.Println(conn2.RemoteAddr().String())

	return conn2, nil
}

func (fakeTcp) SetStreamHandler(processor StreamHandler) {
	panic("656")
}
