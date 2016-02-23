package netmatch

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

// Benchmark data
var (
	networks = []string{
		"62.114.1.255/24", "62.114.2.255/24", "62.114.3.255/24", "62.114.4.255/24", "62.114.5.255/24",
		"62.114.6.255/24", "62.114.7.255/24", "62.114.8.255/24", "62.114.9.255/24", "62.114.10.255/24",
		"62.114.11.255/24", "62.114.12.255/24", "62.114.13.255/24", "62.114.14.255/24", "62.114.15.255/24",
		"62.114.16.255/24", "62.114.17.255/24", "62.114.18.255/24", "62.114.19.255/24", "62.114.20.255/24",
		"62.114.21.255/24", "62.114.22.255/24", "62.114.23.255/24", "62.114.24.255/24", "62.114.25.255/24",
		"62.114.26.255/24", "62.114.27.255/24", "62.114.28.255/24", "62.114.29.255/24", "62.114.30.255/24",
		"62.114.31.255/24", "62.114.32.255/24", "62.114.33.255/24", "62.114.34.255/24", "62.114.35.255/24",
		"62.114.36.255/24", "62.114.37.255/24", "62.114.38.255/24", "62.114.39.255/24", "62.114.40.255/24",
		"62.114.41.255/24", "62.114.42.255/24", "62.114.43.255/24", "62.114.44.255/24", "62.114.45.255/24",
		"62.114.46.255/24", "62.114.47.255/24", "62.114.48.255/24", "62.114.49.255/24", "62.114.50.255/24",
		"62.114.51.255/24", "62.114.52.255/24", "62.114.53.255/24", "62.114.54.255/24", "62.114.55.255/24",
		"62.114.56.255/24", "62.114.57.255/24", "62.114.58.255/24", "62.114.59.255/24", "62.114.60.255/24",
		"62.114.61.255/24", "62.114.62.255/24", "62.114.63.255/24", "62.114.64.255/24", "62.114.65.255/24",
		"62.114.66.255/24", "62.114.67.255/24", "62.114.68.255/24", "62.114.69.255/24", "62.114.70.255/24",
		"62.114.71.255/24", "62.114.72.255/24", "62.114.73.255/24", "62.114.74.255/24", "62.114.75.255/24",
		"62.114.76.255/24", "62.114.77.255/24", "62.114.78.255/24", "62.114.79.255/24", "62.114.80.255/24",
		"62.114.81.255/24", "62.114.82.255/24", "62.114.83.255/24", "62.114.84.255/24", "62.114.85.255/24",
		"62.114.86.255/24", "62.114.87.255/24", "62.114.88.255/24", "62.114.89.255/24", "62.114.90.255/24",
		"62.114.91.255/24", "62.114.92.255/24", "62.114.93.255/24", "62.114.94.255/24", "62.114.95.255/24",
		"62.114.96.255/24", "62.114.97.255/24", "62.114.98.255/24", "62.114.99.255/24", "62.114.100.255/24",
	}

	inNet5   = net.ParseIP("62.114.5.24")
	inNet10  = net.ParseIP("62.114.10.230")
	inNet100 = net.ParseIP("62.114.100.32")
)

func TestKey(t *testing.T) {
	var table = []struct {
		input    net.IP
		expected [16]byte
	}{
		{nil, [16]byte{}},
		{net.ParseIP("0c22:384e:0:0c22:384e::68"), [16]byte{12, 34, 56, 78, 0, 0, 12, 34, 56, 78, 0, 0, 0, 0, 0, 104}},
		{net.ParseIP("12.13.14.15"), [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 12, 13, 14, 15}},       // IPv4 in IPv6 prefix
		{net.ParseIP("12.13.14.15").To4(), [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 12, 13, 14, 15}}, // is equal to the one above, should produce equal output
	}

	for _, tt := range table {
		got := Key(tt.input)
		assert.Equal(t, got, tt.expected)
	}
}

func TestParseNetwork(t *testing.T) {
	var table = []struct {
		input          string
		expectedKey    [16]byte
		expectedLength int
		error          bool
	}{
		{"12.13.14.15/24", [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 12, 13, 14, 0}, 24 + 12*8, false},
		{"0c22:384e:0:0c22:384e::68/110", [16]byte{12, 34, 56, 78, 0, 0, 12, 34, 56, 78, 0, 0, 0, 0, 0, 0}, 110, false},
		{"not_a_subnet", [16]byte{}, 0, true},
	}

	for _, tt := range table {
		key, length, err := ParseNetwork(tt.input)
		if tt.error {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedKey, key)
			assert.Equal(t, tt.expectedLength, length)
		}
	}
}

func TestTrie(t *testing.T) {
	trie := New()
	network := "192.168.1.123/24"
	matchingKey := Key(net.ParseIP("192.168.1.55"))
	nonMatchingKey := Key(net.ParseIP("192.168.2.1"))

	key, length, err := ParseNetwork(network)
	assert.Nil(t, err)

	err = trie.Add(key, length)
	assert.Nil(t, err)

	matches, err := trie.Match(matchingKey)
	assert.Nil(t, err)
	assert.True(t, matches)

	matches, err = trie.Match(nonMatchingKey)
	assert.Nil(t, err)
	assert.False(t, matches)

	err = trie.Remove(key, length)
	assert.Nil(t, err)

	matches, err = trie.Match(matchingKey)
	assert.Nil(t, err)
	assert.False(t, matches)
}

