package models

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/graphql-go/relay"
	"github.com/noh4ck/mvptodo/schema"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

var _ schema.QueryInterface = (*Root)(nil)

var _ schema.MutationInterface = (*Root)(nil)

var _ schema.RelayInterface = (*Root)(nil)

// Convert a slice or array of a specific type to array of interface{}
func ToIntf(s interface{}) []interface{} {
	v := reflect.ValueOf(s)
	// There is no need to check, we want to panic if it's not slice or array
	intf := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		intf[i] = v.Index(i).Interface()
	}
	return intf
}

type Root struct {
}

// Returns the loged in user based on a
// session token in the 'Authorization' header
func (root Root) ViewerQuery(
	ctx context.Context,
) (schema.UserInterface, error) {
	auth := Auth(ctx)

	if auth.Check() == false {
		return nil, nil
	}

	viewer := auth.Viewer()

	return viewer, nil
}

// Returns a connection of the authorized users
// list of todos
func (root Root) TodosQuery(
	ctx context.Context,
	args relay.ConnectionArguments,
) (*relay.Connection, error) {
	return nil, nil
}

// registerUser( input: RegisterUserInput! )
func (root Root) RegisterUserMutation(
	ctx context.Context,
	input schema.RegisterUserInputStruct,
) (schema.RegisterUserPayloadInterface, error) {
	return nil, nil
}

// loginUser( input: LoginUserInput! )
func (root Root) LoginUserMutation(
	ctx context.Context,
	input schema.LoginUserInputStruct,
) (schema.LoginUserPayloadInterface, error) {
	sess := Session(ctx)
	// auth := Auth(ctx)

	db := sess.DB
	user := User{}
	db.Select("*").From(usersTable).
		Where("email = ?", input.Email).
		Load(&user)
	spew.Dump(user)

	if VerifyPassword(*input.Password, user.Password) == false {
		return nil, errors.New("{\"password\": \"Wrong password\"}")
	}

	sessionid := fmt.Sprint(uuid.NewV4())

	session := AuthSession{
		UserID:  user.ID,
		Token:   sessionid,
		Expires: time.Now(),
	}

	_, err := db.InsertInto("sessions").
		Columns("user_id", "token", "expires").
		Record(&session).Exec()

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Could not authenticate user")
	}

	cookie := http.Cookie{
		Name:     "session",
		Value:    sessionid,
		HttpOnly: true,
	}

	http.SetCookie(sess.Response, &cookie)

	return LoginUserPayload{User: user}, nil
}
func (root Root) LogoutUserMutation(
	ctx context.Context,
) (*string, error) {
	return nil, nil
}

// addTodo( input: AddTodoInput! )
func (root Root) AddTodoMutation(
	ctx context.Context,
	input schema.AddTodoInputStruct,
) (schema.AddTodoPayloadInterface, error) {
	auth := Auth(ctx)
	if auth.Check() == false {
		return nil, errors.New("Not loged in")
	}
	db := Session(ctx).DB

	viewer := auth.Viewer()

	todo := Todo{
		UserID:      viewer.ID,
		Title:       *input.Title,
		Description: *input.Description,
		Status:      *input.Status,
	}

	_, err := db.InsertInto(todosTable).
		Columns("user_id", "title", "description", "status").
		Record(&todo).Exec()

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Could not create new todo")
	}

	return AddTodoPayload{Todo: todo}, nil
}

// deleteTodo( input: DeleteTodoInput! )
func (root Root) DeleteTodoMutation(
	ctx context.Context,
	input schema.DeleteTodoInputStruct,
) (schema.DeleteTodoPayloadInterface, error) {
	auth := Auth(ctx)
	if auth.Check() == false {
		return nil, errors.New("Not loged in")
	}

	db := Session(ctx).DB
	viewer := auth.Viewer()

	id := relay.FromGlobalID(*input.TodoID)

	res, err := db.DeleteFrom(todosTable).
		Where("user_id = ?", viewer.ID).
		Where("id = ?", id.ID).Exec()

	affect, _ := res.RowsAffected()
	if err != nil || affect == 0 {
		fmt.Println(err)
		return nil, errors.New("Could not delete the todo")
	}

	return DeleteTodoPayload{TodoID: *input.TodoID}, nil
}

type null struct{}

// updateTodo( input: UpdateTodoInput! )
func (root Root) UpdateTodoMutation(
	ctx context.Context,
	input schema.UpdateTodoInputStruct,
) (schema.UpdateTodoPayloadInterface, error) {
	auth := Auth(ctx)
	if auth.Check() == false {
		return nil, errors.New("Not loged in")
	}

	viewer := auth.Viewer()

	db := Session(ctx).DB
	id := relay.FromGlobalID(*input.TodoID)
	if id == nil {
		return nil, nil
	}

	rows, err := db.Select("1").
		From(todosTable).
		Where("user_id = ?", viewer.ID).
		Where("id = ?", id.ID).Load(&null{})

	if rows == 0 {
		return nil, errors.New("Could not find todo")
	}

	update := db.Update(todosTable)
	if input.Title != nil {
		update.Set("title", *input.Title)
	}
	if input.Description != nil {
		update.Set("description", *input.Description)
	}
	if input.Status != nil {
		update.Set("status", *input.Status)
	}

	_, err = update.Where("id = ?", id.ID).Exec()
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Could not update todo")
	}

	var todo Todo
	db.Select("*").From(todosTable).
		Where("id = ?", id.ID).Load(&todo)

	return UpdateTodoPayload{Todo: todo}, nil
}

func (root Root) ResolveUserNode(ctx context.Context, id string) (schema.UserInterface, error) {
	return nil, nil
}

func (root Root) ResolveTodoNode(ctx context.Context, id string) (schema.TodoInterface, error) {
	return nil, nil
}
