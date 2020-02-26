package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type Traceable struct {
	Stack []uintptr
}

func (t *Traceable) StackTrace() string {
	var stackTrace string

	for _, pc := range t.Stack {
		stackTrace += newTrace(pc).String()
	}

	return stackTrace
}

type Trace struct {
	File           string
	LineNumber     int
	Name           string
	Package        string
	ProgramCounter uintptr
}

func newTrace(pc uintptr) *Trace {
	if pc == 0 {
		return nil
	}

	f := runtime.FuncForPC(pc)
	if f == nil {
		return nil
	}

	t := &Trace{
		Name:           f.Name(),
		ProgramCounter: pc,
	}
	t.File, t.LineNumber = f.FileLine(pc - 1)

	if lastslash := strings.LastIndex(t.Name, "/"); lastslash >= 0 {
		t.Package += t.Name[:lastslash] + "/"
		t.Name = t.Name[lastslash+1:]
	}
	if period := strings.Index(t.Name, "."); period >= 0 {
		t.Package += t.Name[:period]
		t.Name = t.Name[period+1:]
	}

	t.Name = strings.Replace(t.Name, "·", ".", -1)

	return t
}

func (t *Trace) String() string {
	str := fmt.Sprintf("%s:%d (0x%x)\n", t.File, t.LineNumber, t.ProgramCounter)

	return str + fmt.Sprintf("\t%s", t.Name)
}
