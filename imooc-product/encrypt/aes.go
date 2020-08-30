package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// 16，24,32位的字符串分别对于AES-128，AES-192，AES-256 加密方法
var PwdKey = []byte("*.handsomeshop.*")

// PKCS7填充模式
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	// 计算有没有长度缺失，规定加密的字符串必须是长度的整数倍
	padding := blockSize - len(ciphertext)%blockSize
	// Repeat函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字符切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

// 填充的反向操作，删除填充字符串
func PKCS7UnPadding(originData []byte) ([]byte, error)  {
	length := len(originData)
	if length == 0 {
		return nil, errors.New("加密字符串错误")
	}else {
		// 获取填充字符串长度,因为我们封装的 PKCS7Padding 末尾就是使用的需要填充的数量的显示的，所以可以这样取
		unpadding := int(originData[length-1])

		// 截取切片，删除填充字节，并且返回明文
		return originData[:(length-unpadding)], nil
	}
}

// 实现加密
func AesEcrypt(origData []byte, key []byte) ([]byte, error)  {
	// 创建加密实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 获取块的大小
	blockSize := block.BlockSize()
	// 对数据进行填充，满足长度
	origData = PKCS7Padding(origData, blockSize)
	// 采用AES加密算法中的CBC模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 执行加密
	blockMode.CryptBlocks(crypted, origData)

	return crypted, nil
}

// 实现解密
func AesDecrypt(cypted []byte, key []byte) ([]byte, error)  {
	// 创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 获取块的大小
	blockSize := block.BlockSize()
	// 创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))
	// 这个函数也可以用来解密
	blockMode.CryptBlocks(origData, cypted)
	// 去除填充字符串
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}

	return origData, nil

}

// 加密base64
func EnPwdCode(pwd []byte) (string, error)  {
	result, err := AesEcrypt(pwd, PwdKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(result), nil
}

// 解密
func DePwdCode(pwd string) ([]byte, error)  {
	// 解密base64字符串
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)

	if err != nil {
		return nil, err
	}

	// 执行AES解密
	return AesDecrypt(pwdByte, PwdKey)
}

