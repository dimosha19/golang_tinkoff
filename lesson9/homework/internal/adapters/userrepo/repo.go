package userrepo

import (
	validator "github.com/dimosha19/myvalidator"
	"homework9/internal/app"
	myerrors "homework9/internal/errors"
	"homework9/internal/users"
	"sync"
)

type scliceUser struct {
	mx *sync.Mutex
	r  []users.User
}

func New() app.UserRepository {
	mx := sync.Mutex{}
	res := scliceUser{mx: &mx}
	return &res
}

func (p *scliceUser) Add(user users.User) *users.User {
	p.mx.Lock()
	defer p.mx.Unlock()
	(*p).r = append((*p).r, user)
	return &(*p).r[len((*p).r)-1]
}

func (p *scliceUser) Get(userID int64) (*users.User, error) {
	p.mx.Lock()
	defer p.mx.Unlock()
	if userID < p.Size() {
		res := (*p).r[userID]
		return &res, nil
	}
	return nil, myerrors.ErrBadRequest
}

func (p *scliceUser) Size() int64 {
	return int64(len((*p).r))
}

func (p *scliceUser) Update(userID int64, user users.User) (*users.User, error) {
	p.mx.Lock()
	defer p.mx.Unlock()
	err := validator.Validate(user)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	if userID >= p.Size() {
		return nil, myerrors.ErrBadRequest
	}
	(*p).r[userID] = user
	return &user, nil
}
