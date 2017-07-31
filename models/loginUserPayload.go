package models

import (
	"github.com/noh4ck/mvptodo/schema"
	"golang.org/x/net/context"
)

var _ schema.LoginUserPayloadInterface = (*LoginUserPayload)(nil)

type LoginUserPayload struct {
	User User
}

func (loginUserPayload LoginUserPayload) UserField(
	ctx context.Context,
) (schema.UserInterface, error) {
	return loginUserPayload.User, nil
}
