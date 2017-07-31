package models

import (
	"github.com/noh4ck/mvptodo/schema"
	"golang.org/x/net/context"
)

var _ schema.RegisterUserPayloadInterface = (*RegisterUserPayload)(nil)

type RegisterUserPayload struct {
}

func (registerUserPayload RegisterUserPayload) UserField(
	ctx context.Context,
) (schema.UserInterface, error) {
	return nil, nil
}
