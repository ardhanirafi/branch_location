package main

import (
	// L "./lib"
	"github.com/gin-gonic/gin"
	L "github.com/mncbank/api/lib"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	v := r.Group("/mbank/v1/location/")
	{
		//api buat read
		v.GET("inquiry/", L.MbLocInq)
		v.DELETE("inquiry/delete-atm/:id", L.MbDeleteATM)
		v.DELETE("inquiry/delete-bank/:id", L.MbDeleteBank)
		v.POST("inquiry/bank/add", L.MbBankAdd)
		v.POST("inquiry/atm/add", L.MbAtmAdd)
		v.PUT("inquiry/bank/update/:id", L.MbBankUpdate)
		v.PUT("inquiry/atm/update/:id", L.MbAtmUpdate)
		//api buat add, delete, update
	}
	r.Run(":8070") // listen and serve on 0.0.0.0:8070
}

/*
	GET --> untuk khusus get data, jadi gak pake request
	sisanya pake body
*/
