package handler

import (
	base "github.com/SangBejoo/Template/gen/proto"
)

// baseHandler implements the Base service
type baseHandler struct {
	base.UnimplementedBaseServer
}

// NewBaseHandler creates a new baseHandler instance
func NewBaseHandler() *baseHandler {
	return &baseHandler{}
}
