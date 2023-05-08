package userrepo

import (
	validator "github.com/dimosha19/myvalidator"
	"homework9/internal/app"
	myerrors "homework9/internal/errors"
	"homework9/internal/users"
	"sync"
)

type mapUser struct {
	mx  *sync.Mutex
	r   map[int64]users.User
	idx int64
}

func New() app.UserRepository {
	mx := sync.Mutex{}
	res := mapUser{r: make(map[int64]users.User), mx: &mx, idx: 0}
	return &res
}

func (p *mapUser) Add(nickname string, email string) (*users.User, error) {
	res := users.User{ID: p.idx, Nickname: nickname, Email: email}
	p.idx++
	err := validator.Validate(res)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	p.mx.Lock()
	defer p.mx.Unlock()
	(p.r)[res.ID] = res
	a := (p.r)[res.ID]
	return &a, nil
}

func (p *mapUser) Get(userID int64) (*users.User, error) {
	p.mx.Lock()
	defer p.mx.Unlock()
	res, ok := (p.r)[userID]
	if ok {
		return &res, nil
	}
	return nil, myerrors.ErrBadRequest
}

func (p *mapUser) Update(userID int64, nickname string, email string, authorID int64) (*users.User, error) {
	temp := (p.r)[userID]
	if temp.ID != authorID {
		return nil, myerrors.ErrForbidden
	}
	temp.Nickname = nickname
	temp.Email = email
	p.mx.Lock()
	defer p.mx.Unlock()
	err := validator.Validate(temp)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	_, ok := (p.r)[userID]
	if !ok {
		return nil, myerrors.ErrBadRequest
	}
	(p.r)[userID] = temp
	return &temp, nil
}

func (p *mapUser) Idxs() []int64 {
	keys := make([]int64, len(p.r))

	i := 0
	for k := range p.r {
		keys[i] = k
		i++
	}
	return keys
}

func (p *mapUser) Delete(userID int64) bool {
	delete(p.r, userID)
	return true
}
