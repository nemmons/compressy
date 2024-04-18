package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"maps"
	"reflect"
	"testing"
)

func TestGetCharFrequencies(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput map[rune]uint32
	}{
		{
			"aaabbc",
			map[rune]uint32{
				'a': 3,
				'b': 2,
				'c': 1,
			},
		},
		{
			"a\n1!語",
			map[rune]uint32{
				'a':  1,
				'\n': 1,
				'1':  1,
				'!':  1,
				'語':  1,
			},
		},
	}

	for _, testCase := range testCases {
		output := getCharFrequencies(testCase.input)

		for key, expectedValue := range testCase.expectedOutput {
			if output[key] != expectedValue {
				t.Fatalf(`In test case %s: Count for rune '%c' = %v, expected %v`, testCase.input, key, output[key], expectedValue)
			}
		}
	}
}

func TestBuildTreeSimple(t *testing.T) {
	frequencies := map[rune]uint32{
		'a': 1,
		'b': 2,
		'c': 4,
		'語': 5,
	}

	actualTree := buildTree(frequencies)

	expectedTree := InternalNode{
		left: LeafNode{
			char:   '語',
			weight: 5,
		},
		right: InternalNode{
			left: InternalNode{
				left: LeafNode{
					char:   'a',
					weight: 1,
				},
				right: LeafNode{
					char:   'b',
					weight: 2,
				},
				weight: 3,
			},
			right: LeafNode{
				char:   'c',
				weight: 4,
			},
			weight: 7,
		},
		weight: 12,
	}

	equal := reflect.DeepEqual(expectedTree, actualTree)
	if !equal {
		t.Fatalf(`Tree building failed!`)
	}
}

// https://opendsa-server.cs.vt.edu/ODSA/Books/CS3/html/Huffman.html#building-huffman-coding-trees
func TestBuildTree(t *testing.T) {
	frequencies := map[rune]uint32{
		'c': 32,
		'd': 42,
		'e': 120,
		'k': 7,
		'l': 42,
		'm': 24,
		'u': 37,
		'z': 2,
	}

	actualTree := buildTree(frequencies)

	expectedTree := InternalNode{
		left: LeafNode{
			char:   'e',
			weight: 120,
		},
		right: InternalNode{
			left: InternalNode{
				left: LeafNode{
					char:   'u',
					weight: 37,
				},
				right: LeafNode{
					char:   'd',
					weight: 42,
				},
				weight: 79,
			},
			right: InternalNode{
				left: LeafNode{
					char:   'l',
					weight: 42,
				},
				right: InternalNode{
					left: LeafNode{
						char:   'c',
						weight: 32,
					},
					right: InternalNode{
						left: InternalNode{
							left: LeafNode{
								char:   'z',
								weight: 2,
							},
							right: LeafNode{
								char:   'k',
								weight: 7,
							},
							weight: 9,
						},
						right: LeafNode{
							char:   'm',
							weight: 24,
						},
						weight: 33,
					},
					weight: 65,
				},
				weight: 107,
			},
			weight: 186,
		},
		weight: 306,
	}

	equal := reflect.DeepEqual(expectedTree, actualTree)
	if !equal {
		t.Fatalf(`Tree building failed!`)
	}
}

func TestGenerateEncodingsSimple(t *testing.T) {
	frequencies := map[rune]uint32{
		'a': 1,
		'b': 2,
		'c': 4,
		'語': 5,
	}

	tree := buildTree(frequencies)
	encodings := generateEncodings(tree, []bool{})

	expected := map[rune]HuffmanEncoding{
		'a': {
			Value: 4, //0b100
			Bits:  3,
		},
		'b': {
			Value: 5, //0b101
			Bits:  3,
		},
		'c': {
			Value: 3, //0b11
			Bits:  2,
		},
		'語': {
			Value: 0, //0b0
			Bits:  1,
		},
	}

	if !maps.Equal(expected, encodings) {
		t.Fatalf(`Encoding failed!`)
	}

}

func TestEncodingSerialization(t *testing.T) {
	frequencies := getCharFrequencies("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")
	tree := buildTree(frequencies)
	characterEncodings := generateEncodings(tree, []bool{})

	var encodedBytes bytes.Buffer
	encoder := gob.NewEncoder(&encodedBytes)
	err := encoder.Encode(characterEncodings)
	if err != nil {
		t.Fatalf("Encoding serialization failed: %v", err)
	}

	var decodedCharacterEncodings map[rune]HuffmanEncoding
	decoder := gob.NewDecoder(&encodedBytes)
	err = decoder.Decode(&decodedCharacterEncodings)
	if err != nil {
		t.Fatalf("Encoding deserialization failed: %v", err)
	}

	if !reflect.DeepEqual(characterEncodings, decodedCharacterEncodings) {
		t.Fatalf("Encoding serialization/deserialization failed")
	}

}

//func TestPackBits(t *testing.T) {
//	input := []byte{
//		0b00000001,
//		0b00001010,
//		0b00000000,
//		0b00000111,
//	}
//
//	expected := []byte{
//		0b11010011,
//		0b10000000,
//	}
//
//	actual := packBits(input)
//
//	if !reflect.DeepEqual(expected, actual) {
//		panic("fail")
//	}
//}
//
//func TestPackBitsMore(t *testing.T) {
//	input := []byte{
//		0b00000001,
//		0b00001010,
//		0b10001000,
//		0b00000111,
//	}
//
//	expected := []byte{
//		0b11010100,
//		0b01000111,
//	}
//
//	actual := packBits(input)
//
//	if !reflect.DeepEqual(expected, actual) {
//		panic("fail")
//	}
//}
//
//func TestPackBitsMoreAgain(t *testing.T) {
//	input := []byte{
//		0b01000001,
//		0b00001010,
//		0b10001000,
//		0b00000111,
//		0b00100100,
//		0b00000000,
//		0b00000001,
//	}
//
//	expected := []byte{
//		0b10000011,
//		0b01010001,
//		0b00011110,
//		0b01000100,
//	}
//
//	actual := packBits(input)
//
//	if !reflect.DeepEqual(expected, actual) {
//		panic("fail")
//	}
//}

func TestNthBit(t *testing.T) {
	if 1 != nthBit(0b00000001, 8) {
		panic("fail")
	}
	if 0 != nthBit(0b00000001, 7) {
		panic("fail")
	}
	if 1 != nthBit(0b00000011, 7) {
		panic("fail")
	}
	if 1 != nthBit(0b00000010, 7) {
		panic("fail")
	}
	if 1 != nthBit(0b01000000, 2) {
		panic("fail")
	}
	if 0 != nthBit(0b00000010, 4) {
		panic("fail")
	}
	if 1 != nthBit(0b10000000, 1) {
		panic("fail")
	}
}

func TestBuildByte(t *testing.T) {
	input := []bool{true, true, false, true}
	expected := HuffmanEncoding{
		Value: 13,
		Bits:  4,
	}
	actual := buildByte(input)
	if expected != actual {
		fmt.Printf("Found %v, expected %v", actual, expected)
		panic("fail")
	}
}
