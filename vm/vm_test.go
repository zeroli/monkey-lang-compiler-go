package vm

import (
	"fmt"
	"monkey/ast"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1+2", 3},
		{"1-2", -1},
		{"1*2", 2},
		{"4/2", 2},
		{"50/2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2*2*2*2*2", 32},
		{"5*2 + 10", 20},
		{"5 +2 * 10", 25},
		{"5 * (2+10)", 60},
	}
	runVmTests(t, tests)
}

func TestBooleanExpression(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1<2", true},
		{"1>2", false},
		{"1<1", false},
		{"1>1", false},
		{"1==1", true},
		{"1!=1", false},
		{"1==2", false},
		{"1!=2", true},
		{"true==true", true},
		{"false==false", true},
		{"true==false", false},
		{"true!=false", true},
		{"false!=true", true},
		{"(1<2) == true", true},
		{"(1<2) == false", false},
		{"(1>2) == true", false},
		{"(1>2) == false", true},
	}
	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.LastPoppedStackElem()
		testExpectedObject(t, tt.input, tt.expected, stackElem)
	}
}

func testExpectedObject(t *testing.T, input string, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(input, int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(input, bool(expected), actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	}
}

func testIntegerObject(input string, expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("`%s`: object is not Integer. got=%T (%+v)",
			input, actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("`%s`: object has wrong value. got=%d, want=%d",
			input, result.Value, expected)
	}

	return nil
}

func testBooleanObject(input string, expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("`%s`: object is not Boolean. got=%T (%+v)", input, actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("`%s`: object has wrong value. got=%v, want=%v", input, result.Value, expected)
	}
	return nil
}
