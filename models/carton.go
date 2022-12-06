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
