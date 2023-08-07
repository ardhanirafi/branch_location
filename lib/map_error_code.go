package lib

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func GetErrorDesc(errorCD, probOwner string) string {

    wording_i := "Error pada system"
    wording_e := "System error" 
    var finalWording string

    db, err := sql.Open("mysql", "varuser:mnc123@tcp("+dbaddr+":"+dbport+")/mncmbank")
    if err != nil {
        FancyHandleError(err)
        finalWording = wording_i + "\n" + wording_e
        return finalWording
    }
    defer db.Close()

    queryData, err := db.Query("select wording_i, wording_e from mb_error_details where error_code=? and problem_owner=? ",errorCD, probOwner)
    defer queryData.Close()

    if err != nil {
        FancyHandleError(err)
        finalWording = wording_i + "\n" + wording_e
        return finalWording
    }

    if queryData.Next() {
        _ = queryData.Scan(&wording_i, &wording_e)
    }

    finalWording = wording_i + "\n" + wording_e
    return finalWording
}