// The motivation behind this package is that the StructTag implementation shipped
// with Go's standard library is very limited in detecting a malformed StructTag
// and each time StructTag.Get(key) gets called, it results in the StructTag
// being parsed again.
// This package provides a way to parse the StructTag into a structtag.Map.
//
// 	// Example of struct using StructTags to append metadata to fields.
// 	type Server struct {
//		Host string `json:"host" env:"SERVER_HOST" default:"localhost"`
//		Port int    `json:"port" env:"SERVER_PORT" default:"3000"`
//	}
//
package structtag

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	// ErrInvalidSyntax is returned when the StructTag syntax is invalid.
	ErrInvalidSyntax = errors.New("invalid syntax for key value pair")

	// ErrInvalidKey is returned if a key is containing invalid characters or
	// is missing.
	ErrInvalidKey = errors.New("invalid key")

	// ErrInvalidValue is returned if a value is not qouted.
	ErrInvalidValue = errors.New("invalid value")

	// ErrInvalidSeparator is returned if comma is used as separator.
	ErrInvalidSeparator = errors.New("invalid separator, key value pairs should be separated by spaces")
)

// Map is just a map of key value pairs.
type Map map[string]string

// Merge multiple Maps together into a single Tag.
// In case of duplicate keys, the last encountered key will overwrite the existing.
func Merge(maps ...Map) Map {
	for _, t := range maps {
		for k, v := range t {
			maps[0][k] = v
		}
	}

	return maps[0]
}

// StructTag converts the Tag into a StructTag.
func (m Map) StructTag() reflect.StructTag {
	var s string
	for k, v := range m {
		s += fmt.Sprintf(`%s:"%s" `, k, v)
	}
	return reflect.StructTag(strings.TrimSpace(s))
}

// Parse takes a StructTag and parses it into a Tag or returns an error.
// If the given string contains duplicate keys the last key value pair
// will overwrite the previous.
//
// The parsing logic is a slightly modified version of the StructTag.Lookup
// function from the reflect package included in the standard library.
// https://github.com/golang/go/blob/0377f061687771eddfe8de78d6c40e17d6b21a39/src/reflect/type.go#L1132
func Parse(st reflect.StructTag) (Map, error) {
	tags := Map{}

	for st != "" {
		i := 0
		for i < len(st) && st[i] == ' ' {
			i++
		}

		st = st[i:]
		if st == "" {
			break
		}

		i = 0
		for i < len(st) && st[i] > ' ' && st[i] != ':' && st[i] != '"' && st[i] != 0x7f {
			if st[i] == ',' {
				return tags, ErrInvalidSeparator
			}
			i++
		}

		if i == 0 {
			return tags, ErrInvalidKey
		}

		if i+1 >= len(st) || st[i] != ':' {
			return tags, ErrInvalidSyntax
		}

		if st[i+1] != '"' {
			return tags, ErrInvalidValue
		}

		key := string(st[:i])
		st = st[i+1:]

		i = 1
		for i < len(st) && st[i] != '"' {
			if st[i] == '\\' {
				i++
			}
			i++
		}

		if i >= len(st) {
			return tags, ErrInvalidValue
		}

		qvalue := string(st[:i+1])
		st = st[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			return tags, ErrInvalidValue
		}

		tags[key] = value
	}

	return tags, nil
}
