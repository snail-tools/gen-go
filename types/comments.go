package types

import (
	"bytes"
	"strings"
)

// 提前注释信息
//
// @param lines			注释
// @param markers		标记符号，可以为多个。默认为 +  @ 符号
//
// 注释示例:
//   +foo=value1
//   +bar
//   +foo value2
//   +baz="qux"
// return:
//   map[string][]string{"foo":{"value1, "value2"}, "bar": {"true"}, "baz": {"qux"}}
func ExtractCommentTags(lines []string, markers ...byte) (tags map[string][]string, otherLines []string) {
	if len(markers) == 0 {
		markers = []byte{'+', '@'}
	}
	tags = map[string][]string{}
	for _, line := range lines {
		line = strings.Trim(line, " ")

		if !(len(line) != 0 && oneOf(markers, line[0])) {
			otherLines = append(otherLines, line)
			continue
		}

		k, v := splitKV(line[1:])
		tags[k] = append(tags[k], v)
	}

	return
}

func splitKV(line string) (string, string) {
	k := bytes.NewBuffer(nil)
	v := bytes.NewBuffer(nil)

	forValue := false

	for _, c := range line {
		if !forValue && (c == '=' || c == ' ') {
			forValue = true
			continue
		}

		if forValue {
			v.WriteRune(c)
		} else {
			k.WriteRune(c)
		}
	}

	return k.String(), v.String()
}

func oneOf(markers []byte, b byte) bool {
	for _, m := range markers {
		if b == m {
			return true
		}
	}
	return false
}
