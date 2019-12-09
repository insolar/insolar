package args

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrime(t *testing.T) {
	require.Equal(t, 2, Prime(0))
	require.Equal(t, 3, Prime(1))
	require.Equal(t, MaxPrime, Prime(PrimeCount-1))

	// at least for math.MaxUint16
	require.Equal(t, 65521, Primes().Floor(math.MaxUint16))
	require.Equal(t, 65537, Primes().Ceiling(math.MaxUint16))
	require.Equal(t, 65521, OddPrimes().Floor(math.MaxUint16))
	require.Equal(t, 65537, OddPrimes().Ceiling(math.MaxUint16))
	require.Equal(t, 65497, IsolatedOddPrimes().Floor(math.MaxUint16))
	require.Equal(t, 65543, IsolatedOddPrimes().Ceiling(math.MaxUint16))
}

func TestPrimeMaxAndCount(t *testing.T) {
	require.Equal(t, PrimeCount, len(primes)+len(primesOverUint16))
	require.Equal(t, MaxPrime, int(primesOverUint16[len(primesOverUint16)-1]))

	prev := int(primes[0])
	for i := firstOddPrimeIndex; i < PrimeCount; i++ {
		p := Prime(i)
		require.Equal(t, 1, p&1, i)
		require.Less(t, prev, p, i)
		prev = p
	}
}

func TestPrimes(t *testing.T) {
	for i, p := range primes {
		for _, pp := range primes[i+1:] {
			if pp%p == 0 {
				require.Fail(t, "invalid prime", "%d mod %d=%d", pp, p, pp%p)
			}
		}
	}
}

func TestAllPrimes(t *testing.T) {
	pl := Primes()
	for i := pl.Len() - 1; i >= 0; i-- {
		require.Equal(t, Prime(i), pl.Prime(i), i)
	}
}

func TestOddPrimes(t *testing.T) {
	pl := OddPrimes()
	for i := pl.Len() - 1; i >= 0; i-- {
		require.Equal(t, Prime(i+firstOddPrimeIndex), pl.Prime(i), i)
	}
}

func TestIsolatedPrimes(t *testing.T) {
	op := IsolatedOddPrimes()
	n := 0

	for i := 1; i < PrimeCount; i++ {
		p := Prime(i)
		if Prime(i-1) == p-2 || Prime(i+1) == p+2 {
			continue
		}
		require.Equal(t, p, op.Prime(n))
		n++
		if n == op.Len() {
			return
		}
	}
	require.Fail(t, "insufficient number of isolated primes")
}

func TestOddCeilingPrime(t *testing.T) {
	op := OddPrimes()
	require.Equal(t, 3, op.Ceiling(-1))
	require.Equal(t, 3, op.Ceiling(0))
	require.Equal(t, 3, op.Ceiling(1))
	require.Equal(t, 3, op.Ceiling(2))
	require.Equal(t, 3, op.Ceiling(3))

	for i := 1; i < op.Len(); i++ {
		p := Prime(i + firstOddPrimeIndex)
		for j := Prime(i) + 1; j <= p; j++ {
			require.Equal(t, p, op.Ceiling(int(j)))
		}
	}

	require.Equal(t, op.Max(), op.Ceiling(op.Max()))
	require.Equal(t, op.Max(), op.Ceiling(op.Max()+1))
	require.Equal(t, op.Max(), op.Ceiling(MaxPrime))
	require.Equal(t, op.Max(), op.Ceiling(MaxPrime+1))
	//require.Panics(t, func() {
	//	op.CeilingPrime(MaxPrime + 1)
	//})
}

func TestFloorOddPrime(t *testing.T) {
	op := OddPrimes()
	require.Equal(t, op.Min(), op.Floor(-1))
	require.Equal(t, op.Min(), op.Floor(0))
	require.Equal(t, op.Min(), op.Floor(1))
	require.Equal(t, op.Min(), op.Floor(2))
	require.Equal(t, op.Min(), op.Floor(3))

	for i := 0; i < op.Len()-1; i++ {
		p := Prime(i + firstOddPrimeIndex)
		for j := Prime(i+firstOddPrimeIndex) - 1; j >= p; j-- {
			require.Equal(t, p, op.Floor(int(j)))
		}
	}

	require.Equal(t, op.Max(), op.Floor(op.Max()))
	require.Equal(t, op.Max(), op.Floor(op.Max()+1))
	require.Equal(t, op.Max(), op.Floor(MaxPrime))
	require.Equal(t, op.Max(), op.Floor(MaxPrime+1))
}

func TestNearestOddPrime(t *testing.T) {
	op := OddPrimes()
	require.Equal(t, 3, op.Nearest(-1))
	require.Equal(t, 3, op.Nearest(0))
	require.Equal(t, 3, op.Nearest(1))
	require.Equal(t, 3, op.Nearest(2))
	require.Equal(t, 3, op.Nearest(3))
	require.Equal(t, 3, op.Nearest(4))
	require.Equal(t, 5, op.Nearest(5))
	require.Equal(t, 5, op.Nearest(6))
	require.Equal(t, 7, op.Nearest(7))
	require.Equal(t, 7, op.Nearest(8))
	require.Equal(t, 7, op.Nearest(9))
	require.Equal(t, 11, op.Nearest(10))
	require.Equal(t, 11, op.Nearest(11))
	require.Equal(t, 11, op.Nearest(12))
	require.Equal(t, 13, op.Nearest(13))
	require.Equal(t, 13, op.Nearest(14))
	require.Equal(t, 13, op.Nearest(15))
	require.Equal(t, 17, op.Nearest(16))
	require.Equal(t, 17, op.Nearest(17))
	require.Equal(t, 17, op.Nearest(18))
	require.Equal(t, 19, op.Nearest(19))
	require.Equal(t, 19, op.Nearest(20))
	require.Equal(t, 19, op.Nearest(21))
	require.Equal(t, 23, op.Nearest(22))
	require.Equal(t, 23, op.Nearest(23))
	require.Equal(t, 23, op.Nearest(24))
	require.Equal(t, 23, op.Nearest(25))
	require.Equal(t, 23, op.Nearest(26))
	require.Equal(t, 29, op.Nearest(27))
	require.Equal(t, 29, op.Nearest(28))
	require.Equal(t, 29, op.Nearest(29))

	require.Equal(t, op.Max(), op.Nearest(op.Max()))
	require.Equal(t, op.Max(), op.Nearest(op.Max()+1))
	require.Equal(t, op.Max(), op.Nearest(MaxPrime))
	require.Equal(t, op.Max(), op.Nearest(MaxPrime+1))
}
