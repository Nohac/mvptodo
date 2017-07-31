package models

import (
	"fmt"
	"strconv"

	"github.com/graphql-go/relay"
	"github.com/noh4ck/mvptodo/schema"
	"golang.org/x/net/context"
)

var _ schema.UserInterface = (*User)(nil)

var usersTable = "users"

type User struct {
	ID       int64
	Name     string
	Email    string
	Password string
}

func (user User) IdField(
	ctx context.Context,
) (*string, error) {
	id := strconv.FormatInt(user.ID, 10)
	return &id, nil
}
func (user User) NameField(
	ctx context.Context,
) (*string, error) {
	return &user.Name, nil
}
func (user User) EmailField(
	ctx context.Context,
) (*string, error) {
	return &user.Email, nil
}
func (user User) TodosField(
	ctx context.Context,
	args relay.ConnectionArguments,
) (*relay.Connection, error) {
	fmt.Println("Fetching todos")
	db := Session(ctx).DB

	todos := []Todo{}
	db.Select("*").From(todosTable).
		Where("user_id = ?", user.ID).
		Load(&todos)

	todoif := ToIntf(todos)

	return relay.ConnectionFromArray(todoif, args), nil
}
