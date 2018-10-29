package lookup3

import "testing"

var hashWordTests = []struct {
	key       []uint32
	keyLength uint32
	initValue uint32
	hash      uint32
}{
	{[]uint32{1}, 1, 0, 0x72a82a9b},
	{[]uint32{1, 2}, 2, 1, 0x73989811},
	{[]uint32{1, 2, 3}, 3, 0, 0xa46158f5},
	{[]uint32{1, 2, 3, 4}, 4, 1, 0x044ec9ea},
	{[]uint32{1, 2, 3, 4, 5}, 5, 2, 0x39a100d5},
}

var hashWord2Tests = []struct {
	key        []uint32
	keyLength  uint32
	cIn, bIn   uint32
	cOut, bOut uint32
}{
	{[]uint32{1, 2}, 2, 1, 1, 0x301b0127, 0x3ce9fe7e},
	{[]uint32{1, 2, 3}, 3, 0, 0, 0xa46158f5, 0x45915a7e},
	{[]uint32{1, 2, 3, 4}, 4, 1, 0, 0x044ec9ea, 0x729b6663},
	{[]uint32{1, 2, 3, 4, 5}, 5, 0, 1, 0x7489c25b, 0x898e47dd},
}

func TestHashWord(t *testing.T) {
	for _, tt := range hashWordTests {
		h := hash32(tt.initValue, tt.key)
		if h != tt.hash {
			t.Errorf("hash32(%d, %q) => 0x%08x, want 0x%08x\n", tt.initValue, tt.key, h, tt.hash)
		}
	}
}

func TestHashWord2(t *testing.T) {
	for _, tt := range hashWord2Tests {

		c, b := hash64(tt.cIn, tt.bIn, tt.key)
		if c != tt.cOut || b != tt.bOut {
			t.Errorf("hash64(%d, %d, %q) => (0x%08x, 0x%08x), want (0x%08x, 0x%08x)\n",
				tt.cIn, tt.bIn, tt.key, c, b, tt.cOut, tt.bOut)
		}
	}
}
