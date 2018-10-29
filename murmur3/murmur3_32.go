package murmur3

import (
	"unsafe"
)

const (
	c1_32 = 0xcc9e2d51
	c2_32 = 0x1b873593
)

func murmur3_32(seed uint32, b []byte) uint32 {
	h1 := seed
	nblocks := len(b) / 4

	var k1 uint32
	// body
	for i := 0; i < nblocks; i++ {
		k1 = *(*uint32)(unsafe.Pointer(&b[i*4]))
		k1 *= c1_32
		k1 = (k1 << 15) | (k1 >> 17) // rotl32(k1, 15)
		k1 *= c2_32

		h1 ^= k1
		h1 = (h1 << 13) | (h1 >> 19) // rotl32(h1, 13)
		h1 = h1*4 + h1 + 0xe6546b64
	}

	// tail
	tail := b[nblocks*4:]
	k1 = 0
	switch len(b) & 3 {
	case 3:
		k1 ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(tail[0])
		k1 *= c1_32
		k1 = (k1 << 15) | (k1 >> 17) // rotl32(k1, 15)
		k1 *= c2_32
		h1 ^= k1
	}

	h1 ^= uint32(len(b))

	return fmix32(h1)
}

func fmix32(h1 uint32) uint32 {
	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16
	return h1
}
