package pmmpost

import (
	"fmt"
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/config"
	"github.com/npillmayer/gotype/core/config/testadapter"
	"github.com/npillmayer/gotype/core/config/tracing"
)

func TestInit0(t *testing.T) {
	testconfig := testadapter.New()
	config.Initialize(testconfig)
	testconfig.Set("tracinginterpreter", "Info")
	config.DefaultTracing()
	tracing.InterpreterTracer.Infof("Interpreter tracer is alive")
}

func TestCreateInterpreter(t *testing.T) {
	pmmp := NewPMMPostInterpreter(true)
	if pmmp == nil {
		t.Error("Could not create PMMPost interpreter")
	}
}

func TestStatement1(t *testing.T) {
	pmmp := NewPMMPostInterpreter(true)
	pmmp.ParseStatements([]byte("x0=1"))
	v, value := pmmp.Value("x0")
	t.Logf("%s = %v", v, value)
	if !eq(value, 1) {
		t.Errorf("x0 expected to be 1")
	}
	pmmp.ParseStatements([]byte("pair p1"))
	v, value = pmmp.Value("p1")
	t.Logf("%s = %v", v, value)
	if !eq(value, "<unset pair>") {
		t.Errorf("p1 expected to be <unset pair>")
	}
	pmmp.ParseStatements([]byte("p1=origin"))
	v, value = pmmp.Value("p1")
	t.Logf("%s = %v", v, value)
	if !eq(value, "(0,0)") {
		t.Errorf("p1 expected to be (0,0)")
	}
}

// ----------------------------------------------------------------------

func eq(a string, b interface{}) bool {
	return strings.EqualFold(a, fmt.Sprintf("%v", b))
}
