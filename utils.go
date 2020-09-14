package gow

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

//getCurrentDirectory
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		debugPrintError(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Bind is a helper function for given interface object and returns a Gin middleware.
//func Bind(val interface{}) HandlerFunc {
//	value := reflect.ValueOf(val)
//	if value.Kind() == reflect.Ptr {
//		panic(`Bind struct can not be a pointer. Example:
//	Use: gin.Bind(Struct{}) instead of gin.Bind(&Struct{})
//`)
//	}
//	typ := value.Type()
//
//	return func(c *Context) {
//		obj := reflect.New(typ).Interface()
//		if c.Bind(obj) == nil {
//			c。S(BindKey, obj)
//		}
//	}
//}

// WrapF is a helper function for wrapping http.HandlerFunc and returns a Gin middleware.
func WrapF(f http.HandlerFunc) HandlerFunc {
	return func(c *Context) {
		f(c.Writer, c.Request)
	}
}

// WrapH is a helper function for wrapping http.Handler and returns a Gin middleware.
func WrapH(h http.Handler) HandlerFunc {
	return func(c *Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// MarshalXML allows type H to be used with xml.Marshal.
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

func chooseData(custom, wildcard interface{}) interface{} {
	if custom != nil {
		return custom
	}
	if wildcard != nil {
		return wildcard
	}
	panic("negotiation config is invalid")
}

func parseAccept(acceptHeader string) []string {
	parts := strings.Split(acceptHeader, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(strings.Split(part, ";")[0]); part != "" {
			out = append(out, part)
		}
	}
	return out
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// getAddress
func (engine *Engine) getAddress(args ...interface{}) string {
	var (
		host string
		port int
	)
	switch len(args) {
	case 0:
		host, port = getHostAndPort(engine.httpAddr)
	case 1:
		switch arg := args[0].(type) {
		case string:
			host, port = getHostAndPort(args[0].(string))
		case int:
			port = arg
		}
	case 2:
		if arg, ok := args[0].(string); ok {
			host = arg
		}
		if arg, ok := args[1].(int); ok {
			port = arg
		}
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	return addr
}

// getHostAndPort getHostAndPort
func getHostAndPort(addr string) (host string, port int) {
	addrs := strings.Split(addr, ":")
	if len(addrs) == 1 {
		host = ""
		port, _ = strconv.Atoi(addrs[0])
	} else if len(addrs) >= 2 {
		host = addrs[0]
		port, _ = strconv.Atoi(addrs[1])
	}

	return
}
