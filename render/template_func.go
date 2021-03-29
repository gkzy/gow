package render

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"time"
)

func init() {
	templateFuncMap["str2html"] = Str2html
	templateFuncMap["html2str"] = HTML2str
	templateFuncMap["datetimeformat"] = DateTimeFormat
	templateFuncMap["date"] = DateFormat
	templateFuncMap["int_datetimeformat"] = IntDateTimeFormat
	templateFuncMap["int_datetime"] = IntDateTime
	templateFuncMap["int_date"] = IntDate
	templateFuncMap["substr"] = Substr
	templateFuncMap["assets_js"] = AssetsJs
	templateFuncMap["assets_css"] = AssetsCSS
}

//Substr Substr
func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
	}
	return string(bt[start:end])
}

// Str2html str2 to html code
func Str2html(str string) template.HTML {
	return template.HTML(str)
}

// HTML2str
func HTML2str(html string) string {
	re := regexp.MustCompile(`\<[\S\s]+?\>`)
	html = re.ReplaceAllStringFunc(html, strings.ToLower)

	//remove STYLE
	re = regexp.MustCompile(`\<style[\S\s]+?\</style\>`)
	html = re.ReplaceAllString(html, "")

	//remove SCRIPT
	re = regexp.MustCompile(`\<script[\S\s]+?\</script\>`)
	html = re.ReplaceAllString(html, "")

	re = regexp.MustCompile(`\<[\S\s]+?\>`)
	html = re.ReplaceAllString(html, "\n")

	re = regexp.MustCompile(`\s{2,}`)
	html = re.ReplaceAllString(html, "\n")

	return strings.TrimSpace(html)
}

// IntDateTime return datetime string
// format ="YYYY-MM-DD HH:mm:ss"
func IntDateTime(val int64) (ret string) {
	if val < 1 {
		return
	}
	ret = IntDateTimeFormat(val, "YYYY-MM-DD HH:mm:ss")
	return
}

// IntDate 时间戳的日期格式化
func IntDate(val int64) (ret string) {
	if val < 1 {
		return
	}
	ret = IntDateTimeFormat(val, "YYYY-MM-DD")
	return
}

//DateFormat 日期格式化
func DateFormat(t time.Time) string {
	return DateTimeFormat(t, "YYYY-MM-DD")
}

// DateTimeFormat return datetime string
// format ="YYYY-MM-DD HH:mm:ss" or format ="YYYY-MM-DD HH:mm"
func IntDateTimeFormat(val int64, format string) (ret string) {
	if val < 1 {
		return
	}
	tm := time.Unix(val, 0)
	ret = DateTimeFormat(tm, format)
	return
}

// DateTimeFormat 格式化时间显示
// format = "YYYY-MM-DD HH:mm:ss" or "YYYY年MM月DD日"
func DateTimeFormat(t time.Time, format string) string {
	res := strings.Replace(format, "MM", t.Format("01"), -1)
	res = strings.Replace(res, "M", t.Format("1"), -1)
	res = strings.Replace(res, "DD", t.Format("02"), -1)
	res = strings.Replace(res, "D", t.Format("2"), -1)
	res = strings.Replace(res, "YYYY", t.Format("2006"), -1)
	res = strings.Replace(res, "YY", t.Format("06"), -1)
	res = strings.Replace(res, "HH", fmt.Sprintf("%02d", t.Hour()), -1)
	res = strings.Replace(res, "H", fmt.Sprintf("%d", t.Hour()), -1)
	res = strings.Replace(res, "hh", t.Format("03"), -1)
	res = strings.Replace(res, "h", t.Format("3"), -1)
	res = strings.Replace(res, "mm", t.Format("04"), -1)
	res = strings.Replace(res, "m", t.Format("4"), -1)
	res = strings.Replace(res, "ss", t.Format("05"), -1)
	res = strings.Replace(res, "s", t.Format("5"), -1)
	return res
}

// AssetsJs returns script tag with src string.
func AssetsJs(text string) template.HTML {
	text = "<script src=\"" + text + "\"></script>"
	return template.HTML(text)
}

// AssetsCSS returns stylesheet link tag with src string.
func AssetsCSS(text string) template.HTML {
	text = "<link href=\"" + text + "\" rel=\"stylesheet\" />"
	return template.HTML(text)
}
