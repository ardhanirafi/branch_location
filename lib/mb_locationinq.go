package lib

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	// var ReqStruct MbLocationReq
	var ResultElements []MbLocationRespResult

	// if err := c.ShouldBindJSON(&ReqStruct); err != nil {
	// 	c.JSON(400, gin.H{"Error ": err.Error()})
	// 	return
	// }
	// reqJSON, errdec := AESDecryptRequestJSON(c)
	// if errdec != "" {
	// 	return
	// }
	// json.Unmarshal(reqJSON, &ReqStruct)

	// if err := c.ShouldBindBodyWith(&ReqStruct, binding.JSON); err != nil {
	// 	log.Print(err.Error())
	// 	c.JSON(400, gin.H{"Error": err.Error()})
	// 	return
	// }

	db, err := sql.Open("mysql", "root:@tcp("+"localhost"+":"+"3306"+")/mnc")
	if err != nil {
		fmt.Print(err.Error())
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
		c.AbortWithError(500, err)
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
		Message:     "SUKSES BRO",
		MessageCode: "M1",
		Results:     ResultElements,
	}
	// Respjson, _ := json.Marshal(RespVer)

	c.JSON(200, RespVer)
	// c.String(200, string(Respjson))
	// c.JSON(200, gin.H{
	// 	"Content": AESEncrypt(string(Respjson)),
	// })
}

// Inside the lib package

func MbDeleteATM(c *gin.Context) {
	//mengambil nilai dari parameter "id" dari URL dengan menggunakan c.Param("id")
	ID := c.Param("id")
	if ID == "" {
		c.JSON(400, gin.H{"error": "Invalid request data: ID parameter is required"})
		return
	}
	//konversi menjadi integer menggunakan strconv.Atoi(ID)
	idInt, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data: ID parameter must be an integer"})
		return
	}
	/*
		koneksi ke database Parameter pertama adalah tipe database yang digunakan ("mysql"),
		diikuti oleh informasi koneksi seperti username, password, host, dan port.
		Jika terjadi error dalam pembukaan koneksi,
		fungsi akan memberikan respons HTTP dengan status 500 (Internal Server Error) dan mengembalikan error yang terjadi
	*/
	db, err := sql.Open("mysql", "root:@tcp("+"localhost"+":"+"3306"+")/mnc")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	//defer db.Close(). Hal ini memastikan bahwa koneksi database akan ditutup
	defer db.Close()

	_, err = db.Exec("DELETE FROM mb_atm_location WHERE id = ?", idInt)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	/*
		Jika eksekusi query berhasil, fungsi akan memberikan respons HTTP dengan status 200 (OK)
		dan mengembalikan objek JSON yang berisi pesan sukses bahwa data ATM berhasil dihapus.
	*/
	c.JSON(200, gin.H{"message": "ATM data deleted successfully"})
}

func MbDeleteBank(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		c.JSON(400, gin.H{"error": "Invalid request data: ID parameter is required"})
		return
	}

	idInt, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data: ID parameter must be an integer"})
		return
	}

	db, err := gorm.Open(mysql.Open("root:@tcp("+"localhost"+":"+"3306"+")/mnc"), &gorm.Config{})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	var mbBankLocation mb_bank_location
	result := db.Where("id = ?", idInt).Delete(&mbBankLocation)
	if result.Error != nil {
		c.AbortWithError(500, result.Error)
		return
	}

	c.JSON(200, gin.H{"message": "Bank data deleted successfully"})
}

func MbBankAdd(c *gin.Context) {
	var req MbLocationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	dsn := "root:@tcp(" + "localhost" + ":" + "3306" + ")/mnc"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	mbBankLocations := mb_bank_location{
		Utype:     req.Utype,
		Name:      req.Name,
		Address:   req.Addr,
		City:      req.City,
		Postcode:  req.Postcode,
		State:     req.State,
		Country:   req.Country,
		Phone:     req.Phone,
		Fax:       req.Fax,
		Email:     req.Email,
		Bizhour:   req.Bizhour,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		CreatedAt: time.Now(),
	}

	var request map[string]interface{}
	temp, _ := json.Marshal(mbBankLocations)
	json.Unmarshal(temp, &request)

	result := db.Debug().Table("mb_bank_location").Create(&request)
	if result.Error != nil {
		c.AbortWithError(500, result.Error)
		return
	}

	c.JSON(200, gin.H{"message": "Bank data added successfully"})
}

func MbAtmAdd(c *gin.Context) {
	var req MbLocationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentTime := time.Now().Format(time.RFC3339)

	db, err := sql.Open("mysql", "root:@tcp("+"localhost"+":"+"3306"+")/mnc")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT INTO mb_atm_location (
			utype, NAME, address, city, postcode, state, country,
			phone, fax, email, biz_hour, latitude, longitude,
			created, modified
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.Utype, req.Name, req.Addr, req.City, req.Postcode, req.State, req.Country,
		req.Phone, req.Fax, req.Email, req.Bizhour, req.Latitude, req.Longitude, currentTime, nil)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, gin.H{"message": "ATM data added successfully"})
}

