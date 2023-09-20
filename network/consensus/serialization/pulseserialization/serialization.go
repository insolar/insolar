package pulseserialization

import (
	"bytes"
	"encoding/binary"

	"github.com/insolar/insolar/pulse"
)

func Serialize(p pulse.Data) ([]byte, error) {
	buf := bytes.Buffer{}
	if err := binary.Write(&buf, binary.BigEndian, p); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Deserialize(b []byte) (pulse.Data, error) {
	d := pulse.Data{}
	if err := binary.Read(bytes.NewReader(b), binary.BigEndian, &d); err != nil {
		return d, err
	}

	return d, nil
}
