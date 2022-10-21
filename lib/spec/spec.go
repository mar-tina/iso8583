package spec

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"

	"github.com/mar-tina/iso8583/lib/store"
)

//info fmt as below
//field:len:ffmt
//lenpad --> for llvar fields that might indicate their length in the first 2 bytes
func getFieldInfoFromTag(field reflect.StructField) (int, string, error) {
	tags := field.Tag
	info := ""
	lookup := []string{"field", "ln"}
	fieldId, ok := tags.Lookup("field")
	if !ok {
		return 0, "", errors.New("missing required tag: field")
	}

	fieldIdNum, err := strconv.Atoi(fieldId)
	if err != nil {
		return 0, "", errors.New("unusable formate for tag: field")
	}

	for i, key := range lookup {
		val, ok := tags.Lookup(key)
		sep := ""
		if i < len(lookup)-1 {
			sep = ":"
		}
		if ok {
			info += fmt.Sprintf("%v", val) + sep
		}
	}
	return fieldIdNum, info, nil
}

func Register(definitons ...interface{}) error {
	for i := 0; i < len(definitons); i++ {
		fields := reflect.TypeOf(definitons[i])
		secondaryBitmapFields := []int{}
		primaryBitmapFields := []int{}
		specstr := ""

		for i := 0; i < fields.NumField(); i++ {
			fieldNum, fieldInfo, err := getFieldInfoFromTag(fields.Field(i))
			if err != nil {
				return err
			}
			if fieldNum > 68 {
				secondaryBitmapFields = append(secondaryBitmapFields, fieldNum)
			} else {
				primaryBitmapFields = append(primaryBitmapFields, fieldNum)
			}

			log.Printf("fi: %s", fieldInfo)
			specstr += " " + fieldInfo
		}

		buildMap := ""
		// build primary bitmap
		if len(secondaryBitmapFields) > 0 {
			p := buildBitMap(primaryBitmapFields, true, true)
			buildMap += p
			s := buildBitMap(secondaryBitmapFields, false, false)
			buildMap += s
		} else {
			p := buildBitMap(primaryBitmapFields, false, true)
			buildMap += p
		}
		store.Put(fields.Name(), specstr, buildMap)
	}

	return nil
}

func buildBitMap(fields []int, isPrimary, initialSetToOne bool) string {
	var bitmap = ""
	sort.Ints(fields)
	zeroed := []byte{'0'}

	if initialSetToOne {
		bitmap += "1"
	}

	for _, field := range fields {
		if !isPrimary {
			field -= 68
		}

		if field > 0 {
			zeroes := bytes.Repeat(zeroed, field)
			log.Printf("count : %s %d", zeroes, field)
			currenBitmapLength := len(bitmap) + 1
			log.Printf("bitmap [before]: %s", bitmap)
			bitmap += string(zeroes[currenBitmapLength:])
			log.Printf("bitmap [after]: %s", bitmap)

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
	bin, err := strconv.ParseUint(s, 2, 64)
	if err != nil {
		return "error"
	}

	return fmt.Sprintf("%x", bin)
}
