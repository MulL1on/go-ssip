package command

import (
	"encoding/binary"
	"time"
)

type Command struct {
	Type       uint32
	Payload    []byte
	RetryCount int
	Timer      time.Timer
}

func (cmd *Command) Encode() []byte {
	enc := make([]byte, 4+len(cmd.Payload))
	binary.LittleEndian.PutUint32(enc[:4], cmd.Type)
	copy(enc[4:], cmd.Payload)
	return enc
}

func (cmd *Command) Decode(b []byte) {
	cmd.Type = binary.LittleEndian.Uint32(b[:4])
	cmd.Payload = b[4:]
}

func GetCmdType(payload []byte) uint32 {
	return binary.LittleEndian.Uint32(payload[:4])
}

type AckClientIdPayload struct {
	ClientId int64
	UserId   int64
}

func (p *AckClientIdPayload) Encode() []byte {
	enc := make([]byte, 8+8)
	var index int

	index += binary.PutVarint(enc, p.ClientId)
	index += binary.PutVarint(enc[index:], p.UserId)
	return enc[:index]
}

func (p *AckClientIdPayload) Decode(b []byte) {
	var index int
	p.ClientId, index = binary.Varint(b)
	p.UserId, index = binary.Varint(b[index:])
}

type AckMsgPayload struct {
	Seq int64
}

func (p *AckMsgPayload) Encode() []byte {
	enc := make([]byte, 8)
	var index int
	// Seq
	index += binary.PutVarint(enc, p.Seq)
	return enc[:index]
}

func (p *AckMsgPayload) Decode(b []byte) {
	p.Seq, _ = binary.Varint(b)
}
