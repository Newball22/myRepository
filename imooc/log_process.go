package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//整个功能模块的封装
type LogProcess struct {
	rc    chan []byte
	wc    chan *Message
	read  Reader
	write Writer
}

//日志信息
type Message struct {
	TimeLoacal                   time.Time
	BytesSent                    int //流量
	Path, Method, Scheme, Status string
	UpstreamTime, RequestTime    float64
}

//定义一个写入模块的接口，解藕，为了以后更多写入方式的拓展
type Reader interface {
	Read(rc chan []byte)
}

//定义一个读取模块的接口，解藕，为了以后更多读取方式的拓展
type Writer interface {
	Write(wc chan *Message)
}

type ReadFromFile struct {
	path string
}

//读取方式为ReadFromFile
func (r *ReadFromFile) Read(rc chan []byte) {
	//读取文件内容
	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open fail failed err:%s\n", err.Error()))
	}
	f.Seek(0, 2) //把文件的读取指针指向最后
	defer f.Close()
	file := bufio.NewReader(f)
	for {
		line, err := file.ReadBytes('\n') //读取一行的内容
		if err == io.EOF {                //读取到文章末尾,暂时没有数据可读
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadBytes file err:%s\n", err.Error()))
		}
		rc <- line[:len(line)-1] //把最后一个换行符忽略

	}

}

type WriteToInfluxDB struct {
	influxDBDsn string
}

//写入方式为WriteToInfluxDB
func (w *WriteToInfluxDB) Write(wc chan *Message) {
	for value := range wc {
		fmt.Println("写入数据库的数据:", value)
	}

}

//解析模块是为了提取数据中有价值的数据
func (l *LogProcess) Process() {
	//解析模块
	//172.0.0.12 - - [04/Mar/2018:13:49:52 +0000] http GET /foo?query=t HTTP/1.0 200 2133 - KeepAliveClient - 1.005 1.854

	str := `([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`
	r := regexp.MustCompile(str)

	loc, _ := time.LoadLocation("Asia/Shanghai")
	for value := range l.rc {
		ret := r.FindStringSubmatch(string(value))
		log.Println(len(ret))
		if len(ret) != 14 {
			log.Println("FindStringSubmatch failed err:", string(value))
			continue
		}
		message := &Message{}
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], loc)
		if err != nil {
			log.Println("ParseInLocation转换时间戳异常，err:", err.Error(), ret[4])
		}
		message.TimeLoacal = t

		byteSent, _ := strconv.Atoi(ret[8])
		message.BytesSent = byteSent

		//GET /foo?/query=t HTTP/1.0
		reqSli := strings.Split(ret[6], " ")
		if len(reqSli) != 3 {
			log.Println("strings.Split err:", err.Error())
			continue
		}
		message.Method = reqSli[0]

		urlPath, err := url.Parse(reqSli[1])
		if err != nil {
			log.Println("url.Parse failed err:", err.Error())
			continue
		}

		message.Path = urlPath.Path
		message.Scheme = ret[5]
		message.Status = ret[7]

		message.UpstreamTime, _ = strconv.ParseFloat(ret[12], 64)
		message.RequestTime, _ = strconv.ParseFloat(ret[13], 64)
		l.wc <- message

	}

}

func main() {
	r := &ReadFromFile{
		path: "./data/log.txt",
	}
	w := &WriteToInfluxDB{
		influxDBDsn: "userName&passWork",
	}

	lp := &LogProcess{
		rc:    make(chan []byte),
		wc:    make(chan *Message),
		read:  r,
		write: w,
	}
	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)
	time.Sleep(time.Second * 50)
	fmt.Println("Program Over")

}
