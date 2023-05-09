package tests

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	client := getTestClient()
	response, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.UserID)
	assert.Equal(t, response.Data.Nickname, "dimosha")
	assert.Equal(t, response.Data.Email, "dmitriy@mail.ru")
}

func TestGetUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	user, err := client.getUser("0")
	assert.NoError(t, err)
	assert.Equal(t, user.Data.UserID, int64(0))
	assert.Equal(t, user.Data.Nickname, "dimosha")
	assert.Equal(t, user.Data.Email, "dmitriy@mail.ru")
}

func TestGetUserBadId(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	_, err = client.getUser("fg")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	response, err := client.updateUser("0", "D1", "D2", 0)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Nickname, "D1")
	assert.Equal(t, response.Data.Email, "D2")
}

func TestUpdateUserBadId(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	_, err = client.updateUser("dsf", "D1", "D2", 23)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateUserForbidden(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	_, err = client.updateUser("dsf", "D1", "D2", 0)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateUserBadreq(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	_, err = client.updateUser("0", "", "D2", 0)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd(t *testing.T) {
	client := getTestClient()
	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)
	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(0))
	assert.False(t, response.Data.Published)
}

func TestChangeAdStatusBadrequest(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(2, strconv.Itoa(int(response.Data.ID)), true)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestChangeAdStatusBadID(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	_, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, "ljvdb", true)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(0, strconv.Itoa(int(response.Data.ID)), "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
}

func TestUpdateAdStatusBadID(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	_, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(0, "ljvdb", "t1", "t2")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestListAds(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(0, strconv.Itoa(int(response.Data.ID)), true)
	assert.NoError(t, err)

	_, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAds()
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}

func TestListFilterAds(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, strconv.Itoa(int(response.Data.ID)), true)
	assert.NoError(t, err)

	response, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)

	d := time.Now().UTC()
	s := d.Format("02-01-06")

	ads, err := client.listFilterAds("false", s, "-1")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, response.Data.ID)
	assert.Equal(t, ads.Data[0].Title, response.Data.Title)
	assert.Equal(t, ads.Data[0].Text, response.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, response.Data.AuthorID)
	assert.False(t, ads.Data[0].Published)
}

func TestListFilterAdsBadArgs(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, strconv.Itoa(int(response.Data.ID)), true)
	assert.NoError(t, err)

	response, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)

	_, err = client.listFilterAds("sdbvifh", "sdvffdv", "-1")
	assert.ErrorIs(t, err, ErrBadRequest)

	_, err = client.listFilterAds("false", "sdvffdv", "-1")
	assert.ErrorIs(t, err, ErrBadRequest)

	_, err = client.listFilterAds("false", "all", "zfdbsgnmf")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestGetAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	ads, err := client.getAd("0")
	assert.NoError(t, err)
	assert.Equal(t, ads.Data.ID, response.Data.ID)
	assert.Equal(t, ads.Data.Title, response.Data.Title)
	assert.Equal(t, ads.Data.Text, response.Data.Text)
	assert.Equal(t, ads.Data.AuthorID, response.Data.AuthorID)
	assert.False(t, ads.Data.Published)
}

func TestGetAdStatusBadID(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	_, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.getAd("dsjvb")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestDeleteAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	_, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.deleteAdint(0, 0)
	assert.NoError(t, err)
}

func TestDeleteAdForbidden(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	_, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.deleteAdint(0, 1)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestDeleteUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	_, err = client.deleteUser("0")
	assert.NoError(t, err)
}

func TestDeleteUserBadID(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	_, err = client.deleteUser("sdcv")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestDeleteAdBadID(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)
	_, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.deleteAd("ihbdfs", 1)
	assert.ErrorIs(t, err, ErrBadRequest)
}
