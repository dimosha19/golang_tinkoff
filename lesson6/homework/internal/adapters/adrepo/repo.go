package adrepo

import (
	"homework6/internal/ads"
	"homework6/internal/app"
)

type scliceAd []ads.Ad

func New() app.Repository {
	res := scliceAd{}
	return &res
}

func (p *scliceAd) Add(ad ads.Ad) *ads.Ad {
	*p = append(*p, ad)
	return &(*p)[len(*p)-1]
}

func (p *scliceAd) Get(adID int64) *ads.Ad {
	res := (*p)[adID]
	return &res
}

func (p *scliceAd) Size() int64 {
	return int64(len(*p))
}

//Add(ad ads.Ad) (*ads.Ad, error)
//Get(adID int64) (*ads.Ad, error)
//Size() int64
