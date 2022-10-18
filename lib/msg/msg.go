package msg

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mar-tina/iso8583/lib/store"
)

func PackMsg(payload interface{}) (string, error) {
	defn := reflect.TypeOf(payload)
	defVals := reflect.ValueOf(payload)
	sig, bitmap, ok := store.Get(defn.Name())
	if !ok {
		return "", fmt.Errorf("unregistered struct: %s", defn.Name())
	}

	message := sig
	mti := defVals.FieldByName("Mti")
	if mti.String() == "" {
		return "", fmt.Errorf("required field Mti is missing")
	}

	for i := 0; i < defn.NumField(); i++ {
		msg, err := buildmsg(message, bitmap, defn.Field(i).Tag.Get("field"), defn.Field(i).Tag.Get("ln"), defVals.Field(i).String())
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
