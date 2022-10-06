package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "gopkg.in/goracle.v2"
)

type CartonDetail struct {
	RowID     string // ROWID        |AABAE2AAFAAADmiAAA|
	Tagrp     string // TAGRP        |C                 |
	PartNo    string // PARTNO       |7116-5046-02      |
	LotNo     string // LOTNO        |20818026          |
	SerialNo  string // RUNNINGNO    |S2F0501591        |
	LineNo    string // CASEID       |T-1               |
	ReviseNo  int64  // CASENO       |0                 |
	Qty       int64  // STOCKQUANTITY|0                 |
	Shelve    string // SHELVE       |S-PLOUT           |
	IpAddress string // IP_ADDRESS   |192.168.104.120   |
	SiID      string // SIID         |00307             |
	PalletNo  string // PALLETKEY    |-                 |
	InvoiceNo string // INVOICENO    |SI22081902        |
	SiNo      string // SINO         |TIMVOUT           |
}

func PostData(obj *CartonDetail) {

	url := "http://127.0.0.1:4040/api/v1/carton/history"
	method := "POST"
	pData := fmt.Sprintf("row_id=%s&whs=%s&part_no=%s&lot_no=%s&serial_no=%s&die_no=%s&rev_no=%d&qty=%d&shelve=%s&ip_address=%s&emp_id=%s&ref_no=%s&receive_no=%s&description=%s", obj.RowID, obj.Tagrp, obj.PartNo, obj.LotNo, obj.SerialNo, obj.LineNo, obj.ReviseNo, obj.Qty, obj.Shelve, obj.IpAddress, obj.SiID, obj.PalletNo, obj.InvoiceNo, obj.SiNo)
	fmt.Println(pData)
	payload := strings.NewReader(pData)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
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
}

func FetchData(serial_no, ip_address string) {
	username := "expsys"
	password := "expsys"
	host := "192.168.101.215"
	database := "RMW"

	currentTime := time.Now()
	fmt.Println("Starting at : ", currentTime.Format("03:04:05:06 PM"))

	fmt.Println("... Setting up Database Connection")
	db, err := sql.Open("goracle", username+"/"+password+"@"+host+"/"+database)
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

	dbQuery := fmt.Sprintf("SELECT rowid,TAGRP,PARTNO,LOTNO,RUNNINGNO,CASE WHEN CASEID IS NULL THEN '-' ELSE CASEID END CASEID,CASE WHEN CASENO IS NULL THEN 0 ELSE CASENO END CASENO,STOCKQUANTITY,CASE WHEN SHELVE IS NULL THEN '-' ELSE SHELVE END SHELVE,'%s' ip_address,CASE WHEN SIID IS NULL THEN '-' ELSE SIID END SIID,CASE WHEN PALLETKEY IS NULL THEN '-' ELSE PALLETKEY END PALLETKEY,INVOICENO,CASE WHEN SINO IS NULL THEN '-' ELSE SINO END SINO FROM TXP_CARTONDETAILS WHERE RUNNINGNO='%s'", ip_address, serial_no)
	fmt.Println(dbQuery)
	rows, err := db.Query(dbQuery)
	if err != nil {
		fmt.Println(".....Error processing query")
		fmt.Println(err)
		return
	}
	defer rows.Close()

	fmt.Println("... Parsing query results")
	// var tableName string
	var carton CartonDetail
	for rows.Next() {
		rows.Scan(&carton.RowID, &carton.Tagrp, &carton.PartNo, &carton.LotNo, &carton.SerialNo, &carton.LineNo, &carton.ReviseNo, &carton.Qty, &carton.Shelve, &carton.IpAddress, &carton.SiID, &carton.PalletNo, &carton.InvoiceNo, &carton.SiNo)
	}
	go PostData(&carton)
	fmt.Println("... Closing connection")
	fmt.Printf("------------%d-------------------", &carton.RowID)
	finishTime := time.Now()
	fmt.Println("Finished at ", finishTime.Format("03:04:05:06 PM"))
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusCreated).JSON("Hello, world!")
	})

	app.Get("/carton", func(c *fiber.Ctx) error {
		serial_no := c.Query("serial_no")
		ip_address := c.Query("ip_address")
		if serial_no != "" {
			go FetchData(serial_no, ip_address)
			return c.Status(fiber.StatusCreated).JSON(&serial_no)
		}
		return c.Status(fiber.StatusInternalServerError).JSON("error")
	})

	app.Listen(":4000")
}