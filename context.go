package gow

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gkzy/gow/render"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
)

// Context gow context
type Context struct {
	writermem responseWriter
	Request   *http.Request
	Writer    ResponseWriter

	Params   Params
	handlers HandlersChain
	index    int8
	fullPath string

	engine *Engine
	params *Params

	// This mutex protect Keys map
	mu sync.RWMutex

	// Keys is a key/value pair exclusively for the context of each request.
	Keys map[string]interface{}

	Data map[interface{}]interface{}

	// Errors is a list of errors attached to all the handlers/middlewares who used this context.
	Errors errorMsgs

	// Accepted defines a list of manually accepted formats for content negotiation.
	Accepted []string

	// queryCache use url.ParseQuery cached the param query result from c.Request.URL.Query()
	queryCache url.Values

	// formCache use url.ParseQuery cached PostForm contains the parsed form data from POST, PATCH,
	// or PUT body parameters.
	formCache url.Values

	// SameSite allows a server to define a cookie attribute making it impossible for
	// the browser to send this cookie along with cross-site requests.
	sameSite http.SameSite

	//Pager
	Pager *Pager
}

const (
	abortIndex int8 = math.MaxInt8 / 2
)

func (c *Context) reset() {
	c.Writer = &c.writermem
	c.Params = c.Params[0:0]
	c.handlers = nil
	c.index = -1
	c.fullPath = ""
	c.Keys = nil
	c.Errors = c.Errors[0:0]
	c.Accepted = nil
	c.queryCache = nil
	c.formCache = nil
	c.Data = make(map[interface{}]interface{}, 0)
	c.Pager = nil
	*c.params = (*c.params)[0:0]
}

// Handler returns the main handler.
func (c *Context) Handler() HandlerFunc {
	return c.handlers.Last()
}

// HandlerName last handler name
func (c *Context) HandlerName() string {
	return nameOfFunction(c.handlers.Last())
}

