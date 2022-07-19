package iniparser

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var (
	failOutFilePath = filepath.Join(".", "failOutFile.ini")
	successOutFilePath = filepath.Join(".", "successOutFile.ini")
	exampleFilePath = filepath.Join(".", "example.ini")
)

var mapOfSections = map[SectionName]Section{
	SectionName("owner") : {Key("name") : "John Doe",
							Key("organization") : "Acme Inc."},

	SectionName("database") : {Key("server") : "192.0.2.62",
							   Key("port") : "143",
							   Key("file") : "\"payroll.dat\"",},
}

var iniContent = `; last modified 1 April 2001 by John Doe
[owner]
name = John Doe
organization = Acme Inc.

[database]
; use IP address in case network name resolution is not working
server = 192.0.2.62
port = 143
file = "payroll.dat"`

var iniGlobalContent = `; last modified 1 April 2001 by John Doe
name = John Doe
[owner]
organization = Acme Inc.

[database]
; use IP address in case network name resolution is not working
server = 192.0.2.62
port = 143
file = "payroll.dat"`

var ExampleIniContent = `[owner]
name = John Doe

[database]
server = 192.0.2.62
`

var iniErrGlobalProperity = `[owner]
name = John Doe

[]
server = 192.0.2.62
`

var spacedIniContent = `; last modified 1 April 2001 by John Doe
[owner]
name = John Doe
organization= Acme Inc.
[database     ]
; use IP address in case network name resolution is not working
server      = 192.0.2.62     
     port = 143
file = "payroll.dat"`

var emptyLinesIniContent = `; last modified 1 April 2001 by John Doe
[owner]
name = John Doe
organization = Acme Inc.


[database     ]
; use IP address in case network name resolution is not working



server = 192.0.2.62     
port = 143

;

file = "payroll.dat"`

func TestGetSections(t *testing.T) {
	t.Run("get sections", func(t *testing.T) {
		ini := New()
		ini.LoadFromString(iniContent)
		got := ini.GetSections()
		want := mapOfSections

		assertEqualSections(t, got, want)
	})

	t.Run("nil sections", func(t *testing.T) {
		ini := IniParser{}
		got := ini.GetSections()

		assertEqualSections(t, got, nil)
	})
	
}

func ExampleIniParser_GetSections() {
	ini := New()
	stringContent, _ := ini.LoadFromFile(exampleFilePath)
	ini.LoadFromString(stringContent)
	sections := ini.GetSections()
	fmt.Println(sections)
	// Output: map[database:map[file:"payroll.dat" port:143 server:192.0.2.62] owner:map[name:John Doe organization:Acme Inc.]]
}

func TestGetSectionNames(t *testing.T) {
	assertEqualLists := func(t testing.TB, got, want []string) {
		t.Helper()
		if (!reflect.DeepEqual(got, want)) {
			t.Errorf("got %v want %v", got, want)
		}
	}
	t.Run("nil sections", func(t *testing.T) {
		ini := IniParser{}
		got := ini.GetSectionNames()
		want := []string{}
		assertEqualLists(t, got, want)
	})

	t.Run("has sections", func(t *testing.T) {
		ini := New()
		ini.LoadFromString(iniContent)
		got := ini.GetSectionNames()
		want := []string{"database", "owner"}
		assertEqualLists(t, got, want)
	})
}

func ExampleIniParser_GetSectionNames() {
	ini := New()
	stringContent, _ := ini.LoadFromFile(exampleFilePath)
	ini.LoadFromString(stringContent)
	sections := ini.GetSectionNames()
	fmt.Println(sections)
	// Output: [database owner]
}

func TestGet(t *testing.T) {
	t.Run("get value crosponding to key in section", func(t *testing.T) {
		ini := IniParser{sections: map[SectionName]Section{SectionName("owner") : {"name" : "salah"}}}
		got, err := ini.Get(SectionName("owner"), Key("name"))
		want := "salah"
		
		assertNoErrorMsg(t, err)
		assertEqualStrings(t, got, want)
	})
	t.Run("nil sections", func(t *testing.T) {
		ini := IniParser{}
		_, err := ini.Get(SectionName("owner"), Key("name"))
		
		assertErrorMsg(t, err, ErrNullReference)
	})
	t.Run("section doesn't exist", func(t *testing.T) {
		ini := IniParser{sections: map[SectionName]Section{SectionName("owner") : {"name" : "salah"}}}
		_, err := ini.Get(SectionName("employee"), Key("name"))
		
		assertErrorMsg(t, err, ErrSectionNotExist)
	})
	t.Run("key doesn't exist", func(t *testing.T) {
		ini := IniParser{sections: map[SectionName]Section{SectionName("owner") : {"name" : "salah"}}}
		_, err := ini.Get(SectionName("owner"), Key("address"))
		
		assertErrorMsg(t, err, ErrKeyNotExist)
	})
}

