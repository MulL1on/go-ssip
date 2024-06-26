// Code generated by Kitex v0.7.0. DO NOT EDIT.
package groupservice

import (
	server "github.com/cloudwego/kitex/server"
	group "go-ssip/app/common/kitex_gen/group"
)

// NewServer creates a server.Server with the given handler and options.
func NewServer(handler group.GroupService, opts ...server.Option) server.Server {
	var options []server.Option

	options = append(options, opts...)

	svr := server.NewServer(options...)
	if err := svr.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	return svr
}
