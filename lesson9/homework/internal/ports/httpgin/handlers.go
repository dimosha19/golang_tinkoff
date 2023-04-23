package httpgin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"homework9/internal/app"
	myerrors "homework9/internal/errors"
	"net/http"
	"strconv"
	"time"
)

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.ShouldBindJSON(&reqBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		ad, err := a.CreateAd(reqBody.Title, reqBody.Text, reqBody.UserID)
		if err != nil {
			if errors.Is(err, myerrors.ErrBadRequest) {
				c.Status(http.StatusBadRequest)
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
				return
			}
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		ad, err := a.UpdateAdStatus(int64(adID), reqBody.UserID, reqBody.Published)
		if err != nil {
			switch err {
			case myerrors.ErrBadRequest:
				c.Status(http.StatusBadRequest)
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
			case myerrors.ErrForbidden:
				c.Status(http.StatusForbidden)
				c.JSON(http.StatusForbidden, ErrorResponse(err))
			}
			return
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		ad, err := a.UpdateAd(int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)

		if err != nil {
			switch err {
			case myerrors.ErrBadRequest:
				c.Status(http.StatusBadRequest)
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
			case myerrors.ErrForbidden:
				c.Status(http.StatusForbidden)
				c.JSON(http.StatusForbidden, ErrorResponse(err))
			}
			return
		}
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func getAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		pub := c.DefaultQuery("pub", "true")
		if pub != "true" && pub != "false" && pub != "all" {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(myerrors.ErrBadRequest))
			return
		}

		title := c.DefaultQuery("title", "")

		author, err := strconv.Atoi(c.DefaultQuery("author", "-1"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(myerrors.ErrBadRequest))
			return
		}

		date := c.DefaultQuery("date", "all")
		_, err = time.Parse("02-01-06", date)
		if date != "all" && err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(myerrors.ErrBadRequest))
			return
		}

		ad, err := a.GetAds(pub, int64(author), date, title)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdsSuccessResponse(*ad))
	}
}

func getAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		ad, err := a.GetAd(int64(adID))

		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUserRequest
		err := c.ShouldBindJSON(&reqBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		user, err := a.CreateUser(reqBody.Nickname, reqBody.Email)
		if err != nil {
			if errors.Is(err, myerrors.ErrBadRequest) {
				c.Status(http.StatusBadRequest)
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
				return
			}
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

func updateUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateUserRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		user, err := a.UpdateUser(int64(userID), reqBody.Nickname, reqBody.Email, reqBody.UserID)

		if err != nil {
			switch err {
			case myerrors.ErrBadRequest:
				c.Status(http.StatusBadRequest)
				c.JSON(http.StatusBadRequest, ErrorResponse(err))
			case myerrors.ErrForbidden:
				c.Status(http.StatusForbidden)
				c.JSON(http.StatusForbidden, ErrorResponse(err))
			}
			return
		}
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

func getUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		user, err := a.GetUser(int64(userID))

		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

func deleteAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.ShouldBindJSON(&reqBody)
		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		err = a.DeleteAd(int64(adID), reqBody.UserID)

		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, DeleteSuccessResponse(int64(adID)))
	}
}

func deleteUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.ShouldBindJSON(&reqBody)
		adID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		err = a.DeleteUser(int64(adID))

		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, DeleteSuccessResponse(int64(adID)))
	}
}
