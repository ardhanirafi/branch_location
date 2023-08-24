package lib

type MbLocationReq struct {
	//Uuid string `json:"6d24" `
	//Gcn  string `json:"c318" `
	// Reqtype		string`json:"reqtype" binding:"required"`
	Utype     string `json:"utype"`
	Created   string `json:"created,omitempty"`
	Modified  string `json:"modified,omitempty"`
	Name      string `json:"name" binding:"required,min=5,max=20"`
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
