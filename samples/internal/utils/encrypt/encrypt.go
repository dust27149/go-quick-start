package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

// EncryptPKCS1v15 使用平台 RSA 公钥以 PKCS#1 v1.5 填充加密文本，返回 Base64 密文。
func EncryptPKCS1v15(rsaPublicKey, txtToEncrypt string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(rsaPublicKey)
	if err != nil {
		return "", fmt.Errorf("解码 RSA 公钥失败: %w", err)
	}

	// Java 的 X509EncodedKeySpec 对应 Go 中的 PKIX 公钥格式。
	publicKeyAny, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return "", fmt.Errorf("解析 RSA 公钥失败: %w", err)
	}

	publicKey, ok := publicKeyAny.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("公钥不是 RSA 类型")
	}

	// 对应 Java: Cipher.getInstance("RSA/ECB/PKCS1Padding")。
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(txtToEncrypt))
	if err != nil {
		return "", fmt.Errorf("RSA 加密失败: %w", err)
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}
