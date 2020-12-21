package envelope

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/alexykot/cncraft/pkg/envelope/pb"
)

type E interface {
	GetMetaMap() map[string]string
	GetMetaKey(string) (string, bool)
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

type envelope struct {
	proto *pb.Envelope
}

func NewWithHandshake(hs *pb.Handshake, meta map[string]string) E {
	return &envelope{
		proto: &pb.Envelope{
			Meta: meta,
			Message: &pb.Envelope_Handshake{Handshake: hs},
		},
	}
}

func NewEmpty() E {
	return &envelope{}
}

func (e *envelope) GetMetaMap() map[string]string { return e.proto.GetMeta() }
func (e *envelope) GetMetaKey(key string) (string, bool) {
	meta := e.proto.GetMeta()
	if meta == nil {
		return "", false
	}

	val, ok := meta[key]
	return val, ok
}

func (e *envelope) Marshal() ([]byte, error) {
	if e.proto == nil {
		return nil, errors.New("cannot marshal: this envelope is empty")
	}

	bytes, err := proto.Marshal(e.proto)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal envelope to protobuf bytes: %w", err)
	}
	return bytes, nil
}

func (e *envelope) Unmarshal(bytes []byte) error {
	if err := proto.Unmarshal(bytes, e.proto); err != nil {
		return fmt.Errorf("failed to unmarshal protobuf bytes into envelope: %w", err)
	}
	return nil
}
