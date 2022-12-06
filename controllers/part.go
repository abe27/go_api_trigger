package controllers

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/abe27/oracle/api/configs"
	"github.com/abe27/oracle/api/models"
    _ "gopkg.in/goracle.v2"
)

func PostData(RowID *string, obj *models.CartonDetail) {
	method := "POST"
	pData := fmt.Sprintf("row_id=%s&whs=%s&part_no=%s&lot_no=%s&serial_no=%s&die_no=%s&rev_no=%d&qty=%d&shelve=%s&ip_address=%s&emp_id=%s&ref_no=%s&receive_no=%s&description=%s", *RowID, obj.Tagrp, obj.PartNo, obj.LotNo, obj.SerialNo, obj.LineNo, obj.ReviseNo, obj.Qty, obj.Shelve, obj.IpAddress, obj.SiID, obj.PalletNo, obj.InvoiceNo, obj.SiNo)
	payload := strings.NewReader(pData)
	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf("%s/carton/history", configs.REST_URL), payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	fmt.Printf("%s SERIALNO: %s\n", *RowID, obj.SerialNo)
}

func FetchData(frm *models.CartonForm) {
	currentTime := time.Now()
	fmt.Println("Starting at : ", currentTime.Format("03:04:05:06 PM"))

	fmt.Println("... Setting up Database Connection")
	db, err := sql.Open("goracle", configs.USERNAME+"/"+configs.PASSWORD+"@"+configs.HOST+"/"+configs.DATABASE)
	if err != nil {
		fmt.Println("... DB Setup Failed")
		fmt.Println(err)
		return
	}
	defer db.Close()

	fmt.Println("... Opening Database Connection")
	if err = db.Ping(); err != nil {
		fmt.Printf("Error connecting to the database: %s\n", err)
		return
	}
	fmt.Println("... Connected to Database")

	dbQuery := fmt.Sprintf("SELECT TAGRP,PARTNO,LOTNO,RUNNINGNO,CASE WHEN CASEID IS NULL THEN '-' ELSE CASEID END CASEID,CASE WHEN CASENO IS NULL THEN 0 ELSE CASENO END CASENO,STOCKQUANTITY,CASE WHEN SHELVE IS NULL THEN '-' ELSE SHELVE END SHELVE,'%s' ip_address,CASE WHEN SIID IS NULL THEN '-' ELSE SIID END SIID,CASE WHEN PALLETKEY IS NULL THEN '-' ELSE PALLETKEY END PALLETKEY,INVOICENO,CASE WHEN SINO IS NULL THEN '-' ELSE SINO END SINO FROM TXP_CARTONDETAILS WHERE RUNNINGNO='%s'", frm.IpAddress, frm.SerialNo)
	rows, err := db.Query(dbQuery)
	if err != nil {
		fmt.Println(".....Error processing query")
		fmt.Println(err)
		return
	}
	defer rows.Close()

	fmt.Println("... Parsing query results")
	// var tableName string
	var carton models.CartonDetail
	for rows.Next() {
		rows.Scan(&carton.Tagrp, &carton.PartNo, &carton.LotNo, &carton.SerialNo, &carton.LineNo, &carton.ReviseNo, &carton.Qty, &carton.Shelve, &carton.IpAddress, &carton.SiID, &carton.PalletNo, &carton.InvoiceNo, &carton.SiNo)
	}

	fmt.Println("Post RowID: ", frm.RowID)
	PostData(&frm.RowID, &carton)
	fmt.Println("... Closing connection")
	fmt.Printf("------------%d-------------------", &frm.RowID)
	finishTime := time.Now()
	fmt.Println("Finished at ", finishTime.Format("03:04:05:06 PM"))
}

func FetchDataBySerialNo(serialNo string) bool {
	currentTime := time.Now()
	fmt.Println("Starting at : ", currentTime.Format("03:04:05:06 PM"))

	fmt.Println("... Setting up Database Connection")
	db, err := sql.Open("goracle", configs.USERNAME+"/"+configs.PASSWORD+"@"+configs.HOST+"/"+configs.DATABASE)
	if err != nil {
		fmt.Println("... DB Setup Failed")
		fmt.Println(err)
		return false
	}
	defer db.Close()

	fmt.Println("... Opening Database Connection")
	if err = db.Ping(); err != nil {
		fmt.Printf("Error connecting to the database: %s\n", err)
		return false
	}
	fmt.Println("... Connected to Database")

	dbQuery := fmt.Sprintf("SELECT RUNNINGNO FROM TXP_CARTONDETAILS WHERE RUNNINGNO='%s'", serialNo)
	rows, err := db.Query(dbQuery)
	if err != nil {
		fmt.Println(".....Error processing query")
		fmt.Println(err)
		return false
	}
	defer rows.Close()

	fmt.Println("... Parsing query results")
	var serial_no string
	for rows.Next() {
		rows.Scan(&serial_no)
	}

	db.Exec(fmt.Sprintf("UPDATE TXP_CARTONDETAILS SET IS_CHECK=1 WHERE RUNNINGNO='%s'", serial_no))

	fmt.Println("... Closing connection")
	fmt.Printf("------------%s-------------------", serial_no)
	finishTime := time.Now()
	fmt.Println("Finished at ", finishTime.Format("03:04:05:06 PM"))
	return serial_no != ""
}
