package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)
	_, err = client.createUser("dimosha", "dmitriy@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	type Test struct {
		Name      string
		UserID    int64
		AdID      any
		Pub       bool
		ExpectPub bool
		ExpectErr error
	}

	tests := [...]Test{
		{"Switch to published", 0, "0", true, true, nil},
		{"Switch back to unpublished", 0, "0", false, false, nil},
		{"Switch unpublished to unpublished", 0, "0", false, false, nil},
		{"User not exist", 2, "0", false, false, ErrBadRequest},
		{"User not author", 1, "0", false, false, ErrForbidden},
		{"BadID", 0, "dsfdg", false, false, ErrBadRequest},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			response, err = client.changeAdStatus(test.UserID, test.AdID, test.Pub)
			if test.ExpectErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, test.ExpectPub, response.Data.Published)
			} else {
				assert.ErrorIs(t, err, test.ExpectErr)
			}
		})
	}
}
