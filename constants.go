package iniparser

import "path/filepath"

var (
	failOutFilePath = filepath.Join(".", "ini_files", "fail_out_file.ini")
	successOutFilePath = filepath.Join(".", "ini_files", "success_out_file.ini")
	exampleFilePath = filepath.Join(".", "ini_files", "example.ini")
)

const (
	sectionLine = "sectionLine"
	propertyLine = "properityLine"
	commentLine = "commentLine"
	emptyLine = "emptyLine"
	unsportedLine = "unsportedLine"
)

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

var emptySectionNameIniContent = `[owner]
name = John Doe

[]
server = 192.0.2.62
`

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

var mapOfSections = map[SectionName]Section{
	SectionName("owner") : {Key("name") : "John Doe",
							Key("organization") : "Acme Inc."},

	SectionName("database") : {Key("server") : "192.0.2.62",
							   Key("port") : "143",
							   Key("file") : "\"payroll.dat\"",},
}