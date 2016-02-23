# netmatch
[![Build Status](https://travis-ci.org/mrd0ll4r/netmatch.svg?branch=master)](https://travis-ci.org/mrd0ll4r/netmatch)
[![GoDoc](https://godoc.org/github.com/mrd0ll4r/netmatch?status.svg)](https://godoc.org/github.com/mrd0ll4r/netmatch)
[![Go Report Card](https://goreportcard.com/badge/github.com/mrd0ll4r/netmatch)](https://goreportcard.com/report/github.com/mrd0ll4r/netmatch)

Match an IP address against a lot of prefixes in constant time

## What does it do?
It implements a binary Trie-like structure to be used to match IP addresses against prefixes.

Matching takes O(k), where k is the height of the Trie.
The height of the Trie is determined by the longest prefix to be matched.
Because the prefixes can not be longer than some 126 bits, the time requirement is pretty much constant, for more info on performance, check the benchmarks below.

## Usage
Create a new Trie, add prefixes, match against it:

```go
t := netmatch.New()
key, len, err := netmatch.ParseNetwork("192.168.122.255/24")
if err != nil {
    log.Fatal(err)
}

t.Add(key, len)

ipToMatch = net.ParseIP("192.168.122.32")
matched, err := t.Match(netmatch.Key(ipToMatch))
if err != nil {
    log.Fatal(err)
}

//matched is true
```

## Benchmark

```
BenchmarkTrieAdd-4               2000000               649 ns/op
BenchmarkTrieAddRemove-4          200000             11671 ns/op
BenchmarkTrie5Match-4            2000000               628 ns/op
BenchmarkTrie10Match-4           2000000               623 ns/op
BenchmarkTrie100Match-4          2000000               620 ns/op
BenchmarkList5Match-4           10000000               242 ns/op
BenchmarkList10Match-4           3000000               464 ns/op
BenchmarkList100Match-4           300000              4547 ns/op
```

Unsurprisingly, just putting `net.IPNet`s into a slice and iterating over them will take O(n), where n is the number of networks to match against.
The simple Trie implementation should be faster than the slice solution if you have to deal with more than 20 or so networks.
Removing is painfully slow, but the point of this implementation is to have fast matching.