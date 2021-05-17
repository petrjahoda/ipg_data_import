# IPG Data Import Service Changelog

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

Please note, that this project, while following numbering syntax, it DOES NOT
adhere to [Semantic Versioning](http://semver.org/spec/v2.0.0.html) rules.

## Types of changes

* ```Added``` for new features.
* ```Changed``` for changes in existing functionality.
* ```Deprecated``` for soon-to-be removed features.
* ```Removed``` for now removed features.
* ```Fixed``` for any bug fixes.
* ```Security``` in case of vulnerabilities.

## [2021.2.2.17] - 2021-5-17

### Removed
- upx usage from windows binary

### Changed
- directory for csv file on local disk

## [2021.2.2.15] - 2021-5-15

### Changed
- directory for csv file


## [2020.4.2.11] - 2020-11-11

### Added
- when found product group id, prepare time and  scrap percent is updated
- when updating, creating product group, adding cycle from product

## [2020.4.1.26] - 2020-10-26

### Fixed
- fixed leaking goroutine bug when opening sql connections, the right way is this way

## [2020.3.3.19] - 2020-9-19

### Changed
- maps are created with their initial size

## [2020.3.3.18] - 2020-9-18

### Fixed
- when parsing cycle from CSV, proper checking for "9,9" instead of "9.9"

## [2020.3.3.15] - 2020-9-15

### Added
- complete functionality
- windows binary

### Changed
- comparing using map, because it is faster, better with big O(1)
- logging to proper folder when running as windows service

## [2020.3.3.14] - 2020-9-14

### Added
- initial commit
