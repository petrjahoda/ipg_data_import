# IPG Data Import Service

## Description
Go service that downloads user and product data from CSV and updates/creates users and products in Zapsi

* Periodocity of download: 10 minutes
* CSV directory: \\server-zapsi\zapsi_ipg_data
   * user file has to have 'zam' in its name
   * product file has to have 'prod' in its name
   
### User mapping:

|CSV Name|ZapsiName user table|
|------------------|------------------|
|-|FirstName|
|Jméno|Name|
|osobní číslo|Login|
|bckRFID Kód|Rfid changed to decimal (see below)|
|PIN|Barcode|
|PIN|Pin|
|Typ uživatele|UserTypeID|
|-|Email|
|-|Phone|
|1|UserRoleID|
|bckRFID Kód|bckRfid|

### Product mapping:
    
|CSV Name|ZapsiName, product table|
|------------------|------------------|
|název produktu|Name|
|kód produktu|Barcode|
|čas cyklu|Cycle|
|-|IdleFromTime|
|1|ProductStatusID|
|-|Deleted|
|skupina produktů|ProductGroupID|
|kavita|Cavity|
|čas přípravy|-|
|zmetkovitost|-|


© 2020 Petr Jahoda
