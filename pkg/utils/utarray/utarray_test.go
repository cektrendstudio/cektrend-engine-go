package utarray

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHasCommonSubset(t *testing.T) {
	type TestTable struct {
		Name           string
		Origin, Target interface{}
		Expected       bool
	}

	testCases := []TestTable{
		{
			Name:     "array string and array string have subset at least one subset",
			Origin:   []interface{}{"lorem", "ipsum"},
			Target:   []interface{}{"lorem"},
			Expected: true,
		},
		{
			Name:     "array string and array string have subset at least one subset",
			Origin:   []interface{}{"lorem"},
			Target:   []interface{}{"lorem", "ipsum"},
			Expected: true,
		},
		{
			Name:   "array string and array string not have subset at least one subset",
			Origin: []interface{}{"dorem"},
			Target: []interface{}{"lorem", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum",
				"ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum",
				"ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum",
				"ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum",
				"ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum", "ipsum",
			},
			Expected: false,
		},
		{
			Name:     "mix data array and mix data type array not have subset at least one subset",
			Origin:   []interface{}{"dorem", 1},
			Target:   []interface{}{2, "lorem", "ipsum"},
			Expected: false,
		},
		{
			Name:     "mix data array and mix data type array have subset at least one subset",
			Origin:   []interface{}{"dorem", 1},
			Target:   []interface{}{1, "lorem", "ipsum"},
			Expected: true,
		},
		{
			Name:     "array float64 and float64 type have subset at least one subset",
			Origin:   []float64{2.0},
			Target:   2.0,
			Expected: true,
		},
		{
			Name:     "string and string type have subset at least one subset",
			Origin:   "lorem",
			Target:   "lorem",
			Expected: true,
		},
		{
			Name:     "integer and integer type have subset at least one subset",
			Origin:   1,
			Target:   1,
			Expected: true,
		},
		{
			Name:     "integer and integer type have subset at least one subset",
			Origin:   -1,
			Target:   -1,
			Expected: true,
		},
		{
			Name:     "integer and string type not have subset at least one subset",
			Origin:   "-1",
			Target:   -1,
			Expected: false,
		},
		{
			Name:     "array float64 but have more than 1 precision and float64 type have subset at least one subset",
			Origin:   []float64{2.0000000},
			Target:   2.0,
			Expected: true,
		},
		{
			Name:     "array float64 but have more than 1 precision and float64 type no have subset at least one subset",
			Origin:   []float64{2.000000001},
			Target:   2.0,
			Expected: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			result := IsIntersect(test.Origin, test.Target)
			require.Equal(t, test.Expected, result)
		})
	}
}
