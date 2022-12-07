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
	// fmt.Println(configs.USERNAME + "/" + configs.PASSWORD + "@" + configs.HOST + "/" + configs.DATABASE)
	db, err := sql.Open("goracle", configs.USERNAME+"/"+configs.PASSWORD+"@"+configs.HOST+"/"+configs.DATABASE)
	if err != nil {
		fmt.Println("... DB Setup Failed")
		fmt.Println(err)
		r.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&r)
	}
	defer db.Close()

	fmt.Println("... Opening Database Connection")
	if err = db.Ping(); err != nil {
		fmt.Printf("Error connecting to the database: %s\n", err)
		r.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&r)
	}
	fmt.Println("... Connected to Database")

	// dbQuery := fmt.Sprintf("SELECT TAGRP,PARTNO,PARTNO PARTNAME,min(LOTNO) LOTNO,min(CASEID) LINENO, MIN(CASENO) REVISENO,min(SHELVE) SHELVE, min(PALLETKEY) PALLETNO,STOCKQUANTITY QTY,count(PARTNO) CTN,min(SYSDTE) CREATEDAT,max(UPDDTE) UPDATEDAT FROM TXP_CARTONDETAILS WHERE STOCKQUANTITY > 0 AND SHELVE NOT IN ('S-XXX', 'S-PLOUT') AND TAGRP='%s' GROUP BY TAGRP,PARTNO,STOCKQUANTITY ORDER BY PARTNO", tagrp)
	dbQuery := fmt.Sprintf("SELECT TAGRP,PARTNO,min(LOTNO) LOTNO,min(CASEID) LINENO, CASE when MIN(CASENO) IS NULL THEN 0 ELSE MIN(CASENO) END REVISENO,min(SHELVE) SHELVE, min(PALLETKEY) PALLETNO,STOCKQUANTITY QTY,count(PARTNO) CTN,min(SYSDTE) CREATEDAT,max(UPDDTE) UPDATEDAT FROM TXP_CARTONDETAILS WHERE STOCKQUANTITY > 0 AND SHELVE NOT IN ('S-XXX', 'S-PLOUT') AND TAGRP='%s' GROUP BY TAGRP,PARTNO,STOCKQUANTITY ORDER BY PARTNO", tagrp)
	if c.Query("part_no") != "" {
		dbQuery = fmt.Sprintf("SELECT TAGRP,PARTNO,min(LOTNO) LOTNO,min(CASEID) LINENO, CASE when MIN(CASENO) IS NULL THEN 0 ELSE MIN(CASENO) END REVISENO,min(SHELVE) SHELVE, min(PALLETKEY) PALLETNO,STOCKQUANTITY QTY,count(PARTNO) CTN,min(SYSDTE) CREATEDAT,max(UPDDTE) UPDATEDAT FROM TXP_CARTONDETAILS WHERE STOCKQUANTITY > 0 AND SHELVE NOT IN ('S-XXX', 'S-PLOUT') AND TAGRP='%s' AND PARTNO like '%s' GROUP BY TAGRP,PARTNO,STOCKQUANTITY ORDER BY PARTNO", tagrp, "%"+c.Query("part_no")+"%")
	}
	// fmt.Println(dbQuery)
	rows, err := db.Query(dbQuery)
	if err != nil {
		fmt.Println(".....Error processing query")
		r.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&r)
	}
	defer rows.Close()

	fmt.Println("... Parsing query results")
	// var tableName string
	var data []models.Stock
	rnd := 1
	for rows.Next() {
		var r models.Stock
		rows.Scan(&r.Tagrp,
			&r.PartNo,
			&r.LotNo,
			&r.LineNo,
			&r.ReviseNo,
			&r.Shelve,
			&r.PalletNo,
			&r.Qty,
			&r.Ctn,
			&r.CreatedAt,
			&r.UpdatedAt)
		// fmt.Printf("%d %s\n", rnd, r.PartNo)
		r.ID = int64(rnd)
		if r.ReviseNo == "0" {
			r.ReviseNo = "-"
		}
		data = append(data, r)
		rnd++
	}

	// Fetch Data
	r.Message = "Show Stock All"
	r.Data = &data
	return c.Status(fiber.StatusOK).JSON(&r)
}

