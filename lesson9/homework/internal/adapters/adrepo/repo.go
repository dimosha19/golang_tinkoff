package adrepo

import (
	validator "github.com/dimosha19/myvalidator"
	"homework9/internal/ads"
	"homework9/internal/app"
	myerrors "homework9/internal/errors"
	"sync"
	"time"
)

type mapAd struct {
	mx  *sync.Mutex
	r   map[int64]ads.Ad
	idx int64
}

func New() app.AdRepository {
	mx := sync.Mutex{}
	res := mapAd{r: make(map[int64]ads.Ad), mx: &mx, idx: 0}
	return &res
}

func (p *mapAd) Add(title string, text string, userID int64) (*ads.Ad, error) {
	res := ads.Ad{ID: p.idx, Title: title, Text: text, AuthorID: int64(userID), PublishedTime: time.Now().UTC(), UpdateTime: time.Now().UTC()}
	p.idx++
	err := validator.Validate(res)
	if err != nil {
		return nil, err
	}
	p.mx.Lock()
	p.r[res.ID] = res
	p.mx.Unlock()
	a := p.r[res.ID]
	return &a, nil
}

func (p *mapAd) Delete(adID int64, userID int64) bool {
	if p.r[adID].AuthorID != userID {
		return false
	}
	delete(p.r, adID)
	return true
}

func (p *mapAd) Idxs() []int64 {
	keys := make([]int64, len(p.r))

	i := 0
	for k := range p.r {
		keys[i] = k
		i++
	}
	return keys
}

func (p *mapAd) Get(adID int64) (*ads.Ad, error) {
	defer p.mx.Unlock()
	p.mx.Lock()
	res, ok := p.r[adID]
	if ok {
		return &res, nil
	}
	return nil, myerrors.ErrBadRequest
}

func (p *mapAd) Update(adID int64, userID int64, title string, text string, published bool) (*ads.Ad, error) {
	t := (p.r)[adID]
	temp := t
	if temp.AuthorID != userID {
		return nil, myerrors.ErrForbidden
	}
	temp.Published = published
	temp.Title = title
	temp.Text = text
	temp.UpdateTime = time.Now().UTC()
	defer p.mx.Unlock()
	p.mx.Lock()
	err := validator.Validate(temp)
	if err != nil {
		return nil, myerrors.ErrBadRequest
	}
	p.r[adID] = temp
	return &temp, nil
}
