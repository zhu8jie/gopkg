package googleopenrtb

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

// DoubleClickCrypto 用于 DoubleClick Ad Exchange RTB 协议的加密和解密支持
type DoubleClickCrypto struct {
	keys *Keys
}

// 常量定义
const (
	KeyAlgorithm          = "HmacSHA1"
	InitvBase             = 0
	InitvSize             = 16
	InitvTimestampOffset  = 0
	InitvServerIdOffset   = 8
	PayloadBase           = InitvBase + InitvSize
	SignatureSize         = 4
	OverheadSize          = InitvSize + SignatureSize
	counterPageSize       = 20
	counterSections       = 3*256 + 1
	microsPerCurrencyUnit = 1
)

// Keys 持有用于配置 DoubleClick 加密的密钥
type Keys struct {
	encryptionKey []byte
	integrityKey  []byte
}

// NewKeys 创建新的密钥对
func NewKeys(encryptionKey, integrityKey []byte) (*Keys, error) {
	keys := &Keys{
		encryptionKey: encryptionKey,
		integrityKey:  integrityKey,
	}

	// 验证密钥有效性
	if err := keys.validate(); err != nil {
		return nil, err
	}

	return keys, nil
}

// validate 验证密钥有效性
func (k *Keys) validate() error {
	// 测试加密密钥
	if _, err := createHmac(k.encryptionKey); err != nil {
		return fmt.Errorf("invalid encryption key: %v", err)
	}

	// 测试完整性密钥
	if _, err := createHmac(k.integrityKey); err != nil {
		return fmt.Errorf("invalid integrity key: %v", err)
	}

	return nil
}

// GetEncryptionKey 返回加密密钥
func (k *Keys) GetEncryptionKey() []byte {
	return k.encryptionKey
}

// GetIntegrityKey 返回完整性密钥
func (k *Keys) GetIntegrityKey() []byte {
	return k.integrityKey
}

// NewDoubleClickCrypto 使用加密密钥初始化
func NewDoubleClickCrypto(keys *Keys) *DoubleClickCrypto {
	return &DoubleClickCrypto{keys: keys}
}

// decode 从字符串解码为二进制形式（URL安全的base64解码）
func (d *DoubleClickCrypto) decode(data string) ([]byte, error) {
	if data == "" {
		return nil, nil
	}
	return base64.URLEncoding.DecodeString(data)
}

// encode 从二进制形式编码为字符串（URL安全的base64编码）
func (d *DoubleClickCrypto) encode(data []byte) string {
	if data == nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(data)
}

// decrypt 解密数据
func (d *DoubleClickCrypto) decrypt(cipherData []byte) ([]byte, error) {
	if len(cipherData) < OverheadSize {
		return nil, fmt.Errorf("invalid cipherData, %d bytes", len(cipherData))
	}

	// workBytes := initVector || E(payload) || I(signature)
	workBytes := make([]byte, len(cipherData))
	copy(workBytes, cipherData)

	// workBytes := initVector || payload || I(signature)
	d.xorPayloadToHmacPad(workBytes)

	// workBytes := initVector || payload || I'(signature)
	confirmationSignature, err := d.hmacSignature(workBytes)
	if err != nil {
		return nil, err
	}

	// 提取完整性签名
	integritySignature := binary.BigEndian.Uint32(workBytes[len(workBytes)-SignatureSize:])

	// 将确认签名放回数据中
	binary.BigEndian.PutUint32(workBytes[len(workBytes)-SignatureSize:], confirmationSignature)

	if confirmationSignature != integritySignature {
		return nil, fmt.Errorf("signature mismatch: %x vs %x", confirmationSignature, integritySignature)
	}

	return workBytes, nil
}

// encrypt 加密数据
func (d *DoubleClickCrypto) encrypt(plainData []byte) ([]byte, error) {
	if len(plainData) < OverheadSize {
		return nil, fmt.Errorf("invalid plainData, %d bytes", len(plainData))
	}

	// workBytes := initVector || payload || zeros:4
	workBytes := make([]byte, len(plainData))
	copy(workBytes, plainData)

	// workBytes := initVector || payload || I(signature)
	signature, err := d.hmacSignature(workBytes)
	if err != nil {
		return nil, err
	}
	binary.BigEndian.PutUint32(workBytes[len(workBytes)-SignatureSize:], signature)

	// workBytes := initVector || E(payload) || I(signature)
	d.xorPayloadToHmacPad(workBytes)

	return workBytes, nil
}

