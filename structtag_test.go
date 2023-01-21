package structtag_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/quartercastle/structtag"
)

func Example() {
	tags, err := structtag.Parse(`json:"host" env:"SERVER_HOST" default:"localhost"`)

	if err != nil {
		panic(err)
	}

	fmt.Println(tags["env"])
	// Output: SERVER_HOST
}

func ExampleMerge() {
	t1 := structtag.Map{
		"env": "TESTING",
	}

	t2 := structtag.Map{
		"env": "HELLO",
	}

	t := structtag.Merge(t1, t2)
	fmt.Println(t)
	// Output: map[env:HELLO]
}

func ExampleParse() {
	tags, err := structtag.Parse(`env:"SERVER_HOST" default:"localhost"`)
	fmt.Println(tags, err)
	// Output: map[default:localhost env:SERVER_HOST] <nil>
}

func ExampleMap_StructTag() {
	tags, _ := structtag.Parse(`env:"SERVER_HOST"`)
	st := tags.StructTag()
	fmt.Println(st)
	// Output: env:"SERVER_HOST"
}

func TestParse(t *testing.T) {
	tags, _ := structtag.Parse(`env:"SERVER_HOST" default:"localhost"`)

	if v, ok := tags["env"]; !ok || v != "SERVER_HOST" {
		t.Errorf("the key env is not equal to SERVER_HOST; got %s", v)
	}

	if v, ok := tags["default"]; !ok || v != "localhost" {
		t.Errorf("the key default is not equal to localhost; got %s", v)
	}
}

func TestParseErrors(t *testing.T) {
	_, err := structtag.Parse(`invalid syntax`)
	if err != structtag.ErrInvalidSyntax {
		t.Errorf("did not return error %s; got %s", structtag.ErrInvalidSyntax, err)
	}

	cases := []reflect.StructTag{
		`:value`,
		`"":"value"`,
	}

	for _, c := range cases {
		_, err = structtag.Parse(c)
		if err != structtag.ErrInvalidKey {
			t.Errorf("did not return error %s; got %s", structtag.ErrInvalidKey, err)
		}
	}

	cases = []reflect.StructTag{
		`key:value`,
		`key:value"`,
		`key:"value`,
		`key:"value\"`,
		`key:\"value"`,
		`key: ""`,
	}

	for _, c := range cases {
		_, err := structtag.Parse(c)
		if err != structtag.ErrInvalidValue {
			t.Errorf("did not return error %s; got %s", structtag.ErrInvalidValue, err)
		}
	}

	_, err = structtag.Parse(`key:"value", other:"value"`)
	if err != structtag.ErrInvalidSeparator {
		t.Errorf("did not return error %s; got %s", structtag.ErrInvalidSeparator, err)
	}
}
