package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework10/internal/ads"
	"homework10/internal/users"
	"time"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type createUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
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

type adDeleteResponse struct {
	ID      int64  `json:"id"`
	Success string `json:"success"`
}

type adDeleteRequest struct {
	AdID   int64 `json:"ad_id"`
	UserID int64 `json:"user_id"`
}

type userResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
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

type updateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	UserID   int64  `json:"user_id"`
}

func DeleteSuccessResponse(id int64) *gin.H {
	return &gin.H{
		"data": adDeleteResponse{
			ID:      id,
			Success: "Was deleted",
		},
	}
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

func UserSuccessResponse(user *users.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			ID:       user.ID,
			Nickname: user.Nickname,
			Email:    user.Email,
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

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
