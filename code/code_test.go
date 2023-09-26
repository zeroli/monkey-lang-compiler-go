package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpSub, []int{}, []byte{byte(OpSub)}},
		{OpMul, []int{}, []byte{byte(OpMul)}},
		{OpDiv, []int{}, []byte{byte(OpDiv)}},
		{OpPop, []int{}, []byte{byte(OpPop)}},
		{OpTrue, []int{}, []byte{byte(OpTrue)}},
		{OpFalse, []int{}, []byte{byte(OpFalse)}},
		{OpEqual, []int{}, []byte{byte(OpEqual)}},
		{OpNotEqual, []int{}, []byte{byte(OpNotEqual)}},
		{OpGreaterThan, []int{}, []byte{byte(OpGreaterThan)}},
		{OpMinus, []int{}, []byte{byte(OpMinus)}},
		{OpBang, []int{}, []byte{byte(OpBang)}},
		{OpJumpNotTruthy, []int{65534}, []byte{byte(OpJumpNotTruthy), 255, 254}},
		{OpJump, []int{65534}, []byte{byte(OpJump), 255, 254}},
		{OpNull, []int{}, []byte{byte(OpNull)}},
		{OpSetGlobal, []int{1}, []byte{byte(OpSetGlobal), 0, 1}},
		{OpGetGlobal, []int{1}, []byte{byte(OpGetGlobal), 0, 1}},
		{OpArray, []int{1}, []byte{byte(OpArray), 0, 1}},
		{OpHash, []int{10}, []byte{byte(OpHash), 0, 10}},
		{OpIndex, []int{}, []byte{byte(OpIndex)}},
		{OpCall, []int{10}, []byte{byte(OpCall), 10}},
		{OpReturnValue, []int{}, []byte{byte(OpReturnValue)}},
		{OpReturn, []int{}, []byte{byte(OpReturn)}},
		{OpGetLocal, []int{200}, []byte{byte(OpGetLocal), 200}},
		{OpSetLocal, []int{210}, []byte{byte(OpSetLocal), 210}},
		{OpGetBuiltin, []int{100}, []byte{byte(OpGetBuiltin), 100}},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want=%d, got=%d",
				len(tt.expected), len(instruction))
		}

		for i, b := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d",
					i, b, instruction[i])
			}
		}
	}
}

func TestInstructionString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpSub),
		Make(OpMul),
		Make(OpDiv),
		Make(OpPop),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpTrue),
		Make(OpFalse),
		Make(OpEqual),
		Make(OpNotEqual),
		Make(OpGreaterThan),
		Make(OpMinus),
		Make(OpBang),
		Make(OpJumpNotTruthy, 10),
		Make(OpJump, 300),
		Make(OpNull),
		Make(OpSetGlobal, 12),
		Make(OpGetGlobal, 30),
		Make(OpArray, 100),
		Make(OpHash, 3000),
		Make(OpIndex),
		Make(OpCall, 255),
		Make(OpReturnValue),
		Make(OpReturn),
		Make(OpGetLocal, 10),
		Make(OpSetLocal, 20),
		Make(OpGetBuiltin, 100),
	}

	expected := `0000 OpAdd
0001 OpSub
0002 OpMul
0003 OpDiv
0004 OpPop
0005 OpConstant 2
0008 OpConstant 65535
0011 OpTrue
0012 OpFalse
0013 OpEqual
0014 OpNotEqual
0015 OpGreaterThan
0016 OpMinus
0017 OpBang
0018 OpJumpNotTruthy 10
0021 OpJump 300
0024 OpNull
0025 OpSetGlobal 12
0028 OpGetGlobal 30
0031 OpArray 100
0034 OpHash 3000
0037 OpIndex
0038 OpCall 255
0040 OpReturnValue
0041 OpReturn
0042 OpGetLocal 10
0044 OpSetLocal 20
0046 OpGetBuiltin 100
`

	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf("instruction wrongly formatted.\nwant=%q\ngot =%q",
			expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		byteRead int
	}{
		{OpConstant, []int{65535}, 2},
		{OpGetLocal, []int{200}, 1},
		{OpSetLocal, []int{210}, 1},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)
		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		operandsRead, n := ReadOperands(def, instruction[1:])
		if n != tt.byteRead {
			t.Fatalf("n wrong. want=%d, got=%d", tt.byteRead, n)
		}

		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Errorf("operand wrong. want=%d, got=%d", want, operandsRead[i])
			}
		}
	}
}
