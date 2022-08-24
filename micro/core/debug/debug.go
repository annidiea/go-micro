package debug

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	// DebugMode indicates gin mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates gin mode is release.
	ReleaseMode = "release"
	// TestMode indicates gin mode is test.
	TestMode = "test"
)
const (
	debugCode = iota
	releaseCode
	testCode
)

var (
	rpcMode  = debugCode
	modeName = DebugMode
)

var (
	DefaultWriter      io.Writer = os.Stdout
	DefaultErrorWriter io.Writer = os.Stderr
)

var printPrefix = "[debug]"

func SetMode(value string) {
	if value == "" {
		value = DebugMode
	}
	switch value {
	case DebugMode:
		rpcMode = debugCode
	case ReleaseMode:
		rpcMode = releaseCode
	case TestMode:
		rpcMode = testCode
	default:
		panic("shopstar-micro mode unknown: " + value + " (available mode: debug release test)")
	}

	modeName = value
}

func SetPrintPrefix(prefix string) {
	printPrefix = prefix
}

func PrintDirExePos(source string, format string, values ...interface{}) {
	DD("["+source+"]: "+format, values...)
}

func PrintErrDirExePos(source string, err error, format string, values ...interface{}) {
	DEDetails(err, "["+source+"]: "+format, values...)
}

func DD(format string, values ...interface{}) {
	if IsDebugging() {

		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}

		fmt.Fprintf(DefaultWriter, printPrefix+" "+format, values...)
	}
}

func DE(err error) {
	if err != nil {
		if IsDebugging() {
			fmt.Fprintf(DefaultErrorWriter, printPrefix+" [ERROR] %v\n", err)
		}
	}
}

func DEDetails(err error, format string, a ...interface{}) {
	if err != nil {
		if IsDebugging() {
			DD(format+"[ERROR] "+err.Error(), a...)
		}
	}
}

func IsDebugging() bool {
	return rpcMode == debugCode
}
