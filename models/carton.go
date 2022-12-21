package models

import "time"

type CartonDetail struct {
	// RowID     *string // ROWID        |AABAE2AAFAAADmiAAA|
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

type CartonForm struct {
	RowID     string `form:"row_id" json:"row_id"`
	SerialNo  string `form:"serial_no" json:"serial_no"`
	IpAddress string `form:"ip_address" json:"ip_address"`
}

type Stock struct {
	ID        int64     `form:"id" json:"id"`
	Tagrp     string    `form:"tagrp" json:"tagrp"`
	PartNo    string    `form:"part_no" json:"part_no"`
	PartName  string    `form:"part_name" json:"part_name"`
	SerialNo  string    `json:"serial_no"`
	LotNo     string    `json:"lot_no"`
	LineNo    string    `json:"line_no"`
	ReviseNo  string    `json:"revise_no"`
	Shelve    string    `json:"shelve"`
	PalletNo  string    `json:"pallet_no"`
	Qty       float64   `json:"qty"`
	Ctn       float64   `json:"ctn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FrmUpdateStock struct {
	EmpId    string `form:"emp_id" json:"emp_id" binding:"required"`
	SerialNo string `form:"serial_no" json:"serial_no" binding:"required"`
	Shelve   string `form:"shelve" json:"shelve" binding:"required"`
	Qty      int64  `form:"qty" json:"qty"`
	Ctn      int64  `form:"ctn" json:"ctn" binding:"required"`
}

type StockCheck struct {
	Tagrp      string    `json:"tagrp"`    // TAGRP
	PartNo     string    `json:"partno"`   // PARTNO
	PartName   string    `json:"partname"` // PARTNAME
	Total      int64     `json:"total"`    // TOTAL
	Checked    int64     `json:"checked"`  // CHECKED
	NotCheck   int64     `json:"notcheck"` // NOTCHECK
	LastUpdate time.Time `json:"last_update"`
}

type StockCheckDetail struct {
	Tagrp      string    `json:"tagrp"`  //TAGRP
	PartNo     string    `json:"partno"` //PARTNO
	PartName   string    `json:"partname"`
	LotNo      string    `json:"lot_no"`    // LOTNO
	SerialNo   string    `json:"serial_no"` // RUNNINGNO
	Qty        int64     `json:"qty"`       // STOCKQUANTITY
	Shelve     string    `json:"shelve"`    // SHELVE
	PalletNo   string    `json:"pallet_no"` // PALLETKEY
	LineNo     string    `json:"line_no"`
	ReviseNo   string    `json:"revise_no"`
	Checked    int64     `json:"checked_flg"` // STKTAKECHKFLG
	LastUpdate time.Time `json:"last_update"` // UPDDTE
}