func FetchStockByPartNo(c *fiber.Ctx) error {
	var r models.Response
	tagrp := "C"
	if c.Query("tag") != "" {
		tagrp = c.Query("tag")
	}

	// Open Connection
	// fmt.Println(configs.USERNAME + "/" + configs.PASSWORD + "@" + configs.HOST + "/" + configs.DATABASE)
	db, err := sql.Open("goracle", configs.USERNAME+"/"+configs.PASSWORD+"@"+configs.HOST+"/"+configs.DATABASE)
	if err != nil {
		fmt.Println("... DB Setup Failed")
		fmt.Println(err)
		r.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&r)
	}
	defer db.Close()

	fmt.Println("... Opening Database Connection")
	if err = db.Ping(); err != nil {
		fmt.Printf("Error connecting to the database: %s\n", err)
		r.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&r)
	}
	fmt.Println("... Connected to Database")
	filterAll := "AND SHELVE NOT IN ('S-XXX', 'S-PLOUT')"
	if c.Query("is_out") != "1" {
		filterAll = ""
	}
	// dbQuery := fmt.Sprintf("SELECT TAGRP,PARTNO,PARTNO PARTNAME,min(LOTNO) LOTNO,min(CASEID) LINENO, MIN(CASENO) REVISENO,min(SHELVE) SHELVE, min(PALLETKEY) PALLETNO,STOCKQUANTITY QTY,count(PARTNO) CTN,min(SYSDTE) CREATEDAT,max(UPDDTE) UPDATEDAT FROM TXP_CARTONDETAILS WHERE STOCKQUANTITY > 0 AND SHELVE NOT IN ('S-XXX', 'S-PLOUT') AND TAGRP='%s' GROUP BY TAGRP,PARTNO,STOCKQUANTITY ORDER BY PARTNO", tagrp)
	dbQuery := fmt.Sprintf("SELECT TAGRP,PARTNO,LOTNO,RUNNINGNO SERIALNO,CASEID LINENO, CASE when CASENO IS NULL THEN 0 ELSE CASENO END REVISENO,SHELVE,PALLETKEY PALLETNO,STOCKQUANTITY QTY,1 CTN,SYSDTE CREATEDAT,UPDDTE UPDATEDAT FROM TXP_CARTONDETAILS WHERE STOCKQUANTITY > 0 %s AND TAGRP='%s' AND PARTNO='%s' ORDER BY PARTNO,LOTNO,RUNNINGNO,CASEID,CASENO,SYSDTE,UPDDTE", filterAll, tagrp, c.Params("part_no"))
	// fmt.Println(dbQuery)
	rows, err := db.Query(dbQuery)
	if err != nil {
		fmt.Println(".....Error processing query")
		r.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&r)
	}
	defer rows.Close()

	fmt.Println("... Parsing query results")
	// var tableName string
	var data []models.Stock
	rnd := 1
	for rows.Next() {
		var r models.Stock
		rows.Scan(&r.Tagrp,
			&r.PartNo,
			&r.LotNo,
			&r.SerialNo,
			&r.LineNo,
			&r.ReviseNo,
			&r.Shelve,
			&r.PalletNo,
			&r.Qty,
			&r.Ctn,
			&r.CreatedAt,
			&r.UpdatedAt)
		// fmt.Printf("%d %s\n", rnd, r.PartNo)
		r.ID = int64(rnd)
		if r.ReviseNo == "0" {
			r.ReviseNo = "-"
		}
		data = append(data, r)
		rnd++
	}

	// Fetch Data
	r.Message = "Show Stock All"
	r.Data = &data
	return c.Status(fiber.StatusOK).JSON(&r)
}

