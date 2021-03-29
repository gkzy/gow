package core

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"
)

type FontMap struct {
	FontName string
	FileName string
}

// 对接 pdf
type Converter struct {
	pdf         *gopdf.GoPdf // 第三方转换
	atomicCells []string     // 原子单元, 多个单元格最终汇总成PDF文件

	unit  float64    // 单位像素
	fonts []*FontMap // 字体

	linew    float64 // 线宽度(辅助)
	lastFont string  // 最近字体(辅助)
}

// var convert.unit float64 = 2.834645669

// 获取AtomicCells
func (convert *Converter) GetAutomicCells() []string {
	cells := make([]string, len(convert.atomicCells))
	copy(cells, convert.atomicCells)
	return cells
}

// 设置AtomicCells(小心使用)
func (convert *Converter) SetAutomicCells(cells []string) {
	convert.atomicCells = cells
}

// 添加AtomicCell
func (convert *Converter) AddAtomicCell(cell string) {
	if strings.HasPrefix(cell, "F|") {
		if cell == convert.lastFont {
			return
		}

		convert.lastFont = cell
	}

	convert.atomicCells = append(convert.atomicCells, cell)
}

// 从保存的文件解析AtomicCell
func (convert *Converter) ReadFile(fileName string) error {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	text := strings.Replace(string(buf), "\r", "", -1)
	var UTF8_BOM = []byte{239, 187, 191}
	if text[0:3] == string(UTF8_BOM) {
		text = text[3:]
	}
	convert.atomicCells = strings.Split(text, "\n")
	return nil
}

// 将 Text(写入的, 而不是读入的) -> PDF文件
func (convert *Converter) Execute() {
	lines := convert.atomicCells
	for _, line := range lines {
		elements := strings.Split(line, "|")
		switch elements[0] {
		case "P":
			convert.Page(line, elements) // PDF开始, P是标准的格式, P1属于自定义的格式
		case "NP":
			convert.NewPage(line, elements) // 新页面
		case "F":
			convert.Font(line, elements) // 字体
		case "TC":
			convert.TextColor(line, elements) // 文本颜色
		case "LC":
			convert.LineColor(line, elements) // 线条颜色
		case "BC":
			convert.BackgroundColor(line, elements) // 背景颜色
		case "GF", "GS":
			convert.Grey(line, elements)
		case "C", "CL", "CR":
			convert.Cell(line, elements) // 单元格内容
		case "L", "LV", "LH", "LT":
			convert.Line(line, elements) // 行
		case "R":
			convert.Rect(line, elements) // 长方形
		case "O":
			convert.Oval(line, elements) // 椭圆
		case "I":
			convert.Image(line, elements) // 图片
		case "M":
			convert.Margin(line, elements)
		case "EL":
			convert.ExternalLink(line, elements)
		case "ILA":
			convert.InternalLinkAnchor(line, elements)
		case "ILL":
			convert.InternalLinkLink(line, elements)
		default:
			if len(line) > 0 && line[0:1] != "v" {
				fmt.Println("skip:" + line + ":")
			}
		}
	}
}

// 添加字体
func (convert *Converter) AddFont() {
	for _, font := range convert.fonts {
		err := convert.pdf.AddTTFFont(font.FontName, font.FileName)
		if err != nil {
			panic("font file:" + font.FileName + " not found")
		}
	}
}

