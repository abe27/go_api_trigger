package controllers

import (
	"database/sql"
	"fmt"

	"github.com/abe27/oracle/api/configs"
	"github.com/abe27/oracle/api/models"
	"github.com/gofiber/fiber/v2"
	_ "gopkg.in/goracle.v2"
)

func FetchAllStock(c *fiber.Ctx) error {
	var r models.Response
	tagrp := "C"
	if c.Query("tag") != "" {
		tagrp = c.Query("tag")
	}

	// Open Connection
	fmt.Println(configs.USERNAME + "/" + configs.PASSWORD + "@" + configs.HOST + "/" + configs.DATABASE)
	db, err := sql.Open("goracle", configs.USERNAME+"/"+configs.PASSWORD+"@"+configs.HOST+"/"+configs.DATABASE)
	if err != nil {
		fmt.Println("... DB Setup Failed")
		fmt.Println(err)
		panic(err)
	}
	defer db.Close()

	fmt.Println("... Opening Database Connection")
	if err = db.Ping(); err != nil {
		fmt.Printf("Error connecting to the database: %s\n", err)
		panic(err)
	}
	fmt.Println("... Connected to Database")

	dbQuery := fmt.Sprintf("SELECT TAGRP,PARTNO,STOCKQUANTITY,count(PARTNO) ctn FROM TXP_CARTONDETAILS WHERE STOCKQUANTITY > 0 AND TAGRP='%s' GROUP BY TAGRP,PARTNO,STOCKQUANTITY ORDER BY PARTNO", tagrp)
	rows, err := db.Query(dbQuery)
	if err != nil {
		fmt.Println(".....Error processing query")
		panic(err)
	}
	defer rows.Close()

	fmt.Println("... Parsing query results")
	// var tableName string
	var data []models.Stock
	rnd := 1
	for rows.Next() {
		var r models.Stock
		rows.Scan(&r.Tagrp, &r.PartNo, &r.Qty, &r.Ctn)
		fmt.Printf("%d %s\n", rnd, r.PartNo)
		r.ID = int64(rnd)
		data = append(data, r)
		rnd++
	}

	// Fetch Data
	r.Message = "Show Stock All"
	r.Data = &data
	return c.Status(fiber.StatusOK).JSON(&r)
}
