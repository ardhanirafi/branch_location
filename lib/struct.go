package lib

type MbLocationReq struct {
	Uuid				string	`json:"6d24" binding:"required"`
	Gcn					string	`json:"c318" bindind:"required"`
	// Reqtype				string	`json:"reqtype" binding:"required"`
}

type MbLocationResp struct {
	MessageCode			string	`json:"messagecode"`
	Message				string	`json:"message"`
	Success				bool	`json:"success"`
	Results				[]MbLocationRespResult	`json:"results"`
}

type MbLocationRespResult struct {
	Unittype		string	`json:"utype"`
	Name			string	`json:"name"`
	Addr			string	`json:"address"`
	City			string	`json:"city"`
	Postcode		string	`json:"postcode,omitempty"`
	State			string	`json:"state,omitempty"`
	Country			string	`json:"country,omitempty"`
	Phone			string	`json:"phone,omitempty"`
	Fax				string	`json:"fax,omitempty"`
	Email			string	`json:"email,omitempty"`
	Bizhour			string	`json:"bizhour,omitempty"`
	Latitude		string	`json:"latitude"`
	Longitude		string	`json:"longitude"`
}

