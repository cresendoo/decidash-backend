package query

import (
	"fmt"
	"log/slog"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/cresendoo/decidash-backend/pkg/db/encrypt"
)

const (
	tagNameDatabase   = "db"
	tagValueSkipField = "-"
	tagSeparators     = "," // first value must be column name

	columnOptionEncrypt       = "encrypt"
	columnOptionOmitEmpty     = "omitempty"
	columnOptionPrimaryKey    = "primary_key"
	columnOptionAutoIncrement = "auto_increment"

	managedTimeColumnNameCreatedAt = "created_at"
	managedTimeColumnNameUpdatedAt = "updated_at"
	managedTimeColumnNameDeletedAt = "deleted_at"
)

var structInfoMap sync.Map

type structInfo map[string]*columnInfo

type columnInfo struct {
	structIDX int
	name      string

	Encrypted     bool
	OmitEmpty     bool
	PrimaryKey    bool
	AutoIncrement bool

	cipher encrypt.Cipher
}

func getStructInfo(st reflect.Type) structInfo {
	st = indirectType(st)
	stName := st.String()
	info, ok := structInfoMap.Load(stName)
	if !ok {
		newinfo := make(structInfo)
		for idx := 0; idx < st.NumField(); idx++ {
			field := st.Field(idx)
			if tagValue, ok := field.Tag.Lookup(tagNameDatabase); ok {
				if tagValue == tagValueSkipField {
					continue
				}
				ci := &columnInfo{structIDX: idx}
				tags := strings.Split(tagValue, tagSeparators)
				ci.name = tags[0]
				for i := 1; i < len(tags); i++ {
					switch tags[i] {
					case columnOptionAutoIncrement:
						ci.AutoIncrement = true
					case columnOptionEncrypt:
						ci.Encrypted = true
						ci.cipher = encrypt.NewColumnCipher(stName, ci.name)
					case columnOptionPrimaryKey:
						ci.PrimaryKey = true
					case columnOptionOmitEmpty:
						ci.OmitEmpty = true
					}
				}
				newinfo[ci.name] = ci
			}
		}

		// 다른 쓰레드에서 작업이 일어나 이미 Map에 들어가 있다면, 해당 정보를 사용
		// 아닐경우 Map에 추가
		if info, ok = structInfoMap.LoadOrStore(stName, newinfo); ok {
			return info.(structInfo)
		}
	}
	return info.(structInfo)
}

func getStructColumns(v any, options ...getStructColumnsOptionFunc) []string {
	opt := &getStructColumnsOption{}
	for _, f := range options {
		f(opt)
	}

	var st reflect.Type
	var val reflect.Value
	if _, ok := v.(reflect.Type); ok {
		st = v.(reflect.Type)
	} else {
		st = reflect.TypeOf(v)
		val = reflect.Indirect(reflect.ValueOf(v))
	}
	columnInfos := make([]*columnInfo, 0)
	structInfo := getStructInfo(st)
	var pk *columnInfo
	for _, ci := range structInfo {
		if ci.AutoIncrement && opt.excludeAutoIncrement {
			continue
		}
		if ci.PrimaryKey {
			pk = ci
			continue
		}
		if !opt.withManagedTimeColumn {
			if ci.name == managedTimeColumnNameCreatedAt ||
				ci.name == managedTimeColumnNameDeletedAt ||
				ci.name == managedTimeColumnNameUpdatedAt {
				continue
			}
		}
		if !opt.ignoreOmitEmpty &&
			ci.OmitEmpty &&
			val.IsValid() &&
			isEmptyValue(val.Field(ci.structIDX)) {
			continue
		}
		columnInfos = append(columnInfos, ci)
	}

	sort.Sort(byIndex(columnInfos))
	cols := []string{}
	if pk != nil {
		cols = append(cols, pk.name)
	}
	for _, ci := range columnInfos {
		cols = append(cols, ci.name)
	}
	return cols
}

type getStructColumnsOption struct {
	excludeAutoIncrement  bool
	ignoreOmitEmpty       bool
	withManagedTimeColumn bool
}

type getStructColumnsOptionFunc func(*getStructColumnsOption)

// DB에서 관리하는 Time관련 Column은 ReadOnly
func withManagedTimeColumn(flag bool) getStructColumnsOptionFunc {
	return func(o *getStructColumnsOption) {
		o.withManagedTimeColumn = flag
	}
}

// AutoIncrement Column은 Insert, Update 방지
func excludeAutoIncrement(flag bool) getStructColumnsOptionFunc {
	return func(o *getStructColumnsOption) {
		o.excludeAutoIncrement = flag
	}
}

// Select, Update때는 OmitEmpty를 무시
func ignoreOmitEmpty(flag bool) getStructColumnsOptionFunc {
	return func(o *getStructColumnsOption) {
		o.ignoreOmitEmpty = flag
	}
}

func getStructValues(cols []string, v any) []any {
	val := reflect.Indirect(reflect.ValueOf(v))
	structInfo := getStructInfo(val.Type())
	values := make([]any, len(cols))
	for i, col := range cols {
		if ti, ok := structInfo[col]; ok {
			if ti.Encrypted {
				values[i] = encrypt.NewColumn(
					v.(encrypt.EncryptID),
					val.Field(ti.structIDX).Interface(),
					ti.cipher,
				)
			} else {
				values[i] = val.Field(ti.structIDX).Interface()
			}
		}
	}
	slog.Debug(fmt.Sprintf("%+v", values))
	return values
}
