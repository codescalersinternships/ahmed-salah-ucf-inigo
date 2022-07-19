package iniparser

import (
	"fmt"
	"reflect"
)

func ExampleIniParser_LoadFromString() {
	ini := New()
	ini.LoadFromString(iniContent)
	fmt.Println(ini.sections)
	// Output: map[database:map[file:"payroll.dat" port:143 server:192.0.2.62] owner:map[name:John Doe organization:Acme Inc.]]
}
func ExampleIniParser_LoadFromFile() {
	ini := New()
	ini.LoadFromFile(exampleFilePath)
	fmt.Println(ini.sections)
	// output: map[database:map[file:"payroll.dat" port:143 server:192.0.2.62] owner:map[name:John Doe organization:Acme Inc.]]
}

func ExampleIniParser_GetSections() {
	ini := New()
	ini.LoadFromFile(exampleFilePath)
	sections := ini.GetSections()
	fmt.Println(sections)
	// Output: map[database:map[file:"payroll.dat" port:143 server:192.0.2.62] owner:map[name:John Doe organization:Acme Inc.]]
}

func ExampleIniParser_GetSectionNames() {
	ini := New()
	ini.LoadFromFile(exampleFilePath)
	sections := ini.GetSectionNames()
	fmt.Println(sections)
	// Output: [database owner]
}

func ExampleIniParser_Get() {
	ini := New()
	ini.LoadFromFile(exampleFilePath)
	sections, _ := ini.Get("owner", "name")
	fmt.Println(sections)
	// Output: John Doe
}

func ExampleIniParser_Set() {
	ini := New()
	ini.LoadFromFile(exampleFilePath)
	ini.Set("owner", "name", "salah")
	fmt.Println(ini.sections["owner"]["name"])
	// Output: salah
}

func ExampleIniParser_String() {
	ini := New()
	ini.LoadFromString(iniContent)
	stringContent, _ := ini.String()
	oldContentSectionsMap := ini.sections

	newIni := New()
	newIni.LoadFromString(stringContent)
	newContentSectionsMap := newIni.sections
	fmt.Println(reflect.DeepEqual(oldContentSectionsMap, newContentSectionsMap))
	// Output: true
}

func ExampleIniParser_SaveToFile() {
	ini := New()
	ini.SaveToFile(failOutFilePath)
	oldContentSectionsMap := ini.sections
	strContent, _ := ini.LoadFromFile(failOutFilePath)
	ini.LoadFromString(strContent)
	newContentSectionsMap := ini.sections
	fmt.Println(reflect.DeepEqual(oldContentSectionsMap, newContentSectionsMap))
	// Output: true
}