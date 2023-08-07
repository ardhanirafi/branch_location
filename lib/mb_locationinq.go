package lib

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/go-sql-driver/mysql"
)

func MbLocInq(c *gin.Context) {
	var utype string
	var name string
	var address string
	var city string
	var postcode string
	var state string
	var country string
	var phone string
	var fax string
	var email string
	var biz_hour string
	var latitude string
	var longitude string
	var ReqStruct MbLocationReq
	var ResultElements []MbLocationRespResult

	// if err := c.ShouldBindJSON(&ReqStruct); err != nil {
	// 	c.JSON(400, gin.H{"Error ": err.Error()})
	// 	return
	// }
	reqJSON, errdec := AESDecryptRequestJSON(c)
	if errdec != "" {
		return
	}
	json.Unmarshal(reqJSON, &ReqStruct)
	if err := c.ShouldBindBodyWith(&ReqStruct, binding.JSON); err != nil {
		MotionRes := MotionJSON{
			Error: err.Error(),
		}

		jsonRes, _ := json.Marshal(MotionRes)
		encryptRes := AESEncrypt(string(jsonRes))
		c.JSON(400, gin.H{"Content": encryptRes})
		return
	}
	db, err := sql.Open("mysql", "varuser:mnc123@tcp("+dbaddr+":"+dbport+")/mncmbank")
	if err != nil {
		log.Println(err.Error())
	}
	defer db.Close()
	queryget, err := db.Query(`SELECT
				utype,
				NAME,
				address,
				city,
				postcode,
				state,
				country,
				phone,
				fax,
				email,
				biz_hour,
				latitude,
				longitude
			FROM
				mb_atm_location
			UNION ALL
			SELECT
				utype,
				NAME,
				address,
				city,
				postcode,
				state,
				country,
				phone,
				fax,
				email,
				biz_hour,
				latitude,
				longitude
			FROM
				mb_bank_location
			`)
	if err != nil {
		log.Println(err.Error())
	}
	for queryget.Next() {
		queryget.Scan(&utype, &name, &address, &city, &postcode, &state, &country, &phone, &fax, &email, &biz_hour, &latitude, &longitude)
		ResultElements = append(ResultElements, MbLocationRespResult{
			Unittype:  strings.TrimSpace(utype),
			Name:      strings.TrimSpace(name),
			Addr:      strings.TrimSpace(address),
			City:      strings.TrimSpace(city),
			Postcode:  strings.TrimSpace(postcode),
			State:     strings.TrimSpace(state),
			Country:   strings.TrimSpace(country),
			Phone:     strings.TrimSpace(phone),
			Fax:       strings.TrimSpace(fax),
			Email:     strings.TrimSpace(email),
			Bizhour:   strings.TrimSpace(biz_hour),
			Latitude:  strings.TrimSpace(latitude),
			Longitude: strings.TrimSpace(longitude),
		})
	}
	queryget.Close()
	fmt.Println(ResultElements)
	RespVer := MbLocationResp{
		Success:     true,
		Message:     GetErrorDesc("M1", "Mobe"),
		MessageCode: "M1",
		Results:     ResultElements,
	}
	Respjson, _ := json.Marshal(RespVer)
	// c.String(200, string(Respjson))
	c.JSON(200, gin.H{
		"Content": AESEncrypt(string(Respjson)),
	})
}
