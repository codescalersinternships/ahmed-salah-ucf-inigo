// Package INIParser provides functionality for parsing INI format in Go.
package iniparser

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidFilePath = IniParserError("couldn't find the file in the path you provided")
	ErrNullReference = IniParserError("you tried to access object that doesn't exist")
	ErrSectionNotExist = IniParserError("the section you tried to access doesn't exist")
	ErrKeyNotExist = IniParserError("the key you tried to access doesn't exist")
	ErrHasNoData = IniParserError("there is no data yet, you may didn't load data")
)

var (
	rSection, _ = regexp.Compile(`\[.*?\]`)
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
type IniParser struct {
	sections map[SectionName]Section
}

// NewIniParser function create new IniParser object and return it.
func NewIniParser() IniParser{
	return IniParser{map[SectionName]Section{}}
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
	return string(fileContent), err
}

// isCommentLine is a helper predicate that takes a line of string
// and determine if this line is INI comment or not
func isCommentLine(line string) bool {
	firstCharachter, _ := utf8.DecodeRuneInString(line[0:])
	return firstCharachter == ';'
}

// isSectionLine is a helper predicate that takes a line of string
// and determine if this line is INI section or not
func isSectionLine(line string, rSection *regexp.Regexp) bool {
	
	return rSection.MatchString(line)
}

// ParseFieldLine is a helper function that get a line of string and
// parses the line into key and value of type Key and Value respectivly
func parseFieldLine(line string) (Key, string) {
	keyAndValue := strings.Split(line, "=")
	key := Key(strings.Trim(keyAndValue[0], " "))
	value := strings.Trim(keyAndValue[1], " ")
	return key, value
}

// LoadFromString takes iniData of type string as argument
// and loads the data into the object's sections field.
// It's the end-user responsibility to define the sections field
// of type map[SectionName]Section.
// the function returns ErrNullReference error if the user tried
// to Load INI data into IniParser that has sections undefined.
func (i *IniParser) LoadFromString(iniData string) error {
	if i.sections == nil {
		return ErrNullReference
	}
	dataLines := strings.Split(iniData, "\n")
	var sectionName string
	
	for _, line := range dataLines {
		if len(line) > 0 {
			line = strings.Trim(line, " ")
			if isCommentLine(line) {
				continue
			} else if isSectionLine(line, rSection) {
				sectionName = rSection.FindString(line)
				sectionName = strings.TrimLeft(sectionName, " [")
				sectionName = strings.TrimRight(sectionName, " ]")
				i.sections[SectionName(sectionName)] = Section{}
			} else {
				key, value := parseFieldLine(line)
				i.sections[SectionName(sectionName)][key] = value
			}
		}
	}
	return nil
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

