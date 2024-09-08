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

// Note that
//	v.Field(i) == reflect.Zero(v.Field(i).Type())
// and
//	reflect.DeepEqual(v.Field(i), reflect.Zero(v.Field(i).Type()))
// will fail. Doc states:
// « To compare two Values, compare the results of the Interface method.
// Using == on two Values does not compare the underlying values they
// represent. »
//
// However, .Interface() returns false for unexported fields
// (and panic if CanInterface() doesn't guard it:
//	panic: reflect.Value.Interface: cannot return value obtained from
//	unexported field or method)
//
// So, isZeroBasic() only works for public fields ...
func isZeroBasic(v reflect.Value) bool {
	zero := reflect.Zero(v.Type())

	if v.CanInterface() && zero.CanInterface() {
		if v.Comparable() && zero.Comparable() {
			if v.Interface() == zero.Interface() {
				// comment/uncomment to verify
				return true
			}
		}
	}

	return false
}

// ... but since we know (we're doing this already!) we can print
// private fields, we can perform a more  tedious test to get
// a more accurate version
//
// Seems that another option is to rely on unsafe
// (https://stackoverflow.com/a/17982725), but the authors will
// only present unsafe later in the book.
func isZeroBetter(v reflect.Value) bool {
	zero := reflect.Zero(v.Type())

	switch v.Kind() {
	case reflect.Bool:
		return v.Bool() == zero.Bool()

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return v.Int() == zero.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == zero.Uint()

	case reflect.Float32, reflect.Float64:
		return v.Float() == zero.Float()

	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == zero.Complex()

	case reflect.String:
		return v.String() == zero.String()

	case reflect.Array:
		return v.Len() == 0 // ??
	case reflect.Slice, reflect.Ptr, reflect.UnsafePointer:
		return v.IsNil()
//	case reflect.Struct: // hmm.
	}

	return false
}

// But this seems to be the real deal anyway; it even deals with
// nil struct, as demonstrated by the printed os.File
func isReallyZero(v reflect.Value) bool {
	if v.IsValid() {
		return v.IsZero()
	}
	return false
}

var isZero = isReallyZero

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
		n := 0
		for i := 0; i < v.NumField(); i++ {
			if isZero(v.Field(i)) {
				continue
			}

			if n > 0 {
				buf.WriteString("\n"+indent+" ")
			}
			n++
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
		n := 0
		for _, key := range v.MapKeys() {
			if isZero(v.MapIndex(key)) {
				continue
			}
			if n > 0 {
				buf.WriteString("\n"+indent+" ")
			}
			n++
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

type Encoder struct {
	w io.Writer
	escape bool // not implemented
	prefix, indent string
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w,
		true,
		"",
		"",
	}
}

// This is really rough, but should be still fit what's
// expected for this exercise anyway:
//
//	1. we probably would want prettyPrint/encode to take
//	an io.Writer instead of a bytes.Buffer
//
//	2. prettyPrint doesn't manage prefix, and probably doesn't
//	correctly handle indent so well either
//
//	3. HTML escaping is unmanaged.
//
//	4. Eventually, we may need/want a separate Indent function
//	instead of prettyPrint: json/encoding implements and
//	performs indentation in a separate step systematically.
func (enc *Encoder) Encode(v any) error {
	var err error
	var buf bytes.Buffer

	if enc.prefix != "" || enc.indent != "" {
		if err = prettyPrint(&buf, reflect.ValueOf(v), enc.indent); err != nil {
			return err
		}
	} else {
		if err := encode(&buf, reflect.ValueOf(v)); err != nil {
			return err
		}
	}
	_, err = buf.WriteTo(enc.w)
	return err
}

func (enc *Encoder) SetEscapeHTML(on bool) {
	enc.escape = on
}

func (enc *Encoder) SetIndent(prefix, indent string) {
	enc.prefix = prefix
	enc.indent = indent
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
		"foo":                         "",
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
	ys, erry := ppMarshal(strangelove)
	if erry != nil {
		log.Fatal(erry)
	}
	fmt.Println(string(ys))

	enc := NewEncoder(os.Stdout)
	enc.SetIndent("unused-prefix", "") // triggers prettyPrint
	enc.Encode(strangelove)
	enc.Encode(strangelove)
}
