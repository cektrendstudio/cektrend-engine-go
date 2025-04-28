package helper

import (
	"bytes"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string) string {
	password := []byte(pass)

	hash, _ := bcrypt.GenerateFromPassword(password, SALT)

	return string(hash)
}

func ComparePassword(hashed, password []byte) bool {
	h, p := []byte(hashed), []byte(password)

	err := bcrypt.CompareHashAndPassword(h, p)

	return err == nil
}

func CheckPassword(password string) error {
	lengthPassword := len(password)

	if lengthPassword < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if lengthPassword > 20 {
		return errors.New("password must be maximum 20 characters")
	}

	return nil
}

func SplitCamelCase(input string) string {
	var result bytes.Buffer

	for i, char := range input {
		if i > 0 && unicode.IsUpper(char) && unicode.IsLower(rune(input[i-1])) {
			result.WriteRune(' ')
		}
		result.WriteRune(char)
	}

	return result.String()
}

func PrepareDateFilters(startDate, endDate time.Time) (time.Time, time.Time) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	startDateForFilter := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
	endDateForFilter := time.Date(endDate.Year(), endDate.Month(), endDate.Day()+1, 0, 0, 0, 0, loc)
	return startDateForFilter, endDateForFilter
}

func ConvertDateForFilter(startDate, endDate time.Time) (string, string) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	yearStart, monthStart, dayStart := startDate.Date()
	yearEnd, monthEnd, dayEnd := endDate.Date()
	startDateForFilter := time.Date(yearStart, monthStart, dayStart, 0, 0, 0, 0, loc)
	endDateForFilter := time.Date(yearEnd, monthEnd, dayEnd+1, 0, 0, 0, 0, loc)
	startDateStr := startDateForFilter.Format(DATE_FORMAT_YYYY_MM_DD)
	endDateStr := endDateForFilter.Format(DATE_FORMAT_YYYY_MM_DD)
	return startDateStr, endDateStr
}

func RemoveHtmlString(str string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	result := re.ReplaceAllString(str, "")
	return result
}

func ParseDateRange(start, end string) (time.Time, time.Time, error) {
	currentTime := time.Now()

	startDate := currentTime.AddDate(0, 0, -7)
	if start != "" {
		parsedStart, err := time.Parse(DATE_FORMAT_YYYY_MM_DD, start)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		startDate = parsedStart
	}

	endDate := currentTime
	if end != "" {
		parsedEnd, err := time.Parse(DATE_FORMAT_YYYY_MM_DD, end)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		endDate = parsedEnd
	}

	return startDate, endDate, nil
}

func AppEnvIsLoca() bool {
	return os.Getenv("APP_ENV") == "local"
}

func ParsePaginationParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	return page, size
}

func ParseDateFilterYearMonth(date string) (time.Time, error) {
	if date != "" {
		return time.Parse(DATE_FORMAT_YYYY_MM, date)
	}
	return time.Now(), nil
}

func ParseDateTime(format DateFormat, value time.Time) string {
	return value.Format(GoLayout(format))
}

func GoLayout(format DateFormat) string {
	rl := map[string]string{
		"Y":  Year4Digits,
		"y":  Year2Digits,
		"m":  Month2Digits,
		"M":  Month1Digits,
		"d":  Day2Digits,
		"D":  Day1Digits,
		"H":  Hour2Digits,
		"i":  Minute2Digits,
		"s":  Second2Digits,
		".u": Milliseconds,
		"TZ": Timezone,
	}

	for k, v := range rl {
		format = strings.ReplaceAll(format, k, v)
	}

	return format
}
