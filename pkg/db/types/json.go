package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

// JSON is db value type of JSON object
type JSON map[string]any

func (j *JSON) UpdateFields(fields JSON) {
	if *j == nil {
		*j = JSON{}
	}

	for k, v := range fields { // 필드를 순환하면서 비어있지 않은 필드를 업데이트
		// v가 JSON 객체일 경우 재귀호출
		if val, ok := v.(JSON); ok {
			if jsonValue, ok := (*j)[k].(JSON); !ok {
				(*j)[k] = val
			} else {
				jsonValue.UpdateFields(val)
			}
		} else { // 아니면 그냥 업데이트
			(*j)[k] = v
		}
	}
}

func (j *JSON) Scan(value any) (err error) {
	if value == nil {
		*j = nil
		return nil
	}

	*j = make(JSON)
	switch v := value.(type) {
	case string:
		err = json.Unmarshal([]byte(v), j)
		return
	case []byte:
		err = json.Unmarshal(v, j)
		return
	}

	return errorx.New("Can't convert to JSON").With("value", value)
}

func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// JsonRaw for JSON DB field
type JsonRaw json.RawMessage

func (d *JsonRaw) Scan(value any) (err error) {
	if value == nil {
		*d = nil
		return nil
	}

	switch v := value.(type) {
	case string:
		*d = JsonRaw(v)
		return
	case []byte:
		if v == nil {
			return
		}
		*d = append((*d)[0:0], v...)
		return
	}

	return errorx.New("Can't convert to JsonRaw").With("value", value)
}

func (d JsonRaw) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	return []byte(d), nil
}

func (d *JsonRaw) UnmarshalJSON(data []byte) error {
	var msg json.RawMessage
	if err := msg.UnmarshalJSON(data); err != nil {
		return errorx.Wrap(err)
	}
	*d = JsonRaw(msg)
	return nil
}

func (d JsonRaw) MarshalJSON() ([]byte, error) {
	return json.RawMessage(d).MarshalJSON()
}
