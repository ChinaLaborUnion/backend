package payutil

import (
	"bytes"
	"crypto"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	wxconstant "grpc-demo/constants/payment/wxpayment"
	"io/ioutil"
	"net/http"
	"os"
)

// 统一下单接口
func UnifiedOrder(paramMap map[string]interface{}) (payResMap map[string]string, err error) {
	payResMap = make(map[string]string)
	token, err := authorization(http.MethodPost, paramMap, "https://api.mch.weixin.qq.com/v3/pay/transactions/app")

	if err != nil {
		return payResMap, err
	}
	marshal, _ := json.Marshal(paramMap)
	request, err := http.NewRequest(http.MethodPost, "https://api.mch.weixin.qq.com/v3/pay/transactions/app", bytes.NewReader(marshal))
	if err != nil {
		return payResMap, err
	}
	request.Header.Add("Authorization", "WECHATPAY2-SHA256-RSA2048 "+token)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2486.0 Safari/537.36 Edge/13.10586")
	request.Header.Add("Content-type", "application/json;charset='utf-8'")
	request.Header.Add("Accept", "application/json")
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return payResMap, err
	}
	defer func() {
		response.Body.Close()
	}()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return payResMap, err
	}
	if err = json.Unmarshal(bodyBytes, &payResMap); err != nil {
		return payResMap, err
	}
	if payResMap["prepay_id"] == "" {
		return payResMap, errors.New("code:" + payResMap["code"] + "err:" + payResMap["message"])
	}
	return payResMap, nil
}

func AppSign(payResMap map[string]string, nonce, timeStamp string) (payJson string, err error) {
	payMap := make(map[string]string)

	packageStr := "Sign=WXPay"
	payMap["appId"] = wxconstant.AppId
	payMap["partnerid"] = wxconstant.MchId
	payMap["timeStamp"] = fmt.Sprintf("%v", timeStamp)
	payMap["nonceStr"] = nonce
	payMap["package"] = packageStr
	payMap["prepayid"] = payResMap["prepay_id"]
	// 签名
	message := fmt.Sprintf("%s\n%s\n%s\n%s\n", wxconstant.AppId, fmt.Sprintf("%v", timeStamp), nonce, payMap["prepayid"])
	open, err := os.Open("apiclient_key.pem")
	if err != nil {
		return payJson, err
	}
	defer open.Close()
	privateKey, err := ioutil.ReadAll(open)
	if err != nil {
		return payJson, err
	}

	h := sha256.New()
	h.Write([]byte(message))
	sum := h.Sum(nil)

	signBytes, err := signPKCS1v15(sum, privateKey, crypto.SHA256)
	if err != nil {
		return payJson, err
	}
	sign := base64EncodeStr(signBytes)
	return sign, nil
}