// Fungsi untuk mengupdate data bank
func MbBankUpdate(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	idInt, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data: ID parameter must be an integer"})
		return
	}

	var req MbLocationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentTime := time.Now().Format(time.RFC3339)

	db, err := sql.Open("mysql", "root:@tcp("+"localhost"+":"+"3306"+")/mnc")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	defer db.Close()

	updateQuery := "UPDATE mb_bank_location SET"
	values := []interface{}{}

	updateQuery += " modified = ?,"
	values = append(values, currentTime)

	if req.Name != "" {
		updateQuery += " NAME = ?,"
		values = append(values, req.Name)
	}

	if req.Email != "" {
		updateQuery += " email = ?,"
		values = append(values, req.Email)
	}

	if req.Utype != "" {
		updateQuery += " utype = ?,"
		values = append(values, req.Utype)
	}

	if req.Addr != "" {
		updateQuery += " address = ?,"
		values = append(values, req.Addr)
	}

	if req.City != "" {
		updateQuery += " city = ?,"
		values = append(values, req.City)
	}

	if req.Postcode != "" {
		updateQuery += " postcode = ?,"
		values = append(values, req.Postcode)
	}

	if req.State != "" {
		updateQuery += " state = ?,"
		values = append(values, req.State)
	}

	if req.Country != "" {
		updateQuery += " country = ?,"
		values = append(values, req.Country)
	}

	if req.Phone != "" {
		updateQuery += " phone = ?,"
		values = append(values, req.Phone)
	}

	if req.Fax != "" {
		updateQuery += " fax = ?,"
		values = append(values, req.Fax)
	}

	if req.Bizhour != "" {
		updateQuery += " biz_hour = ?,"
		values = append(values, req.Bizhour)
	}

	if req.Latitude != "" {
		updateQuery += " latitude = ?,"
		values = append(values, req.Latitude)
	}

	if req.Longitude != "" {
		updateQuery += " longitude = ?,"
		values = append(values, req.Longitude)
	}

	updateQuery = strings.TrimSuffix(updateQuery, ",")

	updateQuery += " WHERE id = ?"
	values = append(values, idInt)

	_, err = db.Exec(updateQuery, values...)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, gin.H{"message": "Bank data updated successfully"})
}

// Fungsi untuk mengupdate data ATM
func MbAtmUpdate(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	idInt, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data: ID parameter must be an integer"})
		return
	}

	var req MbLocationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentTime := time.Now().Format(time.RFC3339)

	db, err := sql.Open("mysql", "root:@tcp("+"localhost"+":"+"3306"+")/mnc")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	defer db.Close()

	updateQuery := "UPDATE mb_atm_location SET"
	values := []interface{}{}

	updateQuery += " modified = ?,"
	values = append(values, currentTime)

	if req.Name != "" {
		updateQuery += " NAME = ?,"
		values = append(values, req.Name)
	}

	if req.Email != "" {
		updateQuery += " email = ?,"
		values = append(values, req.Email)
	}

	if req.Utype != "" {
		updateQuery += " utype = ?,"
		values = append(values, req.Utype)
	}

	if req.Addr != "" {
		updateQuery += " address = ?,"
		values = append(values, req.Addr)
	}

	if req.City != "" {
		updateQuery += " city = ?,"
		values = append(values, req.City)
	}

	if req.Postcode != "" {
		updateQuery += " postcode = ?,"
		values = append(values, req.Postcode)
	}

	if req.State != "" {
		updateQuery += " state = ?,"
		values = append(values, req.State)
	}

	if req.Country != "" {
		updateQuery += " country = ?,"
		values = append(values, req.Country)
	}

	if req.Phone != "" {
		updateQuery += " phone = ?,"
		values = append(values, req.Phone)
	}

	if req.Fax != "" {
		updateQuery += " fax = ?,"
		values = append(values, req.Fax)
	}

	if req.Bizhour != "" {
		updateQuery += " biz_hour = ?,"
		values = append(values, req.Bizhour)
	}

	if req.Latitude != "" {
		updateQuery += " latitude = ?,"
		values = append(values, req.Latitude)
	}

	if req.Longitude != "" {
		updateQuery += " longitude = ?,"
		values = append(values, req.Longitude)
	}

	// Menghilangkan koma terakhir dari query
	updateQuery = strings.TrimSuffix(updateQuery, ",")

	updateQuery += " WHERE id = ?"
	values = append(values, idInt)

	_, err = db.Exec(updateQuery, values...)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, gin.H{"message": "ATM data updated successfully"})
}
