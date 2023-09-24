package gid

import "github.com/bwmarrin/snowflake"

type IDGenerator struct {
}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

func (g *IDGenerator) CreateUUID() int64 {
	sf, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return sf.Generate().Int64()
}
