package crypts

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesCBC(t *testing.T) {
	key := "AES256Key-32Characters1234567890"
	iv := key[:16]
	data := "ILoveNobleGase"

	cipher, err := EncryptCBC([]byte(key), []byte(iv), []byte(data))
	assert.Nil(t, err)
	assert.Equal(t, "WDq8s1qdHCML8YLhfdmGRw==", cipher.String())

	plain, err := DecryptCBC([]byte(key), []byte(iv), cipher.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, data, string(plain))

	cipher2, err := EncryptCBC([]byte(key), []byte(iv), []byte(data), 32)
	assert.Nil(t, err)
	assert.Equal(t, "vjemH/hxbwNh+WXhkKseCu2GrM4O6bnaaKv59wgkRSE=", cipher2.String())

	plain2, err := DecryptCBC([]byte(key), []byte(iv), cipher2.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, data, string(plain2))
}

func TestAesECB(t *testing.T) {
	key := "AES256Key-32Characters1234567890"
	data := "ILoveNobleGase"

	cipher, err := EncryptECB([]byte(key), []byte(data))
	assert.Nil(t, err)
	assert.Equal(t, "+oYDjdGHY8lK1/sJo750Waw==", cipher.String())

	plain, err := DecryptECB([]byte(key), cipher.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, data, string(plain))

	cipher2, err := EncryptECB([]byte(key), []byte(data), 32)
	assert.Nil(t, err)
	assert.Equal(t, "u0iDWHM8JMnRyJNCiCzKJNib2cOjUrx2FqMjmg3ZTZA=", cipher2.String())

	plain2, err := DecryptECB([]byte(key), cipher.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, data, string(plain2))
}

func TestAesCFB(t *testing.T) {
	key := "AES256Key-32Characters1234567890"
	iv := key[:16]
	data := "ILoveNobleGase"

	cipher, err := EncryptCFB([]byte(key), []byte(iv), []byte(data))
	assert.Nil(t, err)
	assert.Equal(t, "KP7OnZj9J9ONnjn6yA0=", cipher.String())

	plain, err := DecryptCFB([]byte(key), []byte(iv), cipher.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, data, string(plain))
}

func TestAesOFB(t *testing.T) {
	key := "AES256Key-32Characters1234567890"
	iv := key[:16]
	data := "ILoveNobleGase"

	cipher, err := EncryptOFB([]byte(key), []byte(iv), []byte(data))
	assert.Nil(t, err)
	assert.Equal(t, "KP7OnZj9J9ONnjn6yA0=", cipher.String())

	plain, err := DecryptOFB([]byte(key), []byte(iv), cipher.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, data, string(plain))
}

func TestAesCTR(t *testing.T) {
	key := "AES256Key-32Characters1234567890"
	iv := key[:16]
	data := "ILoveNobleGase"

	cipher, err := EncryptCTR([]byte(key), []byte(iv), []byte(data))
	assert.Nil(t, err)
	assert.Equal(t, "KP7OnZj9J9ONnjn6yA0=", cipher.String())

	plain, err := DecryptCTR([]byte(key), []byte(iv), cipher.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, data, string(plain))
}

func TestAesGCM(t *testing.T) {
	key := "AES256Key-32Characters1234567890"
	nonce := key[:12]
	data := "ILoveNobleGase"
	aad := "IIInsomnia"

	cipher, err := EncryptGCM([]byte(key), []byte(nonce), []byte(data), []byte(aad), &GCMOption{})
	assert.Nil(t, err)
	assert.Equal(t, "qciumnROL4U9F0klEKhzE/DngAy/clYUsZGfcafh", cipher.String())
	assert.Equal(t, "qciumnROL4U9F0klEKg=", base64.StdEncoding.EncodeToString(cipher.Data()))
	assert.Equal(t, "cxPw54AMv3JWFLGRn3Gn4Q==", base64.StdEncoding.EncodeToString(cipher.Tag()))

	plain, err := DecryptGCM([]byte(key), []byte(nonce), cipher.Bytes(), []byte(aad), nil)
	assert.Nil(t, err)
	assert.Equal(t, data, string(plain))
}
