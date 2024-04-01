package main

import (
	"fmt"
	"github.com/spaolacci/murmur3"
	"math"
)

// Bloom filter structure
type BloomFilter struct {
	numBits int
	numHashFunctions int
	bitArray []bool
}

func main() {

	fmt.Println("Lovely")
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