package adrepo

import (
	"fmt"
	validator "github.com/dimosha19/myvalidator"
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

func (p *scliceAd) Update(adID int64, ad ads.Ad) (*ads.Ad, error) {
	err := validator.Validate(ad)
	if err != nil {
		return nil, fmt.Errorf("bedrequest")
	}
	*p.Get(adID) = ad
	return &ad, nil
}
