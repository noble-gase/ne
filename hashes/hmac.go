package hashes

import (
	"crypto"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

// HMacSHA1 计算hmac-sha1值
func HMacSHA1(key, str string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// HMacSHA256 计算hmac-sha256值
func HMacSHA256(key, str string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// HMac 计算指定hash算法的hmac值
func HMac(hash crypto.Hash, key, str string) string {
	if !hash.Available() {
		return ""
	}
	h := hmac.New(hash.New, []byte(key))
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