// createInitVector 从组件 (timestamp, serverId) 字段创建初始化向量
// func (d *DoubleClickCrypto) createInitVector(timestamp time.Time, serverId int64) []byte {
// 	initVector := make([]byte, InitvSize)
// 	binary.BigEndian.PutUint64(initVector[InitvTimestampOffset:], millisToSecsAndMicros(timestamp.UnixNano()/int64(time.Millisecond)))
// 	binary.BigEndian.PutUint64(initVector[InitvServerIdOffset:], uint64(serverId))
// 	return initVector
// }

// // createInitVectorFromLong 从长整型时间戳创建初始化向量
// func (d *DoubleClickCrypto) createInitVectorFromLong(timestamp, serverId uint64) []byte {
// 	initVector := make([]byte, InitvSize)
// 	binary.BigEndian.PutUint64(initVector[InitvTimestampOffset:], timestamp)
// 	binary.BigEndian.PutUint64(initVector[InitvServerIdOffset:], serverId)
// 	return initVector
// }

// getTimestamp 从加密或解密数据中获取时间戳字段
// func (d *DoubleClickCrypto) getTimestamp(data []byte) time.Time {
// 	timestamp := binary.BigEndian.Uint64(data[InitvBase+InitvTimestampOffset:])
// 	millis := secsAndMicrosToMillis(timestamp)
// 	return time.Unix(0, millis*int64(time.Millisecond))
// }

// getServerId 从加密或解密数据中获取 serverId 字段
// func (d *DoubleClickCrypto) getServerId(data []byte) int64 {
// 	return int64(binary.BigEndian.Uint64(data[InitvBase+InitvServerIdOffset:]))
// }

// initPlainData 为加密准备明文数据
func (d *DoubleClickCrypto) initPlainData(payloadSize int, initVector []byte) []byte {
	plainData := make([]byte, OverheadSize+payloadSize)

	if initVector == nil {
		// 使用当前时间和随机 serverId
		timestamp := millisToSecsAndMicros(time.Now().UnixNano() / int64(time.Millisecond))
		serverId := randUint64()

		binary.BigEndian.PutUint64(plainData[InitvTimestampOffset:], timestamp)
		binary.BigEndian.PutUint64(plainData[InitvServerIdOffset:], serverId)
	} else {
		copy(plainData[InitvBase:], initVector)
	}

	return plainData
}

// xorPayloadToHmacPad 使用 HMAC pad 对 payload 进行异或操作
func (d *DoubleClickCrypto) xorPayloadToHmacPad(workBytes []byte) {
	payloadSize := len(workBytes) - OverheadSize
	sections := (payloadSize + counterPageSize - 1) / counterPageSize

	if sections > counterSections {
		panic(fmt.Sprintf("Payload is %d bytes, exceeds limit of %d",
			payloadSize, counterPageSize*counterSections))
	}

	pad := make([]byte, counterPageSize+3)
	counterSize := 0

	for section := 0; section < sections; section++ {
		sectionBase := section * counterPageSize
		sectionSize := int(math.Min(float64(payloadSize-sectionBase), float64(counterPageSize)))

		// 生成 HMAC pad
		hmacPad, err := d.generateHmacPad(workBytes[InitvBase:InitvBase+InitvSize], pad, counterSize)
		if err != nil {
			panic(err)
		}

		// 对 payload 进行异或操作
		for i := 0; i < sectionSize; i++ {
			workBytes[PayloadBase+sectionBase+i] ^= hmacPad[i]
		}

		// 清空 pad
		for i := 0; i < counterPageSize; i++ {
			pad[i] = 0
		}

		if counterSize == 0 || pad[counterPageSize+counterSize-1] == 0 {
			counterSize++
		}
	}
}

// generateHmacPad 生成 HMAC pad
func (d *DoubleClickCrypto) generateHmacPad(initVector, pad []byte, counterSize int) ([]byte, error) {
	mac := hmac.New(sha1.New, d.keys.encryptionKey)
	mac.Write(initVector)

	if counterSize != 0 {
		mac.Write(pad[counterPageSize : counterPageSize+counterSize])
	}

	hmacResult := mac.Sum(nil)
	copy(pad, hmacResult)

	return pad[:counterPageSize], nil
}

