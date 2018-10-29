package lookup3

import (
	"unsafe"
)

func rot(n uint32, k uint8) uint32 {
	return (n << k) | (n >> (32 - k))
}

func mix(a, b, c uint32) (uint32, uint32, uint32) {
	a -= c
	a ^= rot(c, 4)
	c += b

	b -= a
	b ^= rot(a, 6)
	a += c

	c -= b
	c ^= rot(b, 8)
	b += a

	a -= c
	a ^= rot(c, 16)
	c += b

	b -= a
	b ^= rot(a, 19)
	a += c

	c -= b
	c ^= rot(b, 4)
	b += a

	return a, b, c
}

func final(a, b, c uint32) (uint32, uint32, uint32) {
	c ^= b
	c -= rot(b, 14)

	a ^= c
	a -= rot(c, 11)

	b ^= a
	b -= rot(a, 25)

	c ^= b
	c -= rot(b, 16)

	a ^= c
	a -= rot(c, 4)

	b ^= a
	b -= rot(a, 14)

	c ^= b
	c -= rot(b, 24)

	return a, b, c
}

func hash32(seed uint32, k []uint32) uint32 {
	var a, b, c uint32
	length := len(k)
	nblocks := length / 3

	a = 0xdeadbeef + (uint32(length) << 2) + seed
	b = a
	c = a

	for i := 0; i < nblocks && length > 3; i++ {
		j := i * 3
		a += k[j+0]
		b += k[j+1]
		c += k[j+2]
		a, b, c = mix(a, b, c)
	}

	tail := k
	if length > 3 {
		tail = k[nblocks*3:]
	}
	switch len(tail) {
	case 3:
		c += tail[2]
		fallthrough
	case 2:
		b += tail[1]
		fallthrough
	case 1:
		a += tail[0]
		a, b, c = final(a, b, c)
	case 0: /* case 0: nothing left to add */
		break
	}

	return c
}

func hash64(seed, moreSeed uint32, k []uint32) (uint32, uint32) {
	var a, b, c uint32
	length := len(k)
	nblocks := length / 3

	a = 0xdeadbeef + (uint32(length) << 2) + seed
	b = a
	c = a + moreSeed

	for i := 0; i < nblocks && length > 3; i++ {
		j := i * 3
		a += k[j+0]
		b += k[j+1]
		c += k[j+2]
		a, b, c = mix(a, b, c)
	}

	tail := k
	if length > 3 {
		tail = k[nblocks*3:]
	}
	switch len(tail) {
	case 3:
		c += tail[2]
		fallthrough
	case 2:
		b += tail[1]
		fallthrough
	case 1:
		a += tail[0]
		a, b, c = final(a, b, c)
	case 0: /* case 0: nothing left to add */
		break
	}

	return c, b
}

