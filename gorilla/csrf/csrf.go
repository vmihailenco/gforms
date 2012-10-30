package csrf

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io"

	"github.com/vmihailenco/gforms"

	"github.com/gorilla/sessions"
)

const (
	fieldName = "_csrf"
)

var (
	ErrCSRFMissing      = errors.New("CSRF token is missing")
	ErrCSRFDoesNotMatch = errors.New("CSRF token does not match")
)

func Token(session *sessions.Session) string {
	tokenI, ok := session.Values[fieldName]
	if !ok {
		return ""
	}
	delete(session.Values, fieldName)

	token, ok := tokenI.(string)
	if !ok {
		return ""
	}
	return token
}

func GenerateToken(session *sessions.Session) string {
	b := make([]byte, 32)
	_, _ = io.ReadFull(rand.Reader, b)
	token := base64.URLEncoding.EncodeToString(b)

	session.Values[fieldName] = token
	return token
}

type Field struct {
	gforms.BaseField
	currentToken string
	newToken     string
	session      *sessions.Session
}

func NewField(session *sessions.Session) *Field {
	f := &Field{
		session:      session,
		currentToken: Token(session),
		newToken:     GenerateToken(session),
	}
	f.SetWidget(gforms.NewHiddenWidget())
	f.SetName(fieldName)

	return f
}

func (f *Field) Value() string {
	return ""
}

func (f *Field) Validate(rawValue interface{}) error {
	value, ok := rawValue.(string)
	if !ok {
		return fmt.Errorf("type %T is not supported")
	}

	if f.currentToken == "" {
		return ErrCSRFMissing
	}
	if value != f.currentToken {
		return ErrCSRFDoesNotMatch
	}

	return nil
}

func (f *Field) SetInitial(initial string) {
	f.newToken = initial
}

func (f *Field) Render(attrs ...string) template.HTML {
	return f.Widget().Render(attrs, f.newToken)
}
