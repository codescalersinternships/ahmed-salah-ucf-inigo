package iniparser

var (
	ErrInvalidFilePath = IniParserError("couldn't find the file in the path you provided")
	ErrNullReference = IniParserError("you tried to access object that doesn't exist")
	ErrSectionNotExist = IniParserError("the section you tried to access doesn't exist")
	ErrKeyNotExist = IniParserError("the key you tried to access doesn't exist")
	ErrHasNoData = IniParserError("there is no data yet, you may didn't load data")
	ErrGlobalProperity = IniParserError("global keys are not allowed")
	ErrEmptySectionName = IniParserError("you should provide sectionName")
	ErrEmptyKey = IniParserError("you should provide key for the properity")
	ErrSyntaxError = IniParserError("syntax error, can't understand this line")
)

type IniParserError string

func (e IniParserError) Error() string {
	return string(e)
}