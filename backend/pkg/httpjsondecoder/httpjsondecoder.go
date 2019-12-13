package httpjsondecoder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
)

type Decoder struct {
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

//registerconverter
func (d *Decoder) Decode(r *http.Request, payload interface{}) error {
	typ := reflect.TypeOf(payload).Elem()
	value := reflect.ValueOf(payload).Elem()

	params := mux.Vars(r) // TODO: replace with FormValues and get rid of mux dependence
	query := r.URL.Query()

	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := value.Field(i)
		jsonTag := fieldType.Tag.Get("json")

		if jsonTag == "" || jsonTag[0] == '-' {
			continue
		}

		paramStringValue, ok := params[jsonTag]

		if ok {
			err := LiteralStore(paramStringValue, fieldValue)
			if err != nil {
				return err
			}
		}

		queryValue, ok := query[jsonTag]
		if ok {
			if len(queryValue) == 1 {
				err := LiteralStore(queryValue[0], fieldValue)
				if err != nil {
					return err
				}
			}
		}
	}

	bytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, payload)
		if err != nil {
			return err
		}
	}

	return nil
}

var numberType = reflect.TypeOf(json.Number(""))

const useNumber bool = false

func LiteralStore(s string, v reflect.Value) error {
	if len(s) == 0 {
		return fmt.Errorf("Empty string given for %v", v.Type())
	}

	{
		vAddr := v
		if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
			// haveAddr = true
			vAddr = v.Addr()
		}

		if vAddr.Type().NumMethod() > 0 && vAddr.CanInterface() {
			if u, ok := vAddr.Interface().(json.Unmarshaler); ok {
				return u.UnmarshalJSON([]byte(s))
			}
			// if !decodingNull {
			// 	if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
			// 		return nil, u, reflect.Value{}
			// 	}
			// }
		}
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil || v.OverflowInt(n) {
			return fmt.Errorf("Failed to parse [%s] in to %v", s, v.Type())
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil || v.OverflowUint(n) {
			return fmt.Errorf("Failed to parse [%s] in to %v", s, v.Type())
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, v.Type().Bits())
		if err != nil || v.OverflowFloat(n) {
			return fmt.Errorf("Failed to parse [%s] in to %v", s, v.Type())
		}
		v.SetFloat(n)
	case reflect.Bool:
		n, err := strconv.ParseBool(s)
		if err != nil {
			return fmt.Errorf("Failed to parse [%s] in to %v", s, v.Type())
		}
		v.SetBool(n)
	default:
		return fmt.Errorf("Unsupported type %v", v.Type())
	}

	return nil
}
