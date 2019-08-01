package ssz_test

import (
	"reflect"
	"testing"

	ssz "github.com/prysmaticlabs/go-ssz"
)

type fork struct {
	PreviousVersion [4]byte
	CurrentVersion  [4]byte
	Epoch           uint64
}

type nestedItem struct {
	Field1 []uint64
	Field2 *fork
	Field3 [3]byte
}

type varItem struct {
	Field2 []uint16
	Field3 []uint16
}

type nestedVarItem struct {
	Field1 []varItem
	Field2 uint64
}

var (
	forkExample = fork{
		PreviousVersion: [4]byte{1, 2, 3, 4},
		CurrentVersion:  [4]byte{5, 6, 7, 8},
		Epoch:           5,
	}
	nestedItemExample = nestedItem{
		Field1: []uint64{1, 2, 3, 4},
		Field2: &forkExample,
		Field3: [3]byte{32, 33, 34},
	}
	nestedVarItemExample = nestedVarItem{
		Field1: []varItem{},
		Field2: 5,
	}
	varItemExample = varItem{
		Field2: []uint16{},
		Field3: []uint16{2, 3},
	}
	varItemAmbiguous = varItem{
		Field3: []uint16{4, 5},
	}
)

func TestMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		input interface{}
		ptr   interface{}
	}{
		// Bool test cases.
		{input: true, ptr: new(bool)},
		{input: false, ptr: new(bool)},
		// Uint8 test cases.
		{input: byte(1), ptr: new(byte)},
		{input: byte(0), ptr: new(byte)},
		// Uint16 test cases.
		{input: uint16(100), ptr: new(uint16)},
		{input: uint16(232), ptr: new(uint16)},
		// Uint32 test cases.
		{input: uint32(1), ptr: new(uint32)},
		{input: uint32(1029391), ptr: new(uint32)},
		// Uint64 test cases.
		{input: uint64(5), ptr: new(uint64)},
		{input: uint64(23929309), ptr: new(uint64)},
		// Byte slice, byte array test cases.
		{input: [8]byte{1, 2, 3, 4, 5, 6, 7, 8}, ptr: new([8]byte)},
		{input: []byte{9, 8, 9, 8}, ptr: new([]byte)},
		// Basic type array test cases.
		{input: [12]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, ptr: new([12]uint64)},
		{input: [100]bool{true, false, true, true}, ptr: new([100]bool)},
		{input: [20]uint16{3, 4, 5}, ptr: new([20]uint16)},
		{input: [20]uint32{4, 5}, ptr: new([20]uint32)},
		{input: [20][2]uint32{{3, 4}, {5}, {8}, {9, 10}}, ptr: new([20][2]uint32)},
		// Basic type slice test cases.
		{input: []uint64{1, 2, 3}, ptr: new([]uint64)},
		{input: []bool{true, false, true, true, true}, ptr: new([]bool)},
		{input: []uint32{0, 0, 0}, ptr: new([]uint32)},
		{input: []uint32{92939, 232, 222}, ptr: new([]uint32)},
		// Struct decoding test cases.
		{input: forkExample, ptr: new(fork)},
		{input: nestedItemExample, ptr: new(nestedItem)},
		{input: nestedVarItemExample, ptr: new(nestedVarItem)},
		{input: varItemExample, ptr: new(varItem)},
		{input: varItemAmbiguous, ptr: new(varItem)},
		// Non-basic type slice/array test cases.
		{input: []fork{forkExample, forkExample}, ptr: new([]fork)},
		{input: [][]uint64{{4, 3, 2}, {1}, {0}}, ptr: new([][]uint64)},
		{input: [][][]uint64{{{1, 2}, {3}}, {{4, 5}}, {{0}}}, ptr: new([][][]uint64)},
		{input: [][3]uint64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, ptr: new([][3]uint64)},
		{input: [3][]uint64{{1, 2}, {4, 5, 6}, {7}}, ptr: new([3][]uint64)},
		{input: [][4]fork{{forkExample, forkExample, forkExample}}, ptr: new([][4]fork)},
		{input: [2]fork{forkExample, forkExample}, ptr: new([2]fork)},
		// Pointer-type test cases.
		{input: &forkExample, ptr: new(fork)},
		{input: &nestedItemExample, ptr: new(nestedItem)},
		{input: []*fork{&forkExample, &forkExample}, ptr: new([]*fork)},
		{input: []*nestedItem{&nestedItemExample, &nestedItemExample}, ptr: new([]*nestedItem)},
		{input: [2]*nestedItem{&nestedItemExample, &nestedItemExample}, ptr: new([2]*nestedItem)},
		{input: [2]*fork{&forkExample, &forkExample}, ptr: new([2]*fork)},
	}
	for _, tt := range tests {
		serializedItem, err := ssz.Marshal(tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if err := ssz.Unmarshal(serializedItem, tt.ptr); err != nil {
			t.Fatal(err)
		}
		output := reflect.ValueOf(tt.ptr)
		inputVal := reflect.ValueOf(tt.input)
		if inputVal.Kind() == reflect.Ptr {
			if !ssz.DeepEqual(output.Interface(), tt.input) {
				t.Errorf("Expected %v, received %v", tt.input, output.Interface())
			}
		} else {
			got := output.Elem().Interface()
			want := tt.input
			if !ssz.DeepEqual(want, got) {
				t.Errorf("Did not unmarshal properly: wanted %v, received %v", tt.input, output.Elem().Interface())
			}
		}
	}
}
