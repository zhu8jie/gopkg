package xsecret

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {

	// sss, err := Encrypt(SCTYPE_AES_ECB, "1000", "123456789abcdefg")
	// if err != nil {
	// 	t.Error(err)
	// }
	// if sss != "Fp599QWzjlqGb0SoX+RFBA==" {
	// 	t.Errorf("%v not equal %v", sss, "Fp599QWzjlqGb0SoX+RFBA==")
	// }

	str := "1000"
	key := "cbecd9671c78bce06a38a8130be2fd7b"
	// iv := "FD2718DD5C312460"
	// fmt.Println(len(key))
	// tmpKey := []byte(key)
	// fmt.Println(len(tmpKey), tmpKey)

	// encryptStr, err := Encrypt(SCTYPE_AES_ECB, str, key)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// 加密方式为:
	// 1.对价格明文采取AES/CBC加密，token线下提供。
	// 2.AES加密之后，对加密字符串进行BASE64编码。
	// 3.对BASE64编码之后的字符串进行URLENCODE。
	// 测试实例：
	// token：exwo0aln9sqgdtdn
	// 加密后密文：CAD3C0EDB517884C04199C87BC2E25A5
	// 解密后价格（单位：分）：500

	// 正整数的底价，采⽤aes cbc加密⽅式 加密并base64 url safe 编码的形式传递，
	// ⽤key的第⼀组块的尺⼨（前16位）做的iv向量，填充⽅式 PKCS5 ，
	// 密钥请联系运营获取 例: key: "cbecd9671c78bce06a38a8130be2fd7b" 加密前: "1000" 加密后:"txshMKoDeDfGxCynkrOwTw=="

	dspWinPriceAesCBC, err := Encrypt(SCTYPE_AES_CBC, str, key)

	// encryptStr, err := (xaes.AESEncryptZeroPadEcb([]byte(str), []byte(key)))
	if err != nil {
		fmt.Println("error: ", err)
	}

	// s := xhex.HexEncode(string(encryptStr))

	fmt.Println(dspWinPriceAesCBC)
}
