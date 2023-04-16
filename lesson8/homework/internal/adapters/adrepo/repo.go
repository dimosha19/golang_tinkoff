package adrepo

import (
	validator "github.com/dimosha19/myvalidator"
	"homework8/internal/ads"
	"homework8/internal/app"
	myerrors "homework8/internal/errors"
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

func (p *scliceAd) Get(adID int64) (*ads.Ad, error) {
	if adID < (*p).Size() {
		res := (*p)[adID]
		return &res, nil
	}
	return nil, myerrors.ErrBadRequest
}

func (p *scliceAd) Size() int64 {
	return int64(len(*p))
}

func (p *scliceAd) Update(adID int64, ad ads.Ad) (*ads.Ad, error) {
	err := validator.Validate(ad)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	get, err := p.Get(adID)
	if err != nil {
		return nil, err
	}
	//*p.Get(adID) = ad
	*get = ad
	return &ad, nil
}