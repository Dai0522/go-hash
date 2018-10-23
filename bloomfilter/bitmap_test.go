package bloomfilter

import (
	"math"
	"testing"
)

func TestLockFreeBitmap_Set(t *testing.T) {
	bits, err := NewLockFreeBitmap(uint64(19200))
	if err != nil {
		panic(err)
	}
	for _, unit := range []struct {
		x        uint64
		expected bool
	}{
		{uint64(1), true},
		{uint64(100), true},
		{uint64(1000), true},
		{uint64(10000), true},
		{uint64(19199), true},
		{uint64(10000), false},
		{uint64(1000), false},
		{uint64(1), false},
	} {
		if ac := bits.Set(unit.x); ac != unit.expected {
			t.Errorf("LockFreeBitmap_Set [%v], actually: [%v]", unit, ac)
		}
	}
}

func TestLockFreeBitmap_Get(t *testing.T) {
	bits, err := NewLockFreeBitmap(uint64(19200))
	if err != nil {
		panic(err)
	}
	for _, unit := range []struct {
		x        uint64
		y        uint64
		expected bool
	}{
		{uint64(0), uint64(1), false},
		{uint64(1), uint64(0), true},
		{uint64(1), uint64(1), true},
		{uint64(1000), uint64(1000), true},
		{uint64(1), uint64(10000), false},
	} {
		bits.Set(unit.x)
		if ac := bits.Get(unit.y); ac != unit.expected {
			t.Errorf("LockFreeBitmap_Get [%v], actually: [%v]", unit, ac)
		}
	}
}

func TestLockFreeBitmap_Merge(t *testing.T) {
	bits, err := NewLockFreeBitmap(uint64(19200))
	if err != nil {
		panic(err)
	}
	tmp, err := NewLockFreeBitmap(uint64(19200))
	if err != nil {
		panic(err)
	}
	tmp.Set(uint64(1))
	tmp.Set(uint64(100))
	tmp.Set(uint64(1000))
	tmp.Set(uint64(10000))
	for _, unit := range []struct {
		x        uint64
		expected bool
	}{
		{uint64(1), false},
		{uint64(100), false},
		{uint64(1000), false},
		{uint64(10000), false},
	} {
		if ac := bits.Get(unit.x); ac != unit.expected {
			t.Errorf("LockFreeBitmap_Merge before Get [%v], actually: [%v]", unit, ac)
		}
	}
	if ac := bits.Merge(&tmp.data); ac != true {
		t.Errorf("LockFreeBitmap_Merge error src cap[%v], actually: [%v]", len(bits.data), ac)
	}
	for _, unit := range []struct {
		x        uint64
		expected bool
	}{
		{uint64(1), true},
		{uint64(100), true},
		{uint64(1000), true},
		{uint64(10000), true},
	} {
		if ac := bits.Get(unit.x); ac != unit.expected {
			t.Errorf("LockFreeBitmap_Merge after Get [%v], actually: [%v]", unit, ac)
		}
	}

	empty, err := NewLockFreeBitmap(uint64(1))
	if err != nil {
		panic(err)
	}
	if ac := bits.Merge(&empty.data); ac != false {
		t.Errorf("LockFreeBitmap_Merge error src cap[%v], actually: [%v]", len(bits.data), ac)
	}
}

func TestLockFreeBitmap_BitSize(t *testing.T) {
	bits, err := NewLockFreeBitmap(uint64(19200))
	if err != nil {
		panic(err)
	}
	if size := bits.BitSize(); size != 19200 {
		t.Errorf("LockFreeBitmap_BitSize [19200], actually: [%v]", size)
	}
}

func TestLockFreeBitmap_BitCount(t *testing.T) {
	bits, err := NewLockFreeBitmap(uint64(19200))
	if err != nil {
		panic(err)
	}
	for _, unit := range []struct {
		x        uint64
		y        uint64
		expected bool
	}{
		{uint64(1), 1, true},
		{uint64(100), 2, true},
		{uint64(1000), 3, true},
		{uint64(10000), 4, true},
		{uint64(19199), 5, true},
		{uint64(10000), 5, true},
		{uint64(1000), 6, false},
		{uint64(1), 7, false},
	} {
		bits.Set(unit.x)
		if ac := bits.BitCount(); (ac == unit.y) != unit.expected {
			t.Errorf("LockFreeBitmap_BitCount [%v], actually: [%v]", unit, ac)
		}
	}
}

func TestLockFreeBitmap_Size(t *testing.T) {
	bits, err := NewLockFreeBitmap(uint64(19200))
	if err != nil {
		panic(err)
	}
	if ac := bits.Size(); ac != (uint32(math.Ceil(19200 / 64))) {
		t.Errorf("LockFreeBitmap_Size [19200], actually: [%v]", ac)
	}
}
