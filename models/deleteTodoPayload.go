package models

import (
	"github.com/noh4ck/mvptodo/schema"
	"golang.org/x/net/context"
)

var _ schema.DeleteTodoPayloadInterface = (*DeleteTodoPayload)(nil)

type DeleteTodoPayload struct {
	TodoID string
}

func (deleteTodoPayload DeleteTodoPayload) DeletedTodoIDField(
	ctx context.Context,
) (*string, error) {
	return &deleteTodoPayload.TodoID, nil
}