// hmacSignature 计算 HMAC 签名
func (d *DoubleClickCrypto) hmacSignature(workBytes []byte) (uint32, error) {
	mac := hmac.New(sha1.New, d.keys.integrityKey)
	mac.Write(workBytes[PayloadBase : len(workBytes)-OverheadSize])
	mac.Write(workBytes[InitvBase : InitvBase+InitvSize])

	hmacResult := mac.Sum(nil)

	// 取前4个字节作为签名
	return binary.BigEndian.Uint32(hmacResult[:4]), nil
}

// dump 调试输出
// func dump(header string, inData, workBytes []byte) string {
// 	timestamp := binary.BigEndian.Uint64(workBytes[InitvBase+InitvTimestampOffset:])
// 	serverId := binary.BigEndian.Uint64(workBytes[InitvBase+InitvServerIdOffset:])

// 	timeObj := time.Unix(0, int64(secsAndMicrosToMillis(timestamp))*int64(time.Millisecond))

// 	return fmt.Sprintf("%s: initVector={timestamp %s, serverId %d}, input =%s, output =%s",
// 		header,
// 		timeObj.Format(time.RFC3339),
// 		serverId,
// 		hex.EncodeToString(inData),
// 		hex.EncodeToString(workBytes))
// }

// createHmac 创建 HMAC 实例
func createHmac(key []byte) ([]byte, error) {
	mac := hmac.New(sha1.New, key)
	return mac.Sum(nil), nil
}

// millisToSecsAndMicros 将毫秒转换为秒和微秒格式
func millisToSecsAndMicros(timestamp int64) uint64 {
	seconds := timestamp / 1000
	micros := (timestamp % 1000) * 1000
	return (uint64(seconds) << 32) | uint64(micros)
}

// secsAndMicrosToMillis 将秒和微秒格式转换为毫秒
// func secsAndMicrosToMillis(secondsAndMicros uint64) int64 {
// 	seconds := secondsAndMicros >> 32
// 	micros := secondsAndMicros & 0xFFFFFFFF
// 	return int64(seconds)*1000 + int64(micros)/1000
// }

// randUint64 生成随机 uint64
func randUint64() uint64 {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		// 如果加密随机数失败，使用时间戳作为后备
		return uint64(time.Now().UnixNano())
	}
	return binary.BigEndian.Uint64(buf)
}

// Price 用于赢取价格加密
type Price struct {
	*DoubleClickCrypto
}

// NewPrice 创建价格加密实例
func NewPrice(keys *Keys) *Price {
	return &Price{DoubleClickCrypto: NewDoubleClickCrypto(keys)}
}

// EncryptPriceMicros 加密赢取价格（微秒）
func (p *Price) EncryptPriceMicros(priceValue uint64, initVector []byte) ([]byte, error) {
	plainData := p.initPlainData(8, initVector)
	binary.BigEndian.PutUint64(plainData[PayloadBase:], priceValue)
	return p.encrypt(plainData)
}

// DecryptPriceMicros 解密赢取价格（微秒）
func (p *Price) DecryptPriceMicros(priceCipher []byte) (uint64, error) {
	if len(priceCipher) != (OverheadSize + 8) {
		return 0, fmt.Errorf("Price is %d bytes, should be %d",
			len(priceCipher), OverheadSize+8)
	}

	plainData, err := p.decrypt(priceCipher)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(plainData[PayloadBase:]), nil
}

// EncodePriceMicros 加密并编码赢取价格（微秒）
func (p *Price) EncodePriceMicros(priceMicros uint64, initVector []byte) (string, error) {
	encrypted, err := p.EncryptPriceMicros(priceMicros, initVector)
	if err != nil {
		return "", err
	}
	return p.encode(encrypted), nil
}

// EncodePriceValue 加密并编码赢取价格（浮点值）
func (p *Price) EncodePriceValue(priceValue float64, initVector []byte) (string, error) {
	priceMicros := uint64(priceValue * microsPerCurrencyUnit)
	return p.EncodePriceMicros(priceMicros, initVector)
}

// DecodePriceMicros 解码并解密赢取价格（微秒）
func (p *Price) DecodePriceMicros(priceCipher string) (uint64, error) {
	if priceCipher == "" {
		return 0, fmt.Errorf("priceCipher cannot be empty")
	}

	decoded, err := p.decode(priceCipher)
	if err != nil {
		return 0, err
	}

	return p.DecryptPriceMicros(decoded)
}

