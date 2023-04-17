package userrepo

import (
	validator "github.com/dimosha19/myvalidator"
	"homework8/internal/app"
	myerrors "homework8/internal/errors"
	"homework8/internal/users"
)

type scliceUser []users.User

func New() app.UserRepository {
	res := scliceUser{}
	//mx := *sync.RWMutex //TODO mutex
	return &res
}

func (p *scliceUser) Add(user users.User) *users.User {
	*p = append(*p, user)
	return &(*p)[len(*p)-1]
}

func (p *scliceUser) Get(userID int64) (*users.User, error) {
	if userID < (*p).Size() {
		res := (*p)[userID]
		return &res, nil
	}
	return nil, myerrors.ErrBadRequest
}

func (p *scliceUser) Size() int64 {
	return int64(len(*p))
}

func (p *scliceUser) Update(userID int64, user users.User) (*users.User, error) {
	err := validator.Validate(user)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	if userID >= (*p).Size() {
		return nil, myerrors.ErrBadRequest
	}
	(*p)[userID] = user
	return &user, nil
}
