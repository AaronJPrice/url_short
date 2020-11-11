package url_store

import (
	"fmt"
	"testing"

	"github.com/aaronjprice/url_short/shortener"
	"github.com/aaronjprice/url_short/testlib"
)

func TestSimple(t *testing.T) {
	m := New(shortener.NewMapHash())
	long := "https://subdomain.domain.tld/path1"

	short := m.Compress(long)
	result, found := m.Expand(short)

	testlib.Assert(t, found)
	testlib.AssertEqual(t, result, long)
}

func TestMany(t *testing.T) {
	m := New(shortener.NewMapHash())
	l1 := "https://subdomain.domain.tld/path1"
	l2 := "https://www.something.com/else"
	l3 := "https://Merry.Pippin.Frodo/Sam"
	l4 := "https://Fiver.Hazel.BigWig/Silver"
	l5 := "abcdefGHIJKL1234567890!@£$%^&*()_+-={}[]:\"|;'\\<>?,./'~`"

	// Use a different order to expand and do some multiple times
	s1 := m.Compress(l1)
	r1, f1 := m.Expand(s1)
	s2 := m.Compress(l2)
	s3 := m.Compress(l3)
	s4 := m.Compress(l4)
	r4, f4 := m.Expand(s4)
	s5 := m.Compress(l5)
	r3, f3 := m.Expand(s3)
	r5, f5 := m.Expand(s5)
	r2, f2 := m.Expand(s2)
	r22, f22 := m.Expand(s2)

	testlib.Assert(t, f1)
	testlib.Assert(t, f2)
	testlib.Assert(t, f22)
	testlib.Assert(t, f3)
	testlib.Assert(t, f4)
	testlib.Assert(t, f5)
	testlib.AssertEqual(t, r1, l1)
	testlib.AssertEqual(t, r2, l2)
	testlib.AssertEqual(t, r22, l2)
	testlib.AssertEqual(t, r3, l3)
	testlib.AssertEqual(t, r4, l4)
	testlib.AssertEqual(t, r5, l5)
}

func TestSimpleConcurrent(t *testing.T) {
	m := New(shortener.NewMapHash())
	long := "https://subdomain.domain.tld/path1"

	ch := make(chan string)

	go func(m *Map) {
		short := m.Compress(long)
		ch <- short
	}(m)

	short := <-ch

	result, found := m.Expand(short)

	testlib.Assert(t, found)
	testlib.AssertEqual(t, result, long)
}

func TestConcurrentRaceCondition(t *testing.T) {
	m := New(shortener.NewMapHash())
	long := "https://subdomain.domain.tld/path1"
	short := m.Compress(long)

	doneCh := make(chan interface{})

	go func(m *Map) {
		// infinite loop
		for i := 0; i >= 0; i++ {
			select {
			case <-doneCh:
				return
			default:
				_ = m.Compress(fmt.Sprint(i))
			}
		}
	}(m)

	for i := 0; i <= 1000; i++ {
		result, found := m.Expand(short)
		testlib.Assert(t, found)
		testlib.AssertEqual(t, result, long)
	}

	close(doneCh)
}
