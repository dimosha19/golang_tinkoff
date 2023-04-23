package app

import (
	"homework9/internal/ads"
	myerrors "homework9/internal/errors"
	"homework9/internal/users"
	"time"
)

type App interface {
	CreateAd(title string, text string, userID int64) (*ads.Ad, error)
	UpdateAdStatus(adID int64, userID int64, published bool) (*ads.Ad, error)
	UpdateAd(adID int64, userID int64, title string, text string) (*ads.Ad, error)
	GetAds(pub string, author int64, date string, title string) (*[]ads.Ad, error)
	GetAd(adID int64) (*ads.Ad, error)
	DeleteAd(adID int64, userID int64) error

	CreateUser(nickname string, email string) (*users.User, error)
	GetUser(userID int64) (*users.User, error)
	UpdateUser(userID int64, nickname string, email string, authorID int64) (*users.User, error)
	DeleteUser(userID int64) error
}

type AdRepository interface {
	Add(title string, text string, userID int64) (*ads.Ad, error)
	Get(adID int64) (*ads.Ad, error)
	Update(adID int64, userID int64, title string, text string, published bool) (*ads.Ad, error)
	Delete(adID int64, userID int64) bool
	Idxs() []int64
}

type UserRepository interface {
	Add(nickname string, email string) (*users.User, error)
	Get(userID int64) (*users.User, error)
	Update(userID int64, nickname string, email string, authorID int64) (*users.User, error)
	Delete(userID int64) bool
	Idxs() []int64
}

func NewApp(adrepo AdRepository, userrepo UserRepository) App {
	return &AppModel{adrepo, userrepo}
}

type AppModel struct {
	adrepo   AdRepository
	userrepo UserRepository
}

func (p *AppModel) CreateAd(title string, text string, userID int64) (*ads.Ad, error) {
	_, err := p.userrepo.Get(userID)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	t, err := p.adrepo.Add(title, text, userID)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	return t, nil
}

func (p *AppModel) DeleteAd(adID int64, userID int64) error {
	if p.adrepo.Delete(adID, userID) {
		return nil
	}
	return myerrors.ErrForbidden
}

func (p *AppModel) UpdateAdStatus(adID int64, userID int64, published bool) (*ads.Ad, error) {
	_, err := p.userrepo.Get(userID)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	t, e := p.adrepo.Get(adID)
	if e != nil {
		return nil, myerrors.ErrBadRequest
	}
	updated, err := p.adrepo.Update(adID, userID, t.Title, t.Text, published)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (p *AppModel) UpdateAd(adID int64, userID int64, title string, text string) (*ads.Ad, error) {
	_, err := p.userrepo.Get(userID)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	t, e := p.adrepo.Get(adID)
	if e != nil {
		return nil, myerrors.ErrBadRequest
	}
	if t.AuthorID != userID {
		return nil, myerrors.ErrForbidden
	}
	updated, err := p.adrepo.Update(adID, userID, title, text, t.Published)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func publishedPredicate(pub string, ad ads.Ad) bool {
	if ad.Published && pub == "true" || !ad.Published && pub == "false" || pub == "all" {
		return true
	}
	return false
}

func authorPredicate(author int64, ad ads.Ad) bool {
	if author == -1 {
		return true
	}
	if author == ad.AuthorID {
		return true
	}
	return false
}

func datePredicate(date string, ad ads.Ad) bool {
	if date == "all" {
		return true
	}
	td, _ := time.Parse("02-01-06", date)
	ady, adm, add := ad.PublishedTime.Date()
	y, m, d := td.Date()
	if y == ady && m == adm && d == add {
		return true
	}
	return false
}

func titlePredicate(title string, ad ads.Ad) bool {
	return ad.Title == title || title == ""
}

func adsPred(pub string, author int64, date string, title string, ad ads.Ad) bool {
	return publishedPredicate(pub, ad) && datePredicate(date, ad) && authorPredicate(author, ad) && titlePredicate(title, ad)
}

func (p *AppModel) GetAds(pub string, author int64, date string, title string) (*[]ads.Ad, error) {
	var res []ads.Ad
	arr := p.adrepo.Idxs()
	for k := range arr {
		t, e := p.adrepo.Get(arr[k])
		if e != nil {
			return nil, myerrors.ErrBadRequest
		}
		if adsPred(pub, author, date, title, *t) {
			res = append(res, *t)
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
	res, err := p.userrepo.Add(nickname, email)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	return res, nil
}

func (p *AppModel) GetUser(userID int64) (*users.User, error) {
	t, e := p.userrepo.Get(userID)
	if e != nil {
		return nil, myerrors.ErrBadRequest
	}
	return t, nil
}

func (p *AppModel) UpdateUser(userID int64, nickname string, email string, authorID int64) (*users.User, error) {
	_, e := p.userrepo.Get(userID)
	if e != nil {
		return nil, myerrors.ErrBadRequest
	}
	updated, err := p.userrepo.Update(userID, nickname, email, authorID)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (p *AppModel) DeleteUser(userID int64) error {
	if p.userrepo.Delete(userID) {
		return nil
	}
	return myerrors.ErrForbidden
}
