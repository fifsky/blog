package aesutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testSecret = "abcdabcdabcdabcd"

func TestAesEncode(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		str, err := AesEncode(testSecret, "123")
		assert.NoError(t, err)
		assert.Equal(t, `R0I+Owj5eBnU5dNFKbtCXw==`, str)
	})

	t.Run("error", func(t *testing.T) {
		errSecret := "abcdabcdabcdabc"
		_, err := AesEncode(errSecret, "123")
		assert.Error(t, err)
	})

}

func TestAesDecode(t *testing.T) {
	t.Run("decode", func(t *testing.T) {
		str := "R0I+Owj5eBnU5dNFKbtCXw=="

		ori, err := AesDecode(testSecret, str)
		assert.NoError(t, err)
		assert.Equal(t, `123`, ori)
	})

	t.Run("decode", func(t *testing.T) {
		str := "R0I+Owj5eBnU5dNFKbtCXw="
		_, err := AesDecode(testSecret, str)
		assert.Error(t, err)
	})

	t.Run("decode", func(t *testing.T) {
		str := "R0I+Owj5eBnU5dNFKbtCXw=="
		errSecret := "abcdabcdabcdabc"
		_, err := AesDecode(errSecret, str)
		assert.Error(t, err)
	})
}