// PDF文件页面的开始
// [P, mm|pt|in, A4, P|L]
// mm|pt|in 表示的尺寸单位, 毫米,像素,英尺
// P|L 表示Portait, Landscape, 表示布局
func (convert *Converter) Page(line string, elements []string) {
	convert.pdf = new(gopdf.GoPdf)

	checkLength(line, elements, 4)
	switch elements[2] {
	/* A0 ~ A5 纸张像素表示
	'A0': [2383.94, 3370.39],
	'A1': [1683.78, 2383.94],
	'A2': [1190.55, 1683.78],
	'A3': [841.89, 1190.55],
	'A4': [595.28, 841.89],
	'A5': [419.53, 595.28],
	*/
	case "A3":
		config := defaultConfigs["A3"]
		convert.setunit(elements[1])
		if elements[3] == "P" {
			convert.Start(config.width, config.height) // 像素
		} else if elements[3] == "L" {
			convert.Start(config.height, config.width)
		} else {
			panic("Page Orientation accept P or L")
		}
	case "A4":
		config := defaultConfigs["A4"]
		convert.setunit(elements[1])
		if elements[3] == "P" {
			convert.Start(config.width, config.height) // 像素
		} else if elements[3] == "L" {
			convert.Start(config.height, config.width)
		} else {
			panic("Page Orientation accept P or L")
		}
	case "A5":
		config := defaultConfigs["A5"]
		convert.setunit(elements[1])
		if elements[3] == "P" {
			convert.Start(config.width, config.height) // 像素
		} else if elements[3] == "L" {
			convert.Start(config.height, config.width)
		} else {
			panic("Page Orientation accept P or L")
		}
	case "LTR":
		config := defaultConfigs["LTR"]
		convert.setunit(elements[1])
		if elements[3] == "P" {
			convert.Start(config.width, config.height) // 像素
		} else if elements[3] == "L" {
			convert.Start(config.height, config.width)
		} else {
			panic("Page Orientation accept P or L")
		}
	default:
		panic("This size not supported yet:" + elements[2])
	}
	convert.AddFont()
	convert.pdf.AddPage()
}

// 单位转换率设置, 基准的像素Pt
func (convert *Converter) setunit(unit string) {
	// 1mm ~ 2.8pt 1in ~ 72pt
	switch unit {
	case "mm":
		convert.unit = 2.834645669
	case "pt":
		convert.unit = 1
	case "in":
		convert.unit = 72
	default:
		panic("This unit is not specified :" + unit)
	}
}

// 构建新的页面
func (convert *Converter) NewPage(line string, elements []string) {
	convert.pdf.AddPage()
}

// 设置PDF文件基本信息(单位,页面大小)
func (convert *Converter) Start(w float64, h float64) {
	convert.pdf.Start(gopdf.Config{
		Unit:     gopdf.Unit_PT,
		PageSize: gopdf.Rect{W: w, H: h},
	}) // 595.28, 841.89 = A4
}

// 设置当前文本使用的字体
// ["", "family", "style", "size"]
// style: "" or "U", ("B", "I")(需要字体本身支持)
func (convert *Converter) Font(line string, elements []string) {
	checkLength(line, elements, 4)
	err := convert.pdf.SetFont(elements[1], elements[2], parseIntPanic(elements[3], line))
	if err != nil {
		panic(err.Error() + " line;" + line)
	}
}

// 设置笔画的灰度 | 设置填充的灰度
// ["GF|GS", grayScale]
// grayScale: 0.0 到 1.0
func (convert *Converter) Grey(line string, elements []string) {
	checkLength(line, elements, 2)
	if elements[0] == "GF" {
		convert.pdf.SetGrayFill(parseFloatPanic(elements[1], line))
	}
	if elements[0] == "GS" {
		convert.pdf.SetGrayStroke(parseFloatPanic(elements[1], line))
	}
}

// 文本颜色
// ["", R, G, B] // RGB文本颜色
func (convert *Converter) TextColor(line string, elements []string) {
	checkLength(line, elements, 4)
	convert.pdf.SetTextColor(uint8(parseIntPanic(elements[1], line)),
		uint8(parseIntPanic(elements[2], line)),
		uint8(parseIntPanic(elements[3], line)))
}

// 画笔颜色
// ["", R, G, B]
func (convert *Converter) LineColor(line string, elements []string) {
	checkLength(line, elements, 4)
	convert.pdf.SetStrokeColor(uint8(parseIntPanic(elements[1], line)),
		uint8(parseIntPanic(elements[2], line)),
		uint8(parseIntPanic(elements[3], line)))
}

