package iniparser

import (
	"bufio"
	"strings"
)

func parse(iniData string) (map[SectionName]Section, error) {
	ini := New()
	var currentSectionName SectionName
	var key Key
	var value string
	var err error
	
	scanner := bufio.NewScanner(strings.NewReader(iniData))

	for scanner.Scan() {
		line := scanner.Text()
		lineType := lineType(line)
		switch lineType {
		case sectionLine:
			currentSectionName, err = parseSection(line)
			if err != nil {
				return 	ini.sections, err
			}
			ini.sections[currentSectionName] = Section{}
		case propertyLine:
			key, value, err = parseProperity(line)
			if err != nil {
				return ini.sections, err
			}
			if currentSectionName == "" {
				return ini.sections, ErrGlobalProperity
			}
			ini.sections[currentSectionName][key] = value
		case commentLine:
		case emptyLine:
			continue

		case unsportedLine:
			return ini.sections, ErrSyntaxError
		}
	}
	return ini.sections, nil
}

func lineType(line string) string {
	if isEmptyLine(line) {
		return emptyLine
	}
	if isSection(line) {
		return sectionLine
	}
	if isProperity(line) {
		return propertyLine
	}
	if isComment(line){
		return commentLine
	}
	
	return unsportedLine
}

func isSection(line string) bool {

	line = strings.TrimSpace(line)

	return line[0] == '[' && line[len(line)-1] == ']' &&
			strings.Count(line, "[") == 1 && strings.Count(line, "]") == 1
}

func isProperity(line string) bool {
	return strings.Count(line, "=") == 1
}

func isComment(line string) bool {
	return line[0] == ';'
}

func isEmptyLine(line string) bool {
	return len(line) == 0
}

func parseSection(sectionLine string) (SectionName, error) {
	if len(sectionLine) == 2 {
		return "", ErrEmptySectionName
	}
	sectionLine = strings.ReplaceAll(sectionLine, " ", "")
	
	sectionName := strings.TrimLeft(sectionLine[1:len(sectionLine)-1], " [")
	sectionName = strings.TrimRight(sectionName, " ]")
	if len(sectionLine) == 2 {
		return "", ErrEmptySectionName
	}

	return SectionName(sectionName), nil
}

func parseProperity(property string) (Key, string, error) {
	sepIdx := strings.Index(property, "=")
	key := property[0:sepIdx]
	if len(key) == 0 {
		return Key(""), "", ErrEmptyKey
	}
	key = strings.TrimSpace(key)
	value := property[sepIdx+1:]
	value = strings.TrimSpace(value)

	return Key(key), value, nil
}