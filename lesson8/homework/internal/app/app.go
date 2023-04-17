package app

import (
	validator "github.com/dimosha19/myvalidator"
	"homework8/internal/ads"
	myerrors "homework8/internal/errors"
	"homework8/internal/users"
	"time"
)

type App interface {
	CreateAd(title string, text string, userID int) (*ads.Ad, error)
	UpdateAdStatus(adID int64, userID int64, published bool) (*ads.Ad, error)
	UpdateAd(adID int64, userID int64, title string, text string) (*ads.Ad, error)
	GetAds(pub string, author int64, date string) (*[]ads.Ad, error)
	GetAd(adID int64) (*ads.Ad, error)

	CreateUser(nickname string, email string) (*users.User, error)
	GetUser(userID int64) (*users.User, error)
	UpdateUser(userID int64, nickname string, email string, authorID int64) (*users.User, error)
}

type AdRepository interface {
	Add(ad ads.Ad) *ads.Ad
	Get(adID int64) (*ads.Ad, error)
	Update(adID int64, ad ads.Ad) (*ads.Ad, error)
	Size() int64
}

type UserRepository interface {
	Add(user users.User) *users.User
	Get(userID int64) (*users.User, error)
	Update(userID int64, user users.User) (*users.User, error)
	Size() int64
}

func NewApp(adrepo AdRepository, userrepo UserRepository) App {
	return &AppModel{adrepo, userrepo}
}

type AppModel struct {
	adrepo   AdRepository
	userrepo UserRepository
}

func (p *AppModel) CreateAd(title string, text string, userID int) (*ads.Ad, error) {
	res := ads.Ad{ID: p.adrepo.Size(), Title: title, Text: text, AuthorID: int64(userID), PublishedTime: time.Now().UTC(), UpdateTime: time.Now().UTC()}
	err := validator.Validate(res)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	t := p.adrepo.Add(res)
	return t, nil
}

func (p *AppModel) UpdateAdStatus(adID int64, userID int64, published bool) (*ads.Ad, error) {
	t, e := p.adrepo.Get(adID)
	if e != nil {
		return nil, myerrors.ErrBadRequest
	}
	temp := *t
	if t.AuthorID != userID {
		return nil, myerrors.ErrForbidden
	}
	temp.Published = published
	temp.UpdateTime = time.Now().UTC()
	updated, err := p.adrepo.Update(adID, temp)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (p *AppModel) UpdateAd(adID int64, userID int64, title string, text string) (*ads.Ad, error) {
	t, e := p.adrepo.Get(adID)
	if e != nil {
		return nil, myerrors.ErrBadRequest
	}
	temp := *t
	if t.AuthorID != userID {
		return nil, myerrors.ErrForbidden
	}
	temp.Title = title
	temp.Text = text
	temp.UpdateTime = time.Now().UTC()
	updated, err := p.adrepo.Update(adID, temp)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (p *AppModel) GetAds(pub string, author int64, date string) (*[]ads.Ad, error) {
	var res []ads.Ad
	var td time.Time
	var err error
	if date != "all" {
		td, err = time.Parse("02-01-06", date)
		if err != nil {
			return nil, myerrors.ErrBadRequest
		}
	}
	for i := int64(0); i < p.adrepo.Size(); i++ {
		t, e := p.adrepo.Get(i)
		if e != nil {
			return nil, myerrors.ErrBadRequest
		}
		if t.Published && pub == "true" || !t.Published && pub == "false" || pub == "all" {
			if date == "all" || td.Equal(t.PublishedTime) {
				if author == -1 || author == t.AuthorID {
					res = append(res, *t)
				}
			}
		}
	}
	return &res, nil
}

func (p *AppModel) GetAd(adID int64) (*ads.Ad, error) {
	t, e := p.adrepo.Get(adID)
	if e != nil {
		return nil, myerrors.ErrBadRequest
	}
	return t, nil
}

func (p *AppModel) CreateUser(nickname string, email string) (*users.User, error) {
	res := users.User{ID: p.userrepo.Size(), Nickname: nickname, Email: email}
	err := validator.Validate(res)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	t := p.userrepo.Add(res)
	return t, nil
}

func (p *AppModel) GetUser(userID int64) (*users.User, error) {
	t, e := p.userrepo.Get(userID)
	if e != nil {
		return nil, myerrors.ErrBadRequest
	}
	return t, nil
}

func (p *AppModel) UpdateUser(userID int64, nickname string, email string, authorID int64) (*users.User, error) {
	t, e := p.userrepo.Get(userID)
	if e != nil {
		return nil, myerrors.ErrBadRequest
	}
	temp := *t
	if t.ID != authorID {
		return nil, myerrors.ErrForbidden
	}
	temp.Nickname = nickname
	temp.Email = email
	updated, err := p.userrepo.Update(userID, temp)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
