package models

import (
	"strconv"

	"github.com/graphql-go/relay"
	"github.com/noh4ck/mvptodo/schema"
	"golang.org/x/net/context"
)

var _ schema.TodoInterface = (*Todo)(nil)

const todosTable = "todos"

type Todo struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	Status      bool
}

func (todo Todo) IdField(
	ctx context.Context,
) (*string, error) {
	id := strconv.FormatInt(todo.ID, 10)
	return &id, nil
}

// The ID of the owner (User)
func (todo Todo) OwnerField(
	ctx context.Context,
) (*string, error) {
	id := strconv.FormatInt(todo.UserID, 10)
	globalid := relay.ToGlobalID("User", id)
	return &globalid, nil
}

// Title of the todo
func (todo Todo) TitleField(
	ctx context.Context,
) (*string, error) {
	return &todo.Title, nil
}

// Descriptoin of the todo
func (todo Todo) DescriptionField(
	ctx context.Context,
) (*string, error) {
	return &todo.Description, nil
}

// The status of the todo (completed or not)
func (todo Todo) StatusField(
	ctx context.Context,
) (*bool, error) {
	return &todo.Status, nil
}
