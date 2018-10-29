package murmur3

func murmur3_64(seed uint32, b []byte) uint64 {
	h1, _ := murmur3_128(seed, b)
	return h1
}
