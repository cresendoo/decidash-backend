package types

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

type IntArray []int

func (arr *IntArray) Scan(value any) (err error) {
	if value == nil {
		*arr = IntArray{}
		return nil
	}

	*arr = IntArray{}
	switch v := value.(type) {
	case string:
		err = json.Unmarshal([]byte(v), arr)
		return
	case []byte:
		err = json.Unmarshal(v, arr)
		return
	}
	return errorx.New("Can't convert to IntArray").With("value", value)
}

func (arr IntArray) Value() (driver.Value, error) {
	if arr == nil {
		return nil, nil
	}
	return json.Marshal(arr)
}

func (arr *IntArray) UnmarshalJSON(b []byte) error {
	var intArr []int
	if err := json.Unmarshal(b, &intArr); err != nil {
		return err
	}
	*arr = intArr
	return nil
}

func (arr IntArray) MarshalJSON() ([]byte, error) {
	if arr == nil {
		return []byte("null"), nil
	}
	return json.Marshal([]int(arr))
}

func (arr IntArray) String() string {
	var str strings.Builder
	str.WriteString("[")
	if len(arr) > 1 {
		str.WriteString(strconv.Itoa(arr[0]))
	}
	for idx := 1; idx < len(arr); idx++ {
		str.WriteString(",")
		str.WriteString(strconv.Itoa(arr[idx]))
	}
	str.WriteString("]")
	return str.String()
}

func (arr IntArray) ToIntArray() []int {
	return arr
}

func (arr *IntArray) Append(numbers ...int) {
	*arr = append(*arr, numbers...)
}
