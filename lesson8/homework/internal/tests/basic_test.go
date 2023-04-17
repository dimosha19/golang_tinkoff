package tests

import (
	"testing"

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

	user, err := client.getUser(int64(0))
	assert.NoError(t, err)
	assert.Equal(t, user.Data.UserID, int64(0))
	assert.Equal(t, user.Data.Nickname, "dimosha")
	assert.Equal(t, user.Data.Email, "dmitriy@mail.ru")
}

func TestUpdateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	response, err = client.updateUser(0, "D1", "D2")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Nickname, "D1")
	assert.Equal(t, response.Data.Email, "D2")
}

func TestCreateAd(t *testing.T) {
	client := getTestClient()
	for i := 0; i < 124; i++ {
		_, err := client.createUser("dimosha", "dmitriy@mail.ru")
		assert.NoError(t, err)
	}
	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(123))
	assert.False(t, response.Data.Published)
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)

	response, err = client.changeAdStatus(0, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(0, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(0, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
}

func TestListAds(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(0, response.Data.ID, true)
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
