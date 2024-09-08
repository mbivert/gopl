package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"io"
	"reflect"
)

func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")

	case reflect.Bool:
		if v.Bool() {
			fmt.Fprintf(buf, "t")
		} else {
			fmt.Fprintf(buf, "nil")
		}

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())

	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%f", v.Float())

	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		fmt.Fprintf(buf, "#C(%f, %f)", real(c), imag(c))

	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())

	case reflect.Ptr:
		return encode(buf, v.Elem())

	case reflect.Array, reflect.Slice: // (value ...)
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')

	case reflect.Struct: // ((name value) ...)
		buf.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')

	case reflect.Map: // ((key value) ...)
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {

			if i > 0 {
				buf.WriteByte(' ')

			}
			buf.WriteByte('(')
			if err := encode(buf, key); err != nil {

				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key)); err != nil {

				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')

	case reflect.Interface: // ("type" value)
		fmt.Fprintf(buf, "(%q ", v.Type().String())
		err := encode(buf, reflect.ValueOf(v.Interface()))
		if err != nil {
			return err
		}
		buf.WriteByte(')')

	case reflect.UnsafePointer:
		fmt.Fprintf(buf, "%p", v.UnsafePointer())

	default: // chan, func
		return fmt.Errorf("unsupported type: %s", v.Type())

	}
	return nil
}

func prettyPrint(buf *bytes.Buffer, v reflect.Value, indent string) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")

	case reflect.Bool:
		if v.Bool() {
			fmt.Fprintf(buf, "t")
		} else {
			fmt.Fprintf(buf, "nil")
		}

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())

	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%f", v.Float())

	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		fmt.Fprintf(buf, "#C(%f, %f)", real(c), imag(c))

	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())

	case reflect.Ptr:
		return prettyPrint(buf, v.Elem(), indent)

	case reflect.Array, reflect.Slice: // (value ...)
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteString("\n"+indent+" ")
			}
			if err := prettyPrint(buf, v.Index(i), indent); err != nil {
				return err
			}
		}
		buf.WriteByte(')')

	case reflect.Struct: // ((name value) ...)
		buf.WriteString("(")
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteString("\n"+indent+" ")
			}
			fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
			ind := indent+" "
			for n := 0; n < len(v.Type().Field(i).Name); n++ {
				ind += " "
			}
			if err := prettyPrint(buf, v.Field(i), ind+" "); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')

	case reflect.Map: // ((key value) ...)
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {

			if i > 0 {
				buf.WriteString("\n"+indent+" ")

			}
			buf.WriteByte('(')
			if err := prettyPrint(buf, key, indent); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := prettyPrint(buf, v.MapIndex(key), indent); err != nil {

				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')

	case reflect.Interface: // ("type" value)
		fmt.Fprintf(buf, "(%q ", v.Type().String())
		ind := indent
		for n := 0; n < len(v.Type().String()); n++ {
			ind += " "
		}
		err := prettyPrint(buf, reflect.ValueOf(v.Interface()), ind)
		if err != nil {
			return err
		}
		buf.WriteByte(')')

	case reflect.UnsafePointer:
		fmt.Fprintf(buf, "%p", v.UnsafePointer())

	default: // chan, func
		return fmt.Errorf("unsupported type: %s (%v)", v.Type(), v)

	}
	return nil
}

// Marshal encodes a Go value in S-expression form.
func ppMarshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := prettyPrint(&buf, reflect.ValueOf(v), ""); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
	Rotation        complex128
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
	// Okay.
	Rotation: 2 + 1i,
	File:     os.Stdout,
}

func main() {
	xs, err := ppMarshal(strangelove)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(xs))
}
