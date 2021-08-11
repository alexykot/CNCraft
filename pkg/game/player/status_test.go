package player

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexykot/cncraft/pkg/envelope/pb"
)

// TestDiggingActionPbEnum asserts that the pb enum and defined constants are aligned at all times.
func TestDiggingActionPbEnum(t *testing.T) {
	for pbName, val := range pb.PlayerDigging_Action_value {
		pbName = strings.ToLower(pbName)
		pbName = strings.ReplaceAll(pbName, "_", " ")
		pbName = strings.Title(pbName)
		pbName = strings.ReplaceAll(pbName, " ", "")

		constVal := DiggingAction(val)
		constName := constVal.String()
		assert.Equal(t, pbName, constName)
	}
}
