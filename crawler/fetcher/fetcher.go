package fetcher

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var rateLimiter = time.Tick(200 * time.Millisecond)

//该模块任务是获取目标URL的网页数据
func Fetch(url string) ([]byte, error) {
	<-rateLimiter
	client := &http.Client{}
	//newUrl := strings.Replace(url, "http://", "https://", 1)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	//cookieTest := "Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1594664194; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1594438856; _efmdata=JhBdr2t6iy7DB%2FgIRFmvSWAUgEhpOW6KvlJ2nnPWQImapuDzpZU5g4xCADfwqcaJoQqIlAl6c6JuMNu3t0uXTUxFEhtHhJ%2FqNoTN5Zp3lgM%3D; _exid=P%2BfV4Jl9OMHfcscNpmjsN%2FU5%2BEe3Ef72OlM6UVdlIt3u8cxS7wsdLFgPGvr%2BNxPb0gQscEVPcLwKjagCCUEMUA%3D%3D; ec=8S6vb77P-1594631624118-977d9f0fe7d251582046124; __channelId=905819%2C0; notificationPreAuthorizeSwitch=50246; sid=2552788a-cb58-46b5-aa6b-0807bb140c8b"
	req.Header.Add("User-Agent", getRanAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	//req.Header.Add("Referer", "https://album.zhenai.com/u/1577319979")
	//req.Header.Add("cookie", cookieTest)
	//req.Host = " www.zhenai.com"
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	//出错处理
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong state code %d", resp.StatusCode)
	}

	//把网页转为utf-8编码
	bodyReader := bufio.NewReader(resp.Body)

	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)

}

func getRanAgent() string {
	agent := [...]string{
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
		"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
		"User-Agent,Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"User-Agent, Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
		"User-Agent,Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := len(agent)
	return agent[r.Intn(len)]

}

//转码工具函数
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error %v\n", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e

}
