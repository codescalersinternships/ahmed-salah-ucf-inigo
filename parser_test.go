package iniparser

import (
	"reflect"
	"testing"
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
		ini := NewIniParser()
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
		ini := NewIniParser()
		ini.LoadFromString(iniContent)
		got := ini.GetSectionNames()
		want := []string{"database", "owner"}
		assertEqualLists(t, got, want)
	})
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

func TestLoadFromFile(t *testing.T) {
	t.Run("valid file path", func(t *testing.T) {
		filePath := "./example.ini"
		ini := NewIniParser()

		got, err := ini.LoadFromFile(filePath)
		want := iniContent

		assertNoErrorMsg(t, err)
		assertEqualStrings(t, got, want)
	})

	t.Run("invalid file path", func(t *testing.T) {
		filepath := "/invalid/file/path/example.ini"
		ini := NewIniParser()
		_, err := ini.LoadFromFile(filepath)

		assertErrorMsg(t, err, ErrInvalidFilePath)
	})

	t.Run("invalid ini file syntax", func(t *testing.T) { // TO-Do
		
	})

}

func TestLoadFromString(t *testing.T) {
	t.Run("nil sections", func(t *testing.T) {
		data := iniContent

		ini := IniParser{}
		err := ini.LoadFromString(data)

		assertErrorMsg(t, err, ErrNullReference)
	})
	t.Run("spaces trimming", func(t *testing.T) {
		data := spacedIniContent

		ini := NewIniParser()
		err := ini.LoadFromString(data)
		got := ini.GetSections()
		want := mapOfSections

		assertNoErrorMsg(t, err)
		assertEqualSections(t, got, want)
	})

	t.Run("empty lines", func(t *testing.T) {
		data := emptyLinesIniContent

		ini := IniParser{sections: map[SectionName]Section{}}
		err := ini.LoadFromString(data)
		got := ini.GetSections()
		want := mapOfSections

		assertNoErrorMsg(t, err)
		assertEqualSections(t, got, want)
	})
}

func assertEqualStrings(t testing.TB, got, want string) {
	t.Helper()
	if (!reflect.DeepEqual(got, want)) {
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
		t.Fatal("expected to get an error.")
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
