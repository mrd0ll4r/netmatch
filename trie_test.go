package netmatch

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestTrie(t *testing.T) {
	trie := New()

	_, ipnet, err := net.ParseCIDR("192.168.1.123/24")
	assert.Nil(t, err)

	size, _ := ipnet.Mask.Size()
	err = trie.Add(Key(ipnet.IP), size+12*8)
	assert.Nil(t, err)

	ip := net.ParseIP("192.168.1.55")
	assert.Nil(t, err)

	matches, err := trie.Match(Key(ip))
	assert.Nil(t, err)
	assert.True(t, matches)

	matches, err = trie.Match(Key(net.ParseIP("192.168.2.1")))
	assert.Nil(t, err)
	assert.False(t, matches)

	err = trie.Remove(Key(ipnet.IP), size+12*8)
	assert.Nil(t, err)

	matches, err = trie.Match(Key(ip))
	assert.Nil(t, err)
	assert.False(t, matches)
}

func TestOverlapping(t *testing.T) {
	var (
		trie      = New()
		inBoth    = net.ParseIP("192.168.1.230")
		inBigOnly = net.ParseIP("192.168.1.20")
	)
	_, bigNet, err := net.ParseCIDR("192.168.1.255/24")
	assert.Nil(t, err)
	bigSize, _ := bigNet.Mask.Size()
	bigSize = bigSize + 12*8

	_, smallNet, err := net.ParseCIDR("192.168.1.255/25")
	assert.Nil(t, err)
	smallSize, _ := smallNet.Mask.Size()
	smallSize = smallSize + 12*8

	matches, err := trie.Match(Key(inBoth))
	assert.Nil(t, err)
	assert.False(t, matches)

	matches, err = trie.Match(Key(inBigOnly))
	assert.Nil(t, err)
	assert.False(t, matches)

	err = trie.Add(Key(smallNet.IP), smallSize)
	assert.Nil(t, err)

	matches, err = trie.Match(Key(inBoth))
	assert.Nil(t, err)
	assert.True(t, matches)

	matches, err = trie.Match(Key(inBigOnly))
	assert.Nil(t, err)
	assert.False(t, matches)

	err = trie.Add(Key(bigNet.IP), bigSize)
	assert.Nil(t, err)

	matches, err = trie.Match(Key(inBoth))
	assert.Nil(t, err)
	assert.True(t, matches)

	matches, err = trie.Match(Key(inBigOnly))
	assert.Nil(t, err)
	assert.True(t, matches)

	err = trie.Remove(Key(bigNet.IP), bigSize)
	assert.Nil(t, err)
	err = trie.Remove(Key(bigNet.IP), bigSize)
	assert.NotNil(t, err)

	matches, err = trie.Match(Key(inBoth))
	assert.Nil(t, err)
	assert.True(t, matches)

	matches, err = trie.Match(Key(inBigOnly))
	assert.Nil(t, err)
	assert.False(t, matches)

	err = trie.Remove(Key(smallNet.IP), smallSize)
	assert.Nil(t, err)
	err = trie.Remove(Key(smallNet.IP), smallSize)
	assert.NotNil(t, err)

	matches, err = trie.Match(Key(inBoth))
	assert.Nil(t, err)
	assert.False(t, matches)

	matches, err = trie.Match(Key(inBigOnly))
	assert.Nil(t, err)
	assert.False(t, matches)

	err = trie.Add(Key(smallNet.IP), smallSize)
	assert.Nil(t, err)

	err = trie.Remove(Key(bigNet.IP), bigSize)
	assert.NotNil(t, err)

	err = trie.Remove(Key(smallNet.IP), smallSize)
	assert.Nil(t, err)
}

