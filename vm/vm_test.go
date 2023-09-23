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
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15/3) * 2 - 10", 50},
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
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!0", false},
		{"!-10", false},
		{"!(if (false) { 5; })", true},
	}
	runVmTests(t, tests)
}

func TestStringExpression(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}
	runVmTests(t, tests)
}

func TestConditions(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else {20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", Null},
		{"if (false) { 10 }", Null},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}
	runVmTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = 1; let two = one + one; one + two", 3},
	}
	runVmTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{`[]`, []int{}},
		{`[1,2,3]`, []int{1, 2, 3}},
		{`[1+3, 3*4, 5+6]`, []int{4, 12, 11}},
	}
	runVmTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{
			"{}", map[object.HashKey]int64{},
		},
		{
			"{1:2, 2:3}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 1}).HashKey(): 2,
				(&object.Integer{Value: 2}).HashKey(): 3,
			},
		},
		{
			"{1+1: 2 * 2, 3 + 3: 4 * 4}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 2}).HashKey(): 4,
				(&object.Integer{Value: 6}).HashKey(): 16,
			},
		},
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
	case string:
		err := testStringObject(input, string(expected), actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case []int:
		err := testArrayObject(input, expected, actual)
		if err != nil {
			t.Errorf("testArrayObject failed: %s", err)
		}
	case map[object.HashKey]int64:
		err := testHashObject(input, expected, actual)
		if err != nil {
			t.Errorf("testHashObject failed: %s", err)
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null: %T (%+v)", actual, actual)
		}
	default:
		t.Errorf("unknown object type: got=%T, want=%T", actual, expected)
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

func testStringObject(input string, expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("`%s`: object is not string. got=%T (%+v)", input, actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("`%s`: object has wrong value. got=%v, want=%v", input, result.Value, expected)
	}
	return nil
}

func testArrayObject(input string, expected []int, actual object.Object) error {
	result, ok := actual.(*object.Array)
	if !ok {
		return fmt.Errorf("`%s`: object is not array. got=%T (%+v)", input, actual, actual)
	}

	if len(result.Elements) != len(expected) {
		return fmt.Errorf("wrong number of elements: want=%d, got=%d",
			len(expected), len(result.Elements))
	}
	for i, expectedElem := range expected {
		err := testIntegerObject(input+fmt.Sprintf("[%d]", i), int64(expectedElem), result.Elements[i])
		if err != nil {
			return fmt.Errorf("testIntegerObject failed: %s", err)
		}
	}
	return nil
}

func testHashObject(input string, expected map[object.HashKey]int64, actual object.Object) error {
	hash, ok := actual.(*object.Hash)
	if !ok {
		return fmt.Errorf("object is not hash, got=%T (%+v)", actual, actual)
	}

	if len(hash.Pairs) != len(expected) {
		return fmt.Errorf("hash has wrong number of pairs. want=%d, got=%d",
			len(expected), len(hash.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := hash.Pairs[expectedKey]
		if !ok {
			return fmt.Errorf("no pair for given key in Pairs")
		}

		err := testIntegerObject(input+fmt.Sprintf("[%v]", expectedKey), expectedValue, pair.Value)
		if err != nil {
			return fmt.Errorf("testIntegerObject failed: %s", err)
		}
	}
	return nil
}
