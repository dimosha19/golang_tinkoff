package httpgin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"homework8/internal/app"
	myerrors "homework8/internal/errors"
	"net/http"
	"strconv"
)

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.ShouldBindJSON(&reqBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.CreateAd(reqBody.Title, reqBody.Text, int(reqBody.UserID))
		if err != nil {
			if errors.Is(err, myerrors.ErrBadRequest) {
				c.Status(http.StatusBadRequest)
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
				return
			}
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
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
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.UpdateAdStatus(int64(adID), reqBody.UserID, reqBody.Published)
		if err != nil {
			switch err {
			case myerrors.ErrBadRequest:
				c.Status(http.StatusBadRequest)
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			case myerrors.ErrForbidden:
				c.Status(http.StatusForbidden)
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
			}
			return
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
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
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.UpdateAd(int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
		if err != nil {
			switch err {
			case myerrors.ErrBadRequest:
				c.Status(http.StatusBadRequest)
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			case myerrors.ErrForbidden:
				c.Status(http.StatusForbidden)
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
			}
			return
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func getAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Server", "matrix")
		ad, err := a.GetAds()
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdsSuccessResponse(*ad))
	}
}