func ExampleIniParser_Get() {
	ini := New()
	stringContent, _ := ini.LoadFromFile(exampleFilePath)
	ini.LoadFromString(stringContent)
	sections, _ := ini.Get("owner", "name")
	fmt.Println(sections)
	// Output: John Doe
}

func TestSet(t *testing.T) {
	t.Run("set value for key in section", func(t *testing.T) {
		ini := IniParser{sections: map[SectionName]Section{SectionName("owner") : {"name" : "salah"}}}
		err := ini.Set(SectionName("owner"), Key("name"), "ahmed")
		got, _ := ini.Get(SectionName("owner"), Key("name"))
		want := "ahmed"
		
		assertNoErrorMsg(t, err)
		assertEqualStrings(t, got, want)
	})
	t.Run("nil sections", func(t *testing.T) {
		ini := IniParser{}
		err := ini.Set(SectionName("owner"), Key("name"), "salah")
		
		assertErrorMsg(t, err, ErrNullReference)
	})
	t.Run("section doesn't exist", func(t *testing.T) {
		ini := IniParser{sections: map[SectionName]Section{SectionName("owner") : {"name" : "salah"}}}
		err := ini.Set(SectionName("employee"), Key("name"), "salah")
		
		assertErrorMsg(t, err, ErrSectionNotExist)
	})
	t.Run("key doesn't exist", func(t *testing.T) {
		ini := IniParser{sections: map[SectionName]Section{SectionName("owner") : {"name" : "salah"}}}
		err := ini.Set(SectionName("owner"), Key("address"), "mahalla")
		
		assertErrorMsg(t, err, ErrKeyNotExist)
	})
}

func ExampleIniParser_Set() {
	ini := New()
	stringContent, _ := ini.LoadFromFile(exampleFilePath)
	ini.LoadFromString(stringContent)
	ini.Set("owner", "name", "salah")
	fmt.Println(ini.sections["owner"]["name"])
	// Output: salah
}

func TestLoadFromFile(t *testing.T) {
	t.Run("valid file path", func(t *testing.T) {
		ini := New()

		got, err := ini.LoadFromFile(exampleFilePath)
		want := iniContent

		assertNoErrorMsg(t, err)
		assertEqualStrings(t, got, want)
	})

	t.Run("invalid file path", func(t *testing.T) {
		filepath := "/invalid/file/path/example.ini"
		ini := New()
		_, err := ini.LoadFromFile(filepath)

		assertErrorMsg(t, err, ErrInvalidFilePath)
	})
}

func ExampleIniParser_LoadFromFile() {
	ini := New()

	stringContent, _ := ini.LoadFromFile(exampleFilePath)
	ini.LoadFromString(stringContent)
	fmt.Println(ini.sections)
	// output: map[database:map[file:"payroll.dat" port:143 server:192.0.2.62] owner:map[name:John Doe organization:Acme Inc.]]
}

func TestLoadFromString(t *testing.T) {
	t.Run("input empty string", func(t *testing.T) {
		ini := New()
		err := ini.LoadFromString("")
		got := ini.sections
		want := map[SectionName]Section{}
		
		assertNoErrorMsg(t, err)
		assertEqualSections(t, got, want)
	})

	t.Run("success parse", func(t *testing.T) {
		ini := New()
		err := ini.LoadFromString(iniContent)
		got := ini.sections
		want := mapOfSections
		
		assertNoErrorMsg(t, err)
		assertEqualSections(t, got, want)	
	})

	t.Run("data contain global content", func(t *testing.T) {
		ini := New()
		err := ini.LoadFromString(iniGlobalContent)
		
		assertErrorMsg(t, err, ErrGlobalProperity)
	})

	t.Run("empty section name", func(t *testing.T) {
		ini := New()
		err := ini.LoadFromString(iniErrGlobalProperity)
		
		assertErrorMsg(t, err, ErrEmptySectionName)
	})

	t.Run("properity with missed key", func(t *testing.T) {
		ini := New()
		err := ini.LoadFromString("[owner]\n=value")
		
		assertErrorMsg(t, err, ErrEmptyKey)
	})

	var testsSpaces = []string{"[owner      ]\nname=salah",
								"[     owner    ]\nname=salah",
								"     [owner    ]    \nname=salah",
								"[owner]\nname    =salah",
								"[owner]\nname    =salah    ",
								"[owner]\n     name    =salah"}
	var contentNoSpaces = "[owner]\nname=salah"

	for _, tt := range testsSpaces {
		t.Run("trim spaces", func(t *testing.T) {
			ini := New()
			ini.LoadFromString(tt)
			got := ini.sections

			newIni := New()
			newIni.LoadFromString(contentNoSpaces)
			want := newIni.sections

			assertEqualSections(t, got, want)
		})
	}

	var testsSyntax = []struct {
		testName string
		content string
	} {
		{"missed section bracket", "owner]\nname=salah"},
		{"multiple property sperators", "[owner]\nname====salah"},
		{"not ini syntax", "{\"name\":\"John\"}"},
	}
	
	// {"owner]\nname=salah", "[owner]\nname====salah", }
	for _, tt := range testsSyntax {
		t.Run("syntax error: " + tt.testName, func(t *testing.T) {
			ini := New()
			err := ini.LoadFromString(tt.content)

			assertErrorMsg(t, err, ErrSyntaxError)
		})
	}

	
}