func TestKey(t *testing.T) {
	var table = []struct {
		input    net.IP
		expected [16]byte
	}{
		{net.ParseIP("0c22:384e:0:0c22:384e::68"), [16]byte{12, 34, 56, 78, 0, 0, 12, 34, 56, 78, 0, 0, 0, 0, 0, 104}},
		{net.ParseIP("12.13.14.15"), [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 12, 13, 14, 15}},       // IPv4 in IPv6 prefix
		{net.ParseIP("12.13.14.15").To4(), [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 12, 13, 14, 15}}, // is equal to the one above, should produce equal output
	}

	for _, tt := range table {
		got := Key(tt.input)
		assert.Equal(t, got, tt.expected)
	}
}

func TestSubnet(t *testing.T) {
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

var (
	_, net1, _  = net.ParseCIDR("192.168.1.0/24")
	_, net2, _  = net.ParseCIDR("192.168.20.255/24")
	_, net3, _  = net.ParseCIDR("10.150.255.255/16")
	_, net4, _  = net.ParseCIDR("10.160.255.255/16")
	_, net5, _  = net.ParseCIDR("172.16.32.255/24")
	_, net6, _  = net.ParseCIDR("172.16.34.255/24")
	_, net7, _  = net.ParseCIDR("11.255.255.255/8")
	_, net8, _  = net.ParseCIDR("12.255.255.255/8")
	_, net9, _  = net.ParseCIDR("243.243.243.255/28")
	_, net10, _ = net.ParseCIDR("142.214.32.255/26")

	_, net11, _ = net.ParseCIDR("62.114.40.255/24")
	_, net12, _ = net.ParseCIDR("62.114.41.255/24")
	_, net13, _ = net.ParseCIDR("62.114.42.255/24")
	_, net14, _ = net.ParseCIDR("62.114.43.255/24")
	_, net15, _ = net.ParseCIDR("62.114.44.255/24")
	_, net16, _ = net.ParseCIDR("62.114.45.255/24")
	_, net17, _ = net.ParseCIDR("62.114.46.255/24")
	_, net18, _ = net.ParseCIDR("62.114.47.255/24")
	_, net19, _ = net.ParseCIDR("62.114.48.255/24")
	_, net20, _ = net.ParseCIDR("62.114.49.255/24")

	_, net21, _ = net.ParseCIDR("62.114.50.255/24")
	_, net22, _ = net.ParseCIDR("62.114.51.255/24")
	_, net23, _ = net.ParseCIDR("62.114.52.255/24")
	_, net24, _ = net.ParseCIDR("62.114.53.255/24")
	_, net25, _ = net.ParseCIDR("62.114.54.255/24")
	_, net26, _ = net.ParseCIDR("62.114.55.255/24")
	_, net27, _ = net.ParseCIDR("62.114.56.255/24")
	_, net28, _ = net.ParseCIDR("62.114.57.255/24")
	_, net29, _ = net.ParseCIDR("62.114.58.255/24")
	_, net30, _ = net.ParseCIDR("62.114.59.255/24")

	_, net31, _ = net.ParseCIDR("62.114.60.255/24")
	_, net32, _ = net.ParseCIDR("62.114.61.255/24")
	_, net33, _ = net.ParseCIDR("62.114.62.255/24")
	_, net34, _ = net.ParseCIDR("62.114.63.255/24")
	_, net35, _ = net.ParseCIDR("62.114.64.255/24")
	_, net36, _ = net.ParseCIDR("62.114.65.255/24")
	_, net37, _ = net.ParseCIDR("62.114.66.255/24")
	_, net38, _ = net.ParseCIDR("62.114.67.255/24")
	_, net39, _ = net.ParseCIDR("62.114.68.255/24")
	_, net40, _ = net.ParseCIDR("62.114.69.255/24")

	_, net41, _ = net.ParseCIDR("62.114.70.255/24")
	_, net42, _ = net.ParseCIDR("62.114.71.255/24")
	_, net43, _ = net.ParseCIDR("62.114.72.255/24")
	_, net44, _ = net.ParseCIDR("62.114.73.255/24")
	_, net45, _ = net.ParseCIDR("62.114.74.255/24")
	_, net46, _ = net.ParseCIDR("62.114.75.255/24")
	_, net47, _ = net.ParseCIDR("62.114.76.255/24")
	_, net48, _ = net.ParseCIDR("62.114.77.255/24")
	_, net49, _ = net.ParseCIDR("62.114.78.255/24")
	_, net50, _ = net.ParseCIDR("62.114.79.255/24")

	_, net51, _ = net.ParseCIDR("62.114.80.255/24")
	_, net52, _ = net.ParseCIDR("62.114.81.255/24")
	_, net53, _ = net.ParseCIDR("62.114.82.255/24")
	_, net54, _ = net.ParseCIDR("62.114.83.255/24")
	_, net55, _ = net.ParseCIDR("62.114.84.255/24")
	_, net56, _ = net.ParseCIDR("62.114.85.255/24")
	_, net57, _ = net.ParseCIDR("62.114.86.255/24")
	_, net58, _ = net.ParseCIDR("62.114.87.255/24")
	_, net59, _ = net.ParseCIDR("62.114.88.255/24")
	_, net60, _ = net.ParseCIDR("62.114.89.255/24")

	_, net61, _ = net.ParseCIDR("62.114.90.255/24")
	_, net62, _ = net.ParseCIDR("62.114.91.255/24")
	_, net63, _ = net.ParseCIDR("62.114.92.255/24")
	_, net64, _ = net.ParseCIDR("62.114.93.255/24")
	_, net65, _ = net.ParseCIDR("62.114.94.255/24")
	_, net66, _ = net.ParseCIDR("62.114.95.255/24")
	_, net67, _ = net.ParseCIDR("62.114.96.255/24")
	_, net68, _ = net.ParseCIDR("62.114.97.255/24")
	_, net69, _ = net.ParseCIDR("62.114.98.255/24")
	_, net70, _ = net.ParseCIDR("62.114.99.255/24")

	_, net71, _ = net.ParseCIDR("62.114.100.255/24")
	_, net72, _ = net.ParseCIDR("62.114.101.255/24")
	_, net73, _ = net.ParseCIDR("62.114.102.255/24")
	_, net74, _ = net.ParseCIDR("62.114.103.255/24")
	_, net75, _ = net.ParseCIDR("62.114.104.255/24")
	_, net76, _ = net.ParseCIDR("62.114.105.255/24")
	_, net77, _ = net.ParseCIDR("62.114.106.255/24")
	_, net78, _ = net.ParseCIDR("62.114.107.255/24")
	_, net79, _ = net.ParseCIDR("62.114.108.255/24")
	_, net80, _ = net.ParseCIDR("62.114.109.255/24")

	_, net81, _ = net.ParseCIDR("62.114.110.255/24")
	_, net82, _ = net.ParseCIDR("62.114.111.255/24")
	_, net83, _ = net.ParseCIDR("62.114.112.255/24")
	_, net84, _ = net.ParseCIDR("62.114.113.255/24")
	_, net85, _ = net.ParseCIDR("62.114.114.255/24")
	_, net86, _ = net.ParseCIDR("62.114.115.255/24")
	_, net87, _ = net.ParseCIDR("62.114.116.255/24")
	_, net88, _ = net.ParseCIDR("62.114.117.255/24")
	_, net89, _ = net.ParseCIDR("62.114.118.255/24")
	_, net90, _ = net.ParseCIDR("62.114.119.255/24")

	_, net91, _  = net.ParseCIDR("62.114.120.255/24")
	_, net92, _  = net.ParseCIDR("62.114.121.255/24")
	_, net93, _  = net.ParseCIDR("62.114.122.255/24")
	_, net94, _  = net.ParseCIDR("62.114.123.255/24")
	_, net95, _  = net.ParseCIDR("62.114.124.255/24")
	_, net96, _  = net.ParseCIDR("62.114.125.255/24")
	_, net97, _  = net.ParseCIDR("62.114.126.255/24")
	_, net98, _  = net.ParseCIDR("62.114.127.255/24")
	_, net99, _  = net.ParseCIDR("62.114.128.255/24")
	_, net100, _ = net.ParseCIDR("62.114.129.255/24")

	inNet5   = net.ParseIP("172.16.32.24")
	inNet10  = net.ParseIP("142.214.32.230")
	inNet100 = net.ParseIP("62.114.129.32")
)

func BenchmarkTrieAdd(b *testing.B) {
	trie := New()
	net1Size, _ := net1.Mask.Size()
	net1Size = net1Size + 12*8
	net1Key := Key(net1.IP)
	for i := 0; i < b.N; i++ {
		trie.Add(net1Key, net1Size)
	}
}

func BenchmarkTrieAddRemove(b *testing.B) {
	trie := New()
	net1Size, _ := net1.Mask.Size()
	net1Size = net1Size + 12*8
	net1Key := Key(net1.IP)
	for i := 0; i < b.N; i++ {
		trie.Add(net1Key, net1Size)
		trie.Remove(net1Key, net1Size)
	}
}

func BenchmarkTrie2Match(b *testing.B) {
	trie := New()
	net1Size, _ := net1.Mask.Size()
	net1Size = net1Size + 12*8
	net10Size, _ := net10.Mask.Size()
	net10Size = net10Size + 12*8

	trie.Add(Key(net1.IP), net1Size)
	trie.Add(Key(net10.IP), net10Size)

	b.ResetTimer()
	inNet10Key := Key(inNet10)
	for i := 0; i < b.N; i++ {
		trie.Match(inNet10Key)
	}
}

func BenchmarkTrie10Match(b *testing.B) {
	trie := New()
	net1Size, _ := net1.Mask.Size()
	net1Size = net1Size + 12*8
	net2Size, _ := net2.Mask.Size()
	net2Size = net2Size + 12*8
	net3Size, _ := net3.Mask.Size()
	net3Size = net3Size + 12*8
	net4Size, _ := net4.Mask.Size()
	net4Size = net4Size + 12*8
	net5Size, _ := net5.Mask.Size()
	net5Size = net5Size + 12*8
	net6Size, _ := net6.Mask.Size()
	net6Size = net6Size + 12*8
	net7Size, _ := net7.Mask.Size()
	net7Size = net7Size + 12*8
	net8Size, _ := net8.Mask.Size()
	net8Size = net8Size + 12*8
	net9Size, _ := net9.Mask.Size()
	net9Size = net9Size + 12*8
	net10Size, _ := net10.Mask.Size()
	net10Size = net10Size + 12*8

	trie.Add(Key(net1.IP), net1Size)
	trie.Add(Key(net2.IP), net2Size)
	trie.Add(Key(net3.IP), net3Size)
	trie.Add(Key(net4.IP), net4Size)
	trie.Add(Key(net5.IP), net5Size)
	trie.Add(Key(net6.IP), net6Size)
	trie.Add(Key(net7.IP), net7Size)
	trie.Add(Key(net8.IP), net8Size)
	trie.Add(Key(net9.IP), net9Size)
	trie.Add(Key(net10.IP), net10Size)

	b.ResetTimer()
	inNet10Key := Key(inNet10)
	for i := 0; i < b.N; i++ {
		trie.Match(inNet10Key)
	}
}

func BenchmarkTrie100Match(b *testing.B) {
	trie := New()
	net1Size, _ := net1.Mask.Size()
	net1Size = net1Size + 12*8
	net2Size, _ := net2.Mask.Size()
	net2Size = net2Size + 12*8
	net3Size, _ := net3.Mask.Size()
	net3Size = net3Size + 12*8
	net4Size, _ := net4.Mask.Size()
	net4Size = net4Size + 12*8
	net5Size, _ := net5.Mask.Size()
	net5Size = net5Size + 12*8
	net6Size, _ := net6.Mask.Size()
	net6Size = net6Size + 12*8
	net7Size, _ := net7.Mask.Size()
	net7Size = net7Size + 12*8
	net8Size, _ := net8.Mask.Size()
	net8Size = net8Size + 12*8
	net9Size, _ := net9.Mask.Size()
	net9Size = net9Size + 12*8
	net10Size, _ := net10.Mask.Size()
	net10Size = net10Size + 12*8

	net11Size, _ := net11.Mask.Size()
	net11Size = net11Size + 12*8
	net12Size, _ := net12.Mask.Size()
	net12Size = net12Size + 12*8
	net13Size, _ := net13.Mask.Size()
	net13Size = net13Size + 12*8
	net14Size, _ := net14.Mask.Size()
	net14Size = net14Size + 12*8
	net15Size, _ := net15.Mask.Size()
	net15Size = net15Size + 12*8
	net16Size, _ := net16.Mask.Size()
	net16Size = net16Size + 12*8
	net17Size, _ := net17.Mask.Size()
	net17Size = net17Size + 12*8
	net18Size, _ := net18.Mask.Size()
	net18Size = net18Size + 12*8
	net19Size, _ := net19.Mask.Size()
	net19Size = net19Size + 12*8
	net20Size, _ := net20.Mask.Size()
	net20Size = net20Size + 12*8

	net21Size, _ := net21.Mask.Size()
	net21Size = net21Size + 12*8
	net22Size, _ := net22.Mask.Size()
	net22Size = net22Size + 12*8
	net23Size, _ := net23.Mask.Size()
	net23Size = net23Size + 12*8
	net24Size, _ := net24.Mask.Size()
	net24Size = net24Size + 12*8
	net25Size, _ := net25.Mask.Size()
	net25Size = net25Size + 12*8
	net26Size, _ := net26.Mask.Size()
	net26Size = net26Size + 12*8
	net27Size, _ := net27.Mask.Size()
	net27Size = net27Size + 12*8
	net28Size, _ := net28.Mask.Size()
	net28Size = net28Size + 12*8
	net29Size, _ := net29.Mask.Size()
	net29Size = net29Size + 12*8
	net30Size, _ := net30.Mask.Size()
	net30Size = net30Size + 12*8

	net31Size, _ := net31.Mask.Size()
	net31Size = net31Size + 12*8
	net32Size, _ := net32.Mask.Size()
	net32Size = net32Size + 12*8
	net33Size, _ := net33.Mask.Size()
	net33Size = net33Size + 12*8
	net34Size, _ := net34.Mask.Size()
	net34Size = net34Size + 12*8
	net35Size, _ := net35.Mask.Size()
	net35Size = net35Size + 12*8
	net36Size, _ := net36.Mask.Size()
	net36Size = net36Size + 12*8
	net37Size, _ := net37.Mask.Size()
	net37Size = net37Size + 12*8
	net38Size, _ := net38.Mask.Size()
	net38Size = net38Size + 12*8
	net39Size, _ := net39.Mask.Size()
	net39Size = net39Size + 12*8
	net40Size, _ := net40.Mask.Size()
	net40Size = net40Size + 12*8

	net41Size, _ := net41.Mask.Size()
	net41Size = net41Size + 12*8
	net42Size, _ := net42.Mask.Size()
	net42Size = net42Size + 12*8
	net43Size, _ := net43.Mask.Size()
	net43Size = net43Size + 12*8
	net44Size, _ := net44.Mask.Size()
	net44Size = net44Size + 12*8
	net45Size, _ := net45.Mask.Size()
	net45Size = net45Size + 12*8
	net46Size, _ := net46.Mask.Size()
	net46Size = net46Size + 12*8
	net47Size, _ := net47.Mask.Size()
	net47Size = net47Size + 12*8
	net48Size, _ := net48.Mask.Size()
	net48Size = net48Size + 12*8
	net49Size, _ := net49.Mask.Size()
	net49Size = net49Size + 12*8
	net50Size, _ := net50.Mask.Size()
	net50Size = net50Size + 12*8

	net51Size, _ := net51.Mask.Size()
	net51Size = net51Size + 12*8
	net52Size, _ := net52.Mask.Size()
	net52Size = net52Size + 12*8
	net53Size, _ := net53.Mask.Size()
	net53Size = net53Size + 12*8
	net54Size, _ := net54.Mask.Size()
	net54Size = net54Size + 12*8
	net55Size, _ := net55.Mask.Size()
	net55Size = net55Size + 12*8
	net56Size, _ := net56.Mask.Size()
	net56Size = net56Size + 12*8
	net57Size, _ := net57.Mask.Size()
	net57Size = net57Size + 12*8
	net58Size, _ := net58.Mask.Size()
	net58Size = net58Size + 12*8
	net59Size, _ := net59.Mask.Size()
	net59Size = net59Size + 12*8
	net60Size, _ := net60.Mask.Size()
	net60Size = net60Size + 12*8

	net61Size, _ := net61.Mask.Size()
	net61Size = net61Size + 12*8
	net62Size, _ := net62.Mask.Size()
	net62Size = net62Size + 12*8
	net63Size, _ := net63.Mask.Size()
	net63Size = net63Size + 12*8
	net64Size, _ := net64.Mask.Size()
	net64Size = net64Size + 12*8
	net65Size, _ := net65.Mask.Size()
	net65Size = net65Size + 12*8
	net66Size, _ := net66.Mask.Size()
	net66Size = net66Size + 12*8
	net67Size, _ := net67.Mask.Size()
	net67Size = net67Size + 12*8
	net68Size, _ := net68.Mask.Size()
	net68Size = net68Size + 12*8
	net69Size, _ := net69.Mask.Size()
	net69Size = net69Size + 12*8
	net70Size, _ := net70.Mask.Size()
	net70Size = net70Size + 12*8

	net71Size, _ := net71.Mask.Size()
	net71Size = net71Size + 12*8
	net72Size, _ := net72.Mask.Size()
	net72Size = net72Size + 12*8
	net73Size, _ := net73.Mask.Size()
	net73Size = net73Size + 12*8
	net74Size, _ := net74.Mask.Size()
	net74Size = net74Size + 12*8
	net75Size, _ := net75.Mask.Size()
	net75Size = net75Size + 12*8
	net76Size, _ := net76.Mask.Size()
	net76Size = net76Size + 12*8
	net77Size, _ := net77.Mask.Size()
	net77Size = net77Size + 12*8
	net78Size, _ := net78.Mask.Size()
	net78Size = net78Size + 12*8
	net79Size, _ := net79.Mask.Size()
	net79Size = net79Size + 12*8
	net80Size, _ := net80.Mask.Size()
	net80Size = net80Size + 12*8

	net81Size, _ := net81.Mask.Size()
	net81Size = net81Size + 12*8
	net82Size, _ := net82.Mask.Size()
	net82Size = net82Size + 12*8
	net83Size, _ := net83.Mask.Size()
	net83Size = net83Size + 12*8
	net84Size, _ := net84.Mask.Size()
	net84Size = net84Size + 12*8
	net85Size, _ := net85.Mask.Size()
	net85Size = net85Size + 12*8
	net86Size, _ := net86.Mask.Size()
	net86Size = net86Size + 12*8
	net87Size, _ := net87.Mask.Size()
	net87Size = net87Size + 12*8
	net88Size, _ := net88.Mask.Size()
	net88Size = net88Size + 12*8
	net89Size, _ := net89.Mask.Size()
	net89Size = net89Size + 12*8
	net90Size, _ := net90.Mask.Size()
	net90Size = net90Size + 12*8

	net91Size, _ := net91.Mask.Size()
	net91Size = net91Size + 12*8
	net92Size, _ := net92.Mask.Size()
	net92Size = net92Size + 12*8
	net93Size, _ := net93.Mask.Size()
	net93Size = net93Size + 12*8
	net94Size, _ := net94.Mask.Size()
	net94Size = net94Size + 12*8
	net95Size, _ := net95.Mask.Size()
	net95Size = net95Size + 12*8
	net96Size, _ := net96.Mask.Size()
	net96Size = net96Size + 12*8
	net97Size, _ := net97.Mask.Size()
	net97Size = net97Size + 12*8
	net98Size, _ := net98.Mask.Size()
	net98Size = net98Size + 12*8
	net99Size, _ := net99.Mask.Size()
	net99Size = net99Size + 12*8
	net100Size, _ := net100.Mask.Size()
	net100Size = net100Size + 12*8

	trie.Add(Key(net1.IP), net1Size)
	trie.Add(Key(net2.IP), net2Size)
	trie.Add(Key(net3.IP), net3Size)
	trie.Add(Key(net4.IP), net4Size)
	trie.Add(Key(net5.IP), net5Size)
	trie.Add(Key(net6.IP), net6Size)
	trie.Add(Key(net7.IP), net7Size)
	trie.Add(Key(net8.IP), net8Size)
	trie.Add(Key(net9.IP), net9Size)
	trie.Add(Key(net10.IP), net10Size)

	trie.Add(Key(net11.IP), net11Size)
	trie.Add(Key(net12.IP), net12Size)
	trie.Add(Key(net13.IP), net13Size)
	trie.Add(Key(net14.IP), net14Size)
	trie.Add(Key(net15.IP), net15Size)
	trie.Add(Key(net16.IP), net16Size)
	trie.Add(Key(net17.IP), net17Size)
	trie.Add(Key(net18.IP), net18Size)
	trie.Add(Key(net19.IP), net19Size)
	trie.Add(Key(net20.IP), net20Size)

	trie.Add(Key(net21.IP), net21Size)
	trie.Add(Key(net22.IP), net22Size)
	trie.Add(Key(net23.IP), net23Size)
	trie.Add(Key(net24.IP), net24Size)
	trie.Add(Key(net25.IP), net25Size)
	trie.Add(Key(net26.IP), net26Size)
	trie.Add(Key(net27.IP), net27Size)
	trie.Add(Key(net28.IP), net28Size)
	trie.Add(Key(net29.IP), net29Size)

	trie.Add(Key(net30.IP), net30Size)
	trie.Add(Key(net31.IP), net31Size)
	trie.Add(Key(net32.IP), net32Size)
	trie.Add(Key(net33.IP), net33Size)
	trie.Add(Key(net34.IP), net34Size)
	trie.Add(Key(net35.IP), net35Size)
	trie.Add(Key(net36.IP), net36Size)
	trie.Add(Key(net37.IP), net37Size)
	trie.Add(Key(net38.IP), net38Size)
	trie.Add(Key(net39.IP), net39Size)

	trie.Add(Key(net40.IP), net40Size)
	trie.Add(Key(net41.IP), net41Size)
	trie.Add(Key(net42.IP), net42Size)
	trie.Add(Key(net43.IP), net43Size)
	trie.Add(Key(net44.IP), net44Size)
	trie.Add(Key(net45.IP), net45Size)
	trie.Add(Key(net46.IP), net46Size)
	trie.Add(Key(net47.IP), net47Size)
	trie.Add(Key(net48.IP), net48Size)
	trie.Add(Key(net49.IP), net49Size)

	trie.Add(Key(net50.IP), net50Size)
	trie.Add(Key(net51.IP), net51Size)
	trie.Add(Key(net52.IP), net52Size)
	trie.Add(Key(net53.IP), net53Size)
	trie.Add(Key(net54.IP), net54Size)
	trie.Add(Key(net55.IP), net55Size)
	trie.Add(Key(net56.IP), net56Size)
	trie.Add(Key(net57.IP), net57Size)
	trie.Add(Key(net58.IP), net58Size)
	trie.Add(Key(net59.IP), net59Size)

	trie.Add(Key(net60.IP), net60Size)
	trie.Add(Key(net61.IP), net61Size)
	trie.Add(Key(net62.IP), net62Size)
	trie.Add(Key(net63.IP), net63Size)
	trie.Add(Key(net64.IP), net64Size)
	trie.Add(Key(net65.IP), net65Size)
	trie.Add(Key(net66.IP), net66Size)
	trie.Add(Key(net67.IP), net67Size)
	trie.Add(Key(net68.IP), net68Size)
	trie.Add(Key(net69.IP), net69Size)

	trie.Add(Key(net70.IP), net70Size)
	trie.Add(Key(net71.IP), net71Size)
	trie.Add(Key(net72.IP), net72Size)
	trie.Add(Key(net73.IP), net73Size)
	trie.Add(Key(net74.IP), net74Size)
	trie.Add(Key(net75.IP), net75Size)
	trie.Add(Key(net76.IP), net76Size)
	trie.Add(Key(net77.IP), net77Size)
	trie.Add(Key(net78.IP), net78Size)
	trie.Add(Key(net79.IP), net79Size)

	trie.Add(Key(net80.IP), net80Size)
	trie.Add(Key(net81.IP), net81Size)
	trie.Add(Key(net82.IP), net82Size)
	trie.Add(Key(net83.IP), net83Size)
	trie.Add(Key(net84.IP), net84Size)
	trie.Add(Key(net85.IP), net85Size)
	trie.Add(Key(net86.IP), net86Size)
	trie.Add(Key(net87.IP), net87Size)
	trie.Add(Key(net88.IP), net88Size)
	trie.Add(Key(net89.IP), net89Size)

	trie.Add(Key(net90.IP), net90Size)
	trie.Add(Key(net91.IP), net91Size)
	trie.Add(Key(net92.IP), net92Size)
	trie.Add(Key(net93.IP), net93Size)
	trie.Add(Key(net94.IP), net94Size)
	trie.Add(Key(net95.IP), net95Size)
	trie.Add(Key(net96.IP), net96Size)
	trie.Add(Key(net97.IP), net97Size)
	trie.Add(Key(net98.IP), net98Size)
	trie.Add(Key(net99.IP), net99Size)

	trie.Add(Key(net100.IP), net100Size)

	b.ResetTimer()
	inNet100Key := Key(inNet100)
	for i := 0; i < b.N; i++ {
		trie.Match(inNet100Key)
	}
}

func BenchmarkList5Match(b *testing.B) {
	nets := []*net.IPNet{net1, net2, net3, net4, net5}

	for i := 0; i < b.N; i++ {
		for _, net := range nets {
			if net.Contains(inNet5) {
				break
			}
		}
	}
}

func BenchmarkList10Match(b *testing.B) {
	nets := []*net.IPNet{net1, net2, net3, net4, net5, net6, net7, net8, net9, net10}

	for i := 0; i < b.N; i++ {
		for _, net := range nets {
			if net.Contains(inNet10) {
				break
			}
		}
	}
}

func BenchmarkList100Match(b *testing.B) {
	nets := []*net.IPNet{
		net1, net2, net3, net4, net5, net6, net7, net8, net9, net10,
		net11, net12, net13, net14, net15, net16, net17, net18, net19, net20,
		net21, net22, net23, net24, net25, net26, net27, net28, net29, net30,
		net31, net32, net33, net34, net35, net36, net37, net38, net39, net40,
		net41, net42, net43, net44, net45, net46, net47, net48, net49, net50,
		net51, net52, net53, net54, net55, net56, net57, net58, net59, net60,
		net61, net62, net63, net64, net65, net66, net67, net68, net69, net70,
		net71, net72, net73, net74, net75, net76, net77, net78, net79, net80,
		net81, net82, net83, net84, net85, net86, net87, net88, net89, net90,
		net91, net92, net93, net94, net95, net96, net97, net98, net99, net100}

	for i := 0; i < b.N; i++ {
		for _, net := range nets {
			if net.Contains(inNet100) {
				break
			}
		}
	}
}
