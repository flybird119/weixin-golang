package db

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	. "github.com/goushuyun/weixin-golang/db"

	"github.com/wothing/log"
)

func TestDBOpe() {
	query := "select store_ids from users"

	var asSlice StringSlice

	rows, err := DB.Query(query)
	if err != nil {
		log.Error(err)
	}

	for rows.Next() {
		log.Debug("Found")
		rows.Scan(&asSlice)
	}

	log.Debugf("\n%+v\n", asSlice)
}

type StringSlice []string

// Implements sql.Scanner for the String slice type
// Scanners take the database value (in this case as a byte slice)
// and sets the value of the type.  Here we cast to a string and
// do a regexp based parse
func (s *StringSlice) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}

	asString := string(asBytes)
	parsed := parseArray(asString)
	(*s) = StringSlice(parsed)

	return nil
}

var (
	// unquoted array values must not contain: (" , \ { } whitespace NULL)
	// and must be at least one char
	unquotedChar  = `[^",\\{}\s(NULL)]`
	unquotedValue = fmt.Sprintf("(%s)+", unquotedChar)

	// quoted array values are surrounded by double quotes, can be any
	// character except " or \, which must be backslash escaped:
	quotedChar  = `[^"\\]|\\"|\\\\`
	quotedValue = fmt.Sprintf("\"(%s)*\"", quotedChar)

	// an array value may be either quoted or unquoted:
	arrayValue = fmt.Sprintf("(?P<value>(%s|%s))", unquotedValue, quotedValue)

	// Array values are separated with a comma IF there is more than one value:
	arrayExp = regexp.MustCompile(fmt.Sprintf("((%s)(,)?)", arrayValue))

	valueIndex int
)

// Find the index of the 'value' named expression
func init() {
	for i, subexp := range arrayExp.SubexpNames() {
		if subexp == "value" {
			valueIndex = i
			break
		}
	}
}

// Parse the output string from the array type.
// Regex used: (((?P<value>(([^",\\{}\s(NULL)])+|"([^"\\]|\\"|\\\\)*")))(,)?)
func parseArray(array string) []string {
	results := make([]string, 0)
	matches := arrayExp.FindAllStringSubmatch(array, -1)
	for _, match := range matches {
		s := match[valueIndex]
		// the string _might_ be wrapped in quotes, so trim them:
		s = strings.Trim(s, "\"")
		results = append(results, s)
	}
	return results
}

// func TestDBOpe() {
// 	query := "insert into users(openid, nickname, avatar, store_ids) values('%s', '%s', '%s', '{%s}')"
//
// 	store_ids := []string{"Wang", "Kai", "forever"}
//
// 	for k, v := range store_ids {
// 		store_ids[k] = fmt.Sprintf("\"%s\"", v)
// 	}
//
// 	store_ids_str := strings.Join(store_ids, ", ")
//
// 	log.Debugf("\n%s\n", store_ids_str)
// 	_, err := DB.Exec(fmt.Sprintf(query, "Wang", "Kai", "Elegant", store_ids_str))
//
// 	if err != nil {
// 		log.Error(err)
// 	}
// }

func SaveAppidToStore(store_id, app_id string) error {
	query := "update store set appid = $1 where id = $2"
	log.Debugf("update store set appid = '%s' where id = '%s'", app_id, store_id)

	_, err := DB.Exec(query, app_id, store_id)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