func ExampleIniParser_LoadFromString() {
	ini := New()
	ini.LoadFromString(iniContent)
	fmt.Println(ini.sections)
	// Output: map[database:map[file:"payroll.dat" port:143 server:192.0.2.62] owner:map[name:John Doe organization:Acme Inc.]]
}

func TestString(t *testing.T) {
	t.Run("nil sections", func(t *testing.T) {
		ini := IniParser{}
		got, err := ini.String()
		want := ""

		assertErrorMsg(t, err, ErrNullReference)
		assertEqualStrings(t, got, want)
	})
	t.Run("has no data yet", func(t *testing.T) {
		ini := New()
		got, err := ini.String()
		want := ""
		assertErrorMsg(t, err, ErrHasNoData)
		assertEqualStrings(t, got, want)
	})
	t.Run("has data", func(t *testing.T) {
		ini := New()
		ini.LoadFromString(iniContent)
		got := ini.sections
		stringContent, err := ini.String()
		ini.LoadFromString(stringContent)
		want := ini.sections
		
		assertNoErrorMsg(t, err)
		assertEqualSections(t, got, want)
	})
}

func ExampleIniParser_String() {
	ini := New()
	ini.LoadFromString(iniContent)
	oldContentSectionsMap := ini.sections
	stringContent, _ := ini.String()
	ini.LoadFromString(stringContent)
	newContentSectionsMap := ini.sections
	fmt.Println(reflect.DeepEqual(oldContentSectionsMap, newContentSectionsMap))
	// Output: true
}

func TestSaveToFile(t *testing.T) {
	t.Run("successful saving", func(t *testing.T) {
		ini := New()
		ini.LoadFromString(iniContent)
		got := ini.SaveToFile(successOutFilePath)
		oldContentSectionsMap := ini.sections
		strContent, _ := ini.LoadFromFile(exampleFilePath)
		ini.LoadFromString(strContent)
		newContentSectionsMap := ini.sections

		assertNoErrorMsg(t, got)
		assertEqualSections(t, oldContentSectionsMap, newContentSectionsMap)

	})
	t.Run("nil sections", func(t *testing.T) {
		ini := IniParser{}
		got := ini.SaveToFile(failOutFilePath)
		
		assertErrorMsg(t, got, ErrNullReference)
	})
	t.Run("has no data", func(t *testing.T) {
		ini := New()
		got := ini.SaveToFile(failOutFilePath)
		oldContentSectionsMap := ini.sections
		strContent, _ := ini.LoadFromFile(failOutFilePath)
		ini.LoadFromString(strContent)
		newContentSectionsMap := ini.sections

		assertErrorMsg(t, got, ErrHasNoData)
		assertEqualSections(t, oldContentSectionsMap, newContentSectionsMap)
	})
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
func assertEqualStrings(t testing.TB, got, want string) {
	t.Helper()
	if (!strings.EqualFold(got, want)) {
		t.Errorf("got %s want %s", got, want)
	}
}

func assertEqualSections(t testing.TB, got, want map[SectionName]Section) {
	t.Helper()
	if (!reflect.DeepEqual(got, want)) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertErrorMsg(t testing.TB, err, want error) {
	t.Helper()
	if err == nil {
		t.Fatal("didn't get error, expected to get an error")
	}

	if err.Error() != want.Error() {
		t.Errorf("got %q want %q", err.Error(), want)
	}
}

func assertNoErrorMsg(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("don't expect to get an error:%q", err.Error())
	}
}
