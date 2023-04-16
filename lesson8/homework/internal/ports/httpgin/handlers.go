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
		}

		ad, err1 := a.CreateAd(reqBody.Title, reqBody.Text, int(reqBody.UserID))
		if err1 != nil {
			if errors.Is(err1, myerrors.ErrBadRequest) {
				c.Status(http.StatusBadRequest)
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			}
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
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
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
		}

		ad, err1 := a.UpdateAdStatus(int64(adID), reqBody.UserID, reqBody.Published)
		if err1 != nil {
			switch err1 {
			case myerrors.ErrBadRequest:
				c.Status(http.StatusBadRequest)
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			case myerrors.ErrForbidden:
				c.Status(http.StatusForbidden)
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
			}
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
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
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
		}

		ad, err1 := a.UpdateAd(int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
		if err1 != nil {
			switch err1 {
			case myerrors.ErrBadRequest:
				c.Status(http.StatusBadRequest)
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			case myerrors.ErrForbidden:
				c.Status(http.StatusForbidden)
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
			}
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

//func getAd(a app.App) fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		adID, err := c.ParamsInt("ad_id")
//		if err != nil {
//			c.Status(http.StatusBadRequest)
//			return c.JSON(AdErrorResponse(err))
//		}
//
//		ad, err := a.GetAd(int64(adID))
//		if err != nil {
//			return err
//		}
//
//		return c.JSON(AdSuccessResponse(ad))
//	}
//}
