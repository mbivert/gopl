package main

import (
	"fmt"
	"reflect"
	"strconv"
)

const maxDepth = 42

func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	fmt.Printf(display(name, reflect.ValueOf(x), "", 0))
}

// formatAtom formats a value without inspecting its internal structure.
func formatAtom(v reflect.Value, indent string, depth int) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	case reflect.Struct:
		fallthrough
	case reflect.Array:
		return display(v.Type().String(), v, indent+"\t", depth+1)
	default: // reflect.Interface
		return v.Type().String() + " value"
	}
}

func display(path string, v reflect.Value, indent string, depth int) string {
	if depth >= maxDepth {
		return "...\n"
	}

	switch v.Kind() {
	case reflect.Invalid:
		return fmt.Sprintf("%s%s = invalid\n", path, indent)
	case reflect.Slice, reflect.Array:
		s := ""
		for i := 0; i < v.Len(); i++ {
			s += display(fmt.Sprintf("%s%s[%d]", indent, path, i), v.Index(i), indent, depth+1)
		}
		return s
	case reflect.Struct:
		s := ""
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s%s.%s", indent, path, v.Type().Field(i).Name)
			s += display(fieldPath, v.Field(i), indent, depth+1)
		}
		return s
	case reflect.Map:
		s := ""
		for _, key := range v.MapKeys() {
			// meh
			nl := ""
			if key.Kind() == reflect.Struct || key.Kind() == reflect.Array {
				nl = "\n"
			}
			s += display(fmt.Sprintf("%s%s[%s%s]", indent, path, nl,
				formatAtom(key, indent, depth)), v.MapIndex(key), indent, depth+1)
		}
		return s
	case reflect.Ptr:
		if v.IsNil() {
			return fmt.Sprintf("%s%s = nil\n", indent, path)
		} else {
			return display(fmt.Sprintf("%s(*%s)", indent, path), v.Elem(), indent, depth+1)
		}
	case reflect.Interface:
		if v.IsNil() {
			return fmt.Sprintf("%s%s = nil\n", indent, path)
		} else {
			s := fmt.Sprintf("%s%s.type = %s\n", indent, path, v.Elem().Type())
			return s + display(path+".value", v.Elem(), indent, depth+1)
		}
	default: // basic types, channels, funcs
		return fmt.Sprintf("%s%s = %s\n", indent, path, formatAtom(v, indent, depth))
	}
}

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
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
}

type Person struct {
	Name string
	Age  int
}

// a struct that points to itself
type Cycle struct {
	Value int
	Tail  *Cycle
}

// We can't recycle Movie here, as struct{} may only
// be used as key when they are comparable, i.e. its
// fields must be "scalar" types or comparable structs
// (maps are the issue)
//
// We could have used *Movie, but that's a different exercise.
var known = map[Person]bool{
	Person{"Peter Sellers", 54}:              false,
	Person{"Dennis MacAlistair Ritchie", 70}: true,
}

// similarly, we can't use slices as keys, but arrays are okay.
var sums = map[[3]int]int{
	[3]int{1, 2, 3}: 1 + 2 + 3,
	[3]int{4, 5, 6}: 4 + 5 + 6,
}

func main() {
	Display("strangelove", strangelove)
	Display("known", known)
	Display("sums", sums)
	var c Cycle
	c = Cycle{42, &c}
	Display("c", c)
}

