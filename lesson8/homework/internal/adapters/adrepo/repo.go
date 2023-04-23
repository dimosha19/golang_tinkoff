package adrepo

import (
	validator "github.com/dimosha19/myvalidator"
	"homework8/internal/ads"
	"homework8/internal/app"
	myerrors "homework8/internal/errors"
	"sync"
)

type scliceAd struct {
	mx *sync.Mutex
	r  []ads.Ad
}

func New() app.AdRepository {
	mx := sync.Mutex{}
	res := scliceAd{mx: &mx}
	return &res
}

func (p *scliceAd) Add(ad ads.Ad) *ads.Ad {
	p.mx.Lock()
	(*p).r = append((*p).r, ad)
	p.mx.Unlock()
	return &((*p).r)[len((*p).r)-1] // это в рамочку
}

func (p *scliceAd) Get(adID int64) (*ads.Ad, error) {
	defer p.mx.Unlock()
	if adID < p.Size() {
		p.mx.Lock()
		res := (*p).r[adID]
		return &res, nil
	}
	return nil, myerrors.ErrBadRequest
}

func (p *scliceAd) Size() int64 {
	res := int64(len((*p).r))
	return res
}

func (p *scliceAd) Update(adID int64, ad ads.Ad) (*ads.Ad, error) {
	defer p.mx.Unlock()
	p.mx.Lock()
	err := validator.Validate(ad)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	if adID >= p.Size() {
		return nil, myerrors.ErrBadRequest
	}
	(*p).r[adID] = ad
	//*p.Get(adID) = ad
	return &ad, nil
}
