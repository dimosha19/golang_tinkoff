package httpgin

import (
	"errors"
	"homework8/internal/app"
	myerrors "homework8/internal/errors"
	"net/http"
)

// Метод для создания объявления (ad)
func createAd(a app.App) {
	var reqBody createAdRequest
	err := c.BodyParser(&reqBody)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.JSON(AdErrorResponse(err))
	}

	ad, err1 := a.CreateAd(reqBody.Title, reqBody.Text, int(reqBody.UserID))
	if err1 != nil {
		if errors.Is(err1, myerrors.ErrBadRequest) {
			c.Status(http.StatusBadRequest)
		}
		c.JSON(AdErrorResponse(err1))
	}

	// TODO: вызов логики, например, CreateAd(c.Context(), reqBody.Title, reqBody.Text, reqBody.UserID)
	// TODO: метод должен возвращать AdSuccessResponse или ошибку.

	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.JSON(AdErrorResponse(err))
	}
	c.JSON(AdSuccessResponse(ad))
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody changeAdStatusRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, err1 := a.UpdateAdStatus(int64(adID), reqBody.UserID, reqBody.Published)
		if err1 != nil {
			switch err1 {
			case myerrors.ErrBadRequest:
				c.Status(http.StatusBadRequest)
			case myerrors.ErrForbidden:
				c.Status(http.StatusForbidden)
			}
			return c.JSON(AdErrorResponse(err1))
		}

		// TODO: вызов логики ChangeAdStatus(c.Context(), int64(adID), reqBody.UserID, reqBody.Published)
		// TODO: метод должен возвращать AdSuccessResponse или ошибку.

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(AdErrorResponse(err))
		}

		return c.JSON(AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody updateAdRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, err1 := a.UpdateAd(int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
		if err1 != nil {
			switch err1 {
			case myerrors.ErrBadRequest:
				c.Status(http.StatusBadRequest)
			case myerrors.ErrForbidden:
				c.Status(http.StatusForbidden)
			}
			return c.JSON(AdErrorResponse(err1))
		}

		// TODO: вызов логики, например, UpdateAd(c.Context(), int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
		// TODO: метод должен возвращать AdSuccessResponse или ошибку.

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(AdErrorResponse(err))
		}

		return c.JSON(AdSuccessResponse(ad))
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
