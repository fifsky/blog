package aesutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

const (
	NOPADDING = iota
	PKCS5PADDING
)

// PKCS5补位
func PKCSPadding(text []byte, blockSize int) []byte {
	padding := blockSize - len(text)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(text, padtext...)
}

// 去除PKCS5补位
func PKCSUnPadding(text []byte) ([]byte, error) {
	length := len(text)

	if length == 0 {
		return text, nil
	}

	padtext := int(text[length-1])
	if length < padtext {
		return nil, errors.New("unpadding length error")
	}
	return text[:(length - padtext)], nil
}

// 补位方法
func Padding(text []byte, padding int) ([]byte, error) {
	switch padding {
	case NOPADDING:
		if len(text)%8 != 0 {
			return nil, errors.New("text length invalid")
		}
	case PKCS5PADDING:
		return PKCSPadding(text, 8), nil
	default:
		return nil, errors.New("padding type error")
	}

	return text, nil
}

// 去除补位方法
func UnPadding(text []byte, padding int) ([]byte, error) {
	switch padding {
	case NOPADDING:
		if len(text)%8 != 0 {
			return nil, errors.New("text length invalid")
		}
	case PKCS5PADDING:
		return PKCSUnPadding(text)
	default:
		return nil, errors.New("padding type error")
	}
	return text, nil
}

// AesEncrypt AES 加密
func AesEncrypt(key, origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCSPadding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AesDecrypt AES 解密
func AesDecrypt(key, crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	return PKCSUnPadding(origData)
}

func AesEncode(token, data string) (string, error) {
	s, err := AesEncrypt([]byte(token), []byte(data))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(s), nil
}

func AesDecode(token, data string) (string, error) {
	bt, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		return "", err
	}

	s, err := AesDecrypt([]byte(token), bt)

	if err != nil {
		return "", err
	}
	return string(s), nil
}