// FullPath returns a matched route full path. For not found routes
// returns an empty string.
//     router.GET("/user/:id", func(c *gin.Context) {
//         c.FullPath() == "/user/:id" // true
//     })
func (c *Context) FullPath() string {
	return c.fullPath
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// IsProd return bool
//	是否运行在生产环境下
func (c *Context) IsProd() bool {
	return c.engine.RunMode == ProdMode
}

// GetRunMode return app run mode string
//	return dev or prod
func (c *Context) GetRunMode() string {
	return c.engine.RunMode
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort abort handler
func (c *Context) Abort() {
	c.index = abortIndex
}

// StopRun
func (c *Context) StopRun() {
	panic(stopRun)
}

func (c *Context) Error(err error) *Error {
	if err == nil {
		panic("err is nil")
	}

	parsedError, ok := err.(*Error)
	if !ok {
		parsedError = &Error{
			Err:  err,
			Type: ErrorTypePrivate,
		}
	}

	c.Errors = append(c.Errors, parsedError)
	return parsedError
}

/************************************/
/************ INPUT DATA ************/
/************************************/

// IsWeChat return is WeChat user-agent
func (c *Context) IsWeChat() bool {
	return strings.Contains(strings.ToLower(c.UserAgent()), strings.ToLower("MicroMessenger"))
}

// IsAjax return is ajax request
func (c *Context) IsAjax() bool {
	return c.GetHeader("X-Requested-With") == "XMLHttpRequest"
}

// IsWebsocket return is websocket request
func (c *Context) IsWebsocket() bool {
	if strings.Contains(strings.ToLower(c.GetHeader("Connection")), "upgrade") &&
		strings.EqualFold(c.GetHeader("Upgrade"), "websocket") {
		return true
	}
	return false
}

// Referer return request referer
func (c *Context) Referer() string {
	return c.Request.Referer()
}

// Host return request host string
func (c *Context) Host() string {
	return c.Request.Host
}

// Param returns the value of the URL param.
// It is a shortcut for c.Params.ByName(key)
//     router.GET("/user/{id}, func(c *gin.Context) {
//         // a GET request to /user/john
//         id := c.Param("id") // id == "john"
//     })
func (c *Context) Param(key string) string {
	return c.Params.ByName(key)
}

// ParamInt  return the value of the URL param
func (c *Context) ParamInt(key string) (int, error) {
	v := c.Param(key)
	return strconv.Atoi(v)
}

//  ParamInt64  return the value of the URL param
func (c *Context) ParamInt64(key string) (int64, error) {
	v := c.Param(key)
	return strconv.ParseInt(v, 10, 64)
}

// UserAgent get useragent
func (c *Context) UserAgent() string {
	return c.GetHeader("User-Agent")
}

// Query return query string
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// Form return request.FormValue key
func (c *Context) Form(key string) string {
	return c.Request.FormValue(key)
}

// input
func (c *Context) input() url.Values {
	if c.Request.Form == nil {
		c.Request.ParseForm()
	}
	return c.Request.Form
}

// formValue formValue
func (c *Context) formValue(key string) string {
	if v := c.Form(key); v != "" {
		return v
	}
	if c.Request.Form == nil {
		c.Request.ParseForm()
	}
	return c.Request.Form.Get(key)
}

// GetString 按key返回字串值，可以设置default值
func (c *Context) GetString(key string, def ...string) string {
	if v := c.formValue(key); v != "" {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// GetStrings GetStrings
func (c *Context) GetStrings(key string, def ...[]string) []string {
	var defaultDef []string
	if len(def) > 0 {
		defaultDef = def[0]
	}

	if v := c.input(); v == nil {
		return defaultDef
	} else if kv := v[key]; len(kv) > 0 {
		return kv
	}
	return defaultDef
}

// GetInt return int
func (c *Context) GetInt(key string, def ...int) (int, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.Atoi(v)
}

// GetInt8 GetInt8
//	-128~127
func (c *Context) GetInt8(key string, def ...int8) (int8, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseInt(v, 10, 8)
	return int8(i64), err
}

//GetUint8 GetUint8
//	0~255
func (c *Context) GetUint8(key string, def ...uint8) (uint8, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseUint(v, 10, 8)
	return uint8(i64), err
}

// GetInt16 GetInt16
//	-32768~32767
func (c *Context) GetInt16(key string, def ...int16) (int16, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseInt(v, 10, 16)
	return int16(i64), err
}

// GetUint16 GetUint16
//	0~65535
func (c *Context) GetUint16(key string, def ...uint16) (uint16, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseUint(v, 10, 16)
	return uint16(i64), err
}

//GetInt32 GetInt32
//	-2147483648~2147483647
func (c *Context) GetInt32(key string, def ...int32) (int32, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseInt(v, 10, 32)
	return int32(i64), err
}

// GetUint32 GetUint32
//	0~4294967295
func (c *Context) GetUint32(key string, def ...uint32) (uint32, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseUint(v, 10, 32)
	return uint32(i64), err
}

// GetInt64 GetInt64
//	-9223372036854775808~9223372036854775807
func (c *Context) GetInt64(key string, def ...int64) (int64, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseInt(v, 10, 64)
}

// GetUint64 GetUint64
//	0~18446744073709551615
func (c *Context) GetUint64(key string, def ...uint64) (uint64, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseUint(v, 10, 64)
	return uint64(i64), err
}

// GetFloat64 GetFloat64
func (c *Context) GetFloat64(key string, def ...float64) (float64, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseFloat(v, 64)
}

//GetBool GetBool
func (c *Context) GetBool(key string, def ...bool) (bool, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseBool(v)
}

/************************************/
/******** UPLOAD********/
/************************************/

// GetFile get single file from request
func (c *Context) GetFile(key string) (multipart.File, *multipart.FileHeader, error) {
	if c.Request.MultipartForm == nil {
		if err := c.Request.ParseMultipartForm(c.engine.MaxMultipartMemory); err != nil {
			return nil, nil, err
		}
	}
	return c.Request.FormFile(key)
}

// GetFiles get files from request
func (c *Context) GetFiles(key string) ([]*multipart.FileHeader, error) {
	if files, ok := c.Request.MultipartForm.File[key]; ok {
		return files, nil
	}
	return nil, http.ErrMissingFile
}

// SaveToFile saves uploaded file to new path.
//	upload the file and save it on the server
//	c.SaveToFile("file","./upload/1.jpg")
func (c *Context) SaveToFile(fromFile, toFile string) error {
	file, _, err := c.Request.FormFile(fromFile)
	if err != nil {
		return err
	}
	defer file.Close()
	f, err := os.OpenFile(toFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, file)
	return nil
}

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// SetKey is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) SetKey(key string, value interface{}) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}

	c.Keys[key] = value
	c.mu.Unlock()
}

// GetKey returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) GetKey(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.GetKey(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// KeyString returns the value associated with the key as a string.
func (c *Context) KeyString(key string) (s string) {
	if val, ok := c.GetKey(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// KeyBool returns the value associated with the key as a boolean.
func (c *Context) KeyBool(key string) (b bool) {
	if val, ok := c.GetKey(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// KeyInt returns the value associated with the key as an integer.
func (c *Context) KeyInt(key string) (i int) {
	if val, ok := c.GetKey(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// KeyInt64 returns the value associated with the key as an integer.
func (c *Context) KeyInt64(key string) (i64 int64) {
	if val, ok := c.GetKey(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// KeyUint returns the value associated with the key as an unsigned integer.
func (c *Context) KeyUint(key string) (ui uint) {
	if val, ok := c.GetKey(key); ok && val != nil {
		ui, _ = val.(uint)
	}
	return
}

// KeyUint64 returns the value associated with the key as an unsigned integer.
func (c *Context) KeyUint64(key string) (ui64 uint64) {
	if val, ok := c.GetKey(key); ok && val != nil {
		ui64, _ = val.(uint64)
	}
	return
}

// KeyFloat64 returns the value associated with the key as a float64.
func (c *Context) KeyFloat64(key string) (f64 float64) {
	if val, ok := c.GetKey(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// KeyTime returns the value associated with the key as time.
func (c *Context) KeyTime(key string) (t time.Time) {
	if val, ok := c.GetKey(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// KeyDuration returns the value associated with the key as a duration.
func (c *Context) KeyDuration(key string) (d time.Duration) {
	if val, ok := c.GetKey(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// KeyStringSlice returns the value associated with the key as a slice of strings.
func (c *Context) KeyStringSlice(key string) (ss []string) {
	if val, ok := c.GetKey(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// KeyStringMap returns the value associated with the key as a map of interfaces.
func (c *Context) KeyStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.GetKey(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

// KeyStringMapString returns the value associated with the key as a map of strings.
func (c *Context) KeyStringMapString(key string) (sms map[string]string) {
	if val, ok := c.GetKey(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// KeyStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Context) KeyStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.GetKey(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}

/************************************/
/******** RESPONSE RENDERING ********/
/************************************/

// Header Header
func (c *Context) Header(key, value string) {
	if value == "" {
		c.Writer.Header().Del(key)
		return
	}
	c.Writer.Header().Set(key, value)
}

// GetHeader returns value from request headers.
func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

// SetSameSite with cookie
func (c *Context) SetSameSite(samesite http.SameSite) {
	c.sameSite = samesite
}

// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: c.sameSite,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

// GetCookie returns the named cookie provided in the request or
// ErrNoCookie if not found. And return the named cookie is unescaped.
// If multiple cookies match the given name, only one cookie will
// be returned.
func (c *Context) GetCookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

// Redirect http redirect
// like : 301 302 ...
func (c *Context) Redirect(code int, url string) {
	c.Writer.WriteHeader(code)
	http.Redirect(c.Writer, c.Request, url, code)
}

// Body request body
func (c *Context) Body() []byte {
	if c.Request.Body == nil {
		return []byte{}
	}
	var body []byte
	body, _ = ioutil.ReadAll(c.Request.Body)

	c.Request.Body.Close()
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body
}

// Render render html
func (c *Context) Render(statusCode int, r render.Render) {
	c.Writer.WriteHeader(statusCode)
	if !bodyAllowedForStatus(statusCode) {
		r.WriteContentType(c.Writer)
		c.Writer.WriteHeader(statusCode)
		return
	}
	if err := r.Render(c.Writer); err != nil {
		fmt.Println(err)
	}
}

// ServerHTML render html page
func (c *Context) ServerHTML(statusCode int, name string, data ...interface{}) {
	//未设置 AutoRender时，不渲染模板
	if !c.engine.AutoRender {
		c.ServerString(404, string(default404Body))
		return
	}
	var v interface{}
	if len(data) > 0 {
		v = data[0]
	} else {
		v = c.Data
	}
	render := render.HTMLRender{}.Instance(c.engine.viewsPath, name, c.engine.FuncMap, c.engine.delims, c.engine.AutoRender, c.engine.RunMode, v)
	c.Render(statusCode, render)
}

// HTML render html page
// 		When inputting data, use the value of data, otherwise use c.Data
func (c *Context) HTML(name string, data ...interface{}) {
	if len(data) > 0 {
		v := data[0]
		c.ServerHTML(http.StatusOK, name, v)
		return
	}
	c.ServerHTML(http.StatusOK, name)
}

// ServerString write string into the response body
func (c *Context) ServerString(code int, msg string) {
	if code < 0 {
		code = http.StatusOK
	}
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Status(code)
	c.Writer.Write([]byte(msg))
}

// String write string into the response
func (c *Context) String(msg string) {
	c.ServerString(http.StatusOK, msg)
}

// ServerYAML serializes the given struct as YAML into the response body.
func (c *Context) ServerYAML(code int, data interface{}) {
	if code < 0 {
		code = http.StatusOK
	}
	c.Header("Content-Type", "application/x-yaml; charset=utf-8")
	c.Status(code)

	bytes, err := yaml.Marshal(data)
	if err != nil {
		c.Header("Content-Type", "")
		c.ServerString(http.StatusServiceUnavailable, err.Error())
	}
	c.Writer.Write(bytes)
}

// YAML serializes the given struct as YAML into the response body
func (c *Context) YAML(data interface{}) {
	c.ServerYAML(http.StatusOK, data)
}

// ServerJSON serializes the given struct as JSON into the response body.
func (c *Context) ServerJSON(code int, data interface{}) {
	if code < 0 {
		code = http.StatusOK
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if c.engine.RunMode == DevMode {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(data); err != nil {
		c.Header("Content-Type", "")
		c.ServerString(http.StatusServiceUnavailable, err.Error())
	}
}

// JSON serializes the given struct as JSON into the response body
func (c *Context) JSON(data interface{}) {
	c.ServerJSON(http.StatusOK, data)
}

// ServerJSONP write data by jsonp format
func (c *Context) ServerJSONP(code int, callback string, data interface{}) {
	if code < 0 {
		code = http.StatusOK
	}
	c.Header("Content-Type", "application/javascript; charset=utf-8")
	c.Status(code)

	bytes, err := json.Marshal(data)
	if err != nil {
		c.Header("Content-Type", "")
		c.ServerString(http.StatusServiceUnavailable, err.Error())
	}
	c.Writer.Write([]byte(callback + "{"))
	c.Writer.Write(bytes)
	c.Writer.Write([]byte(");"))
}

// JSONP write date by jsonp format
func (c *Context) JSONP(callback string, data interface{}) {
	c.ServerJSONP(http.StatusOK, callback, data)
}

// ServerXML write data by xml format
func (c *Context) ServerXML(code int, data interface{}) {
	if code < 0 {
		code = http.StatusOK
	}
	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.Status(code)
	encoder := xml.NewEncoder(c.Writer)
	if err := encoder.Encode(data); err != nil {
		c.Header("Content-Type", "")
		c.ServerString(http.StatusServiceUnavailable, err.Error())
	}
}

// XML response xml doc
func (c *Context) XML(data interface{}) {
	c.ServerXML(http.StatusOK, data)
}

// GetIP return k8s Cluster ip
//	if ip=="" return "10.10.10.2"
func (c *Context) GetIP() (ip string) {
	//H5服务器端返传递的IP
	ip = c.GetHeader("ip")
	if ip == "" {
		ip = c.GetHeader("X-Original-Forwarded-For")
	}
	if ip == "" {
		ip = c.GetHeader("Remote-Host")
	}
	if ip == "" {
		ip = c.GetHeader("X-Real-IP")
	}
	if ip == "" {
		ip = c.ClientIP()
	}
	if ip == "" {
		ip = "10.10.10.2"
	}
	return ip
}

// ClientIP get client ip address
func (c *Context) ClientIP() (ip string) {
	addr := c.Request.RemoteAddr
	str := strings.Split(addr, ":")
	if len(str) > 1 {
		ip = str[0]
	}
	return
}

// Status sets the HTTP response code.
func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

// AbortWithStatus calls `Abort()` and writes the headers with the specified status code.
// For example, a failed attempt to authenticate a request could use: context.AbortWithStatus(401).
func (c *Context) AbortWithStatus(code int) {
	c.Status(code)
	c.Writer.WriteHeaderNow()
	c.Abort()
}

// File writes the specified file into the body stream in a efficient way.
func (c *Context) File(filepath string) {
	http.ServeFile(c.Writer, c.Request, filepath)
}

// FileFromFS writes the specified file from http.FileSytem into the body stream in an efficient way.
func (c *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		c.Request.URL.Path = old
	}(c.Request.URL.Path)

	c.Request.URL.Path = filepath

	http.FileServer(fs).ServeHTTP(c.Writer, c.Request)
}

// FileAttachment writes the specified file into the body stream in an efficient way
// On the client side, the file will typically be downloaded with the given filename
func (c *Context) FileAttachment(filepath, filename string) {
	c.Header("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	http.ServeFile(c.Writer, c.Request, filepath)
}

// Download download data
func (c *Context) Download(data []byte) {
	c.Header("Content-Type", "application/octet-stream; charset=utf-8")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(data)
}

//============private method=============

func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}
