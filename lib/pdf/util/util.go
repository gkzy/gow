package util

import (
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

func AddComma(s string) string {
	if strings.Index(s, ".") == -1 {
		return addCommaSub(s)
	}
	ss := strings.Split(s, ".")
	ss[0] = addCommaSub(ss[0])
	return ss[0] + "." + ss[1]
}

func addCommaSub(s string) string {
	res := ""
	if len(s) < 4 {
		return s
	}
	pos := len(s) % 3
	if pos > 0 {
		res += s[0:pos] + ","
	}
	for i := pos; i < len(s); i += 3 {
		res += s[i : i+3]
		//fmt.Printf("pos %v \n", i)
		if i < len(s)-3 {
			res += ","
		}
	}
	return res
}

func ReadTextFile(filename string, colno int) []interface{} {
	res, _ := ioutil.ReadFile(filename)
	lines := strings.Split(string(res), "\n")
	list := make([]interface{}, 0, 100)
	for _, line := range lines {
		line = strings.Replace(line, "\r", "", -1)
		cols := strings.Split(line, "\t")
		if len(cols) < colno {
			continue
		}
		list = append(list, cols)
	}
	return list
}

func Ftoa(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func Btoa(b uint8) string {
	val := int(b)
	return strconv.Itoa(val)
}

func IsEmpty(object interface{}) bool {
	if object == nil {
		return true
	}

	objValue := reflect.ValueOf(object)
	switch objValue.Kind() {
	// collection types are empty when they have no element
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		deref := objValue.Elem().Interface()
		return IsEmpty(deref)
		// for all other types, compare against the zero value
	default:
		zero := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(object, zero.Interface())
	}
}

func CheckColor(color string) string {
	color = strings.Replace(color, " ", "", -1)
	rgb := strings.Split(color, ",")
	if len(rgb) != 3 {
		panic("the color err")
	}

	for i := range rgb {
		value, err := strconv.Atoi(rgb[i])
		if err != nil {
			panic(err)
		}
		if value < 0 || value > 255 {
			panic("the R,G,B value error")
		}
	}

	return color
}

func GetColorRGB(color string) (r, g, b int) {
	color = CheckColor(color)
	rgb := strings.Split(color, ",")
	return Atoi(rgb[0]), Atoi(rgb[1]), Atoi(rgb[2])
}

func Atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
