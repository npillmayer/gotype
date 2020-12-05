package pmmpost

import (
	"fmt"
	"strings"
	"testing"

	"github.com/npillmayer/gotype/backend/gfx"
	"github.com/npillmayer/gotype/core/arithmetic"
	"github.com/npillmayer/gotype/core/path"
	"github.com/npillmayer/schuko/testadapter"
	"github.com/npillmayer/schuko/tracing"
)

func TestInit0(t *testing.T) {
	testconfig := testadapter.New()
	config.Initialize(testconfig)
	testconfig.Set("tracinginterpreter", "Info")
	testconfig.Set("tracinggraphics", "Info")
	config.DefaultTracing()
	tracing.InterpreterTracer.Infof("Interpreter tracer is alive!")
}

func TestCreateInterpreter(t *testing.T) {
	pmmp := NewPMMPostInterpreter(true, nil)
	if pmmp == nil {
		t.Error("Could not create PMMPost interpreter")
	}
}

func TestStatement1(t *testing.T) {
	pmmp := NewPMMPostInterpreter(true, pictureCallback)
	check(t, pmmp.ParseStatements([]byte("x0=1;")))
	v, value := pmmp.ValueString("x0")
	t.Logf("%s = %v", v, value)
	if !eq(value, 1) {
		t.Errorf("x0 expected to be 1")
	}
	v, value = pmmp.ValueString("p1") // tag p is predeclared as type pair
	t.Logf("%s = %v", v, value)
	if !eq(value, "<unset pair>") {
		t.Errorf("p1 expected to be <unset pair>")
	}
	check(t, pmmp.ParseStatements([]byte("p1=origin;")))
	v, value = pmmp.ValueString("p1")
	t.Logf("%s = %v", v, value)
	if !eq(value, "(0,0)") {
		t.Errorf("p1 expected to be (0,0)")
	}
}

func TestPathAssign(t *testing.T) {
	pmmp := NewPMMPostInterpreter(true, pictureCallback)
	check(t, pmmp.ParseStatements([]byte("P := origin .. (1,1);")))
	varval, err := pmmp.ValueOf("P")
	if err != nil {
		t.Errorf("cannot get value of P: %s", err.Error())
	}
	p := varval.AsPath()
	pstr := path.PathAsString(p, p.Controls)
	t.Logf("path value of %s = %v", varval.VariableFullName, pstr)
}

func TestPicture(t *testing.T) {
	pmmp := NewPMMPostInterpreter(true, pictureCallback)
	shipoutOk = false
	check(t, pmmp.ParseStatements([]byte(`
    	beginfig("testpic", 3cm, 3cm);
			P := origin .. (1,1);
			draw P;
		endfig;`)))
	varval, err := pmmp.ValueOf("P")
	if err != nil {
		t.Errorf("cannot get value of P: %s", err.Error())
	}
	p := varval.AsPath()
	if p.N() != 2 {
		t.Errorf("path should be of length 2, is %d", p.N())
	}
	pstr := path.PathAsString(p, p.Controls)
	t.Logf("path value of %s = %v", varval.VariableFullName, pstr)
	if !shipoutOk {
		t.Error("no picture shipped out")
	}
}

func TestProg1(t *testing.T) {
	pmmp := NewPMMPostInterpreter(true, pictureCallback)
	//tracing.InterpreterTracer.SetTraceLevel(tracing.LevelDebug)
	check(t, pmmp.ParseStatements([]byte(`
    	pair a, b, c;
        xpart a=xpart c; ypart a=ypart b;
        b=(1,1);
        c=(2,2);
        show a;
	`)))
	varval, err := pmmp.ValueOf("a")
	if err != nil {
		t.Errorf("cannot get value of a: %s", err.Error())
	}
	a := varval.AsPair()
	if !a.YPart().Equals(arithmetic.ConstOne) {
		t.Errorf("a should be (2,1), is %v", a)
	}
}

func TestSyntaxError(t *testing.T) {
	pmmp := NewPMMPostInterpreter(true, pictureCallback)
	tracing.InterpreterTracer.SetTraceLevel(tracing.LevelDebug)
	errs := pmmp.ParseStatements([]byte("hello world;"))
	if len(errs) == 0 {
		t.Error("syntax error not detected")
	}
}

// ----------------------------------------------------------------------

func eq(a string, b interface{}) bool {
	return strings.EqualFold(a, fmt.Sprintf("%v", b))
}

func check(t *testing.T, errs []error) {
	if errs != nil {
		for _, err := range errs {
			t.Logf(">>> %s", err.Error())
		}
		t.Fail()
	}
}

var shipoutOk bool

func pictureCallback(pic *gfx.Picture) {
	tracing.GraphicsTracer.Infof("SHIPPING PICTURE '%s'\n", pic.Name)
	shipoutOk = true
}
