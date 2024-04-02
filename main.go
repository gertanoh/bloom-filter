package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/spaolacci/murmur3"
)

// Bloom filter structure
type BloomFilter struct {
	numBits          uint32
	numHashFunctions uint16
	bitArray         []bool
}

var (
	buildTime string
	version   string
)

const bloomFilterFilePath = "words.bf"

func main() {

	buildFlag := ""
	flag.StringVar(&buildFlag, "build", "", "file.txt")
	flag.Parse()

	bf := NewBloomFilter(100000, 0.00001)
	if buildFlag != "" {
		err := FillDataInFilter(bf, buildFlag)
		if err != nil {
			fmt.Println("Error filling bloom filter:", err)
		}
		return
	} else {
		bf, err := LoadBloomFilter(bloomFilterFilePath)
		if err != nil {
			fmt.Println("Error in loading Bloom Filter in memory")
			return
		}
		for _, word := range flag.Args() {
			// Perform spell check using Bloom filter
			if !bf.Contains(word) {
				fmt.Println(word + " is misspelled")
			}
		}
	}

}

func FillDataInFilter(filter *BloomFilter, filepath string) error {
	// read each line of the file and add it to the filter
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		for _, word := range words {
			filter.Add(word)
		}
	}
	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	saveBloomFilter(filter, bloomFilterFilePath)
	return nil
}

// saveBloomFilter saves the Bloom filter to disk with a header
func saveBloomFilter(bf *BloomFilter, filePath string) error {
	// Create or truncate the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header to file
	header := make([]byte, 12)                 // 4 bytes for identifier, 2 bytes for version number, 2 bytes for number of hash functions, 4 bytes for number of bits
	copy(header[:4], "CCBF")                   // Identifier
	binary.BigEndian.PutUint16(header[4:6], 1) // Version number (example: 1)
	binary.BigEndian.PutUint16(header[6:8], uint16(bf.numHashFunctions))
	binary.BigEndian.PutUint32(header[8:12], bf.numBits)
	if _, err := file.Write(header); err != nil {
		return err
	}

	numBytes := (bf.numBits + 7) / 8 // Round up to the nearest whole number of bytes

	// Convert the bit array to a byte slice
	bitArrayBytes := make([]byte, numBytes)
	var i uint32
	for i = 0; i < bf.numBits; i++ {
		if bf.bitArray[i] {
			bitArrayBytes[i/8] |= 1 << (i % 8)
		}
	}

	// Write bit array to file
	if _, err := file.Write(bitArrayBytes); err != nil {
		return err
	}

	return nil
}

// Load bloom filter in memory
func LoadBloomFilter(filePath string) (*BloomFilter, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the header
	header := make([]byte, 12)
	if _, err := file.Read(header); err != nil {
		return nil, err
	}

	// Check the header
	if string(header[:4]) != "CCBF" {
		return nil, fmt.Errorf("invalid header")
	}
	version := binary.BigEndian.Uint16(header[4:6])
	numHashFunctions := binary.BigEndian.Uint16(header[6:8])
	numBits := binary.BigEndian.Uint32(header[8:12])
	if version != 1 {
		return nil, fmt.Errorf("unsupported version")
	}

	if numHashFunctions == 0 {
		return nil, fmt.Errorf("invalid numHashFunctions")
	}

	if numBits == 0 {
		return nil, fmt.Errorf("invalid numBits")
	}

	// Read the bit array
	bitArrayBytes := make([]byte, (numBits+7)/8) // Round up to the nearest whole number of bytes
	if _, err := file.Read(bitArrayBytes); err != nil {
		return nil, err
	}

	// Convert the bit array bytes to a boolean array
	bitArray := make([]bool, numBits)
	var i uint32
	for i = 0; i < numBits; i++ {
		bitArray[i] = (bitArrayBytes[i/8] & (1 << (i % 8))) != 0
	}

	bf := &BloomFilter{
		numBits:          numBits,
		numHashFunctions: numHashFunctions,
		bitArray:         bitArray,
	}
	return bf, nil
}

// calculateNumBits calculates the number of bits needed in the Bloom filter
func calculateNumBits(numItems uint32, falsePositiveProb float64) uint32 {
	return uint32(-(float64(numItems) * math.Log(falsePositiveProb)) / (math.Log(2) * math.Log(2)))
}

// calculateNumHashFunctions calculates the number of hash functions needed in the Bloom filter
func calculateNumHashFunctions(numBits, numItems uint32) uint16 {
	return uint16(float64(numBits) / float64(numItems) * math.Log(2))
}

// BloomFilter constructor
func NewBloomFilter(numItems uint32, falsePositiveErrorRate float64) *BloomFilter {
	bf := &BloomFilter{}
	bf.numBits = calculateNumBits(numItems, falsePositiveErrorRate)
	bf.numHashFunctions = calculateNumHashFunctions(bf.numBits, numItems)
	bf.bitArray = make([]bool, bf.numBits)

	return bf
}

func (bf *BloomFilter) Add(item string) {
	var i uint16
	for i = 0; i < bf.numHashFunctions; i++ {
		hashIndex := murmurHash(item, i) % bf.numBits
		bf.bitArray[hashIndex] = true
	}
}

func (bf *BloomFilter) Contains(item string) bool {
	var i uint16
	for i = 0; i < bf.numHashFunctions; i++ {
		hashIndex := murmurHash(item, i) % bf.numBits
		if !bf.bitArray[hashIndex] {
			return false
		}
	}
	return true
}

// murmurHash generates a hash value for the given item and seed using MurmurHash3
func murmurHash(item string, seed uint16) uint32 {
	hash := murmur3.New32WithSeed(uint32(seed))
	hash.Write([]byte(item))
	return hash.Sum32()
}
