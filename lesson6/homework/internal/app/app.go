package app

import (
	"fmt"
	"homework6/internal/ads"
	"homework6/internal/validator"
)

type App interface {
	CreateAd(title string, text string, userID int) (*ads.Ad, error)
	UpdateAdStatus(adID int64, userID int64, published bool) (*ads.Ad, error)
	UpdateAd(adID int64, userID int64, title string, text string) (*ads.Ad, error)
}

type Repository interface {
	Add(ad ads.Ad) *ads.Ad
	Get(adID int64) *ads.Ad
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
		return nil, fmt.Errorf("bedrequest")
	}
	t := p.repo.Add(res)
	return t, nil
}

func (p *AppModel) UpdateAdStatus(adID int64, userID int64, published bool) (*ads.Ad, error) {
	t := p.repo.Get(adID)
	temp := *t
	if t.AuthorID != userID {
		return nil, fmt.Errorf("forbidden")
	}
	temp.Published = published
	err := validator.Validate(temp)
	if err != nil {
		return nil, fmt.Errorf("bedrequest")
	}
	t = &temp
	return t, nil
}

func (p *AppModel) UpdateAd(adID int64, userID int64, title string, text string) (*ads.Ad, error) {
	t := p.repo.Get(adID)
	temp := *t
	if t.AuthorID != userID {
		return nil, fmt.Errorf("forbidden")
	}
	temp.Title = title
	temp.Text = text
	err := validator.Validate(temp)
	if err != nil {
		return nil, fmt.Errorf("bedrequest")
	}
	t = &temp
	return t, nil
}
