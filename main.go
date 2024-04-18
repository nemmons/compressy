package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"maps"
	"os"
	"strings"
)

const PsuedoEOF = "■" //mark the end of encoded test, making it easier to discard leftover bits afterward
const Divider = "⍭"   //to divide the encoded file into two sections: frequency map and encoded text

func main() {
	mode := os.Args[1] //['compress'/'decompress']

	var filename string
	var defaultFilename string
	if mode == "decompress" {
		defaultFilename = "compressed.txt"
	} else {
		defaultFilename = "input.txt"
	}
	flag.StringVar(&filename, "file", defaultFilename, "The name of the file containing the text to be processed")
	flag.Parse()

	input, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if mode == "compress" {
		err = compress(string(input) + PsuedoEOF)
		if err != nil {
			panic(err)
		}
	} else if mode == "decompress" {
		err = decompress(string(input))
		if err != nil {
			panic(err)
		}
	} else {
		panic(fmt.Sprintf("Unknown Mode %s", mode))
	}
}

func compress(input string) error {
	frequencies := getCharFrequencies(input)

	tree := buildTree(frequencies)

	//fmt.Printf("%#v\n", frequencies)
	var characterEncodings map[rune]HuffmanEncoding
	characterEncodings = generateEncodings(tree, []bool{})
	//fmt.Printf("%#v\n", characterEncodings)

	var output bytes.Buffer
	encoder := gob.NewEncoder(&output)
	err := encoder.Encode(frequencies)
	if err != nil {
		return err
	}

	compressedText := encodeText(input, characterEncodings)

	//fmt.Printf("%#v\n", compressedText)

	output.WriteString(Divider)
	output.Write(compressedText)

	err = os.WriteFile("compressed.txt", output.Bytes(), 0666)
	if err != nil {
		return err
	}

	return nil
}

func encodeText(input string, encodings map[rune]HuffmanEncoding) []byte {
	encodedText := make([]HuffmanEncoding, len(input))
	for i, char := range input {
		encodedText[i] = encodings[char]
	}
	return packBits(encodedText)
}

func decompress(input string) error {
	frequencies, message, found := strings.Cut(input, Divider)
	if !found {
		panic("Could not parse input file - missing separator between encodings and message!")
	}

	var characterFrequencies map[rune]uint32
	decoder := gob.NewDecoder(strings.NewReader(frequencies))
	err := decoder.Decode(&characterFrequencies)
	if err != nil {
		return err

	}

	//fmt.Printf("%#v\n", []byte(message))
	result := decode(message, characterFrequencies)

	//fmt.Printf("%#v\n", result)

	trimmed, _, found := strings.Cut(result, PsuedoEOF)
	if !found {
		panic("Could not parse translated output- missing psuedo-EOF character!")
	}

	//fmt.Println(trimmed)
	var output bytes.Buffer
	output.WriteString(trimmed)
	err = os.WriteFile("decompressed.txt", output.Bytes(), 0666)
	if err != nil {
		panic(err)
	}
	return nil
}

func decode(message string, frequencies map[rune]uint32) string {
	tree := buildTree(frequencies) //TODO - lift this up, probably

	//TODO - this is not performant, probably because we're copying the tree so much? Maybe this should be flattened into an array where each element points to its children, for faster traversal?
	workingNode := tree

	var output string

	for _, code := range []byte(message) {
		for i := range 8 {
			if nthBit(code, i+1) == 0 {
				workingNode = workingNode.(InternalNode).left
			} else {
				workingNode = workingNode.(InternalNode).right
			}
			if workingNode.isLeaf() {
				output += string(workingNode.(LeafNode).char)
				workingNode = tree
			}
		}
	}
	return output
}

func getCharFrequencies(input string) map[rune]uint32 {
	results := make(map[rune]uint32)
	for _, val := range input {
		count, ok := results[val]
		if ok {
			results[val] = count + 1
		} else {
			results[val] = 1
		}
	}

	return results
}

func generateEncodings(node TreeNode, path []bool) map[rune]HuffmanEncoding {
	if node.isLeaf() {
		return map[rune]HuffmanEncoding{
			node.getVal(): buildByte(path),
		}
	}
	left := generateEncodings(node.(InternalNode).left, append(path, false))
	right := generateEncodings(node.(InternalNode).right, append(path, true))
	maps.Copy(left, right)
	return left
}

func buildByte(path []bool) HuffmanEncoding {
	var wip uint32
	for _, val := range path {
		wip = wip << 1
		if val {
			wip |= 1
		}
	}

	return HuffmanEncoding{
		Value: wip,
		Bits:  uint8(len(path)),
	}
}

func packBits(encodedText []HuffmanEncoding) []byte {
	var packedBytes []byte
	var buffer uint64
	var bufferLen uint8

	for _, encoding := range encodedText {
		buffer = (buffer << encoding.Bits) | uint64(encoding.Value)
		bufferLen += encoding.Bits

		for bufferLen >= 8 {
			bufferLen -= 8
			packedBytes = append(packedBytes, byte(buffer>>bufferLen))
		}
	}

	if bufferLen > 0 {
		packedBytes = append(packedBytes, byte(buffer<<(8-bufferLen)))
	}

	return packedBytes
}

func nthBit(number byte, n int) byte {
	if (1<<(8-n))&number > 0 {
		return 1
	}
	return 0
}

type HuffmanEncoding struct {
	Value uint32
	Bits  uint8
}
