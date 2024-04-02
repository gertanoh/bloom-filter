# Bloom Filter simple Implementation

## What is it ?
A Bloom filter is a probalistic data structure, space efficient, that is used to test whether an element is present in a set or not.

Why do we need it ?

Most of us are familiar with a classic algo problem where one uses a set to check the presence of an element.
However, this solution is not memory efficient as the space complexity is O(n), thus scale linear with the size of set. That's where Bloom filters come into play.
Because of it probalistic nature, it is not error free, it comes with false positives, however, there are no false negatives.
When the filter returns no, it is a firm no. For yes, it is possibly yes.
For large dataset, where storing all the elements is not possible in memory, bloom is a huge gain.
It is used a lot in storage. LSM-based storage uses it to check if a key is inside a memtable before checking inside the memtable.

## How does it work ?
A bloom filter is comprised of a array of bit and some hash functions.
Each item is hashed to determine which bit of the array to set for this item.
Multiple hash functions are used to reduce collisions and reduce the error rate.
The two parameters m (number of bits in the array) and n (number of hash functions) are central to provide good accuracy with low error rate.
Check out these two pages to see two pages on how to pick them :

 * [Wikipedia](https://en.wikipedia.org/wiki/Tar_(computing)#File_format)
 * [Bloom Filter Calculator](https://hur.st/bloomfilter/)


## Example application
Spell checker coding challenges from [John Crickett](https://codingchallenges.fyi/challenges/challenge-bloom)

Create a spellchecker application from a list of words on your language of choice.

If you are on Unix, you can use cat /usr/share/dict/words >> dict.txt to create your list of words.

### Build

* Build
make build
* Create the filter
./ccspellcheck -build dict.txt

* Run
```
./ccspellcheck hello1 adventure lov love "infelicity's4"
hello1 is misspelled
lov is misspelled
infelicity's4 is misspelled
```