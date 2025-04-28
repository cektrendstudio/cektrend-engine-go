package uttime

import (
	"fmt"
	"testing"
	"time"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"

	"github.com/stretchr/testify/assert"
)

func setTimezone(zone string) serror.SError {
	var loc *time.Location
	loc, errx := GetTimezone(zone)
	if errx != nil {
		errx.AddComments("while get timezone")
		return errx
	}

	time.Local = loc
	return errx
}

func TestParseFromString(t *testing.T) {
	errx := setTimezone("+7")
	if errx != nil {
		t.Fatal(errx.String())
		t.FailNow()
	}

	most := func(tim time.Time, err error) time.Time {
		if err == nil {
			return tim
		}

		t.Fatal(err.Error())
		return time.Now()
	}

	vals := map[string]time.Time{
		"2020-01-01 01:00:00":   most(time.Parse(time.RFC3339, "2020-01-01T01:00:00+07:00")),
		"2020-01-01 01:00:00 Z": most(time.Parse(time.RFC3339, "2020-01-01T01:00:00Z")),
		"2020-01-01T01:00:00":   most(time.Parse(time.RFC3339, "2020-01-01T01:00:00+07:00")),
		"2020-01-01T01:00:00Z":  most(time.Parse(time.RFC3339, "2020-01-01T01:00:00Z")),
		"2020-01-01":            most(time.Parse(time.RFC3339, "2020-01-01T00:00:00+07:00")),
		"2020-01-01Z":           most(time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")),
	}

	for k, v := range vals {
		tim, errx := ParseFromString(k)
		if errx != nil {
			errx.AddComments("while parse from string")

			t.Fatal(errx.String())
			t.FailNow()
		}

		assert.Equal(t, v, tim, fmt.Sprintf("from %v", k))
	}
}

func TestParseUTCFromString(t *testing.T) {
	errx := setTimezone("+7")
	if errx != nil {
		t.Fatal(errx.String())
		t.FailNow()
	}

	most := func(tim time.Time, err error) time.Time {
		if err == nil {
			return tim
		}

		t.Fatal(err.Error())
		return time.Now()
	}

	vals := map[string]time.Time{
		"2020-01-01 01:00:00":   most(time.Parse(time.RFC3339, "2020-01-01T01:00:00Z")),
		"2020-01-01 01:00:00 Z": most(time.Parse(time.RFC3339, "2020-01-01T01:00:00Z")),
		"2020-01-01T01:00:00":   most(time.Parse(time.RFC3339, "2020-01-01T01:00:00Z")),
		"2020-01-01T01:00:00Z":  most(time.Parse(time.RFC3339, "2020-01-01T01:00:00Z")),
		"2020-01-01":            most(time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")),
		"2020-01-01Z":           most(time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")),
	}

	for k, v := range vals {
		tim, errx := ParseUTCFromString(k)
		if errx != nil {
			errx.AddComments("while parse utc from string")

			t.Fatal(errx.String())
			t.FailNow()
		}

		assert.Equal(t, v, tim, fmt.Sprintf("from %v", k))
	}
}