func FetchStockByShelve(c *fiber.Ctx) error {
	var r models.Response
	tagrp := "C"
	if c.Query("tag") != "" {
		tagrp = c.Query("tag")
	}

	// Open Connection
	// fmt.Println(configs.USERNAME + "/" + configs.PASSWORD + "@" + configs.HOST + "/" + configs.DATABASE)
	db, err := sql.Open("goracle", configs.USERNAME+"/"+configs.PASSWORD+"@"+configs.HOST+"/"+configs.DATABASE)
	if err != nil {
		fmt.Println("... DB Setup Failed")
		fmt.Println(err)
		r.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&r)
	}
	defer db.Close()

	fmt.Println("... Opening Database Connection")
	if err = db.Ping(); err != nil {
		fmt.Printf("Error connecting to the database: %s\n", err)
		r.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&r)
	}
	fmt.Println("... Connected to Database")
	// dbQuery := fmt.Sprintf("SELECT TAGRP,PARTNO,PARTNO PARTNAME,min(LOTNO) LOTNO,min(CASEID) LINENO, MIN(CASENO) REVISENO,min(SHELVE) SHELVE, min(PALLETKEY) PALLETNO,STOCKQUANTITY QTY,count(PARTNO) CTN,min(SYSDTE) CREATEDAT,max(UPDDTE) UPDATEDAT FROM TXP_CARTONDETAILS WHERE STOCKQUANTITY > 0 AND SHELVE NOT IN ('S-XXX', 'S-PLOUT') AND TAGRP='%s' GROUP BY TAGRP,PARTNO,STOCKQUANTITY ORDER BY PARTNO", tagrp)
	dbQuery := fmt.Sprintf("SELECT TAGRP,PARTNO,LOTNO,RUNNINGNO SERIALNO,CASEID LINENO, CASE when CASENO IS NULL THEN 0 ELSE CASENO END REVISENO,SHELVE,PALLETKEY PALLETNO,STOCKQUANTITY QTY,1 CTN,SYSDTE CREATEDAT,UPDDTE UPDATEDAT FROM TXP_CARTONDETAILS WHERE TAGRP='%s' AND SHELVE in ('%s') ORDER BY PARTNO,LOTNO,RUNNINGNO,CASEID,CASENO,SYSDTE,UPDDTE", tagrp, c.Params("shelve_no"))
	// fmt.Println(dbQuery)
	rows, err := db.Query(dbQuery)
	if err != nil {
		fmt.Println(".....Error processing query")
		r.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&r)
	}
	defer rows.Close()

	fmt.Println("... Parsing query results")
	// var tableName string
	var data []models.Stock
	rnd := 1
	for rows.Next() {
		var r models.Stock
		rows.Scan(&r.Tagrp,
			&r.PartNo,
			&r.LotNo,
			&r.SerialNo,
			&r.LineNo,
			&r.ReviseNo,
			&r.Shelve,
			&r.PalletNo,
			&r.Qty,
			&r.Ctn,
			&r.CreatedAt,
			&r.UpdatedAt)
		// fmt.Printf("%d %s\n", rnd, r.PartNo)
		r.ID = int64(rnd)
		if r.ReviseNo == "0" {
			r.ReviseNo = "-"
		}
		data = append(data, r)
		rnd++
	}

	// Fetch Data
	r.Message = "Show Stock All"
	r.Data = &data
	return c.Status(fiber.StatusOK).JSON(&r)
}

