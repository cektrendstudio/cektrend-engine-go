package uttime

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"

	"github.com/stretchr/testify/assert"
)

func TestTimeParse(t *testing.T) {
	errx := setTimezone("+7")
	if errx != nil {
		t.Fatal(errx.String())
		t.FailNow()
	}

	vals := map[string]Time{
		`"2020-01-01 01:00:00"`:       Time(MostParse("2020-01-01 01:00:00 +07:00")),
		`"2020-02-12T18:35:11+07:00"`: Time(MostParse("2020-02-12T18:35:11+07:00")),
		`"2020-02-12T18:35:11Z"`:      Time(MostParse("2020-02-13T01:35:11+07:00")),
		`"2020-05-03"`:                Time(MostParse("2020-05-03")),
		`"2020-01-01+07:00"`:          Time(MostParse("2020-01-01+07:00")),
		`"20200503130001"`:            Time(MostParse("2020-05-03 13:00:01")),
		`"20200503"`:                  Time(MostParse("2020-05-03")),
		`1597038426375`:               Time(MostParse("2020-08-10 12:47:06")),
	}

	for k, v := range vals {
		var vx Time
		err := json.Unmarshal([]byte(k), &vx)
		if err != nil {
			if errx, ok := err.(serror.SError); ok {
				t.Fatal(errx.String())
				t.FailNow()
				return
			}

			t.Fatal(err.Error())
			t.FailNow()
		}

		assert.Equal(t, v, vx, fmt.Sprintf("from %v", k))
	}
}

func TestTimeFormat(t *testing.T) {
	errx := setTimezone("+7")
	if errx != nil {
		t.Fatal(errx.String())
		t.FailNow()
	}

	vals := map[Time]string{
		Time(MostParse("2020-01-01 01:00:00")):                               `"2020-01-01T01:00:00+07:00"`,
		Time(MostParseForceTimezone("2020-02-12 18:35:11", "Asia/Shanghai")): `"2020-02-12T18:35:11+08:00"`,
	}

	for k, v := range vals {
		byt, err := json.Marshal(k)
		if err != nil {
			if errx, ok := err.(serror.SError); ok {
				t.Fatal(errx.String())
				t.FailNow()
				return
			}

			t.Fatal(err.Error())
			t.FailNow()
		}

		assert.Equal(t, v, string(byt), fmt.Sprintf("from %v", k))
	}
}

func TestDateParse(t *testing.T) {
	errx := setTimezone("+7")
	if errx != nil {
		t.Fatal(errx.String())
		t.FailNow()
	}

	vals := map[string]Date{
		`"2020-01-01 01:00:00"`:       Date(MostParse("2020-01-01 01:00:00 +07:00")),
		`"2020-02-12T18:35:11+07:00"`: Date(MostParse("2020-02-12T18:35:11+07:00")),
		`"2020-02-12T18:35:11Z"`:      Date(MostParse("2020-02-13T01:35:11+07:00")),
		`"2020-05-03"`:                Date(MostParse("2020-05-03")),
		`"2020-01-01+07:00"`:          Date(MostParse("2020-01-01+07:00")),
		`"20200503130001"`:            Date(MostParse("2020-05-03 13:00:01")),
		`"20200503"`:                  Date(MostParse("2020-05-03")),
		`1597038426375`:               Date(MostParse("2020-08-10 12:47:06")),
	}

	for k, v := range vals {
		var vx Date
		err := json.Unmarshal([]byte(k), &vx)
		if err != nil {
			if errx, ok := err.(serror.SError); ok {
				t.Fatal(errx.String())
				t.FailNow()
				return
			}

			t.Fatal(err.Error())
			t.FailNow()
		}

		assert.Equal(t, v, vx, fmt.Sprintf("from %v", k))
	}
}

func TestDateFormat(t *testing.T) {
	errx := setTimezone("+7")
	if errx != nil {
		t.Fatal(errx.String())
		t.FailNow()
	}

	vals := map[Date]string{
		Date(MostParse("2020-01-01 01:00:00")):                               `"2020-01-01"`,
		Date(MostParseForceTimezone("2020-02-12 18:35:11", "Asia/Shanghai")): `"2020-02-12"`,
	}

	for k, v := range vals {
		byt, err := json.Marshal(k)
		if err != nil {
			if errx, ok := err.(serror.SError); ok {
				t.Fatal(errx.String())
				t.FailNow()
				return
			}

			t.Fatal(err.Error())
			t.FailNow()
		}

		assert.Equal(t, v, string(byt), fmt.Sprintf("from %v", k))
	}
}
