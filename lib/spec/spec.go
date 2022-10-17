package spec

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type entry struct {
	value  string
	bitmap string
}

var signatureStore = struct {
	sync.RWMutex
	m map[string]entry
}{m: make(map[string]entry)}

func put(key string, value entry) {
	log.Printf("put: %s, %v", key, value)
	signatureStore.Lock()
	signatureStore.m[key] = value
	signatureStore.Unlock()
}

func get(key string) (entry, bool) {
	val, ok := signatureStore.m[key]
	return val, ok
}

type Iso8583 interface {
}

func getFieldNum(tags reflect.StructTag, key string) (int, bool) {
	f := tags.Get(key)
	val, err := strconv.Atoi(f)
	if err != nil {
		return 0, false
	}

	return val, true
}

func Register(definitons ...interface{}) error {
	for i := 0; i < len(definitons); i++ {
		defn := definitons[i]
		fields := reflect.TypeOf(defn)
		secondaryBitmapFields := []int{}
		primaryBitmapFields := []int{}
		defname := fields.Name()
		buildStr := ""

		for i := 0; i < fields.NumField(); i++ {
			tags := fields.Field(i).Tag
			fLen, okfLen := getFieldNum(tags, "ln")
			fieldNum, okfNum := getFieldNum(tags, "field")
			if !okfLen || !okfNum {
				return fmt.Errorf("invalid field tag for: %s", fields.Field(i).Name)
			}
			if fieldNum > 68 {
				secondaryBitmapFields = append(secondaryBitmapFields, fieldNum)
			} else {
				primaryBitmapFields = append(primaryBitmapFields, fieldNum)
			}

			buildStr += " " + fmt.Sprintf("%d", fieldNum) + ":" + fmt.Sprintf("%d", fLen)
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
		put(defname, entry{
			value:  buildStr,
			bitmap: buildMap,
		})
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
	bin, err := strconv.ParseUint(s, 2, 64)
	if err != nil {
		return "error"
	}

	return fmt.Sprintf("%x", bin)
}

func PackMsg[K Iso8583](payload K) (string, error) {
	defn := reflect.TypeOf(payload)
	defVals := reflect.ValueOf(payload)
	sig, ok := get(defn.Name())
	if !ok {
		return "", fmt.Errorf("unregistered struct: %s", defn.Name())
	}

	message := sig.value
	mti := defVals.FieldByName("Mti")
	if mti.String() == "" {
		return "", fmt.Errorf("required field Mti is missing")
	}

	for i := 0; i < defn.NumField(); i++ {
		msg, err := buildmsg(message, sig.bitmap, defn.Field(i).Tag.Get("field"), defn.Field(i).Tag.Get("ln"), defVals.Field(i).String())
		if err != nil {
			return "", err
		}
		message = msg
	}

	return strings.ReplaceAll(message, " ", ""), nil
}

func buildmsg(sig, bmap string, key, length, value string) (string, error) {
	ln, err := strconv.Atoi(length)
	if err != nil {
		return "", fmt.Errorf("failed to convert len to decimal: %s", err)
	}
	if key == "0" {
		value += bmap
	}
	return strings.Replace(sig, key+":"+length, value, ln), nil
}

func PackNetMsg() {

}
