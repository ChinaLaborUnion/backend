package payutil

import (
	"encoding/json"
	"errors"
	"fmt"
	wxconstant "grpc-demo/constants/payment/wxpayment"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// 获取公钥
const publicKeyUrl = "https://api.mch.weixin.qq.com/v3/certificates"

type TokenResponse struct {
	Data []TokenResponseData `json:"data"`
}
type TokenResponseData struct {
	EffectiveTime      string             `json:"effective_time"`
	EncryptCertificate EncryptCertificate `json:"encrypt_certificate"`
	ExpireTime         string             `json:"expire_time"`
	SerialNo           string             `json:"serial_no"`
}
type EncryptCertificate struct {
	Algorithm      string `json:"algorithm"`
	AssociatedData string `json:"associated_data"`
	Ciphertext     string `json:"ciphertext"`
	Nonce          string `json:"nonce"`
}

var publicSyncMap sync.Map

// 获取公钥
func getPublicKey() (key string, err error) {
	var prepareTime int64 = 24 * 3600 * 3 // 证书提前三天过期旧证书，获取新证书
	nowTime := time.Now().Unix()
	// 读取公钥缓存数据
	cacheValueKey := fmt.Sprintf("app_id:%s:public_key:value", wxconstant.AppId)
	cacheExpireTimeKey := fmt.Sprintf("app_id:%s:public_key:expire_time", wxconstant.AppId)
	cacheValue, keyValueOk := publicSyncMap.Load(cacheValueKey)
	cacheExpireTime, expireTimeOk := publicSyncMap.Load(cacheExpireTimeKey)
	if keyValueOk && expireTimeOk {
		// 格式化时间
		local, _ := time.LoadLocation("Local")
		location, _ := time.ParseInLocation(time.RFC3339, cacheExpireTime.(string), local)
		// 判断是否过期，证书没有过期直接返回
		if location.Unix()-prepareTime > nowTime {
			return cacheValue.(string), nil
		}
	}
	token, err := authorization(http.MethodGet, nil, publicKeyUrl)
	if err != nil {
		return key, err
	}
	request, err := http.NewRequest(http.MethodGet, publicKeyUrl, nil)
	if err != nil {
		return key, err
	}
	request.Header.Add("Authorization", "WECHATPAY2-SHA256-RSA2048 "+token)
	request.Header.Add("User-Agent", "用户代理(https://zh.wikipedia.org/wiki/User_agent)")
	request.Header.Add("Content-type", "application/json;charset='utf-8'")
	request.Header.Add("Accept", "application/json")
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return key, err
	}
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return key, err
	}
	//fmt.Println(string(bodyBytes))
	var tokenResponse TokenResponse
	if err = json.Unmarshal(bodyBytes, &tokenResponse); err != nil {
		return key, err
	}
	for _, encryptCertificate := range tokenResponse.Data {
		// 格式化时间
		local, _ := time.LoadLocation("Local")
		location, err := time.ParseInLocation(time.RFC3339, encryptCertificate.ExpireTime, local)
		if err != nil {
			return key, err
		}
		// 判断是否过期，证书没有过期直接返回
		if location.Unix()-prepareTime > nowTime {
			decryptBytes, err := DecryptGCM(wxconstant.AesKey, encryptCertificate.EncryptCertificate.Nonce, encryptCertificate.EncryptCertificate.Ciphertext,
				encryptCertificate.EncryptCertificate.AssociatedData)
			if err != nil {
				return key, err
			}
			key = string(decryptBytes)
			publicSyncMap.Store(cacheValueKey, key)
			publicSyncMap.Store(cacheExpireTimeKey, encryptCertificate.ExpireTime)
			return key, nil
		}
	}
	return key, errors.New("get public key error")
}