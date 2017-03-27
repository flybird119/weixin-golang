/**
 * Copyright 2015-2016, Wothing Co., Ltd.
 * All rights reserved.
 *
 * Created by Elvizlai on 2016/04/07 17:34
 */

package misc

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"goushuyun/errs"
)

func FazzyQuery(value string) string {
	var fazzy_value = "%"
	for _, char := range value {
		char := fmt.Sprintf("%c", char)
		if char != " " {
			fazzy_value += (char + "%")
		}
	}
	return fazzy_value
}

func Md5String(objs ...interface{}) string {
	text := ""
	for i := range objs {
		text += fmt.Sprint(objs[i])
	}

	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func Contains(array []string, element string) bool {
	for i := range array {
		if array[i] == element {
			return true
		}
	}
	return false
}

var reg = regexp.MustCompile("86-1\\d{10}")

func MobileFormat(mobile string) ([]string, error) {
	if !reg.MatchString(mobile) {
		return []string{}, errs.NewError(errs.ErrMobileFormat, `The mobile should match 86-1\d{10}`)
	}

	s := strings.Split(mobile, "-")
	if _, err := strconv.Atoi(s[0]); err != nil {
		return []string{}, errs.NewError(errs.ErrMobileFormat, `mobile should 86-156xxxx`)
	}
	if _, err := strconv.Atoi(s[1]); err != nil {
		return []string{}, errs.NewError(errs.ErrMobileFormat, `mobile should 86-156xxxx`)
	}

	return []string{s[0], s[1]}, nil
}

func SuperPrint(x interface{}) string {
	buff := bytes.NewBuffer([]byte{})
	if err := encode(buff, reflect.ValueOf(x)); err != nil {
		return err.Error()
	}
	return buff.String()
}

func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())

	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())

	case reflect.Bool:
		fmt.Fprintf(buf, "%t", v.Bool())

	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%g", v.Float())

	case reflect.Ptr:
		buf.WriteByte('&')
		return encode(buf, v.Elem())

	case reflect.Array, reflect.Slice:
		buf.WriteString(v.Type().String())
		buf.WriteByte('{')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteString(", ")
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte('}')

	case reflect.Struct:
		buf.WriteString(v.Type().String())
		buf.WriteByte('{')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteString(", ")
			}
			fmt.Fprintf(buf, "%s:", v.Type().Field(i).Name)
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
		}
		buf.WriteByte('}')

	case reflect.Map:
		buf.WriteString(v.Type().String())
		buf.WriteByte('{')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteString(", ")
			}
			if err := encode(buf, key); err != nil {
				return err
			}
			buf.WriteByte(':')
			if err := encode(buf, v.MapIndex(key)); err != nil {
				return err
			}
		}
		buf.WriteByte('}')

	case reflect.Interface:
		return encode(buf, v.Elem())

	default: // complex, chan, func
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}
