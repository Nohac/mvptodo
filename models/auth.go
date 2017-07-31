package models

import (
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type AuthType struct {
	CookieName string
	viewer     *User
	session    *SessionType
}

type AuthSession struct {
	UserID  int64
	Token   string
	Expires time.Time
}

type AuthCredentials map[string]string

type AuthorizedViewer interface {
	GetID()
	Object() interface{}
}

type Authenticatable interface {
	Authenticate(cookie *http.Cookie) bool
	CreateSession(creds AuthCredentials) string
	DestroySession() bool
	Viewer() AuthorizedViewer
}

func NewAuth(s *SessionType) *AuthType {
	return &AuthType{session: s}
}

func PasswordHash(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("Could not hash password")
	}

	return string(pw), nil
}

func VerifyPassword(pw string, pwh string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwh), []byte(pw))
	if err != nil {
		return false
	}

	return true
}

func (auth *AuthType) Authenticate() bool {
	cookie, err := auth.session.Request.Cookie("session")
	if err != nil {
		return false
	}

	token := cookie.Value
	db := auth.session.DB

	var viewer User
	var userid struct {
		UserID int64
	}

	row, err := db.
		Select("token, user_id").
		From("sessions").
		Where("token = ?", token).
		Load(&userid)

	if err != nil || row == 0 {
		return false
	}

	_, err = db.
		Select("*").
		From("users").
		Where("id = ?", userid.UserID).
		Load(&viewer)

	auth.viewer = &viewer

	return true
}

func (auth *AuthType) SignOut() bool {
	if auth.Check() == false {
		return false
	}

	viewer := auth.Viewer()
	cookie, err := auth.session.Request.Cookie("session")

	token := cookie.Value
	db := auth.session.DB

	res, err := db.DeleteFrom(sessionsTable).
		Where("user_id = ?", viewer.ID).
		Where("token = ?", token).
		Exec()

	if err != nil {
		return false
	}

	affect, _ := res.RowsAffected()
	if affect > 0 {
		return true
	}

	return false
}

func (auth *AuthType) Viewer() *User {
	return auth.viewer
}

func (auth *AuthType) Check() bool {
	if auth.viewer != nil {
		return true
	}

	return false
}

func Auth(ctx context.Context) *AuthType {
	val := ctx.Value("auth")
	if val == nil {
		panic("`auth` is not provided in the context")
	}

	session, ok := val.(*AuthType)
	if ok == false {
		panic("Could not cast the provided `auth` to AuthType")
	}

	return session
}
