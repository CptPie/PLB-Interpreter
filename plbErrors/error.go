package plbErrors

import (
	"bytes"
	"fmt"
)

// PLBError is a custom error type for PLB errors
type PLBError struct {
	ErrorCode  string // Error code, e.g. E501
	Message    string // Error message, e.g. Invalid token type
	File       string // Filename where the error occurred
	LineNumber int    // LineNumber where the error occurred
	Column     int    // Column (number of character in the line) where the error occurred
	LineText   string // Literal line contents where the error occurred to display to the user
}

// NewPLBError creates a new PLBError
func NewPLBError(errorCode string, message string, file string, line int, column int, lineText string) *PLBError {
	return &PLBError{
		ErrorCode:  errorCode,
		Message:    message,
		File:       file,
		LineNumber: line,
		Column:     column,
		LineText:   lineText,
	}
}

// Error returns a string representation of the plbErrors
func (e *PLBError) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Error %s: %s\n", e.ErrorCode, e.Message))
	buffer.WriteString(fmt.Sprintf("Location: %s %d:%d\n", e.File, e.LineNumber, e.Column))
	buffer.WriteString(fmt.Sprintf("%s\n", e.LineText))
	buffer.WriteString(fmt.Sprintf("%s^\n", bytes.Repeat([]byte(" "), e.Column-1)))
	return buffer.String()
}
