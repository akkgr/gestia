package main

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/context"
)

type TokenAuth struct {
	handler             http.Handler
	store               TokenStore
	getter              TokenGetter
	UnauthorizedHandler http.HandlerFunc
}

type TokenGetter interface {
	GetTokenFromRequest(req *http.Request) string
}

type TokenStore interface {
	CheckToken(token string) (Token, error)
}

type Token interface {
	IsExpired() bool
	fmt.Stringer
	ClaimGetter
}

type ClaimSetter interface {
	SetClaim(string, interface{}) ClaimSetter
}

type ClaimGetter interface {
	Claims(string) interface{}
}

func DefaultUnauthorizedHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(401)
	fmt.Fprint(w, "unauthorized")
}

type QueryStringTokenGetter struct {
	Parameter string
}

func (q QueryStringTokenGetter) GetTokenFromRequest(req *http.Request) string {
	return req.URL.Query().Get(q.Parameter)
}

func NewQueryStringTokenGetter(parameter string) *QueryStringTokenGetter {
	return &QueryStringTokenGetter{
		Parameter: parameter,
	}
}

type BearerGetter struct{}

func (b *BearerGetter) GetTokenFromRequest(req *http.Request) string {
	authStr := req.Header.Get("Authorization")
	if !strings.HasPrefix(authStr, "Bearer ") {
		return ""
	}
	return authStr[7:]
}

/*
	Returns a TokenAuth object implemting Handler interface

	if a handler is given it proxies the request to the handler

	if a unauthorizedHandler is provided, unauthorized requests will be handled by this HandlerFunc,
	otherwise a default unauthorized handler is used.

	store is the TokenStore that stores and verify the tokens
*/
func NewTokenAuth(handler http.Handler, unauthorizedHandler http.HandlerFunc, store TokenStore, getter TokenGetter) *TokenAuth {
	t := &TokenAuth{
		handler:             handler,
		store:               store,
		getter:              getter,
		UnauthorizedHandler: unauthorizedHandler,
	}
	if t.getter == nil {
		t.getter = &BearerGetter{} //NewQueryStringTokenGetter("token")
	}
	if t.UnauthorizedHandler == nil {
		t.UnauthorizedHandler = DefaultUnauthorizedHandler
	}
	return t
}

/* wrap a HandlerFunc to be authenticated */
func (t *TokenAuth) HandleFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token, err := t.Authenticate(req)
		if err != nil {
			t.UnauthorizedHandler.ServeHTTP(w, req)
			return
		}
		context.Set(req, "token", token)
		handlerFunc.ServeHTTP(w, req)
	}
}

func (t *TokenAuth) Authenticate(req *http.Request) (Token, error) {
	strToken := t.getter.GetTokenFromRequest(req)
	if strToken == "" {
		return nil, errors.New("token required")
	}
	token, err := t.store.CheckToken(strToken)
	if err != nil {
		return nil, errors.New("Invalid token")
	}
	return token, nil
}

/* implement Handler */
func (t *TokenAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token, err := t.Authenticate(req)
	if err != nil {
		t.UnauthorizedHandler.ServeHTTP(w, req)
		return
	}
	context.Set(req, "token", token)
	t.handler.ServeHTTP(w, req)
	context.Clear(req)
}

func Get(req *http.Request) Token {
	return context.Get(req, "token").(Token)
}

//Memory Store
type MemoryTokenStore struct {
	tokens   map[string]*MemoryToken
	idTokens map[string]*MemoryToken
	salt     string
}

type MemoryToken struct {
	ExpireAt time.Time
	Token    string
	Id       string
}

func (t *MemoryToken) IsExpired() bool {
	return time.Now().After(t.ExpireAt)
}

func (t *MemoryToken) String() string {
	return t.Token
}

/* lookup 'exp' or 'id' */
func (t *MemoryToken) Claims(key string) interface{} {
	switch key {
	case "exp":
		return t.ExpireAt
	case "id":
		return t.Id
	default:
		return nil
	}
}

func (s *MemoryTokenStore) generateToken(id string) []byte {
	hash := sha1.New()
	now := time.Now()
	timeStr := now.Format(time.ANSIC)
	hash.Write([]byte(timeStr))
	hash.Write([]byte(id))
	hash.Write([]byte("salt"))
	return hash.Sum(nil)
}

/* returns a new token with specific id */
func (s *MemoryTokenStore) NewToken(id interface{}) *MemoryToken {
	strId := id.(string)
	bToken := s.generateToken(strId)
	strToken := base64.URLEncoding.EncodeToString(bToken)
	t := &MemoryToken{
		ExpireAt: time.Now().Add(time.Minute * 30),
		Token:    strToken,
		Id:       strId,
	}
	oldT, ok := s.idTokens[strId]
	if ok {
		delete(s.tokens, oldT.Token)
	}
	s.tokens[strToken] = t
	s.idTokens[strId] = t
	return t
}

/* Create a new memory store */
func NewMemoryTokenStore(salt string) *MemoryTokenStore {
	return &MemoryTokenStore{
		salt:     salt,
		tokens:   make(map[string]*MemoryToken),
		idTokens: make(map[string]*MemoryToken),
	}

}

func (s *MemoryTokenStore) CheckToken(strToken string) (Token, error) {
	t, ok := s.tokens[strToken]
	if !ok {
		return nil, errors.New("Failed to authenticate")
	}
	if t.ExpireAt.Before(time.Now()) {
		delete(s.tokens, strToken)
		return nil, errors.New("Token expired")
	}
	return t, nil
}
