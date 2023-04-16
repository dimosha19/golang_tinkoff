package app

import (
	validator "github.com/dimosha19/myvalidator"
	"homework8/internal/ads"
	myerrors "homework8/internal/errors"
	"time"
)

type App interface {
	CreateAd(title string, text string, userID int) (*ads.Ad, error)
	UpdateAdStatus(adID int64, userID int64, published bool) (*ads.Ad, error)
	UpdateAd(adID int64, userID int64, title string, text string) (*ads.Ad, error)
	GetAds(pub string, author int64, date string) (*[]ads.Ad, error)
	GetAd(adID int64) (*ads.Ad, error)
}

type Repository interface {
	Add(ad ads.Ad) *ads.Ad
	Get(adID int64) (*ads.Ad, error)
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
	res := ads.Ad{ID: p.repo.Size(), Title: title, Text: text, AuthorID: int64(userID), PublishedTime: time.Now().UTC(), UpdateTime: time.Now().UTC()}
	err := validator.Validate(res)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	t := p.repo.Add(res)
	return t, nil
}

func (p *AppModel) UpdateAdStatus(adID int64, userID int64, published bool) (*ads.Ad, error) {
	t, e := p.repo.Get(adID)
	if e != nil {
		return nil, myerrors.ErrBadRequest
	}
	temp := *t
	if t.AuthorID != userID {
		return nil, myerrors.ErrForbidden
	}
	temp.Published = published
	temp.UpdateTime = time.Now().UTC()
	updated, err := p.repo.Update(adID, temp)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (p *AppModel) UpdateAd(adID int64, userID int64, title string, text string) (*ads.Ad, error) {
	t, e := p.repo.Get(adID)
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
	updated, err := p.repo.Update(adID, temp)
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
	for i := int64(0); i < p.repo.Size(); i++ {
		t, e := p.repo.Get(i)
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
	t, e := p.repo.Get(adID)
	if e != nil {
		return nil, myerrors.ErrBadRequest
	}
	return t, nil
}
