package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"homework10/internal/ads"
	"homework10/internal/app"
	"homework10/internal/tests/mocks"
	"homework10/internal/users"
	"testing"
)

func Test_AddUser(t *testing.T) {
	AdRepo := mocks.NewAdRepository(t)

	UserRepo := &mocks.UserRepository{}
	UserRepo.On("Add", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(&users.User{ID: 0, Nickname: "name", Email: "mail"}, nil).
		Once()

	service := app.NewApp(AdRepo, UserRepo)
	res, err := service.CreateUser("name", "mail")
	assert.Equal(t, res.ID, int64(0))
	assert.Nil(t, err)
}

func Test_GetUser(t *testing.T) {
	AdRepo := mocks.NewAdRepository(t)
	UserRepo := mocks.NewUserRepository(t)

	UserRepo.On("Get", mock.AnythingOfType("int64")).
		Return(&users.User{ID: 0, Nickname: "name", Email: "mail"}, nil).
		Once()

	service := app.NewApp(AdRepo, UserRepo)
	res, err := service.GetUser(0)
	assert.Equal(t, res.ID, int64(0))
	assert.Nil(t, err)
}

func Test_UpdateUser(t *testing.T) {
	AdRepo := mocks.NewAdRepository(t)
	UserRepo := mocks.NewUserRepository(t)
	UserRepo.On("Update", mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
		Return(&users.User{ID: 0, Nickname: "name1", Email: "mail1"}, nil).
		Once()
	UserRepo.On("Get", mock.AnythingOfType("int64")).
		Return(&users.User{ID: 0, Nickname: "name", Email: "mail"}, nil).
		Once()

	service := app.NewApp(AdRepo, UserRepo)
	res, err := service.UpdateUser(0, "name1", "mail1", 0)
	assert.Equal(t, res.ID, int64(0))
	assert.Nil(t, err)
}

func Test_DeleteUser(t *testing.T) {
	AdRepo := mocks.NewAdRepository(t)
	UserRepo := mocks.NewUserRepository(t)
	UserRepo.On("Delete", mock.AnythingOfType("int64")).
		Return(true).
		Once()

	service := app.NewApp(AdRepo, UserRepo)
	err := service.DeleteUser(0)
	assert.Nil(t, err)
}

////////////////////____AD____////////////////////////

func Test_AddAd(t *testing.T) {
	AdRepo := &mocks.AdRepository{}
	AdRepo.On("Add", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
		Return(&ads.Ad{
			ID:    0,
			Title: "title",
			Text:  "text",
		}, nil).
		Once()

	UserRepo := &mocks.UserRepository{}
	UserRepo.On("Get", mock.AnythingOfType("int64")).
		Return(&users.User{ID: 1, Nickname: "name", Email: "mail"}, nil).
		Once()

	service := app.NewApp(AdRepo, UserRepo)
	res, err := service.CreateAd("title", "text", 1)
	assert.Equal(t, res.ID, int64(0))
	assert.Nil(t, err)
}

func Test_GetAd(t *testing.T) {
	AdRepo := &mocks.AdRepository{}
	AdRepo.On("Get", mock.AnythingOfType("int64")).
		Return(&ads.Ad{
			ID:    0,
			Title: "title",
			Text:  "text",
		}, nil).
		Once()

	UserRepo := &mocks.UserRepository{}
	UserRepo.On("Get", mock.AnythingOfType("int64")).
		Return(&users.User{ID: 0, Nickname: "name", Email: "mail"}, nil).
		Once()

	service := app.NewApp(AdRepo, UserRepo)
	res, err := service.GetAd(0)
	assert.Equal(t, res.ID, int64(0))
	assert.Nil(t, err)
}

func Test_GetAds(t *testing.T) {
	AdRepo := &mocks.AdRepository{}
	AdRepo.On("Get", mock.AnythingOfType("int64")).
		Return(&ads.Ad{
			ID:    0,
			Title: "title",
			Text:  "text",
		}, nil).
		Once()
	AdRepo.On("Idxs").
		Return([]int64{0}).
		Once()

	UserRepo := &mocks.UserRepository{}
	UserRepo.On("Get", mock.AnythingOfType("int64")).
		Return(&users.User{ID: 0, Nickname: "name", Email: "mail"}, nil).
		Once()

	service := app.NewApp(AdRepo, UserRepo)
	res, err := service.GetAds("all", -1, "all", "")
	assert.Equal(t, (*res)[0].ID, int64(0))
	assert.Nil(t, err)
}

func Test_UpdateAd(t *testing.T) {
	AdRepo := &mocks.AdRepository{}
	AdRepo.On("Get", mock.AnythingOfType("int64")).
		Return(&ads.Ad{
			ID:        0,
			Title:     "title",
			Text:      "text",
			Published: true,
		}, nil).
		Once()
	AdRepo.On("Update", mock.AnythingOfType("int64"), mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("bool")).
		Return(&ads.Ad{
			ID:        0,
			Title:     "t1",
			Text:      "t2",
			Published: true,
		}, nil).
		Once()

	UserRepo := &mocks.UserRepository{}
	UserRepo.On("Get", mock.AnythingOfType("int64")).
		Return(&users.User{ID: 0, Nickname: "name", Email: "mail"}, nil).
		Once()

	service := app.NewApp(AdRepo, UserRepo)
	res, err := service.UpdateAd(0, 0, "t1", "t2")
	assert.Equal(t, res.ID, int64(0))
	assert.Nil(t, err)
}

func Test_DeleteAd(t *testing.T) {
	AdRepo := &mocks.AdRepository{}
	AdRepo.On("Delete", mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
		Return(true).
		Once()

	UserRepo := &mocks.UserRepository{}
	UserRepo.On("Get", mock.AnythingOfType("int64")).
		Return(&users.User{ID: 0, Nickname: "name", Email: "mail"}, nil).
		Once()

	service := app.NewApp(AdRepo, UserRepo)
	err := service.DeleteAd(0, 0)
	assert.Nil(t, err)
}
