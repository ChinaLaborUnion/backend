package payutil

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	wxconstant "grpc-demo/constants/payment/wxpayment"
	util "grpc-demo/utils/hash"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// 对消息的散列值进行数字签名
func signPKCS1v15(msg, privateKey []byte, hashType crypto.Hash) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key decode error")
	}
	pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("parse private key error")
	}
	key, ok := pri.(*rsa.PrivateKey)
	if ok == false {
		return nil, errors.New("private key format error")
	}
	sign, err := rsa.SignPKCS1v15(rand.Reader, key, hashType, msg)

	if err != nil {
		return nil, errors.New("sign error")
	}
	return sign, nil
}

// base编码
func base64EncodeStr(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}
func base64DecodeStr(src string) string {
	sDec, _ := base64.StdEncoding.DecodeString(src)
	return string(sDec)
}

func authorization(method string, paramMap map[string]interface{}, rawUrl string) (token string, err error) {
	var body string
	if len(paramMap) != 0 {
		paramJsonBytes, err := json.Marshal(paramMap)
		if err != nil {
			return token, err
		}
		body = string(paramJsonBytes)
	}
	urlPart, err := url.Parse(rawUrl)
	if err != nil {
		return token, err
	}
	canonicalUrl := urlPart.RequestURI()
	timestamp := time.Now().Unix()
	nonce := util.GetRandomString(32)
	message := fmt.Sprintf("%s\n%s\n%d\n%s\n%s\n", method, canonicalUrl, timestamp, nonce, body)

	open, err := os.Open("apiclient_key.pem")
	if err != nil {
		return token, err
	}
	defer open.Close()
	privateKey, err := ioutil.ReadAll(open)
	if err != nil {
		return token, err
	}

	h := sha256.New()
	h.Write([]byte(message))
	sum := h.Sum(nil)

	signBytes, err := signPKCS1v15(sum, privateKey, crypto.SHA256)

	if err != nil {
		return token, err
	}
	sign := base64EncodeStr(signBytes)

	token = fmt.Sprintf("mchid=\"%s\",nonce_str=\"%s\",timestamp=\"%d\",serial_no=\"%s\",signature=\"%s\"",
		wxconstant.MchId, nonce, timestamp, wxconstant.PrivateSerialNo, sign)

	return token, nil
}

func DecryptGCM(aesKey, nonceV, ciphertextV, additionalDataV string) ([]byte, error) {
	key := []byte(aesKey)
	nonce := []byte(nonceV)
	additionalData := []byte(additionalDataV)
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextV)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, additionalData)
	if err != nil {
		return nil, err
	}
	return plaintext, err
}

//验证数字签名
func VerifyRsaSign(msg []byte, sign []byte, publicStr []byte, hashType crypto.Hash) bool {
	//pem解码
	block, _ := pem.Decode(publicStr)
	//x509解码
	publicKeyInterface, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}
	publicKey := publicKeyInterface.PublicKey.(*rsa.PublicKey)
	//验证数字签名
	err = rsa.VerifyPKCS1v15(publicKey, hashType, msg, sign) //crypto.SHA1
	return err == nil
}
// 验证签名
func notifyValidate(timeStamp ,nonce,rawPost,signature string) (bool, error) {
	signature = base64DecodeStr(signature)
	message := fmt.Sprintf("%s\n%s\n%s\n", timeStamp, nonce, rawPost)
	publicKey, err := getPublicKey()
	if err != nil {
		return false, err
	}

	h := sha256.New()
	h.Write([]byte(message))
	sum := h.Sum(nil)

	return VerifyRsaSign(sum, []byte(signature), []byte(publicKey), crypto.SHA256), nil
}

type NotifyResponse struct {
	ID         string    `json:"id"`
	CreateTime string   `json:"create_time"`
	EventType  string    `json:"event_type"`
	Resource NotifyResource `json:"resource"`
	summary    string 	`json:"summary"`
}
type NotifyResource struct {
	Algorithm   string 		`json:"algorithm"`
	Ciphertext  string `json:"ciphertext"`
	AssociatedData string `json:"associated_data"`
	OriginalType   string 	`json:"original_type"`
	Nonce   string `json:"nonce"`
}

func NotifyDecrypt(rawPost string) (decrypt string, err error) {
	var notifyResponse NotifyResponse
	if err = json.Unmarshal([]byte(rawPost), &notifyResponse); err != nil {
		return decrypt, err
	}
	decryptBytes, err := DecryptGCM(wxconstant.AesKey, notifyResponse.Resource.Nonce, notifyResponse.Resource.Ciphertext,
		notifyResponse.Resource.AssociatedData)
	if err != nil {
		return decrypt, err
	}
	decrypt = string(decryptBytes)
	return decrypt, nil
}

func ParseNotify(req *http.Request) (notifyRsp *NotifyResponse, err error) {
	notifyRsp = new(NotifyResponse)
	if err = json.NewDecoder(req.Body).Decode(notifyRsp); err != nil {
		return nil, fmt.Errorf("json.NewDecoder.Decode：%w", err)
	}
	return
}