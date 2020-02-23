package config

import (
	"bytes"
	"log"
	"os"
	"regexp"
	"strings"
)

type IfDo struct {
	seg string
	fn  func(l *Line)
}

type Seg struct {
	word []string
}

type Line struct {
	indent int
	word   []string
}

type View struct {
	body []*Line
}

type Config struct {
	body []*Line
}

func ReadLine(line string) Line {
	reg := regexp.MustCompilePOSIX("^[ ]*")
	indent := reg.FindString(line)
	l := strings.Replace(line, indent, "", 1)
	s := new(string)
	s = &l
	for {
		if strings.Index(*s, "  ") > 0 {
			l := strings.Replace(*s, "  ", " ", -1)
			s = &l
		} else {
			break
		}
	}
	return Line{len(indent), strings.Split(*s, " ")}
}

func NewSeg(seg string) Seg {
	reg := regexp.MustCompilePOSIX("^[ ]*")
	indent := reg.FindString(seg)
	l := strings.Replace(seg, indent, "", 1)
	s := new(string)
	s = &l
	for {
		if strings.Index(*s, "  ") > 0 {
			l := strings.Replace(*s, "  ", " ", -1)
			s = &l
		} else {
			break
		}
	}
	return Seg{strings.Split(*s, " ")}
}

func NewView() *View {
	v := new(View)
	v.body = make([]*Line, 0)
	return v
}

func Switch(i ...interface{}) []IfDo {
	id := make([]IfDo, 0)
	if len(i)%2 != 0 {
		return id
	}
	for ii, o := range i {
		if ii%2 == 0 {
			seg, ok1 := o.(string)
			fn, ok2 := i[ii+1].(func(l *Line))
			if ok1 && ok2 {
				id = append(id, IfDo{seg, fn})
			}

		}
	}

	return id
}

func (l *Line) Copy() *Line {
	_l := new(Line)
	_l.word = make([]string, 0)
	_l.word = append(_l.word, l.word...)
	return _l
}
func (l *Line) String() string {
	indent := ""
	for i := 1; i < l.indent; i++ {
		indent += " "
	}
	str := indent
	for _, s := range l.word {
		str += " " + s
	}
	return str
}

func (l *Line) IfBegin(seg string, fn func(l *Line)) bool {
	s := NewSeg(seg)
	if len(l.word) < len(s.word) {
		return false
	}
	for i, w := range s.word {
		if w == "*" {
			continue
		}
		if l.word[i] != w {
			return false
		}
	}
	fn(l)
	return true
}

func (l *Line) IfInclude(word string, fn func(l *Line)) {
	if strings.Index(l.String(), word) > -1 {
		fn(l)
	}
}

func (l *Line) GetIndex(index int) string {
	if index >= len(l.word) {
		return ""
	}
	return l.word[index]
}

func (l *Line) GetNext(word string) string {
	for i, w := range l.word {
		if w == word {
			if i+1 < len(l.word) {
				return l.word[i+1]
			}
			return ""
		}
	}
	return ""
}

func (l *Line) DeleteWord(word string) {
	for i, w := range l.word {
		if w == word {
			_w := make([]string, 0)
			_w = append(_w, l.word[:i]...)
			_w = append(_w, l.word[i+1:]...)
			l.word = _w
		}
	}

}

func (l *Line) DeleteWordNext(word string) {
	for i, w := range l.word {
		if w == word && i+1 < len(l.word) {
			_w := make([]string, 0)
			_w = append(_w, l.word[:i+1]...)
			_w = append(_w, l.word[i+2:]...)
			l.word = _w
		}
	}

}

func (l *Line) DeleteWordByIndex(index int) {
	if index < len(l.word) {
		w := make([]string, 0)
		w = append(w, l.word[:index]...)
		w = append(w, l.word[index+1:]...)
		l.word = w
	}
}

func (l *Line) ModifWord(old string, new string) {
	for i, w := range l.word {
		if w == old {
			end := l.word[i+1:]
			l.word = append(l.word[:i], new)
			l.word = append(l.word, end...)
			return
		}
	}
}

func (l *Line) ModifWordNext(word string, new string) {
	for i, w := range l.word {
		if w == word && i+1 < len(l.word) {
			end := l.word[i+2:]
			l.word = append(l.word[:i+1], new)
			l.word = append(l.word, end...)
			return
		}
	}
}

