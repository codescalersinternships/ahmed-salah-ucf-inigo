# INI Parser

INI Parser is a simple Go package that provides functionality for parsing INI format

## Overview

An INI Parser is a package that read/write INI data from/to file streams and strings.

It can also do the following:

1. List all section names
2. serialize convert into a dictionary/map
3. gets the value of key key in section
4. sets a key in section section_name to value value


## How to use parser?

First make the parser object:
```go
iniParser := IniParser.New()
```

To parse from a file:
```go
iniParser.LoadFromFile("file.ini")
```

Or from a string:
```go
iniText := `[section]
domain = wikipedia.org
[section.subsection]
foo = bar`
iniParser.LoadFromString(iniText)
```

You can get the sections names of your parser object:
```go
sliceOfStrings := iniParser.GetSectionNames()
```

You can get your ini file as a map:
```go
mapOfSections := iniParser.GetSections()
```

You can set a property:
```go
iniParser.Set("sectionName", "key", "value")
```

You can get a property:
```go
value := iniParser.Get("sectionName", "key")
```

You can return the iniParser object as string:
```go
parserAsString := iniParser.String()
```

You can save your object in a file:
```go
iniParser.SaveToFile("filePath")
```