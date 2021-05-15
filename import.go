package main

import (
	"database/sql"
	"encoding/csv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func importData() {
	timer := time.Now()
	logInfo("MAIN", "Importing process started")
	zapsiUsers, zapsiProducts, downloadedFromZapsi := downloadDataFromZapsi()
	csvUsers, csvProducts, downloadedFromCsvFile := downloadDataFromCsvFile()
	if downloadedFromZapsi && downloadedFromCsvFile {
		logInfo("MAIN", "Zapsi Users: "+strconv.Itoa(len(zapsiUsers)))
		logInfo("MAIN", "Zapsi Products: "+strconv.Itoa(len(zapsiProducts)))
		logInfo("MAIN", "CSV Users: "+strconv.Itoa(len(csvUsers)))
		logInfo("MAIN", "CSV Products: "+strconv.Itoa(len(csvProducts)))
		updatedUsers, createdUsers := processUsers(zapsiUsers, csvUsers)
		updatedProducts, createdProducts := processProducts(zapsiProducts, csvProducts)
		logInfo("MAIN", "Updated users: "+strconv.Itoa(updatedUsers))
		logInfo("MAIN", "Created users: "+strconv.Itoa(createdUsers))
		logInfo("MAIN", "Updated products: "+strconv.Itoa(updatedProducts))
		logInfo("MAIN", "Created products: "+strconv.Itoa(createdProducts))
	}
	logInfo("MAIN", "Importing process complete, time elapsed: "+time.Since(timer).String())
}

func processProducts(zapsiProducts map[string]product, csvProducts []csvProduct) (int, int) {
	timer := time.Now()
	logInfo("MAIN", "Processing products")
	updated := 0
	created := 0
	for _, csvProduct := range csvProducts {
		if serviceRunning {
			_, productInZapsi := zapsiProducts[csvProduct.kodProduktu]
			if productInZapsi {
				updateProductInZapsi(csvProduct)
				updated++
			} else {
				createProductInZapsi(csvProduct)
				created++
			}
		}
	}
	logInfo("MAIN", "Products processed, time elapsed: "+time.Since(timer).String())
	return updated, created
}

func createProductInZapsi(csvProduct csvProduct) {
	logInfo("MAIN", csvProduct.nazevProduktu+": Product does not exist in Zapsi, creating...")
	productGroupId := getProductGroupId(csvProduct)
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}

	cavityAsInt, err := strconv.Atoi(csvProduct.kavita)
	if err != nil {
		cavityAsInt = 0
	}
	cycleAsInt, err := strconv.Atoi(strings.ReplaceAll(csvProduct.casCyklu, ",", "."))
	var cycleAsFloat float64
	if err != nil {
		cycleAsFloat, err = strconv.ParseFloat(strings.ReplaceAll(csvProduct.casCyklu, ",", "."), 64)
		if err != nil {
			cycleAsFloat = 0.0
		}
	} else {
		cycleAsFloat = float64(cycleAsInt)
	}
	var product product
	product.Name = csvProduct.nazevProduktu
	product.Barcode = csvProduct.kodProduktu
	product.Cycle = cycleAsFloat
	product.ProductGroupID = productGroupId
	product.ProductStatusID = 1
	product.Cavity = cavityAsInt
	db.Save(&product)
}

func updateProductInZapsi(csvProduct csvProduct) {
	productGroupId := getProductGroupId(csvProduct)
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}
	cavityAsInt, err := strconv.Atoi(csvProduct.kavita)
	if err != nil {
		cavityAsInt = 0
	}
	cycleAsInt, err := strconv.Atoi(strings.ReplaceAll(csvProduct.casCyklu, ",", "."))
	var cycleAsFloat float64
	if err != nil {
		cycleAsFloat, err = strconv.ParseFloat(strings.ReplaceAll(csvProduct.casCyklu, ",", "."), 64)
		if err != nil {
			cycleAsFloat = 0.0
		}
	} else {
		cycleAsFloat = float64(cycleAsInt)
	}

	db.Model(&product{}).Where(user{Barcode: csvProduct.kodProduktu}).Updates(product{
		Name:           csvProduct.nazevProduktu,
		Cycle:          cycleAsFloat,
		Cavity:         cavityAsInt,
		ProductGroupID: productGroupId,
	})
}

