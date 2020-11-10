package main

import "database/sql"

type user struct {
	OID        int           `gorm:"primary_key;column:OID"`
	Login      string        `gorm:"column:Login"`
	Password   string        `gorm:"column:Password"`
	Name       string        `gorm:"column:Name"`
	FirstName  string        `gorm:"column:FirstName"`
	Rfid       string        `gorm:"column:Rfid"`
	Barcode    string        `gorm:"column:Barcode"`
	Pin        string        `gorm:"column:Pin"`
	Function   string        `gorm:"column:Function"`
	UserTypeID sql.NullInt32 `gorm:"column:UserTypeID"`
	UserRoleID sql.NullInt32 `gorm:"column:UserRoleID"`
	Email      string        `gorm:"column:Email"`
	Phone      string        `gorm:"column:Phone"`
	BckRfid    string        `gorm:"column:BckRfid"`
}

func (user) TableName() string {
	return "user"
}

type product struct {
	OID             int     `gorm:"primary_key;column:OID"`
	Name            string  `gorm:"column:Name"`
	Barcode         string  `gorm:"column:Barcode"`
	Cycle           float64 `gorm:"column:Cycle"`
	IdleFromTime    int     `gorm:"column:IdleFromTime"`
	ProductStatusID int     `gorm:"column:ProductStatusID"`
	Deleted         int     `gorm:"column:Deleted"`
	ProductGroupID  int     `gorm:"column:ProductGroupID"`
	Cavity          int     `gorm:"column:Cavity"`
}

func (product) TableName() string {
	return "product"
}

type productGroup struct {
	OID          int     `gorm:"primary_key;column:OID"`
	Name         string  `gorm:"column:Name"`
	Cycle        float64 `gorm:"column:Cycle"`
	PrepareTime  float64 `gorm:"column:PrepareTime"`
	ScrapPercent float64 `gorm:"column:ScrapPercent"`
}

func (productGroup) TableName() string {
	return "product_group"
}

type csvUser struct {
	jmeno        string
	osobniCislo  string
	rfidKod      string
	pin          string
	typUzivatele string
}

type csvProduct struct {
	nazevProduktu   string
	kodProduktu     string
	skupinaProduktu string
	casCyklu        string
	kavita          string
	casPripravy     string
	zmetkovitost    string
}
