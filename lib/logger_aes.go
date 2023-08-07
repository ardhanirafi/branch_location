package lib

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	"runtime"
	// "fmt"

	"github.com/gin-gonic/gin"
	"github.com/sbecker/gin-api-demo/util"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

var (
	Graylog_Host = os.Getenv("Graylog_Host")
	Graylog_Port_Traffic = os.Getenv("Graylog_Port_Traffic")
	Graylog_Port_Error = os.Getenv("Graylog_Port_Error")
)

type LogJSON struct {
	ClientIP string `json:"clientip"`
	TimeStamp string `json:"timestamp"`
	Method string `json:"method"`
	Path string `json:"path"`
	StatusCode int `json:"statuscode"`
	Latency float64 `json:"latency"`
	Request string `json:"request"`
	Response string `json:"response"`
	//Request []Request `json:"request"`
	//Response []Response `json:"response"`
}

type bodyLogWriter struct {
    gin.ResponseWriter
    body *bytes.Buffer
}

const (
	reqFormat = `"request":`
	resFormat = `"response":`

	gelfStr = "gelf.NewWriter: %s"
)

func FancyHandleError(errCode error) (b bool) {
	
	if errCode != nil {
		gelfWriter, err := gelf.NewTCPWriter(Graylog_Host+":"+Graylog_Port_Error)

		if err != nil {
			log.Fatalf(gelfStr, err)
		}else{
			// log to both stderr and graylog2
			log.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
			log.SetFlags(0)
		}
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, fn, line, _ := runtime.Caller(1)

		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, errCode.Error())
		b = true
	}
	return
}

func Loggertogray(whattolog string) {
	gelfWriter, err := gelf.NewTCPWriter(Graylog_Host+":"+Graylog_Port_Traffic)
	if err != nil {
		log.Fatalf(gelfStr, err)
	}
	log.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
	log.SetFlags(0)
	log.Printf(whattolog)
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
    w.body.Write(b)
    return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
    w.body.WriteString(s)
    return w.ResponseWriter.WriteString(s)
}

func JSONLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start Timer
		start := time.Now()

		var bodyBytes []byte
        if c.Copy().Request.Body != nil {
          bodyBytes, _ = ioutil.ReadAll(c.Copy().Request.Body)
		}
		
		// Get Request Data
      	c.Copy().Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// Get Response Data
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
        c.Writer = blw
		c.Next()
		
		// Stop Timer
		duration := util.GetDurationInMillseconds(start)

		reqJSON := string(bodyBytes)

		rawReqJSON := reqJSON
		rawResJSON := blw.body.String()

		byteReq := []byte(rawReqJSON)
		byteRes := []byte(rawResJSON)
		
		var encryptReq MotionJSON
		var decryptRes MotionJSON
		var STRCTReq map[string]interface{}
		var STRCTRes map[string]interface{}
	
		json.Unmarshal(byteReq, &encryptReq)
		json.Unmarshal(byteRes, &decryptRes)

		

		// fmt.Println("biji : ",string(AESDecrypt(STRCTReq.Content)))
		jsonReq := AESDecrypt(encryptReq.Content)
		buffReq := new(bytes.Buffer)
		json.Compact(buffReq, jsonReq)
		// JSONReq, _ :=  json.MarshalIndent(decryptReq, "", "  ")
		// decryptedReq := buffReq.String()

		var resKey string
		if decryptRes.Error != "" {
			resKey = decryptRes.Error
		}else{
			resKey = decryptRes.Content
		}

		jsonRes := AESDecrypt(resKey)
		buffRes := new(bytes.Buffer)
		json.Compact(buffRes, jsonRes)
		// JSONRes, _ :=  json.MarshalIndent(decryptRes, "", "  ")
		// decryptedRes := buffRes.String()

		if string(jsonReq) == ""{
			json.Unmarshal(byteReq, &STRCTReq)
		}else{
			checkJSON := isJSON(string(jsonReq))
			switch checkJSON {
				case false :
					json.Unmarshal(byteReq, &STRCTReq)
				default :
				json.Unmarshal(buffReq.Bytes(), &STRCTReq)
			}
		}
		
		json.Unmarshal(buffRes.Bytes(), &STRCTRes)

		JSONReq, _ :=  json.MarshalIndent(STRCTReq, "", "  ")
		StringJSONReq := string(JSONReq)

		JSONRes, _ :=  json.MarshalIndent(STRCTRes, "", "  ")
		StringJSONRes := string(JSONRes)

		LogJSON := LogJSON{
			ClientIP: util.GetClientIP(c),
			TimeStamp: start.Format("2006-01-02 15:04:05"),
			Method: c.Request.Method,
			Path: c.Request.RequestURI,
			StatusCode: c.Writer.Status(),
			Latency: duration,
			Request:"",
			Response:"",
			//Request: string(bodyBytes),
			//Response: blw.body.String(),
		}

		if c.Writer.Status() >= 500 {
			//entry.Error(c.Errors.String())
			LogJSON.Response = c.Errors.String()
		}

		LogWrite, _ := json.MarshalIndent(LogJSON, "", "  ")
		LogWriteString := string(LogWrite)

		/*LogWriteString = strings.Replace(LogWriteString, reqFormat + ` ""`, reqFormat + "["+StringJSONReq+"]", -1)
		LogWriteString = strings.Replace(LogWriteString, resFormat + ` ""`, resFormat + "["+StringJSONRes+"]", -1)*/

		LogWriteString = strings.Replace(LogWriteString, reqFormat+` ""`, reqFormat+StringJSONReq, -1)
		LogWriteString = strings.Replace(LogWriteString, resFormat+` ""`, resFormat+StringJSONRes, -1)

		gelfWriter, err := gelf.NewTCPWriter(Graylog_Host+":"+Graylog_Port_Traffic)
		if err != nil {
			log.Fatalf(gelfStr, err)
		}

		// log to both stderr and graylog2
		log.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
		log.SetFlags(0)
		log.Printf(LogWriteString)
	}
}