func getProductGroupId(csvProduct csvProduct) int {
	prepareTimeAsFloat, err := strconv.ParseFloat(csvProduct.casPripravy, 64)
	if err != nil {
		logError("MAIN", "Problem parsing cavity for "+csvProduct.nazevProduktu+": "+csvProduct.kavita)
		prepareTimeAsFloat = 0.0
	}
	scrapPercentAsFloat, err := strconv.ParseFloat(csvProduct.zmetkovitost, 64)
	if err != nil {
		logError("MAIN", "Problem parsing cavity for "+csvProduct.nazevProduktu+": "+csvProduct.kavita)
		scrapPercentAsFloat = 0.0
	}
	cycleAsInt, err := strconv.Atoi(strings.ReplaceAll(csvProduct.casCyklu, ",", "."))
	var cycleAsFloat float64
	if err != nil {
		cycleAsFloat, err = strconv.ParseFloat(strings.ReplaceAll(csvProduct.casCyklu, ",", "."), 64)
		if err != nil {
			cycleAsFloat = 0.0
		}
	} else {
		cycleAsFloat = float64(cycleAsInt)
	}
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return 1
	}
	var existingProductGroup productGroup
	db.Where("Name like ?", csvProduct.skupinaProduktu).Find(&existingProductGroup)
	if existingProductGroup.OID > 0 {
		logInfo("MAIN", "Updating product group "+csvProduct.skupinaProduktu)
		existingProductGroup.PrepareTime = prepareTimeAsFloat
		existingProductGroup.ScrapPercent = scrapPercentAsFloat
		existingProductGroup.Cycle = cycleAsFloat
		db.Save(&existingProductGroup)
		return existingProductGroup.OID
	}
	logInfo("MAIN", "Product group "+csvProduct.skupinaProduktu+" does not exist, creating ...")
	var newProductGroup productGroup
	newProductGroup.Name = csvProduct.skupinaProduktu
	newProductGroup.PrepareTime = prepareTimeAsFloat
	newProductGroup.ScrapPercent = scrapPercentAsFloat
	existingProductGroup.Cycle = cycleAsFloat
	db.Save(&newProductGroup)
	var brandNewProductGroup productGroup
	db.Where("Name like ?", csvProduct.skupinaProduktu).Find(&brandNewProductGroup)
	if brandNewProductGroup.OID > 0 {
		return brandNewProductGroup.OID
	}
	return 1
}

func processUsers(zapsiUsers map[string]user, csvUsers []csvUser) (int, int) {
	timer := time.Now()
	logInfo("MAIN", "Processing users")
	updated := 0
	created := 0
	for _, csvUser := range csvUsers {
		if serviceRunning {
			_, userInZapsi := zapsiUsers[csvUser.osobniCislo]
			if userInZapsi {
				updateUserInZapsi(csvUser, zapsiUsers[csvUser.osobniCislo])
				updated++
			} else {
				createUserInZapsi(csvUser)
				created++
			}
		}
	}
	logInfo("MAIN", "Users processed, time elapsed: "+time.Since(timer).String())
	return updated, created
}

func updateUserInZapsi(csvUser csvUser, zapsiUser user) {
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}
	userTypeIdToInsert := getUserTypeID(csvUser)
	rfidInDecimal := getRfidInHex(csvUser.rfidKod)
	db.Model(&user{}).Where(user{Login: zapsiUser.Login}).Updates(user{
		Name:       csvUser.jmeno,
		Login:      csvUser.osobniCislo,
		Rfid:       rfidInDecimal,
		Pin:        csvUser.pin,
		Barcode:    csvUser.pin,
		UserTypeID: sql.NullInt32{Int32: int32(userTypeIdToInsert), Valid: true},
		BckRfid:    csvUser.rfidKod,
	})
}

func getRfidInHex(rfid string) string {
	trimmedRfid := rfid[0:8]
	originalData := "0123456789ABCDEF"
	replacementData := []rune("084C2A6E195D3B7F")
	rfidConverted := ""
	for _, element := range trimmedRfid {
		index := strings.IndexRune(originalData, element)
		rfidConverted += string(replacementData[index])
	}
	result, err := strconv.ParseUint(rfidConverted, 16, 32)
	if err != nil {
		logError("MAIN", "Problem converting rfid to decimal: "+err.Error())
	}
	resultAsString := strconv.FormatUint(result, 10)

	resultLength := len(resultAsString)
	difference := 10 - resultLength
	for i := 0; i < difference; i++ {
		resultAsString = "0" + resultAsString
	}
	return resultAsString
}

