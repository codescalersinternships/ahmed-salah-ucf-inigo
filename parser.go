// Package INIParser provides functionality for parsing INI format in Go.
package iniparser

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	sectionLine = "sectionLine"
	propertyLine = "properityLine"
	commentLine = "commentLine"
	emptyLine = "emptyLine"
	unsportedLine = "unsportedLine"
)

var (
	ErrInvalidFilePath = IniParserError("couldn't find the file in the path you provided")
	ErrNullReference = IniParserError("you tried to access object that doesn't exist")
	ErrSectionNotExist = IniParserError("the section you tried to access doesn't exist")
	ErrKeyNotExist = IniParserError("the key you tried to access doesn't exist")
	ErrHasNoData = IniParserError("there is no data yet, you may didn't load data")
	ErrGlobalProperity = IniParserError("global keys are not allowed")
	ErrEmptySectionName = IniParserError("you should provide sectionName")
	ErrSyntaxError = IniParserError("syntax error, can't understand this line")
	ErrEmptyKey = IniParserError("you should provide key for the properity")

)

type IniParserError string

func (e IniParserError) Error() string {
	return string(e)
}


type (
	// SectionName is the type of keys in sections map for IniParser struct
	SectionName string
	// Key is the type of the keys for INI fields
	Key string
	// Section is the type of values for sections in IniParser
	Section map[Key]string
)

// IniParser is the type that represent INI file structure and methods
// INI content is represented as a map in which keys are section names
// and values are maps of keys and values from the ini properties.
type IniParser struct {
	sections map[SectionName]Section
}

// New function create new IniParser object and return it.
func New() *IniParser{
	return &IniParser{map[SectionName]Section{}}
}


// GetSections return map of sections
func (i *IniParser) GetSections() (sections map[SectionName]Section) {
	sections = i.sections
	return
}

// GetSectionNames is a function that returns a slice
// of all section names in the IniParser object
func (i *IniParser) GetSectionNames () ([]string) {
	sectionNamesList := []string{}
	for sectionName := range i.sections {
		sectionNamesList = append(sectionNamesList, string(sectionName))
	}
	sort.Strings(sectionNamesList)
	return sectionNamesList
}


// Get function gets the section name of type SectionName and the key
// of type Key and return the Value associated with that key that has
// type Value.
// The function returns err == nil if the returned successfully.
// 						err == ErrNullReference if sections is not defined.
// 						err == ErrSectionNotExist if no section with name sectionName.
// 						err == ErrKeyNotExist if no key with name key.
func (i *IniParser) Get(sectionName SectionName, key Key) (string, error) {
	if i.sections == nil {
		return "", ErrNullReference
	}
	if _, ok := i.sections[sectionName]; !ok {
		return "", ErrSectionNotExist
	}
	value, ok := i.sections[sectionName][key]
	if !ok {
		return "", ErrKeyNotExist
	}
	return value, nil
}

func (i *IniParser) Set(sectionName SectionName, key Key, value string) error{
	if i.sections == nil {
		return ErrNullReference
	}
	if _, ok := i.sections[sectionName]; !ok {
		return ErrSectionNotExist
	}
	
	if _, ok := i.sections[sectionName][key]; !ok {
		return ErrKeyNotExist
	}

	i.sections[sectionName][key] = value
	return nil
}

// LoadFromFile get filePath as argument and returns the file content as a string
// A successful call returns err == nil, and non-successful call returns an error
// of type ErrInvalidFilePath
func (i *IniParser) LoadFromFile(filePath string) (string, error) {
	
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", ErrInvalidFilePath
	}
	i.sections, err = parse(string(fileContent))
	return string(fileContent), err
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
	return len(line) == 0 || line[0] == '\n'
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

func parseSection(sectionLine string) (SectionName, error) {
	if len(sectionLine) == 2 {
		return "", ErrEmptySectionName
	}
	sectionName := strings.TrimLeft(sectionLine[1:len(sectionLine)-1], " [")
	sectionName = strings.TrimRight(sectionName, " ]")
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

func parse(iniData string) (map[SectionName]Section, error) {
	ini := New()
	var currentSectionName SectionName
	var key Key
	var value string
	var err error

	dataLines := strings.Split(iniData, "\n")

	for _, line := range dataLines {
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

// LoadFromString takes iniData of type string as argument
// and loads the data into the object's sections field.
// It's the end-user responsibility to define the sections field
// of type map[SectionName]Section.
// the function returns ErrNullReference error if the user tried
// to Load INI data into IniParser that has sections undefined.
func (i *IniParser) LoadFromString(iniData string) (err error) {
	i.sections, err = parse(iniData)
	
	return err
}

func (i *IniParser) String() (string, error) {
	if (i.sections == nil) {
		return "", ErrNullReference
	}
	if len(i.sections) == 0 {
		return "", ErrHasNoData
	}
	var result string
	for SectionName, section := range i.sections {
		result += fmt.Sprintf("[%s]\n", SectionName)
		for name, value := range section {
			result += fmt.Sprintf("%v = %s\n", name, value)
		}
		
	}
	return result, nil
}

func (i *IniParser) SaveToFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil{
		file.Close()
		return err
	}
	defer file.Close()

	content, err1 := i.String()
	file.WriteString(content)
	if err1 == ErrNullReference {
		err = err1
	} else if err1 == ErrHasNoData {
		err = err1
		file.WriteString(content)
	}
	return err
}

