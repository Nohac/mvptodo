package models

import (
	"github.com/noh4ck/mvptodo/schema"
	"golang.org/x/net/context"
)

var _ schema.AddTodoPayloadInterface = (*AddTodoPayload)(nil)

type AddTodoPayload struct {
	Todo Todo
}

func (addTodoPayload AddTodoPayload) TodoField(
	ctx context.Context,
) (schema.TodoInterface, error) {
	return addTodoPayload.Todo, nil
}
