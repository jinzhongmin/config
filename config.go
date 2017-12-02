package config

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//Config ..
type Config struct {
	lines []*Line
}

//View ..
type View struct {
	lines []*Line
}

//Line ..
type Line struct {
	tab  int
	args []string
}

func read(file string) string {
	fs, err := os.Open(file)
	if err != nil {
	}
	bt, err := ioutil.ReadAll(fs)
	if err != nil {
	}
	return string(bt)
}
func clean(text string) string {
	text = strings.Replace(text, "\r\n[", ";;;[;;;", -1)
	text = strings.Replace(text, "\r\n#", ";;;#;;;", -1)

	for i := 20; i > 0; i-- {
		text = strings.Replace(text, "\r\n"+blank(i), ";;;"+strconv.Itoa(i)+";", -1)
		if i > 1 {
			text = strings.Replace(text, blank(i), "", -1)
		}
	}

	text = strings.Replace(text, "\r\n", "", -1)

	text = strings.Replace(text, ";;;[;;;", "\r\n[", -1)
	text = strings.Replace(text, ";;;#;;;", "\r\n#", -1)

	for i := 20; i > 0; i-- {
		text = strings.Replace(text, ";;;"+strconv.Itoa(i)+";", "\r\n"+blank(i), -1)
	}

	return text
}
func blank(num int) string {
	str := ""
	for i := 0; i < num; i++ {
		str += " "
	}
	return str
}

//New ..
func New(fs string) *Config {
	config := new(Config)
	text := clean(read(fs))

	lines := strings.Split(text, "\r\n")
	for i := range lines {
		config.lines = append(config.lines, NewLine(lines[i]))
	}

	return config
}

//Lines ..
func (config *Config) Lines(match []string) []*Line {
	lines := make([]*Line, 0)

	for i := range config.lines {
		line := config.lines[i]

		if len(line.args) < len(match) {
			continue
		}

		flag := true
		for ii := range match {
			if line.args[ii] == match[ii] {
				continue
			}
			flag = false
		}

		if flag == true {
			lines = append(lines, line)
		}
	}

	return lines
}

//ViewsAuto ..
func (config *Config) ViewsAuto(match []string) []*View {
	views := make([]*View, 0)

	for i := 0; i < len(config.lines); i++ {
		lines := make([]*Line, 0)
		line := config.lines[i]

		if len(line.args) < len(match) {
			continue
		}

		flag := true
		for ii := range match {
			if line.args[ii] == match[ii] {
				continue
			}
			flag = false
			break
		}

		if flag == true { //开始匹配
			lines = append(lines, line)
			i++

			tabs := line.Tabs()

			for ; i < len(config.lines); i++ {
				if config.lines[i].Tabs() > tabs {
					lines = append(lines, config.lines[i])
				} else {
					view := new(View)
					view.lines = lines
					views = append(views, view)
					i = i - 1
					break
				}
			}
		}
	}
	return views
}

//Views ..
func (config *Config) Views(begin []string, end []string) []*View {
	views := make([]*View, 0)

	for i := 0; i < len(config.lines); i++ {
		lines := make([]*Line, 0)
		line := config.lines[i]

		if len(line.args) < len(begin) {
			continue
		}

		flag := true
		for ii := range begin {
			if line.args[ii] == begin[ii] {
				continue
			}
			flag = false
			break
		}

		if flag == true { //开始匹配
			lines = append(lines, line)

			for i = i + 1; i < len(config.lines); i++ {
				line = config.lines[i]

				if len(line.args) < len(begin) {
					lines = append(lines, line)
					continue
				}

				for ii := range end {
					if line.args[ii] == end[ii] {
						continue
					}
					flag = false
					break
				}
				lines = append(lines, config.lines[i])

				if flag == true {
					view := new(View)
					view.lines = lines
					views = append(views, view)
					i = i - 1
					break
				}

				flag = true
			}
		}
	}
	return views
}

//LinesAll ..
func (view *View) LinesAll() []*Line {
	return view.lines
}

//Lines ..
func (view *View) Lines(match []string) []*Line {
	lines := make([]*Line, 0)

	for i := range view.lines {
		line := view.lines[i]

		if len(line.args) < len(match) {
			continue
		}

		flag := true
		for ii := range match {
			if line.args[ii] == match[ii] {
				continue
			}
			flag = false
		}

		if flag == true {
			lines = append(lines, line)
		}
	}

	return lines
}

//NewLine ..
func NewLine(line string) *Line {
	tab := 0
	for i := 20; i > 0; i-- {
		index := strings.Index(line, blank(i))
		if index == 0 {
			tab = i
			break
		}
	}

	line = strings.Replace(line, blank(tab), "", 1)
	args := strings.Split(line, " ")

	l := new(Line)
	l.tab = tab
	l.args = args

	return l
}

//NextS ..
func (line *Line) NextS(arg string) string {
	for i := range line.args {
		if line.args[i] == arg && len(line.args) > i+1 {
			return line.args[i]
		}
	}
	return ""
}

//NextI ..
func (line *Line) NextI(arg string) int {
	for i := range line.args {
		if line.args[i] == arg && len(line.args)-1 > i {
			r, err := strconv.Atoi(line.args[i])
			if err != nil {
				return r
			}
		}
	}
	return -1
}

//Indexs ..
func (line *Line) Indexs(indexs []int) []string {
	args := make([]string, 0)
	for i := range indexs {
		args = append(args, line.args[indexs[i]])
	}

	return args
}

//IndexStr ..
func (line *Line) IndexStr(index int) string {
	return line.args[index]
}

//IndexInt ..
func (line *Line) IndexInt(index int) int {
	i, _ := strconv.Atoi(line.args[index])
	return i
}

//Len ..
func (line *Line) Len() int {
	return len(line.args)
}

//Tabs ..
func (line *Line) Tabs() int {
	return line.tab
}
