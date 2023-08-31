package lib

import (
	"gorm.io/gorm"
)

type MbLocationReq struct {
	//Uuid string `json:"6d24" `
	//Gcn  string `json:"c318" `
	// Reqtype		string`json:"reqtype" binding:"required"`
	Utype     string `json:"utype"`
	Created   string `json:"created,omitempty"`
	Modified  string `json:"modified,omitempty"`
	Name      string `json:"name"`
	Addr      string `json:"address"`
	City      string `json:"city"`
	Postcode  string `json:"postcode,omitempty"`
	State     string `json:"state,omitempty"`
	Country   string `json:"country,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Fax       string `json:"fax,omitempty"`
	Email     string `json:"email,omitempty"`
	Bizhour   string `json:"bizhour,omitempty"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type MbLocationResp struct {
	MessageCode string                 `json:"messagecode"`
	Message     string                 `json:"message"`
	Success     bool                   `json:"success"`
	Results     []MbLocationRespResult `json:"results"`
}

type MbLocationRespResult struct {
	Unittype  string `json:"utype"`
	Name      string `json:"name"`
	Addr      string `json:"address"`
	City      string `json:"city"`
	Postcode  string `json:"postcode,omitempty"`
	State     string `json:"state,omitempty"`
	Country   string `json:"country,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Fax       string `json:"fax,omitempty"`
	Email     string `json:"email,omitempty"`
	Bizhour   string `json:"bizhour,omitempty"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type mb_bank_location struct {
	ID        uint        `gorm:"primaryKey"`
	Utype     string      `json:"utype"`
	Name      string      `json:"name"`
	Address   string      `json:"address"`
	City      string      `json:"city"`
	Postcode  string      `json:"postcode,omitempty"`
	State     string      `json:"state,omitempty"`
	Country   string      `json:"country,omitempty"`
	Phone     string      `json:"phone,omitempty"`
	Fax       string      `json:"fax,omitempty"`
	Email     string      `json:"email,omitempty"`
	Bizhour   string      `json:"biz_hour,omitempty"`
	Latitude  string      `json:"latitude"`
	Longitude string      `json:"longitude"`
	CreatedAt interface{} `json:"created"`
}

type mb_atm_location struct {
	ID        uint        `gorm:"primaryKey"`
	Utype     string      `json:"utype"`
	Name      string      `json:"name"`
	Address   string      `json:"address"`
	City      string      `json:"city"`
	Postcode  string      `json:"postcode,omitempty"`
	State     string      `json:"state,omitempty"`
	Country   string      `json:"country,omitempty"`
	Phone     string      `json:"phone,omitempty"`
	Fax       string      `json:"fax,omitempty"`
	Email     string      `json:"email,omitempty"`
	Bizhour   string      `json:"biz_hour,omitempty"`
	Latitude  string      `json:"latitude"`
	Longitude string      `json:"longitude"`
	CreatedAt interface{} `json:"created"`
}

type MbBankLocationsDelete struct {
	ID        uint   `gorm:"primaryKey"`
	Utype     string `json:"utype"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Postcode  string `json:"postcode,omitempty"`
	State     string `json:"state,omitempty"`
	Country   string `json:"country,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Fax       string `json:"fax,omitempty"`
	Email     string `json:"email,omitempty"`
	Bizhour   string `json:"biz_hour,omitempty"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type MbAtmLocationsDelete struct {
	ID        uint   `gorm:"primaryKey"`
	Utype     string `json:"utype"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Postcode  string `json:"postcode,omitempty"`
	State     string `json:"state,omitempty"`
	Country   string `json:"country,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Fax       string `json:"fax,omitempty"`
	Email     string `json:"email,omitempty"`
	Bizhour   string `json:"biz_hour,omitempty"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type MbBankLocation struct {
	gorm.Model
	Name      string `json:"name"`
	Email     string `json:"email"`
	Utype     string `json:"utype"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Postcode  string `json:"postcode"`
	State     string `json:"state"`
	Country   string `json:"country"`
	Phone     string `json:"phone"`
	Fax       string `json:"fax"`
	BizHour   string `json:"bizhour"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Modified  string `json:"modified"`
}

type MbAtmLocation struct {
	gorm.Model
	Name      string `json:"name"`
	Email     string `json:"email"`
	Utype     string `json:"utype"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Postcode  string `json:"postcode"`
	State     string `json:"state"`
	Country   string `json:"country"`
	Phone     string `json:"phone"`
	Fax       string `json:"fax"`
	BizHour   string `json:"bizhour"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Modified  string `json:"modified"`
}

func (MbBankLocation) TableName() string {
	return banklocation
}

func (MbAtmLocation) TableName() string {
	return atmlocation
}
