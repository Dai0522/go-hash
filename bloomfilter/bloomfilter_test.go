package bloomfilter

import (
	"testing"
)

func TestBloomFilter(t *testing.T) {
	write, err := New(1000, 0.0001)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	testCases := []struct {
		set  int64
		get  int64
		want bool
	}{
		{123123, 123123, true},
		{987654321, 987654321, true},
		{0, 0, true},
		{1, 1, true},
		{-110, -110, true},
		{-98765321, -98771231, false},
		{-987621, -1, false},
		{-2, 100, false},
	}
	for _, c := range testCases {
		result := write.PutUint64(uint64(c.set))
		if !result {
			t.Errorf("bloomfilter PutUint64 failed: set[%d] got[%v]", c.set, result)
			t.FailNow()
		}

		result = write.MightContainUint64(uint64(c.get))
		if result != c.want {
			t.Errorf("bloomfilter MightContainUint64 failed: get[%d] want[%v] got[%v]", c.get, c.want, result)
			t.FailNow()
		}
	}

	b := write.Serialized()
	t.Logf("BloomFilter serialized capacity[1000] fpp[0.0001] length[%d]", len(b))
	read, err := Load(b)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for _, c := range testCases {
		result := read.MightContainUint64(uint64(c.get))
		if result != c.want {
			t.Errorf("bloomfilter MightContainUint64 failed: get[%d] want[%v] got[%v]", c.get, c.want, result)
			t.FailNow()
		}
	}
}
