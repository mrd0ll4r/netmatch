# netmatch
[![Build Status](https://travis-ci.org/mrd0ll4r/netmatch.svg?branch=master)](https://travis-ci.org/mrd0ll4r/netmatch)

Match an IP address against a lot of prefixes in constant time

## What does it do?
It stupidly implements a binary Trie-like structure to be used to match IP addresses against prefixes.
Matching takes O(k), where k is the height of the Trie.
The height of the Trie is determined by the longest prefix to be matched.
Because the prefixes can not be longer than some 126 bits, the time requirement is pretty much constant, for more info on performance, check the benchmarks below.

## Usage
Create a new Trie, add prefixes, match against it:

```go
t := netmatch.New()
key,len,err := netmatch.Subnet("192.168.122.255/24")
if err != nil {
    panic(err)
}

t.Add(key,len)

ipToMatch = net.ParseIP("192.168.122.32")
matched, err := t.Match(netmatch.Key(ipToMatch))
if err != nil {
    panic(err)
}

//matched is true
```

## Benchmark

```
BenchmarkTrieAdd-8               2000000               737 ns/op
BenchmarkTrieAddRemove-8          200000             10553 ns/op
BenchmarkTrie2Match-8            2000000               820 ns/op
BenchmarkTrie10Match-8           2000000               813 ns/op
BenchmarkTrie100Match-8          2000000               683 ns/op
BenchmarkList5Match-8           10000000               225 ns/op
BenchmarkList10Match-8           3000000               434 ns/op
BenchmarkList100Match-8           300000              4428 ns/op
```

Unsurprisingly, just putting `net.IPNet`s into a slice and iterating over them will take O(n), where n is the number of subnets to match against.
The stupid Trie implementation should be faster than the slice solution if you have to deal with more than 20 subnets or so.
Interestingly, a Trie with 100 disjoint subnets performs better than a Trie with less entries.
Removing is painfully slow, but the point of this implementation is to have fast matching.