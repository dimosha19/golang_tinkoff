package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/ads"
	"time"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type adResponse struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Text          string    `json:"text"`
	AuthorID      int64     `json:"author_id"`
	Published     bool      `json:"published"`
	PublishedTime time.Time `json:"published_time"`
	UpdatedTime   time.Time `json:"updated_time"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

func AdSuccessResponse(ad *ads.Ad) *gin.H {
	return &gin.H{
		"data": adResponse{
			ID:            ad.ID,
			Title:         ad.Title,
			Text:          ad.Text,
			AuthorID:      ad.AuthorID,
			Published:     ad.Published,
			UpdatedTime:   ad.UpdateTime,
			PublishedTime: ad.PublishedTime,
		},
		"error": nil,
	}
}

func AdsSuccessResponse(a []ads.Ad) *gin.H {
	var response []adResponse
	for i := range a {
		response = append(response, adResponse{
			ID:            a[i].ID,
			Title:         a[i].Title,
			Text:          a[i].Text,
			AuthorID:      a[i].AuthorID,
			Published:     a[i].Published,
			PublishedTime: a[i].PublishedTime,
			UpdatedTime:   a[i].UpdateTime,
		})
	}

	return &gin.H{
		"data":  response,
		"error": nil,
	}
}

func AdErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
