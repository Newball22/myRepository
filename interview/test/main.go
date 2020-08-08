package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Session struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	SessionKey    string `json:"session_key"`
	SessionSecret string `json:"session_secret"`
	Scope         string `json:"scope"`
	ExpiresIn     int    `json:"expires_in"`
}

type Result struct {
	LogId     int64       `json:"log_id"`
	Direction int         `json:"direction"`
	Wrn       int         `json:"words_result_num"`
	Wr        WordsResult `json:"words_result"`
}

type WordsResult struct {
	Model        Words `json:"品牌型号"`
	IssueDate    Words `json:"发证日期"`
	UseCharacter Words `json:"使用性质"`
	EngineNo     Words `json:"发动机号码"`
	PlateNo      Words `json:"号牌号码"`
	Owner        Words `json:"所有人"`
	Address      Words `json:"住址"`
	RegisterDate Words `json:"注册日期"`
	Vin          Words `json:"车辆识别代号"`
	VehicleType  Words `json:"车辆类型"`
}

type Words struct {
	Words string `json:"words"`
}

type VehicleLicense struct {
	PlateNo      string
	VehicleType  string
	Owner        string
	Address      string
	UseCharacter string
	Model        string
	Vin          string
	EngineNo     string
	RegisterDate string
	IssueDate    string
}

const (
	clientID     = "0oVZhEvStHPnzkEv8GQwkdKX"         // API Key
	clientSecret = "rp6oOMvjzsbaUqtwY9dTXtZ5UGNbIPF9" // Secret Key
)

func main() {

	//url := "https://aip.baidubce.com/rest/2.0/ocr/v1/vehicle_license"
	file, _ := os.Open("1.JPG")
	resp, _ := ioutil.ReadAll(file)
	defer file.Close()
	obj := hex.EncodeToString(resp)
	fmt.Println(obj)
	//access_token := "24.0c0c3dbba47c2656698aa5824082a25b.2592000.1598899040.282335-21739964"
	//requestUrl := url + "?access_token=" + access_token
	//contentType := "application/x-www-form-urlencoded"
	//data, _ := http.Post(requestUrl, contentType, resp)
	//fmt.Println(data)

	//token, err := accessToken(clientID, clientSecret)
	//if err != nil {
	//	fmt.Println("accessToken()", err)
	//	return
	//}
	//fmt.Println("Token:", *token)
	//
	//img, err := ioutil.ReadFile("1.JPG")
	//if err != nil {
	//	fmt.Println("ioutil.ReadFile()", err)
	//	return
	//}
	//var vl VehicleLicense
	//err = ImageToText(*token, img, &vl)
	//if err != nil {
	//	fmt.Println("ImageToText()", err)
	//	return
	//}
	//fmt.Println("PlateNo: ", vl.PlateNo)
	//fmt.Println("VIN: ", vl.Vin)

}

func accessToken(id string, secret string) (token *string, err error) {
	apiURL := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s", id, secret)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var session Session
	err = json.Unmarshal(body, &session)
	if err != nil {
		return nil, err
	}
	token = &session.AccessToken
	return token, nil
}

func ImageToText(token string, image []byte, vl *VehicleLicense) error {
	apiURL := fmt.Sprintf("https://aip.baidubce.com/rest/2.0/ocr/v1/vehicle_license?access_token=%s", token)
	param := "image=" + url.QueryEscape(base64.StdEncoding.EncodeToString(image))
	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(param))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	convert(result.Wr, vl)
	return nil
}

func convert(wr WordsResult, vl *VehicleLicense) {
	vl.PlateNo = wr.PlateNo.Words
	vl.VehicleType = wr.VehicleType.Words
	vl.Owner = wr.Owner.Words
	vl.Address = wr.Address.Words
	vl.UseCharacter = wr.UseCharacter.Words
	vl.Model = wr.Model.Words
	vl.Vin = wr.Vin.Words
	vl.EngineNo = wr.EngineNo.Words
	vl.RegisterDate = wr.RegisterDate.Words
	vl.IssueDate = wr.IssueDate.Words
	return
}
