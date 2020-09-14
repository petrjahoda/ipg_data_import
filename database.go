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
}

func (user) TableName() string {
	return "user"
}

type product struct {
}

func (product) TableName() string {
	return "product"
}

type csvUser struct {
}

type csvProduct struct {
}
