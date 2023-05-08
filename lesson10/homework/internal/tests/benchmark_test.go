package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Benchmark100users(b *testing.B) {
	client := getTestClient()
	for i := 0; i < 100; i++ {
		_, err := client.createUser("dimosha", "dmitriy@mail.ru")
		assert.NoError(b, err)
	}
}
