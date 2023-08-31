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

// ConnectToDatabase membangun dan membuka koneksi database MySQL menggunakan parameter yang diberikan.
// Mengembalikan instance *gorm.DB yang siap digunakan, atau error jika koneksi gagal.
func ConnectToDatabase(dbusn, dbaddr, dbport, dbtable string) (*gorm.DB, error) {
	dsn := dbusn + ":@tcp(" + dbaddr + ":" + dbport + ")/" + dbtable
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Inside the lib package

func MbDeleteATM(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		c.JSON(400, gin.H{"error": "Invalid request data: ID parameter is required"})
		return
	}

	// Initialize the GORM DB connection
	db, err := ConnectToDatabase(dbusn, dbaddr, dbport, dbtable)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	var mbAtmLocation MbAtmLocationsDelete

	// Delete the bank location entry
	result := db.Table("mb_atm_location").Where("id = ?", ID).Delete(&mbAtmLocation)

	if result.Error != nil {
		c.AbortWithError(500, result.Error)
		return
	}

	c.JSON(200, gin.H{"message": "Bank data deleted successfully"})
}

func MbDeleteBank(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		c.JSON(400, gin.H{"error": "Invalid request data: ID parameter is required"})
		return
	}

	// idInt, err := strconv.Atoi(ID)
	// if err != nil {
	// 	c.JSON(400, gin.H{"error": "Invalid request data: ID parameter must be an integer"})
	// 	return
	// }

	// Initialize the GORM DB connection
	db, err := ConnectToDatabase(dbusn, dbaddr, dbport, dbtable)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	var mbBankLocation MbBankLocationsDelete

	// Find the bank location entry by ID
	// result := db.Table("mb_bank_location").First(&mbBankLocation, idInt)

	// if result.Error != nil {
	// 	c.AbortWithError(500, result.Error)
	// 	return
	// }

	// Delete the bank location entry
	result := db.Table("mb_bank_location").Where("id = ?", ID).Delete(&mbBankLocation)

	if result.Error != nil {
		c.AbortWithError(500, result.Error)
		return
	}

	c.JSON(200, gin.H{"message": "Bank data deleted successfully"})
}

func MbBankAdd(c *gin.Context) {
	// Mengambil data JSON dari permintaan dan melakukan binding ke variabel req
	var req MbLocationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Membangun string DSN untuk koneksi ke basis data MySQL
	db, err := ConnectToDatabase(dbusn, dbaddr, dbport, dbtable)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// Membentuk data lokasi bank baru
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

	// Mengubah data menjadi format JSON
	var request map[string]interface{}
	temp, _ := json.Marshal(mbBankLocations)
	json.Unmarshal(temp, &request)

	// Memasukkan data ke dalam basis data
	result := db.Debug().Table("mb_bank_location").Create(&request)
	if result.Error != nil {
		c.AbortWithError(500, result.Error)
		return
	}

	// Mengirim respons sukses
	c.JSON(200, gin.H{"message": "Bank data added successfully"})
}

func MbAtmAdd(c *gin.Context) {
	// Mengambil data JSON dari permintaan dan melakukan binding ke variabel req
	var req MbLocationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db, err := ConnectToDatabase(dbusn, dbaddr, dbport, dbtable)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// Membentuk data lokasi ATM baru
	mbAtmLocations := mb_atm_location{
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

	// Mengubah data menjadi format JSON
	var request map[string]interface{}
	temp, _ := json.Marshal(mbAtmLocations)
	json.Unmarshal(temp, &request)

	// Mencetak data dalam format JSON
	fmt.Println(request)

	// Memasukkan data ke dalam basis data
	result := db.Debug().Table("mb_atm_location").Create(&request)
	if result.Error != nil {
		c.AbortWithError(500, result.Error)
		return
	}

	// Mengirim respons sukses
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
	fmt.Println(dbusn, dbaddr, dbport, dbtable)
	// Inisialisasi koneksi database menggunakan GORM
	db, err := ConnectToDatabase(dbusn, dbaddr, dbport, dbtable)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	var bankLocation MbBankLocation
	if err := db.Unscoped().Where("id = ?", idInt).First(&bankLocation).Error; err != nil {
		c.AbortWithError(500, err)
		return
	}

	// Update nilai-nilai yang diperlukan
	if req.Name != "" {
		bankLocation.Name = req.Name
	}

	if req.Addr != "" {
		bankLocation.Address = req.Addr
	}

	if req.Email != "" {
		bankLocation.Email = req.Email
	}

	if req.Utype != "" {
		bankLocation.Utype = req.Utype
	}

	if req.City != "" {
		bankLocation.City = req.City
	}

	if req.Postcode != "" {
		bankLocation.Postcode = req.Postcode
	}

	if req.State != "" {
		bankLocation.State = req.State
	}

	if req.Country != "" {
		bankLocation.Country = req.Country
	}

	if req.Phone != "" {
		bankLocation.Phone = req.Phone
	}

	if req.Fax != "" {
		bankLocation.Fax = req.Fax
	}

	if req.Bizhour != "" {
		bankLocation.BizHour = req.Bizhour
	}

	if req.Latitude != "" {
		bankLocation.Latitude = req.Latitude
	}

	if req.Longitude != "" {
		bankLocation.Longitude = req.Longitude
	}
	bankLocation.Modified = time.Now().Format(time.RFC3339)

	// Lakukan perubahan pada database
	if err := db.Unscoped().Omit("created_at", "updated_at", "deleted_at").Save(&bankLocation).Error; err != nil {
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
	fmt.Println(dbusn, dbaddr, dbport, dbtable)
	// Inisialisasi koneksi database menggunakan GORM
	db, err := ConnectToDatabase(dbusn, dbaddr, dbport, dbtable)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	var atmLocation MbAtmLocation
	if err := db.Unscoped().Where("id = ?", idInt).Omit("modified", "created_at", "updated_at", "deleted_at").First(&atmLocation).Error; err != nil {
		c.AbortWithError(500, err)
		return
	}

	// Update nilai-nilai yang diperlukan
	if req.Name != "" {
		atmLocation.Name = req.Name
	}

	if req.Addr != "" {
		atmLocation.Address = req.Addr
	}

	if req.Email != "" {
		atmLocation.Email = req.Email
	}

	if req.Utype != "" {
		atmLocation.Utype = req.Utype
	}

	if req.City != "" {
		atmLocation.City = req.City
	}

	if req.Postcode != "" {
		atmLocation.Postcode = req.Postcode
	}

	if req.State != "" {
		atmLocation.State = req.State
	}

	if req.Country != "" {
		atmLocation.Country = req.Country
	}

	if req.Phone != "" {
		atmLocation.Phone = req.Phone
	}

	if req.Fax != "" {
		atmLocation.Fax = req.Fax
	}

	if req.Bizhour != "" {
		atmLocation.BizHour = req.Bizhour
	}

	if req.Latitude != "" {
		atmLocation.Latitude = req.Latitude
	}

	if req.Longitude != "" {
		atmLocation.Longitude = req.Longitude
	}

	atmLocation.Modified = time.Now().Format(time.RFC3339)

	// Lakukan perubahan pada database
	if err := db.Unscoped().Omit("created_at", "updated_at", "deleted_at").Save(&atmLocation).Error; err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, gin.H{"message": "ATM data updated successfully"})
}
