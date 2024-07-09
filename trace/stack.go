package trace

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"runtime"

	"github.com/alioth-center/infrastructure/utils/values"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// Stack generates a formatted stack trace of the calling goroutine.
// The stack trace starts from the function that calls Stack and ascends up the call stack.
// Each entry in the stack trace includes the file name, line number, and the address of the call.
// If the function can read the source file, it also includes the line of source code at the call site.
// The skip parameter allows skipping a number of stack frames to start the trace at a higher level in the call stack.
//
// Parameters:
//
//	skip (int): The number of stack frames to skip before starting the stack trace.
//
// Returns:
//
//	[]byte: A byte slice containing the formatted stack trace.
func Stack(skip int) []byte {
	var (
		lines    [][]byte            // Lines of the source file being processed
		lastFile string              // The last file processed, to avoid reading the same file multiple times
		buf      = new(bytes.Buffer) // Buffer to hold the formatted stack trace
	)

	for i := skip; ; i++ { // Infinite loop, will break when runtime.Caller returns false
		pc, file, line, ok := runtime.Caller(i) // Get caller info at the current level
		if !ok {                                // If no more callers, break the loop
			break
		}
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc) // Write file, line, and PC to buffer
		if file != lastFile {                                     // If this is a new file
			data, err := os.ReadFile(file) // Read the file
			if err != nil {                // If there's an error reading the file, skip to the next iteration
				continue
			}
			lines = bytes.Split(data, []byte{'\n'}) // Split the file into lines
			lastFile = file                         // Update lastFile to the current file
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line)) // Write the function name and source line to buffer
	}
	return buf.Bytes() // Return the formatted stack trace
}

// source returns a trimmed line of source code from a slice of lines based on the specified line number.
// It is used to extract the specific line of code that corresponds to a stack trace entry.
// If the line number is out of range, it returns a placeholder indicating unknown source.
//
// Parameters:
//
//	lines ([][]byte): A slice of byte slices, each representing a line of source code from a file.
//	n     (int):      The line number to retrieve from the lines slice, adjusted to be zero-indexed.
//
// Returns:
//
//	[]byte: The trimmed line of source code at the specified line number, or a placeholder if out of range.
func source(lines [][]byte, n int) []byte {
	n--
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function retrieves the name of the function at the specified program counter (PC) address.
// It formats the function name by trimming the package path and replacing center dots with periods.
// If the function cannot be found, it returns a placeholder indicating unknown function.
//
// Parameters:
//
//	pc (uintptr): The program counter address of the function to retrieve the name for.
//
// Returns:
//
//	[]byte: The name of the function at the specified PC, formatted and trimmed, or a placeholder if not found.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.ReplaceAll(name, centerDot, dot)
	return name
}

// Caller returns a string representing the file and line number of the caller
// of the function that called Caller, adjusted by the skip parameter. The skip
// parameter is used to ascend the call stack by a number of frames, allowing
// the retrieval of callers further up the call stack. The resulting string is
// in the format "file:line", where "file" is the full path to the source file
// containing the caller and "line" is the line number within that file.
//
// Parameters:
//
//	skip (int): The number of stack frames to ascend, with 0 identifying the
//	            caller of Caller, 1 identifying the caller's caller, and so on.
//
// Returns:
//
//	caller (string): A string in the format "file:line" indicating the source
//	                 file and line number of the caller adjusted by skip.
func Caller(skip int) (caller string) {
	_, file, line, _ := runtime.Caller(2 + skip)
	return values.BuildStrings(file, ":", values.IntToString(line))
}

// FunctionLocation returns a string representing the file and line number
// of the location where the provided function is defined. The resulting
// string is in the format "file:line", where "file" is the full path to
// the source file containing the function definition and "line" is the
// line number within that file.
//
// Parameters:
//
//	fn (any): The function whose location is to be determined.
//
// Returns:
//
//	location (string): A string in the format "file:line" indicating the
//	                   source file and line number of the function definition.
func FunctionLocation(fn any) (location string) {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		return "provided argument is not a function"
	}

	pc := val.Pointer()
	fnInfo := runtime.FuncForPC(pc)
	if fnInfo == nil {
		return "unable to retrieve function information"
	}

	file, line := fnInfo.FileLine(pc)
	return values.BuildStrings(file, ":", values.IntToString(line))
}