func TestOverlapping(t *testing.T) {
	var (
		trie         = New()
		inBothKey    = Key(net.ParseIP("192.168.1.230"))
		inBigOnlyKey = Key(net.ParseIP("192.168.1.20"))
		smallNet     = "192.168.1.255/25"
		bigNet       = "192.168.1.255/24"
	)
	bigKey, bigLength, err := ParseNetwork(bigNet)
	assert.Nil(t, err)

	smallKey, smallLength, err := ParseNetwork(smallNet)
	assert.Nil(t, err)

	matches, err := trie.Match(inBothKey)
	assert.Nil(t, err)
	assert.False(t, matches)

	matches, err = trie.Match(inBigOnlyKey)
	assert.Nil(t, err)
	assert.False(t, matches)

	err = trie.Add(smallKey, smallLength)
	assert.Nil(t, err)

	matches, err = trie.Match(inBothKey)
	assert.Nil(t, err)
	assert.True(t, matches)

	matches, err = trie.Match(inBigOnlyKey)
	assert.Nil(t, err)
	assert.False(t, matches)

	err = trie.Add(bigKey, bigLength)
	assert.Nil(t, err)

	matches, err = trie.Match(inBothKey)
	assert.Nil(t, err)
	assert.True(t, matches)

	matches, err = trie.Match(inBigOnlyKey)
	assert.Nil(t, err)
	assert.True(t, matches)

	err = trie.Remove(bigKey, bigLength)
	assert.Nil(t, err)
	err = trie.Remove(bigKey, bigLength)
	assert.NotNil(t, err)

	matches, err = trie.Match(inBothKey)
	assert.Nil(t, err)
	assert.True(t, matches)

	matches, err = trie.Match(inBigOnlyKey)
	assert.Nil(t, err)
	assert.False(t, matches)

	err = trie.Remove(smallKey, smallLength)
	assert.Nil(t, err)
	err = trie.Remove(smallKey, smallLength)
	assert.NotNil(t, err)

	matches, err = trie.Match(inBothKey)
	assert.Nil(t, err)
	assert.False(t, matches)

	matches, err = trie.Match(inBigOnlyKey)
	assert.Nil(t, err)
	assert.False(t, matches)

	err = trie.Add(smallKey, smallLength)
	assert.Nil(t, err)

	err = trie.Remove(bigKey, bigLength)
	assert.NotNil(t, err)

	err = trie.Remove(smallKey, smallLength)
	assert.Nil(t, err)
}

func BenchmarkTrieAdd(b *testing.B) {
	trie := New()
	key, length, err := ParseNetwork(networks[0])
	assert.Nil(b, err)

	for i := 0; i < b.N; i++ {
		trie.Add(key, length)
	}
}

func BenchmarkTrieAddRemove(b *testing.B) {
	trie := New()
	key, length, err := ParseNetwork(networks[0])
	assert.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Add(key, length)
		trie.Remove(key, length)
	}
}

func BenchmarkTrie5Match(b *testing.B) {
	trie := New()
	inNet5Key := Key(inNet5)

	for i := 0; i < 5; i++ {
		key, length, err := ParseNetwork(networks[i])
		assert.Nil(b, err)
		trie.Add(key, length)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Match(inNet5Key)
	}
}

func BenchmarkTrie10Match(b *testing.B) {
	trie := New()
	inNet10Key := Key(inNet10)

	for i := 0; i < 10; i++ {
		key, length, err := ParseNetwork(networks[i])
		assert.Nil(b, err)
		trie.Add(key, length)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Match(inNet10Key)
	}
}

func BenchmarkTrie100Match(b *testing.B) {
	trie := New()
	inNet100Key := Key(inNet100)

	for i := 0; i < 100; i++ {
		key, length, err := ParseNetwork(networks[i])
		assert.Nil(b, err)
		trie.Add(key, length)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Match(inNet100Key)
	}
}

func BenchmarkList5Match(b *testing.B) {
	nets := make([]*net.IPNet, 5)

	for i := 0; i < 5; i++ {
		_, network, err := net.ParseCIDR(networks[i])
		assert.Nil(b, err)
		nets[i] = network
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, network := range nets {
			if network.Contains(inNet5) {
				break
			}
		}
	}
}

func BenchmarkList10Match(b *testing.B) {
	nets := make([]*net.IPNet, 10)

	for i := 0; i < 10; i++ {
		_, network, err := net.ParseCIDR(networks[i])
		assert.Nil(b, err)
		nets[i] = network
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, network := range nets {
			if network.Contains(inNet10) {
				break
			}
		}
	}
}

func BenchmarkList100Match(b *testing.B) {
	nets := make([]*net.IPNet, 100)

	for i := 0; i < 100; i++ {
		_, network, err := net.ParseCIDR(networks[i])
		assert.Nil(b, err)
		nets[i] = network
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, network := range nets {
			if network.Contains(inNet100) {
				break
			}
		}
	}
}
