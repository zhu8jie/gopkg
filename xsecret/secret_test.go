package xsecret

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"testing"
)

func AesEncrypt(str string, secretKey string) (string, error) {
	// 将输入字符串转换为字节数组
	toEncryptArray := []byte(str)
	key := []byte(secretKey)

	// 创建并初始化RijndaelManaged等价物
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 使用PKCS7填充
	padding := aes.BlockSize - len(toEncryptArray)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	toEncryptArray = append(toEncryptArray, padtext...)

	// ECB模式不需要IV，创建一个新的cipher.BlockMode
	ciphertext := make([]byte, len(toEncryptArray))
	mode := cipher.NewCBCEncrypter(block, bytes.Repeat([]byte{0}, aes.BlockSize)) // 由于ECB模式不使用IV，此处仅为示例而设

	mode.CryptBlocks(ciphertext, toEncryptArray)

	// 对结果进行Base64编码，并替换特殊字符
	tmpstr := base64.StdEncoding.EncodeToString(ciphertext)
	tmpstr = removeBase64Padding(tmpstr)
	tmpstr = replaceBase64SpecialChars(tmpstr)

	return tmpstr, nil
}

func removeBase64Padding(b64Str string) string {
	return b64Str[:len(b64Str)-len(b64Str)%4]
}

func replaceBase64SpecialChars(b64Str string) string {
	b64Str = string(bytes.ReplaceAll([]byte(b64Str), []byte("+"), []byte("-"))[:])
	b64Str = string(bytes.ReplaceAll([]byte(b64Str), []byte("/"), []byte("_"))[:])
	b64Str = string(bytes.ReplaceAll([]byte(b64Str), []byte("="), []byte("")[:]))
	return string(b64Str)
}

// func main() {
//     encrypted, err := AesEncrypt("135", "zTqHtwVgno07JCdOyGRxh6o5rkIKYhAI")
//     if err != nil {
//        fmt.Println("Error:", err)
//     } else {
//        fmt.Println("Encrypted:", encrypted)
//     }
// }

func TestEncrypt(t *testing.T) {
	// var dspWinPrice float64 = 135
	// dspWinPriceAesECB, _ := Encrypt(SCTYPE_AES_ECB, "135", "zTqHtwVgno07JCdOyGRxh6o5rkIKYhAI")

	// b64str := xaes.AesEncryptECB([]byte("135"), []byte("zTqHtwVgno07JCdOyGRxh6o5rkIKYhAI"))
	// bt := xbase64.Base64Encrypt(b64str)
	// str, _ := AesEncrypt("135", "zTqHtwVgno07JCdOyGRxh6o5rkIKYhAI")
	// // str := strings.ReplaceAll(dspWinPriceAesECB, "+", "-")
	// // str = strings.ReplaceAll(str, "/", "_")
	// // str = strings.ReplaceAll(str, "=", "")

	// fmt.Println(str, " === ", "A7rGUjZ7Hxr-qLPIRNh2qA")
	//A7rGUjZ7Hxr-qLPIRNh2qA

}
