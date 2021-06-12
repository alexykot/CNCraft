package envelope

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/alexykot/cncraft/pkg/envelope/pb"
)

type E struct {
	pb.Envelope
}

func (e *E) GetMetaMap() map[string]string { return e.GetMeta() }
func (e *E) GetMetaKey(key string) (string, bool) {
	meta := e.GetMeta()
	if meta == nil {
		return "", false
	}

	val, ok := meta[key]
	return val, ok
}

func (e *E) Marshal() (bytes []byte, err error) {
	if e.Message == nil && e.ShardEvent == nil {
		return nil, errors.New("cannot marshal: this E is empty")
	}

	if bytes, err = proto.Marshal(e); err != nil {
		return nil, fmt.Errorf("failed to marshal E to protobuf bytes: %w", err)
	}

	return bytes, nil
}

func (e *E) Unmarshal(bytes []byte) error {
	if err := proto.Unmarshal(bytes, e); err != nil {
		return fmt.Errorf("failed to unmarshal protobuf bytes into E: %w", err)
	}
	return nil
}