func (convert *Converter) BackgroundColor(line string, elements []string) {
	checkLength(line, elements, 9)

	//convert.pdf.SetLineWidth(0)               // 宽带最小
	convert.pdf.SetStrokeColor(255, 255, 255) // 白色线条

	convert.pdf.SetFillColor(uint8(parseIntPanic(elements[5], line)),
		uint8(parseIntPanic(elements[6], line)),
		uint8(parseIntPanic(elements[7], line))) // 设置填充颜色

	convert.pdf.RectFromUpperLeftWithStyle(parseFloatPanic(elements[1], line)*convert.unit,
		parseFloatPanic(elements[2], line)*convert.unit,
		parseFloatPanic(elements[3], line)*convert.unit,
		parseFloatPanic(elements[4], line)*convert.unit, "F")

	convert.pdf.SetFillColor(0, 0, 0) // 颜色恢复
	convert.pdf.SetStrokeColor(0, 0, 0)

	convert.pdf.SetLineType("solid")

	x := parseFloatPanic(elements[1], line) * convert.unit
	y := parseFloatPanic(elements[2], line) * convert.unit
	w := parseFloatPanic(elements[3], line) * convert.unit
	h := parseFloatPanic(elements[4], line) * convert.unit

	lines := elements[8] //  LEFT,TOP,RIGHT,BOTTOM
	if lines[0] == '1' {
		convert.pdf.Line(x, y, x, y+h)
	}
	if lines[1] == '1' {
		convert.pdf.Line(x, y, x+w, y)
	}
	if lines[2] == '1' {
		convert.pdf.Line(x+w, y, x+w, y+h)
	}
	if lines[3] == '1' {
		convert.pdf.Line(x, y+h, x+w, y+h)
	}
}

// 椭圆
// ["", x1, y1, x2, y2]
func (convert *Converter) Oval(line string, elements []string) {
	checkLength(line, elements, 5)
	convert.pdf.Oval(parseFloatPanic(elements[1], line)*convert.unit,
		parseFloatPanic(elements[2], line)*convert.unit,
		parseFloatPanic(elements[3], line)*convert.unit,
		parseFloatPanic(elements[4], line)*convert.unit)
}

// 长方形
// ["R", x1, y1, x2, y2]
func (convert *Converter) Rect(line string, eles []string) {
	checkLength(line, eles, 5)
	adj := convert.linew * convert.unit * 0.5
	convert.pdf.Line(
		parseFloatPanic(eles[1], line)*convert.unit,
		parseFloatPanic(eles[2], line)*convert.unit+adj,
		parseFloatPanic(eles[3], line)*convert.unit+adj*2,
		parseFloatPanic(eles[2], line)*convert.unit+adj)

	convert.pdf.Line(
		parseFloatPanic(eles[1], line)*convert.unit+adj,
		parseFloatPanic(eles[2], line)*convert.unit,
		parseFloatPanic(eles[1], line)*convert.unit+adj,
		parseFloatPanic(eles[4], line)*convert.unit+adj*2)

	convert.pdf.Line(
		parseFloatPanic(eles[1], line)*convert.unit,
		parseFloatPanic(eles[4], line)*convert.unit+adj,
		parseFloatPanic(eles[3], line)*convert.unit+adj*2,
		parseFloatPanic(eles[4], line)*convert.unit+adj)

	convert.pdf.Line(
		parseFloatPanic(eles[3], line)*convert.unit+adj,
		parseFloatPanic(eles[2], line)*convert.unit,
		parseFloatPanic(eles[3], line)*convert.unit+adj,
		parseFloatPanic(eles[4], line)*convert.unit+adj*2)
}

