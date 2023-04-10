package app

import (
	"github.com/dimosha19/myvalidator"
	"homework6/internal/ads"
	myerrors "homework6/internal/errors"
)

type App interface {
	CreateAd(title string, text string, userID int) (*ads.Ad, error)
	UpdateAdStatus(adID int64, userID int64, published bool) (*ads.Ad, error)
	UpdateAd(adID int64, userID int64, title string, text string) (*ads.Ad, error)
}

type Repository interface {
	Add(ad ads.Ad) *ads.Ad
	Get(adID int64) *ads.Ad
	Update(adID int64, ad ads.Ad) (*ads.Ad, error)
	Size() int64
}

func NewApp(repo Repository) App {
	return &AppModel{repo}
}

type AppModel struct {
	repo Repository
}

func (p *AppModel) CreateAd(title string, text string, userID int) (*ads.Ad, error) {
	res := ads.Ad{ID: p.repo.Size(), Title: title, Text: text, AuthorID: int64(userID)}
	err := validator.Validate(res)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	t := p.repo.Add(res)
	return t, nil
}

func (p *AppModel) UpdateAdStatus(adID int64, userID int64, published bool) (*ads.Ad, error) {
	t := p.repo.Get(adID)
	temp := *t
	if t.AuthorID != userID {
		return nil, myerrors.ErrForbidden
	}
	temp.Published = published
	updated, err := p.repo.Update(adID, temp)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (p *AppModel) UpdateAd(adID int64, userID int64, title string, text string) (*ads.Ad, error) {
	t := p.repo.Get(adID)
	temp := *t
	if t.AuthorID != userID {
		return nil, myerrors.ErrForbidden
	}
	temp.Title = title
	temp.Text = text
	updated, err := p.repo.Update(adID, temp)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