func hashLittle32(seed uint32, k []byte) uint32 {
	var a, b, c uint32
	length := len(k)

	a = 0xdeadbeef + uint32(length) + seed
	b = a
	c = a

	u := *(*uint32)(unsafe.Pointer(&k[0]))
	if u&3 == 0 {
		nblocks := length / 12
		for i := 0; i < nblocks && length > 12; i++ {
			j := i * 12
			a += *((*uint32)(unsafe.Pointer(&k[j])))
			b += *((*uint32)(unsafe.Pointer(&k[j+4])))
			c += *((*uint32)(unsafe.Pointer(&k[j+8])))
			a, b, c = mix(a, b, c)
		}

		tail := k
		if length > 12 {
			tail = k[nblocks*12:]
		}
		switch len(tail) {
		case 12:
			c += *(*uint32)(unsafe.Pointer(&tail[8]))
			b += *(*uint32)(unsafe.Pointer(&tail[4]))
			a += *(*uint32)(unsafe.Pointer(&tail[0]))
		case 11:
			c += uint32(tail[10]) << 16
			fallthrough
		case 10:
			c += uint32(tail[9]) << 8
			fallthrough
		case 9:
			c += uint32(tail[8])
			fallthrough
		case 8:
			b += *(*uint32)(unsafe.Pointer(&tail[4]))
			a += *(*uint32)(unsafe.Pointer(&tail[0]))
		case 7:
			b += uint32(tail[6]) << 16
			fallthrough
		case 6:
			b += uint32(tail[5]) << 8
			fallthrough
		case 5:
			b += uint32(tail[4])
			fallthrough
		case 4:
			a += *(*uint32)(unsafe.Pointer(&tail[0]))
		case 3:
			a += uint32(tail[2]) << 16
			fallthrough
		case 2:
			a += uint32(tail[1]) << 8
			fallthrough
		case 1:
			a += uint32(tail[0])
		case 0:
			return c
		}
	} else if u&1 == 0 {
		nblocks := length / 12
		for i := 0; i < nblocks && length > 12; i++ {
			j := i * 12
			a += uint32(*(*uint16)(unsafe.Pointer(&k[j]))) + (uint32(*(*uint16)(unsafe.Pointer(&k[j+2]))) << 16)
			b += uint32(*(*uint16)(unsafe.Pointer(&k[j+4]))) + (uint32(*(*uint16)(unsafe.Pointer(&k[j+6]))) << 16)
			c += uint32(*(*uint16)(unsafe.Pointer(&k[j+8]))) + (uint32(*(*uint16)(unsafe.Pointer(&k[j+10]))) << 16)
			a, b, c = mix(a, b, c)
		}

		tail := k
		if length > 12 {
			tail = k[nblocks*12:]
		}
		switch len(tail) {
		case 12:
			c += uint32(*(*uint16)(unsafe.Pointer(&tail[8]))) + (uint32(*(*uint16)(unsafe.Pointer(&tail[10]))) << 16)
			b += uint32(*(*uint16)(unsafe.Pointer(&tail[4]))) + (uint32(*(*uint16)(unsafe.Pointer(&tail[6]))) << 16)
			a += uint32(*(*uint16)(unsafe.Pointer(&tail[0]))) + (uint32(*(*uint16)(unsafe.Pointer(&tail[2]))) << 16)
		case 11:
			c += uint32(tail[10]) << 16
			fallthrough
		case 10:
			c += uint32(*(*uint16)(unsafe.Pointer(&tail[8])))
			b += uint32(*(*uint16)(unsafe.Pointer(&tail[4]))) + (uint32(*(*uint16)(unsafe.Pointer(&tail[6]))) << 16)
			a += uint32(*(*uint16)(unsafe.Pointer(&tail[0]))) + (uint32(*(*uint16)(unsafe.Pointer(&tail[2]))) << 16)
		case 9:
			c += uint32(tail[8]) << 16
			fallthrough
		case 8:
			b += uint32(*(*uint16)(unsafe.Pointer(&tail[4]))) + (uint32(*(*uint16)(unsafe.Pointer(&tail[6]))) << 16)
			a += uint32(*(*uint16)(unsafe.Pointer(&tail[0]))) + (uint32(*(*uint16)(unsafe.Pointer(&tail[2]))) << 16)
		case 7:
			b += uint32(tail[6]) << 16
			fallthrough
		case 6:
			b += uint32(*(*uint16)(unsafe.Pointer(&tail[4])))
			a += uint32(*(*uint16)(unsafe.Pointer(&tail[0]))) + (uint32(*(*uint16)(unsafe.Pointer(&tail[2]))) << 16)
		case 5:
			b += uint32(tail[6]) << 16
			fallthrough
		case 4:
			a += uint32(*(*uint16)(unsafe.Pointer(&tail[0]))) + (uint32(*(*uint16)(unsafe.Pointer(&tail[2]))) << 16)
		case 3:
			a += uint32(tail[2]) << 16
			fallthrough
		case 2:
			a += uint32(*(*uint16)(unsafe.Pointer(&tail[0]))) + (uint32(*(*uint16)(unsafe.Pointer(&tail[2]))) << 16)
		case 1:
			a += uint32(tail[0]) << 16
		case 0:
			return c
		}
	} else {
		nblocks := length / 12
		for i := 0; i < nblocks && length > 12; i++ {
			a += uint32(uint8(k[i]))
			a += uint32(uint8(k[i+1])) << 8
			a += uint32(uint8(k[i+2])) << 16
			a += uint32(uint8(k[i+3])) << 24
			b += uint32(uint8(k[i+4]))
			b += uint32(uint8(k[i+5])) << 8
			b += uint32(uint8(k[i+6])) << 16
			b += uint32(uint8(k[i+7])) << 24
			c += uint32(uint8(k[i+8]))
			c += uint32(uint8(k[i+9])) << 8
			c += uint32(uint8(k[i+10])) << 16
			c += uint32(uint8(k[i+11])) << 24
			a, b, c = mix(a, b, c)
		}

		tail := k
		if length > 12 {
			tail = k[nblocks*12:]
		}
		switch length {
		case 12:
			c += uint32(uint8(tail[11])) << 24
		case 11:
			c += uint32(uint8(tail[10])) << 16
		case 10:
			c += uint32(uint8(tail[9])) << 8
		case 9:
			c += uint32(uint8(tail[8]))
		case 8:
			b += uint32(uint8(tail[7])) << 24
		case 7:
			b += uint32(uint8(tail[6])) << 16
		case 6:
			b += uint32(uint8(tail[5])) << 8
		case 5:
			b += uint32(uint8(tail[4]))
		case 4:
			a += uint32(uint8(tail[3])) << 24
		case 3:
			a += uint32(uint8(tail[2])) << 16
		case 2:
			a += uint32(uint8(tail[1])) << 8
		case 1:
			a += uint32(uint8(tail[0]))
		case 0:
			return c
		}
	}

	a, b, c = final(a, b, c)
	return c
}