// 图片
// ["I", path, x, y, x1, y2]
func (convert *Converter) Image(line string, elements []string) {
	checkLength(line, elements, 6)
	r := new(gopdf.Rect)
	r.W = parseFloatPanic(elements[4], line)*convert.unit - parseFloatPanic(elements[2], line)*convert.unit
	r.H = parseFloatPanic(elements[5], line)*convert.unit - parseFloatPanic(elements[3], line)*convert.unit

	convert.pdf.Image(
		elements[1],
		parseFloatPanic(elements[2], line)*convert.unit,
		parseFloatPanic(elements[3], line)*convert.unit,
		r,
	)
}

// 线
// ["L", x1, y1, x2, y2] 两点之间的线
// ["LH", x1, y1, x2] 水平线
// ["LV", x1, y2, y2] 垂直线
// ["LT", "dashed|dotted|straight", w] 虚线,点,直线
func (convert *Converter) Line(line string, elements []string) {
	switch elements[0] {
	case "L":
		checkLength(line, elements, 5)
		convert.pdf.Line(
			parseFloatPanic(elements[1], line)*convert.unit,
			parseFloatPanic(elements[2], line)*convert.unit,
			parseFloatPanic(elements[3], line)*convert.unit,
			parseFloatPanic(elements[4], line)*convert.unit,
		)
	case "LH":
		checkLength(line, elements, 4)
		convert.pdf.Line(
			parseFloatPanic(elements[1], line)*convert.unit,
			parseFloatPanic(elements[2], line)*convert.unit,
			parseFloatPanic(elements[3], line)*convert.unit,
			parseFloatPanic(elements[2], line)*convert.unit,
		)
	case "LV":
		checkLength(line, elements, 4)
		convert.pdf.Line(
			parseFloatPanic(elements[1], line)*convert.unit,
			parseFloatPanic(elements[2], line)*convert.unit,
			parseFloatPanic(elements[1], line)*convert.unit,
			parseFloatPanic(elements[3], line)*convert.unit,
		)
	case "LT":
		checkLength(line, elements, 3)
		lineType := elements[1]
		if lineType == "" {
			lineType = "straight"
		}
		convert.pdf.SetLineType(lineType)
		convert.linew = parseFloatPanic(elements[2], line)
		convert.pdf.SetLineWidth(convert.linew * convert.unit)
	}
}

// 单元格
// ["C", family, size, x, y, content] // 从(x,y) 位置开始写入content
// ["CL", x, y, content] // 从(x,y) 位置开始写入content
// ["CR", x, y, w, content] // 从右往左写入w长度的内容
func (convert *Converter) Cell(line string, elements []string) {
	switch elements[0] {
	case "C":
		checkLength(line, elements, 6)
		err := convert.pdf.SetFont(elements[1], "", parseIntPanic(elements[2], line))
		if err != nil {
			panic(err.Error() + " line;" + line)
		}
		convert.setPosition(elements[3], elements[4], line)
		convert.pdf.Cell(nil, elements[5])
	case "CL":
		checkLength(line, elements, 4)
		convert.setPosition(elements[1], elements[2], line)
		convert.pdf.Cell(nil, elements[3])
	case "CR":
		checkLength(line, elements, 5)
		tw, err := convert.pdf.MeasureTextWidth(elements[4])
		if err != nil {
			panic(err.Error() + " line;" + line)
		}
		x := parseFloatPanic(elements[1], line) * convert.unit
		y := parseFloatPanic(elements[2], line) * convert.unit
		w := parseFloatPanic(elements[3], line) * convert.unit
		finalx := x + w - tw
		convert.pdf.SetX(finalx)
		convert.pdf.SetY(y)
		convert.pdf.Cell(nil, elements[4])
	}
}

func (convert *Converter) setPosition(x string, y string, line string) {
	convert.pdf.SetX(parseFloatPanic(x, line) * convert.unit)
	convert.pdf.SetY(parseFloatPanic(y, line) * convert.unit)
}

