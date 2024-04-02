package main

import (
	"fmt"
	"github.com/spaolacci/murmur3"
	"math"
	"flag"
	"strings"
	"bufio"
	"os"
)

// Bloom filter structure
type BloomFilter struct {
	numBits int
	numHashFunctions int
	bitArray []bool
}

func main() {

	fmt.Println("Lovely")

	buildFlag := ""
	flag.StringVar(&buildFlag, "build", "", "file.txt")
	flag.Parse()

	bf := NewBloomFilter(102401, 0.01)
	if  buildFlag != "" {
		err := FillDataInFilter(bf, buildFlag)
		if err != nil {
			fmt.Println("Error filling bloom filter:", err)
		}
		return
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
	header := make([]byte, 12) // 4 bytes for identifier, 2 bytes for version number, 2 bytes for number of hash functions, 4 bytes for number of bits
	copy(header[:4], "CCBF")    // Identifier
	binary.BigEndian.PutUint16(header[4:6], 1) // Version number (example: 1)
	binary.BigEndian.PutUint16(header[6:8], uint16(bf.NumHashFunctions))
	binary.BigEndian.PutUint32(header[8:12], uint32(bf.NumBits))
	if _, err := file.Write(header); err != nil {
		return err
	}

	// Write bit array to file
	for _, bit := range bf.BitArray {
		var b byte
		if bit {
			b = 1
		}
		if _, err := file.Write([]byte{b}); err != nil {
			return err
		}
	}

	return nil
}

// calculateNumBits calculates the number of bits needed in the Bloom filter
func calculateNumBits(numItems int, falsePositiveProb float64) int {
    return int(-(float64(numItems) * math.Log(falsePositiveProb)) / (math.Log(2) * math.Log(2)))
}

// calculateNumHashFunctions calculates the number of hash functions needed in the Bloom filter
func calculateNumHashFunctions(numBits, numItems int) int {
    return int(float64(numBits) / float64(numItems) * math.Log(2))
}

// BloomFilter constructor
func NewBloomFilter(numItems int, falsePositiveErrorRate float64) *BloomFilter{
	bf := &BloomFilter{}
	bf.numBits = calculateNumBits(numItems, falsePositiveErrorRate)
	bf.numHashFunctions = calculateNumHashFunctions(bf.numBits, numItems)
	bf.bitArray = make([]bool, bf.numBits)

	return bf
}

func (bf *BloomFilter) Add(item string) {
	for i := 0; i < bf.numHashFunctions; i++ {
		hashIndex := murmurHash(item, i) % bf.numBits
		bf.bitArray[hashIndex] = true
	}
}

func (bf *BloomFilter) Contains(item string) bool {
	for i := 0; i < bf.numHashFunctions; i++ {
		hashIndex := murmurHash(item, i) % bf.numBits
		if bf.bitArray[hashIndex] == false {
			return false
		}
	}
	return true
}

// murmurHash generates a hash value for the given item and seed using MurmurHash3
func murmurHash(item string, seed int) int {
    hash := murmur3.New32WithSeed(uint32(seed))
    hash.Write([]byte(item))
    return int(hash.Sum32())
}