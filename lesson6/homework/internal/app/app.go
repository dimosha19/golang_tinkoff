package app

import (
	"fmt"
	"homework6/internal/ads"
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
	t := p.repo.Add(res)
	return t, nil
}

func (p *AppModel) UpdateAdStatus(adID int64, userID int64, published bool) (*ads.Ad, error) {
	t := p.repo.Get(adID)
	t.Published = published
	if t.AuthorID != userID {
		return &ads.Ad{}, fmt.Errorf("forbidden")
	}
	return t, nil
}

func (p *AppModel) UpdateAd(adID int64, userID int64, title string, text string) (*ads.Ad, error) {
	t := p.repo.Get(adID)
	if t.AuthorID != userID {
		return &ads.Ad{}, fmt.Errorf("forbidden")
	}
	t.Text = text
	t.Title = title
	return t, nil
}
