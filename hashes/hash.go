package hashes

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

// MD5 计算md5值
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA1 计算sha1值
func SHA1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256 计算sha256值
func SHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Hash 计算指定算法的hash值
func Hash(hash crypto.Hash, str string) string {
	if !hash.Available() {
		return ""
	}
	h := hash.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
