package models

import (
	"github.com/noh4ck/mvptodo/schema"
	"golang.org/x/net/context"
)

var _ schema.UpdateTodoPayloadInterface = (*UpdateTodoPayload)(nil)

type UpdateTodoPayload struct {
	Todo Todo
}

func (updateTodoPayload UpdateTodoPayload) TodoField(
	ctx context.Context,
) (schema.TodoInterface, error) {
	return updateTodoPayload.Todo, nil
}
