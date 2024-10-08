package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"io"
	"reflect"
	"slices"
	"strings"
)

type valueKey struct {
	value, key reflect.Value
}

func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("null")

	case reflect.Bool:
		if v.Bool() {
			fmt.Fprintf(buf, "true")
		} else {
			fmt.Fprintf(buf, "false")
		}

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())

	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%v", v.Float())


	// Not supported by json/encoding
	// case reflect.Complex64, reflect.Complex128:
		// c := v.Complex()
		// fmt.Fprintf(buf, "#C(%f, %f)", real(c), imag(c))

	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())

	case reflect.Ptr:
		return encode(buf, v.Elem())

	case reflect.Array, reflect.Slice:
		buf.WriteByte('[')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(']')

	case reflect.Struct:
		buf.WriteByte('{')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprintf(buf, "\"%s\":", v.Type().Field(i).Name)
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
		}
		buf.WriteByte('}')

	case reflect.Map:
		// We need to sort the keys (json/encoding does)
		vks := make([]valueKey, v.Len())
		iter := v.MapRange()

		// encoding/json is smarter than that
		for i:= 0; iter.Next(); i++ {
			vks[i] = valueKey{iter.Value(), iter.Key()}
		}

		slices.SortFunc(vks, func(i, j valueKey) int {
			return strings.Compare(i.key.String(), j.key.String())
		})
		buf.WriteByte('{')
		for i, vk := range vks {
			if i > 0 {
				buf.WriteByte(',')

			}
			if err := encode(buf, vk.key); err != nil {
				return err
			}
			buf.WriteByte(':')
			if err := encode(buf, vk.value); err != nil {
				return err
			}
		}
		buf.WriteByte('}')

	// Again, unmanaged by encoding/json
	case reflect.Interface:
		buf.WriteString("{}")
/*
		fmt.Fprintf(buf, "(%q ", v.Type().String())
		err := encode(buf, reflect.ValueOf(v.Interface()))
		if err != nil {
			return err
		}
		buf.WriteByte(')')
*/
	case reflect.UnsafePointer:
		fmt.Fprintf(buf, "%p", v.UnsafePointer())

	default: // chan, func
		return fmt.Errorf("unsupported type: %s", v.Type())

	}
	return nil
}

// Marshal encodes a Go value in S-expression form.
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
	Score           float64
	File            io.WriteCloser
}

var strangelove = Movie{
	Title:    "Dr. Strangelove",
	Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
	Year:     1964,
	Color:    false,
	Actor: map[string]string{
		"Dr. Strangelove":            "Peter Sellers",
		"Grp. Capt. Lionel Mandrake": "Peter Sellers",
		"Pres. Merkin Muffley":       "Peter Sellers",
		"Gen. Buck Turgidson":        "George C. Scott",
		"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
		`Maj. T.J. "King" Kong`:      "Slim Pickens",
	},
	Oscars: []string{
		"Best Actor (Nomin.)",
		"Best Adapted Screenplay (Nomin.)",
		"Best Director (Nomin.)",
		"Best Picture (Nomin.)",
	},
	Score: 7.8,
	File:     os.Stdout,
}

func main() {
	xs, err := Marshal(strangelove)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(xs))
}
