package main

import (
	"encoding/csv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

func importData() {
	timer := time.Now()
	logInfo("MAIN", "Importing process started")
	zapsiUsers, zapsiProducts, downloadedFromZapsi := downloadDataFromZapsi()
	csvUsers, csvProducts, downloadedFromCsvFile := downloadDataFromCsvFile()
	if downloadedFromZapsi && downloadedFromCsvFile {
		sort.Slice(zapsiUsers, func(i, j int) bool { return zapsiUsers[i].Login <= zapsiUsers[j].Login })
		logInfo("MAIN", "Zapsi Users: "+strconv.Itoa(len(zapsiUsers)))
		updateUsers(zapsiUsers, csvUsers)
		updateProducts(zapsiProducts, csvProducts)
	}
	logInfo("MAIN", "Importing process complete, time elapsed: "+time.Since(timer).String())
}

func updateProducts(zapsiProducts []product, csvProducts []csvProduct) {
	timer := time.Now()
	logInfo("MAIN", "Updating products")
	for _, csvProduct := range csvProducts {
		if serviceRunning {
			index, userInZapsi := BinarySearchProduct(zapsiProducts, csvProduct)
			if userInZapsi {
				UpdateProductInZapsi(csvProduct, zapsiProducts[index])
			} else {
				CreateProductInZapsi(csvProduct)
			}
		}
	}
	logInfo("MAIN", "Products updated, time elapsed: "+time.Since(timer).String())
}

func BinarySearchProduct(zapsiProducts []product, csvProduct csvProduct) (interface{}, interface{}) {
	index := sort.Search(len(zapsiProducts), func(i int) bool { return zapsiProducts[i].Login >= csvProduct.Cislo })
	userInZapsi := index < len(zapsiProducts) && zapsiProducts[index].Login == csvProduct.Cislo
	return index, userInZapsi
}

func updateUsers(zapsiUsers []user, csvUsers []csvUser) {
	timer := time.Now()
	logInfo("MAIN", "Updating users")
	for _, csvUser := range csvUsers {
		if serviceRunning {
			index, userInZapsi := BinarySearchUser(zapsiUsers, csvUser)
			if userInZapsi {
				UpdateUserInZapsi(csvUser, zapsiUsers[index])
			} else {
				CreateUserInZapsi(csvUser)
			}
		}
	}
	logInfo("MAIN", "Users updated, time elapsed: "+time.Since(timer).String())
}

//func UpdateUserInZapsi(heliosUser hvw_Zamestnanci, zapsiUser user) {
//	logInfo("MAIN", heliosUser.Jmeno+" "+heliosUser.Prijmeni+": User exists in Zapsi, updating to "+heliosUser.EVOLoginZam)
//	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
//	if err != nil {
//		logError("MAIN", "Problem opening database: "+err.Error())
//		return
//	}
//	sqlDB, err := db.DB()
//	defer sqlDB.Close()
//
//	var userTypeIdToInsert int
//	updateUserType := false
//	if heliosUser.Serizovac {
//		userTypeIdToInsert = 2
//		updateUserType = true
//	}
//	db.Model(&user{}).Where(user{Login: zapsiUser.Login}).Updates(user{
//		Name:       heliosUser.Prijmeni,
//		FirstName:  heliosUser.Jmeno,
//		Rfid:       heliosUser.EVOLoginZam,
//		Barcode:    heliosUser.EVOLoginZam,
//		Pin:        heliosUser.EVOLoginZam,
//		UserTypeID: sql.NullInt32{Int32: int32(userTypeIdToInsert), Valid: updateUserType},
//	})
//}
//
//func CreateUserInZapsi(heliosUser hvw_Zamestnanci) {
//	logInfo("MAIN", heliosUser.Jmeno+" "+heliosUser.Prijmeni+": User does not exist in Zapsi, creating...")
//	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
//	if err != nil {
//		logError("MAIN", "Problem opening database: "+err.Error())
//		return
//	}
//	sqlDB, err := db.DB()
//	defer sqlDB.Close()
//	var user user
//	user.Login = heliosUser.Cislo
//	user.FirstName = heliosUser.Jmeno
//	user.Name = heliosUser.Prijmeni
//	user.Rfid = heliosUser.EVOLoginZam
//	user.Barcode = heliosUser.EVOLoginZam
//	user.Pin = heliosUser.EVOLoginZam
//	if heliosUser.Serizovac {
//		user.UserTypeID = sql.NullInt32{Int32: 2, Valid: true}
//	} else {
//		user.UserTypeID = sql.NullInt32{Int32: 1, Valid: true}
//	}
//	user.UserRoleID = sql.NullInt32{Int32: 2, Valid: true}
//	db.Save(&user)
//	return
//}
//
func BinarySearchUser(zapsiUsers []user, csvUser csvUser) (int, bool) {
	index := sort.Search(len(zapsiUsers), func(i int) bool { return zapsiUsers[i].Login >= heliosUser.Cislo })
	userInZapsi := index < len(zapsiUsers) && zapsiUsers[index].Login == heliosUser.Cislo
	return index, userInZapsi
}

func downloadDataFromCsvFile() ([]csvUser, []csvProduct, bool) {
	timer := time.Now()
	logInfo("MAIN", "Downloading data from Csv")
	var files []string
	root := "E:/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		println(err.Error())
	}
	for _, file := range files {
		csvFile, err := os.Open(file)
		if err != nil {
			println("Cannot open file: " + err.Error())
		} else {
			r := csv.NewReader(csvFile)
			r.Comma = ';'
			for {
				record, err := r.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					println("Cannot read file: " + err.Error())
					break
				}
				for _, data := range record {
					print(data + " ")
				}
				println("")
			}
		}
	}
	logInfo("MAIN", "Csv data downloaded, time elapsed: "+time.Since(timer).String())
	return []csvUser{}, []csvProduct{}, true
}

func downloadDataFromZapsi() ([]user, []product, bool) {
	timer := time.Now()
	logInfo("MAIN", "Downloading data from Zapsi")
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return []user{}, []product{}, false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var users []user
	db.Find(&users)
	var products []product
	db.Find(&products)
	logInfo("MAIN", "Zapsi data downloaded, time elapsed: "+time.Since(timer).String())
	return users, products, true
}