func getUserTypeID(csvUser csvUser) int {
	userTypeIdToInsert := 1
	if strings.Contains(csvUser.typUzivatele, "Mistr") {
		userTypeIdToInsert = 3
	} else if strings.Contains(csvUser.typUzivatele, "Seřizovač") {
		userTypeIdToInsert = 2
	} else if strings.Contains(csvUser.typUzivatele, "Údržbář") {
		userTypeIdToInsert = 4
	} else if strings.Contains(csvUser.typUzivatele, "Supervizor") {
		userTypeIdToInsert = 5
	}
	return userTypeIdToInsert
}

func createUserInZapsi(csvUser csvUser) {
	logInfo("MAIN", csvUser.jmeno+": User does not exist in Zapsi, creating...")
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}
	userTypeIdToInsert := getUserTypeID(csvUser)
	rfidInDecimal := getRfidInHex(csvUser.rfidKod)
	var user user
	user.Login = csvUser.osobniCislo
	user.Name = csvUser.jmeno
	user.Rfid = rfidInDecimal
	user.Barcode = csvUser.pin
	user.Pin = csvUser.pin
	user.UserTypeID = sql.NullInt32{Int32: int32(userTypeIdToInsert), Valid: true}
	user.UserRoleID = sql.NullInt32{Int32: 2, Valid: true}
	user.BckRfid = csvUser.rfidKod
	db.Save(&user)
}

func downloadDataFromCsvFile() ([]csvUser, []csvProduct, bool) {
	timer := time.Now()
	logInfo("MAIN", "Downloading data from Csv")
	var files []string
	root := "\\\\zapsi.ipg.local\\zapsi_ipg_data"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		logError("MAIN", err.Error())
	}
	var csvUsers []csvUser
	var csvProducts []csvProduct
	for _, file := range files {
		if strings.Contains(file, "zam") {
			logInfo("MAIN", "Processing users from "+file)
			csvFile, err := os.Open(file)
			if err != nil {
				logError("MAIN", "Cannot open file: "+err.Error())
			} else {
				r := csv.NewReader(csvFile)
				r.Comma = ';'
				skipFirstLine := true
				for {
					record, err := r.Read()
					if err == io.EOF {
						break
					}
					if err != nil {
						logError("MAIN", "Cannot read file: "+err.Error())
						break
					}
					if skipFirstLine {
						logInfo("MAIN", "Skipping first line containing header")
						skipFirstLine = false
					} else {
						var csvUser csvUser
						if len(record) > 4 {
							if len(record[1]) > 1 {
								csvUser.jmeno = record[0]
								csvUser.osobniCislo = record[1]
								csvUser.rfidKod = record[2]
								csvUser.pin = record[3]
								csvUser.typUzivatele = record[4]
								csvUsers = append(csvUsers, csvUser)
							}
						}
					}
				}
			}

		} else if strings.Contains(file, "prod") {
			logInfo("MAIN", "Processing products from "+file)
			csvFile, err := os.Open(file)
			if err != nil {
				logError("MAIN", "Cannot open file: "+err.Error())
			} else {
				r := csv.NewReader(csvFile)
				r.Comma = ';'
				skipFirstLine := true
				for {
					record, err := r.Read()
					if err == io.EOF {
						break
					}
					if err != nil {
						logError("MAIN", "Cannot read file: "+err.Error())
						break
					}
					if skipFirstLine {
						logInfo("MAIN", "Skipping first line containing header")
						skipFirstLine = false
					} else {
						var csvProduct csvProduct
						if len(record) > 6 {
							csvProduct.nazevProduktu = record[0]
							csvProduct.kodProduktu = record[1]
							csvProduct.skupinaProduktu = record[2]
							csvProduct.casCyklu = record[3]
							csvProduct.kavita = record[4]
							csvProduct.casPripravy = record[5]
							csvProduct.zmetkovitost = record[6]
							csvProducts = append(csvProducts, csvProduct)
						}
					}
				}
			}
		}
	}
	logInfo("MAIN", "Csv data downloaded, time elapsed: "+time.Since(timer).String())
	return csvUsers, csvProducts, true
}

func downloadDataFromZapsi() (map[string]user, map[string]product, bool) {
	timer := time.Now()
	logInfo("MAIN", "Downloading data from Zapsi")
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return nil, nil, false
	}
	var users []user
	var products []product
	db.Find(&users)
	db.Find(&products)
	returnProducts := make(map[string]product, len(products))
	returnUsers := make(map[string]user, len(users))
	for _, product := range products {
		returnProducts[product.Barcode] = product
	}
	for _, user := range users {
		returnUsers[user.Login] = user
	}
	logInfo("MAIN", "Zapsi data downloaded, time elapsed: "+time.Since(timer).String())
	return returnUsers, returnProducts, true
}