func (l *Line) ModifWordByIndex(index int, new string) {
	if index < len(l.word) {
		end := l.word[index+1:]
		l.word = append(l.word[:index], new)
		l.word = append(l.word, end...)
	}
}

func (v *View) Copy() *View {
	_v := new(View)
	_v.body = make([]*Line, 0)
	v.RangeLine(func(l *Line) {
		_v.body = append(_v.body, l.Copy())
	})

	return _v
}

func (v *View) String() string {
	str := ""
	for _, l := range v.body {
		str += l.String() + "\n"
	}
	return str
}

func (v *View) AppendLine(l *Line) {
	v.body = append(v.body, l)
}

func (v *View) GetLineByIndex(index int) *Line {
	return v.body[index]
}

func (v *View) RangeLine(fn func(l *Line)) {
	for i, _ := range v.body {
		fn(v.body[i])
	}
}

func (v *View) RangeLineIfBegin(word string, fn func(l *Line)) {
	for i, _ := range v.body {
		v.body[i].IfBegin(word, fn)
	}
}

func (v *View) RangeLineIfInclude(word string, fn func(l *Line)) {
	for i, _ := range v.body {
		v.body[i].IfInclude(word, fn)
	}
}

func (v *View) RangeLineSwitchBegin(sw []IfDo) {
	for i, _ := range v.body {
		for _, ii := range sw {
			v.body[i].IfBegin(ii.seg, ii.fn)
		}
	}
}

func (v *View) RangeLineSwitchInclude(sw []IfDo) {
	for i, _ := range v.body {
		for _, ii := range sw {
			v.body[i].IfInclude(ii.seg, ii.fn)
		}
	}
}

func LoadConfig(path string) *Config {
	fs, err := os.Open(path)
	if err != nil {
		log.Panicln(err)
	}

	buf := bytes.Buffer{}
	buf.ReadFrom(fs)

	cfg := buf.String()
	cfgls := strings.Split(cfg, "\r\n")

	ls := make([]*Line, 0)
	for _, l := range cfgls {
		_l := ReadLine(l)
		ls = append(ls, &_l)
	}

	c := new(Config)
	c.body = ls
	return c
}

func (c *Config) GetViews(seg string) []*View {
	view := make([]*View, 0)

	lv := new(View)
	li := 0
	for _, l := range c.body {
		if l.IfBegin(seg, func(l *Line) {
			if len(lv.body) > 0 {
				lv.body[0].IfBegin(seg, func(l *Line) {
					view = append(view, lv)
				})
			}

			lv = new(View)
			lv.body = append(lv.body, l)
			li = l.indent
		}) {
			continue
		}

		if l.indent < li {
			if len(lv.body) > 0 {
				lv.body[0].IfBegin(seg, func(l *Line) {
					view = append(view, lv)
				})
			}

			lv = new(View)
			continue
		}
		if l.indent >= li {
			li = l.indent
		}
		lv.body = append(lv.body, l)
	}

	return view
}

func (c *Config) RangeViews(seg string, fn func(v *View)) {
	vs := c.GetViews(seg)
	for _, v := range vs {
		fn(v)
	}
}

func (c *Config) RangeLinesBegin(seg string, fn func(l *Line)) {
	for _, l := range c.body {
		l.IfBegin(seg, fn)
	}
}

func (c *Config) RangeLinesInclude(seg string, fn func(l *Line)) {
	for _, l := range c.body {
		l.IfInclude(seg, fn)
	}
}

func (c *Config) RangeLineSwitchBegin(sw []IfDo) {
	for i, _ := range c.body {
		for _, ii := range sw {
			c.body[i].IfBegin(ii.seg, ii.fn)
		}
	}
}

func (c *Config) RangeLineSwitchInclude(sw []IfDo) {
	for i, _ := range c.body {
		for _, ii := range sw {
			c.body[i].IfInclude(ii.seg, ii.fn)
		}
	}
}

func (c *Config) GetLinesViewBegin(seg string) *View {
	v := new(View)
	c.RangeLinesBegin(seg, func(l *Line) {
		v.body = append(v.body, l)
	})
	return v
}

func (c *Config) GetLinesViewInclude(seg string) *View {
	v := new(View)
	c.RangeLinesInclude(seg, func(l *Line) {
		v.body = append(v.body, l)
	})
	return v
}
