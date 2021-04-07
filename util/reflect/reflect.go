package reflect

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrInvalidStruct = errors.New("invalid struct specified")
	ErrInvalidValue  = errors.New("invalid value specified")
)

func MergeMap(dst interface{}, mp map[string]interface{}, tags []string) error {
	var err error
	var sval reflect.Value
	var fname string

	dviface := reflect.ValueOf(dst)
	if dviface.Kind() == reflect.Ptr {
		dviface = dviface.Elem()
	}

	if dviface.Kind() != reflect.Struct {
		return ErrInvalidStruct
	}

	dtype := dviface.Type()
	for idx := 0; idx < dtype.NumField(); idx++ {
		dfld := dtype.Field(idx)
		dval := dviface.Field(idx)
		if !dval.CanSet() || len(dfld.PkgPath) != 0 || !dval.IsValid() {
			continue
		}

		fname = ""
		for _, tname := range tags {
			tvalue, ok := dfld.Tag.Lookup(tname)
			if !ok {
				continue
			}

			tpart := strings.Split(tvalue, ",")
			switch tname {
			case "protobuf":
				fname = tpart[3][5:]
			default:
				fname = tpart[0]
			}

			if fname != "" {
				break
			}
		}

		if fname == "" {
			fname = strings.ToLower(dfld.Name)
		}

		val, ok := mp[fname]
		if !ok {
			continue
		}

		sval = reflect.ValueOf(val)

		switch dval.Kind() {
		case reflect.Bool:
			err = mergeBool(dval, sval)
		case reflect.String:
			err = mergeString(dval, sval)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err = mergeInt(dval, sval)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			err = mergeUint(dval, sval)
		case reflect.Float32, reflect.Float64:
			err = mergeFloat(dval, sval)
		case reflect.Struct:
			mp, ok := sval.Interface().(map[string]interface{})
			if !ok {
				return ErrInvalidValue
			}
			err = MergeMap(dval.Interface(), mp, tags)
			/*
				  case reflect.Interface:
							  err = d.decodeBasic(name, input, outVal)
							case reflect.Map:
								err = mergeMap(dval, sval)
							case reflect.Ptr:
								err = mergePtr(dval, sval)
							case reflect.Slice:
								err = mergeSlice(dval, sval)
							case reflect.Array:
								err = mergeArray(dval, sval)
			*/
		default:
			err = ErrInvalidValue
		}
	}

	if err != nil {
		err = fmt.Errorf("err: %v key %v invalid val %v", err, fname, sval.Interface())
	}

	return err
}

func mergeBool(dval reflect.Value, sval reflect.Value) error {
	switch sval.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch sval.Int() {
		case 1:
			dval.SetBool(true)
		case 0:
			dval.SetBool(false)
		default:
			return ErrInvalidValue
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch sval.Uint() {
		case 1:
			dval.SetBool(true)
		case 0:
			dval.SetBool(false)
		default:
			return ErrInvalidValue
		}
	case reflect.Float32, reflect.Float64:
		switch sval.Float() {
		case 1:
			dval.SetBool(true)
		case 0:
			dval.SetBool(false)
		default:
			return ErrInvalidValue
		}
	case reflect.Bool:
		dval.SetBool(sval.Bool())
	case reflect.String:
		switch sval.String() {
		case "t", "T", "true", "TRUE", "True", "1", "yes":
			dval.SetBool(true)
		case "f", "F", "false", "FALSE", "False", "0", "no":
			dval.SetBool(false)
		default:
			return ErrInvalidValue
		}
	default:
		return ErrInvalidValue
	}
	return nil
}

func mergeString(dval reflect.Value, sval reflect.Value) error {
	switch sval.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dval.SetString(strconv.FormatInt(sval.Int(), sval.Type().Bits()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dval.SetString(strconv.FormatUint(sval.Uint(), sval.Type().Bits()))
	case reflect.Float32, reflect.Float64:
		dval.SetString(strconv.FormatFloat(sval.Float(), 'f', -1, sval.Type().Bits()))
	case reflect.Bool:
		switch sval.Bool() {
		case true:
			dval.SetString(strconv.FormatBool(true))
		case false:
			dval.SetString(strconv.FormatBool(false))
		}
	case reflect.String:
		dval.SetString(sval.String())
	default:
		return ErrInvalidValue
	}
	return nil
}

func mergeInt(dval reflect.Value, sval reflect.Value) error {
	switch sval.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dval.SetInt(sval.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dval.SetInt(int64(sval.Uint()))
	case reflect.Float32, reflect.Float64:
		dval.SetInt(int64(sval.Float()))
	case reflect.Bool:
		switch sval.Bool() {
		case true:
			dval.SetInt(1)
		case false:
			dval.SetInt(0)
		}
	case reflect.String:
		l, err := strconv.ParseInt(sval.String(), 0, dval.Type().Bits())
		if err != nil {
			return err
		}
		dval.SetInt(l)
	default:
		return ErrInvalidValue
	}
	return nil
}

func mergeUint(dval reflect.Value, sval reflect.Value) error {
	switch sval.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dval.SetUint(uint64(sval.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dval.SetUint(sval.Uint())
	case reflect.Float32, reflect.Float64:
		dval.SetUint(uint64(sval.Float()))
	case reflect.Bool:
		switch sval.Bool() {
		case true:
			dval.SetUint(1)
		case false:
			dval.SetUint(0)
		}
	case reflect.String:
		l, err := strconv.ParseUint(sval.String(), 0, dval.Type().Bits())
		if err != nil {
			return err
		}
		dval.SetUint(l)
	default:
		return ErrInvalidValue
	}
	return nil
}

func mergeFloat(dval reflect.Value, sval reflect.Value) error {
	switch sval.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dval.SetFloat(float64(sval.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dval.SetFloat(float64(sval.Uint()))
	case reflect.Float32, reflect.Float64:
		dval.SetFloat(sval.Float())
	case reflect.Bool:
		switch sval.Bool() {
		case true:
			dval.SetFloat(1)
		case false:
			dval.SetFloat(0)
		}
	case reflect.String:
		l, err := strconv.ParseFloat(sval.String(), dval.Type().Bits())
		if err != nil {
			return err
		}
		dval.SetFloat(l)
	default:
		return ErrInvalidValue
	}
	return nil
}
