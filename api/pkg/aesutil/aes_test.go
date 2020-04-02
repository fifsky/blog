package aesutil

import (
	"fmt"
	"net/url"
	"testing"
)

func TestAesEncode(t *testing.T) {
	secret := "Cif$kyL!1024@iLU"

	str, err := AesEncode(secret, "123")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(url.QueryEscape(str))

	ori, err := AesDecode(secret, str)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ori)
}

func TestAesDecode(t *testing.T) {
	secret := "Cif$kyL!1024@iLU"

	str := "ug2+QKXIrvkMN/g7zKIwFg=="

	ori, err := AesDecode(secret, str)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ori)
}