// 外部链接
// ["EL", x, y, w, h, content, link] // 从(x,y)开始写入content,并添加外链接
func (convert *Converter) ExternalLink(line string, elements []string) {
	checkLength(line, elements, 7)

	x, y := parseFloatPanic(elements[1], line), parseFloatPanic(elements[2], line)
	w, h := parseFloatPanic(elements[3], line), parseFloatPanic(elements[4], line)

	convert.pdf.SetX(x)
	convert.pdf.SetY(y)

	convert.pdf.Text(elements[5])
	y1 := y
	if y-h > 0 {
		y1 = y - h
	}
	convert.pdf.AddExternalLink(elements[6], x, y1, w, h)

	convert.pdf.SetX(x + w)
	convert.pdf.SetY(y)
}

// 内部链接, 锚点
// ["ILA", x, y, w, h, content, anchor]
func (convert *Converter) InternalLinkAnchor(line string, elements []string) {
	checkLength(line, elements, 7)

	x, y := parseFloatPanic(elements[1], line), parseFloatPanic(elements[2], line)
	w, h := parseFloatPanic(elements[3], line), parseFloatPanic(elements[4], line)
	convert.pdf.SetX(x)
	convert.pdf.SetY(y)

	convert.pdf.Text(elements[5])
	y1 := y
	if y-h > 0 {
		y1 = y - h
	}
	convert.pdf.AddInternalLink(elements[6], x, y1, w, h)

	convert.pdf.SetX(x + w)
	convert.pdf.SetY(y)
}

// 内部链接, 链接
// ["ILL", x, y, w  content, anchor]
func (convert *Converter) InternalLinkLink(line string, elements []string) {
	checkLength(line, elements, 6)

	convert.pdf.SetX(parseFloatPanic(elements[1], line))
	convert.pdf.SetY(parseFloatPanic(elements[2], line))

	convert.pdf.Text(elements[4])
	convert.pdf.SetAnchor(elements[5])

	convert.pdf.SetX(parseFloatPanic(elements[1], line) + parseFloatPanic(elements[3], line))
	convert.pdf.SetY(parseFloatPanic(elements[2], line))
}

// 辅助方法
func (convert *Converter) Margin(line string, eles []string) {
	checkLength(line, eles, 3)
	top := parseFloatPanic(eles[1], line)
	left := parseFloatPanic(eles[2], line)
	if top != 0.0 {
		convert.pdf.SetTopMargin(top)
	}

	if left != 0.0 {
		convert.pdf.SetLeftMargin(left)
	}
}

func (convert *Converter) GetXY() (x, y float64) {
	return convert.pdf.GetX(), convert.pdf.GetY()
}

func (convert *Converter) MeasureTextWidth(text string) float64 {
	w, err := convert.pdf.MeasureTextWidth(text)
	if err != nil {
		panic(err)
	}

	return w
}

func (convert *Converter) SetFont(family, style string, size int) {
	convert.pdf.SetFont(family, style, size)
}

func (convert *Converter) NoCompression() {
	convert.pdf.SetNoCompression()
}

func (convert *Converter) WritePdf(filepath string) {
	convert.pdf.WritePdf(filepath)

}

func (convert *Converter) CompressLevel(level int) {
	convert.pdf.SetCompressLevel(level)
}

func (convert *Converter) GetBytesPdf() (ret []byte) {
	return convert.pdf.GetBytesPdf()
}

func checkLength(line string, eles []string, no int) {
	if len(eles) < no {
		panic("Column short:" + line)
	}
}

func parseIntPanic(num string, line string) int {
	i, err := strconv.Atoi(num)
	if err != nil {
		panic(num + " not Integer :" + line)
	}
	return i
}

func parseFloatPanic(num string, line string) float64 {
	if num == "" {
		return 0
	}
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		panic(num + " not Numeric :" + line)
	}
	return f
}