// DecodePriceValue 解码并解密赢取价格（浮点值）
func (p *Price) DecodePriceValue(priceCipher string) (float64, error) {
	micros, err := p.DecodePriceMicros(priceCipher)
	if err != nil {
		return 0, err
	}
	return float64(micros) / float64(microsPerCurrencyUnit), nil
}

// AdId 用于广告ID加密
type AdId struct {
	*DoubleClickCrypto
}

// NewAdId 创建广告ID加密实例
func NewAdId(keys *Keys) *AdId {
	return &AdId{DoubleClickCrypto: NewDoubleClickCrypto(keys)}
}

// EncryptAdId 加密广告ID
func (a *AdId) EncryptAdId(adidPlain, initVector []byte) ([]byte, error) {
	if len(adidPlain) != 16 {
		return nil, fmt.Errorf("AdId is %d bytes, should be 16", len(adidPlain))
	}

	plainData := a.initPlainData(16, initVector)
	copy(plainData[PayloadBase:], adidPlain)
	return a.encrypt(plainData)
}

// DecryptAdId 解密广告ID
func (a *AdId) DecryptAdId(adidCipher []byte) ([]byte, error) {
	if len(adidCipher) != (OverheadSize + 16) {
		return nil, fmt.Errorf("AdId is %d bytes, should be %d",
			len(adidCipher), OverheadSize+16)
	}

	plainData, err := a.decrypt(adidCipher)
	if err != nil {
		return nil, err
	}

	return plainData[PayloadBase : len(plainData)-SignatureSize], nil
}

// Idfa 用于IDFA加密
type Idfa struct {
	*DoubleClickCrypto
}

// NewIdfa 创建IDFA加密实例
func NewIdfa(keys *Keys) *Idfa {
	return &Idfa{DoubleClickCrypto: NewDoubleClickCrypto(keys)}
}

// EncryptIdfa 加密IDFA
func (i *Idfa) EncryptIdfa(idfaPlain, initVector []byte) ([]byte, error) {
	plainData := i.initPlainData(len(idfaPlain), initVector)
	copy(plainData[PayloadBase:], idfaPlain)
	return i.encrypt(plainData)
}

// DecryptIdfa 解密IDFA
func (i *Idfa) DecryptIdfa(idfaCipher []byte) ([]byte, error) {
	plainData, err := i.decrypt(idfaCipher)
	if err != nil {
		return nil, err
	}
	return plainData[PayloadBase : len(plainData)-SignatureSize], nil
}

// EncodeIdfa 加密并编码IDFA
func (i *Idfa) EncodeIdfa(idfaPlain, initVector []byte) (string, error) {
	encrypted, err := i.EncryptIdfa(idfaPlain, initVector)
	if err != nil {
		return "", err
	}
	return i.encode(encrypted), nil
}

// DecodeIdfa 解码并解密IDFA
func (i *Idfa) DecodeIdfa(idfaCipher string) ([]byte, error) {
	decoded, err := i.decode(idfaCipher)
	if err != nil {
		return nil, err
	}
	return i.DecryptIdfa(decoded)
}

// Hyperlocal 用于超本地信息加密
type Hyperlocal struct {
	*DoubleClickCrypto
}

// NewHyperlocal 创建超本地信息加密实例
func NewHyperlocal(keys *Keys) *Hyperlocal {
	return &Hyperlocal{DoubleClickCrypto: NewDoubleClickCrypto(keys)}
}

// EncryptHyperlocal 加密超本地信息
func (h *Hyperlocal) EncryptHyperlocal(hyperlocalPlain, initVector []byte) ([]byte, error) {
	plainData := h.initPlainData(len(hyperlocalPlain), initVector)
	copy(plainData[PayloadBase:], hyperlocalPlain)
	return h.encrypt(plainData)
}

// DecryptHyperlocal 解密超本地信息
func (h *Hyperlocal) DecryptHyperlocal(hyperlocalCipher []byte) ([]byte, error) {
	plainData, err := h.decrypt(hyperlocalCipher)
	if err != nil {
		return nil, err
	}
	return plainData[PayloadBase : len(plainData)-SignatureSize], nil
}
