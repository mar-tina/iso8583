package spec

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"sync"
)

type fieldMapping struct {
	Field     int
	Fieldname string
	Len       int
	Pos       int
	Value     string
}

var signatureStore = struct {
	sync.RWMutex
	m map[string]fieldMapping
}{m: make(map[string]fieldMapping)}

// func put(key string, value fieldMapping) {
// 	signatureStore.Lock()
// 	signatureStore.m[key] = value
// }

type Iso8583 interface {
}

func Register[K Iso8583](definitons ...K) error {
	for _, defn := range definitons {
		fields := reflect.TypeOf(defn)
		secondaryBitmapFields := []int{}
		primaryBitmapFields := []int{}
		allFields := []fieldMapping{}

		for i := 0; i < fields.NumField(); i++ {
			tags := fields.Field(i).Tag
			fieldLen := tags.Get("ln")
			fieldTag := tags.Get("field")
			fieldLenNum, err1 := strconv.Atoi(fieldLen)
			fieldTagNum, err := strconv.Atoi(fieldTag)
			if err != nil || err1 != nil {
				return fmt.Errorf("invalid field tag for: %s", fields.Field(i).Name)
			}

			if fieldTagNum > 68 {
				secondaryBitmapFields = append(secondaryBitmapFields, fieldTagNum)
			} else {
				primaryBitmapFields = append(primaryBitmapFields, fieldTagNum)
			}

			allFields = append(allFields, fieldMapping{
				Pos:       i,
				Len:       fieldLenNum,
				Fieldname: tags.Get("json"),
				Field:     fieldTagNum,
			})
		}

		// build primary bitmap
		if len(secondaryBitmapFields) > 0 {
			primary := buildBitMap(primaryBitmapFields, true, true)
			allFields = append(allFields, fieldMapping{
				Field: 1,
				Len:   16,
				Value: primary,
				Pos:   1,
			})
			secodary := buildBitMap(secondaryBitmapFields, false, false)
			allFields = append(allFields, fieldMapping{
				Field: 1,
				Len:   16,
				Value: secodary,
				Pos:   2,
			})

		} else {
			primary := buildBitMap(primaryBitmapFields, false, true)
			allFields = append(allFields, fieldMapping{
				Field: 1,
				Len:   16,
				Value: primary,
				Pos:   1,
			})
		}
	}

	return nil
}

// initialSetToOne defautls to the initial bitmap field is set to '1'
func buildBitMap(fields []int, isPrimary, initialSetToOne bool) string {
	var bitmap = ""
	sort.Ints(fields)
	zeroed := []byte{'0'}

	if initialSetToOne {
		bitmap += "1"
	}
	if !isPrimary {
		bitmap += "0"
	}

	for _, field := range fields {
		if !isPrimary {
			field -= 68
		}
		if field > 0 {
			zeroes := bytes.Repeat(zeroed, field)
			currenBitmapLength := len(bitmap) + 1
			bitmap += string(zeroes[currenBitmapLength:])
			bitmap += "1"
		}
	}

	if len(bitmap) < 64 {
		zeroes := bytes.Repeat(zeroed, 64)
		bitmap += string(zeroes[len(bitmap):])
	}

	return parseBinToHex(bitmap)
}

func parseBinToHex(s string) string {
	ui, err := strconv.ParseUint(s, 2, 64)
	if err != nil {
		return "error"
	}

	return fmt.Sprintf("%x", ui)
}

func PackMsg(defn interface{}) {
	log.Printf("defn: %v", defn)
}

func PackNetMsg() {

}
