// Package INIParser provides functionality for parsing INI format in Go.
package iniparser

import (
	"fmt"
	"os"
	"sort"
)


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

// LoadFromString takes iniData of type string as argument
// and loads the data into the object's sections field.
// The function returns err == ErrGlobalProperity if file contains global properties.
//			err == ErrEmptySectionName if section line has no name i.e: [ ].
// 			err == ErrEmptyKey if properity has no key
// 			err == ErrSyntaxError if there is any unsupported format
func (i *IniParser) LoadFromString(iniData string) (err error) {
	i.sections, err = parse(iniData)
	
	return err
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
// 			err == ErrNullReference if sections is not defined.
// 			err == ErrSectionNotExist if no section with name sectionName.
// 			err == ErrKeyNotExist if no key with name key.
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


// Set function get section name and key and value, and set the key in the
// given section with given value.
// It returns err == ErrNullReference if the user try to access undefined sections
// err == ErrSectionNotExist if the sectionName doesn't exist
// err == ErrKeyNotExist if the key doesn't exist
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

// String function converts the IniParser object into string type
// and returns that string.
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

// SaveToFile get filePath and save the data in sections
// field into the file at that filePath.
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

