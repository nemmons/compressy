# Compressy

A solution for a Huffman text compression [coding challenge](https://codingchallenges.fyi/challenges/challenge-huffman) 
in golang.

## Overview

A text file can be compressed into a smaller binary format by removing redundant bits. We place each unique character 
from the input text into a binary tree, sorted by frequency (with more frequent characters placed higher), and then encode
each character based on the path taken to traverse the tree to reach it (left = 0, right = 1). The frequency distribution is
stored alongside the encoded text to allow the tree to be reconstructed, and then decoding this format is essentially
a matter of reading the stream of bits and traversing the tree, fetching characters from the leaf nodes encountered.

### Example
Input Text: `aaaaaaeeebbcc` 

A '■' character is appended to the end, as a 'psuedo-eof' marker, to make it easier to cut off leftover bits when decoding later.

Frequencies:
- a: 6
- e: 3
- b: 2
- c: 2
- ■: 1

Tree:

*Leaf nodes show the character and its frequency, internal nodes show the sum of child frequencies)*
```text
       14
    /      \
6(a)        8
          /   \
      3(e)      5
              /   \
            2(c)    3
                   /  \
                1(■)  2(b)
```
encodings:

- a: 0 
- e: 10 
- c: 110 
- b: 1111
- ■: 1110 

Encoded Text: 001010010000111111111101101110

Split Into Bytes (with trailing padding): 00101001 00001111 11111101 10111000

So, we've taken 13 characters and encoded them into 4 bytes. 


## Usage

```shell
go build compressy
./compressy compress
./compressy decompress
diff -as input.txt decompressed.txt
```

## Todo

- [ ] Improve Overview (Blog post?)
- [ ] Clean up code
  - [ ] Improve clarity  
  - [ ] Figure out best practices for code splitting across files in golang
- [ ] Improve test coverage
  - [ ] better unit tests
  - [ ] e2e
- [ ] Improve performance (decoding speed)
  - [ ] Target is [Les Misérables](https://www.gutenberg.org/files/135/135-0.txt) 