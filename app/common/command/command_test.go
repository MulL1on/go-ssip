package command

import (
	"github.com/stretchr/testify/assert"
	"go-ssip/app/common/consts"
	"testing"
)

func TestCommand(t *testing.T) {
	cmd := &Command{
		Type:    consts.CommandTypeGetMsg,
		Payload: []byte("foo"),
	}
	enc := cmd.Encode()
	assert.NotNil(t, enc)

	newCmd := &Command{}
	newCmd.Decode(enc)
	assert.Equal(t, cmd, newCmd)
}

func TestAckMsgPayload(t *testing.T) {
	payload := &AckMsgPayload{
		Seq: 1,
	}
	enc := payload.Encode()
	assert.NotNil(t, enc)
	newPayload := &AckMsgPayload{}
	newPayload.Decode(enc)
	assert.Equal(t, payload, newPayload)
}

func TestAckClientIdPayload(t *testing.T) {
	payload := &AckClientIdPayload{
		ClientId: 1,
		UserId:   1705062900332236800,
	}
	enc := payload.Encode()
	assert.NotNil(t, enc)
	newPayload := &AckClientIdPayload{}
	newPayload.Decode(enc)
	assert.Equal(t, payload, newPayload)
}
