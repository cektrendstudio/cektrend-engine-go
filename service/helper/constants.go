package helper

import (
	"time"
)

var JakartaLoc *time.Location

func init() {
	JakartaLoc, _ = time.LoadLocation(JAKARTA_TIME_LOC)
}

// DateFormat date format
type DateFormat = string

const (
	// Utilities
	JAKARTA_TIME_LOC = "Asia/Jakarta"
	ERROR_STRING     = "ERROR:"
	SALT             = 8

	// Date Format
	DATE_FORMAT_LONG            = "2006-01-02 15:04:05 -0700 MST"
	DATE_FORMAT_YYYY_MM         = "2006-01"
	DATE_FORMAT_YYYY_MM_DD      = "2006-01-02"
	DATE_FORMAT_YYYY_MM_DD_TIME = "2006-01-02 15:04:05"

	// Error Messages
	FILTER_NOT_ALLOWED       = "filter not allowed"
	ROLE_NOT_ALLOWED         = "your role is not allowed to use this API"
	DIRECT_LINE_NOT_APPROVED = "Direct line not yet approval this request"
	APPLICATION_ERROR        = "Application Error, please contact the dev team for further inquiry"

	// Year4Digits for Years in 4 digits
	Year4Digits = "2006"

	// Year2Digits for Years in 2 digits
	Year2Digits = "06"

	// Month2Digits for Months in 2 digits
	Month2Digits = "01"

	// Month1Digits for Months in 1 digits
	Month1Digits = "1"

	// Day2Digits for Days in 2 digits
	Day2Digits = "02"

	// Day1Digits for Days in 1 digits
	Day1Digits = "2"

	// Hour2Digits for Hours in 2 digits
	Hour2Digits = "15"

	// Minute2Digits for Minutes in 2 digits
	Minute2Digits = "04"

	// Second2Digits for Second in 2 digits
	Second2Digits = "05"

	// Milliseconds for Time Milliseconds
	Milliseconds = ".999999999"

	// Timezone for Timezone Location
	Timezone = "MST"

	// DefaultDateFormat for Default Date Format
	DefaultDateFormat DateFormat = "Y-m-d"

	// DefaultDateWithTimezoneFormat for Default Date Format with Timezone
	DefaultDateWithTimezoneFormat DateFormat = "Y-m-d TZ"

	// INDateFormat for Indonesian Date Format
	INDateFormat DateFormat = "d-m-Y"

	// DefaultDateTimeFormat for Default Date Time Format
	DefaultDateTimeFormat DateFormat = "Y-m-d H:i:s"

	// INDateTimeFormat for Indonesian Date Time Format
	INDateTimeFormat DateFormat = "d-m-Y H:i:s"

	// DefaultDateTimeWithTimezoneFormat for Default Date Time Format With Timezone
	DefaultDateTimeWithTimezoneFormat DateFormat = "Y-m-d H:i:s TZ"

	// DefaultTimeFormat for Default Time Format
	DefaultTimeFormat DateFormat = "H:i:s"

	// DefaultDateTimeFormat for Default Date Time Format With Milliseconds
	DefaultDateTimeWithMillisecondsFormat DateFormat = "Y-m-d H:i:s.u"
)