func FetchStockBySerialNo(c *fiber.Ctx) error {
	var r models.Response
	// Open Connection
	// fmt.Println(configs.USERNAME + "/" + configs.PASSWORD + "@" + configs.HOST + "/" + configs.DATABASE)
	db, err := sql.Open("goracle", configs.USERNAME+"/"+configs.PASSWORD+"@"+configs.HOST+"/"+configs.DATABASE)
	if err != nil {
		fmt.Println("... DB Setup Failed")
		fmt.Println(err)
		r.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&r)
	}
	defer db.Close()

	fmt.Println("... Opening Database Connection")
	if err = db.Ping(); err != nil {
		fmt.Printf("Error connecting to the database: %s\n", err)
		r.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&r)
	}
	fmt.Println("... Connected to Database")
	// dbQuery := fmt.Sprintf("SELECT TAGRP,PARTNO,PARTNO PARTNAME,min(LOTNO) LOTNO,min(CASEID) LINENO, MIN(CASENO) REVISENO,min(SHELVE) SHELVE, min(PALLETKEY) PALLETNO,STOCKQUANTITY QTY,count(PARTNO) CTN,min(SYSDTE) CREATEDAT,max(UPDDTE) UPDATEDAT FROM TXP_CARTONDETAILS WHERE STOCKQUANTITY > 0 AND SHELVE NOT IN ('S-XXX', 'S-PLOUT') AND TAGRP='%s' GROUP BY TAGRP,PARTNO,STOCKQUANTITY ORDER BY PARTNO", tagrp)
	dbQuery := fmt.Sprintf("SELECT TAGRP,PARTNO,LOTNO,RUNNINGNO SERIALNO,CASEID LINENO, CASE when CASENO IS NULL THEN 0 ELSE CASENO END REVISENO,SHELVE,PALLETKEY PALLETNO,STOCKQUANTITY QTY,1 CTN,SYSDTE CREATEDAT,UPDDTE UPDATEDAT FROM TXP_CARTONDETAILS WHERE RUNNINGNO LIKE '%s' ORDER BY PARTNO,LOTNO,RUNNINGNO,CASEID,CASENO,SYSDTE,UPDDTE", "%"+c.Params("serial_no")+"%")
	// fmt.Println(dbQuery)
	rows, err := db.Query(dbQuery)
	if err != nil {
		fmt.Println(".....Error processing query")
		r.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&r)
	}
	defer rows.Close()

	fmt.Println("... Parsing query results")
	// var tableName string
	var data []models.Stock
	rnd := 1
	for rows.Next() {
		var r models.Stock
		rows.Scan(&r.Tagrp,
			&r.PartNo,
			&r.LotNo,
			&r.SerialNo,
			&r.LineNo,
			&r.ReviseNo,
			&r.Shelve,
			&r.PalletNo,
			&r.Qty,
			&r.Ctn,
			&r.CreatedAt,
			&r.UpdatedAt)
		// fmt.Printf("%d %s\n", rnd, r.PartNo)
		r.ID = int64(rnd)
		if r.ReviseNo == "0" {
			r.ReviseNo = "-"
		}
		data = append(data, r)
		rnd++
	}

	// Fetch Data
	r.Message = "Show Stock All"
	r.Data = &data
	return c.Status(fiber.StatusOK).JSON(&r)
}

func UpdateStockBySerialNo(c *fiber.Ctx) error {
	var r models.Response
	var frm models.FrmUpdateStock
	if err := c.BodyParser(&frm); err != nil {
		r.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&r)
	}

	db, err := sql.Open("goracle", configs.USERNAME+"/"+configs.PASSWORD+"@"+configs.HOST+"/"+configs.DATABASE)
	if err != nil {
		fmt.Println("... DB Setup Failed")
		r.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&r)
	}
	defer db.Close()

	strExecute := fmt.Sprintf("UPDATE TXP_CARTONDETAILS SET STOCKQUANTITY=0,SHELVE='S-PLOUT',SIDTE=sysdate,SINO='TIMVOUT',SIID='%s' WHERE RUNNINGNO='%s'", frm.EmpId, frm.SerialNo)
	if frm.Ctn > 0 {
		strExecute = fmt.Sprintf("UPDATE TXP_CARTONDETAILS SET STOCKQUANTITY=RECEIVINGQUANTITY,SHELVE='SNON',SIDTE=null,SINO=null,SIID='%s' WHERE RUNNINGNO='%s'", frm.EmpId, frm.SerialNo)
	}
	if _, err := db.Exec(strExecute); err != nil {
		r.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&r)
	}

	if frm.Ctn > 0 {
		rows, err := db.Query(fmt.Sprintf("SELECT RECEIVINGQUANTITY FROM TXP_CARTONDETAILS WHERE RUNNINGNO='%s'", frm.SerialNo))
		if err != nil {
			fmt.Println(".....Error processing query")
			r.Message = err.Error()
			return c.Status(fiber.StatusBadRequest).JSON(&r)
		}
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&frm.Qty)
			frm.Ctn = 1
		}
	}

	fmt.Println(strExecute)
	r.Data = &frm
	return c.Status(fiber.StatusOK).JSON(&r)
}
