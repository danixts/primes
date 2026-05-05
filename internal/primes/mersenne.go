package primes

import "math/big"

type MersennePrime struct {
	P     int
	Value *big.Int
}

func MersennePrimes(maxExponent int) []MersennePrime {
	if maxExponent < 2 {
		return []MersennePrime{}
	}

	exponents := Sieve(maxExponent)
	result := make([]MersennePrime, 0, len(exponents))
	for _, p := range exponents {
		if IsMersennePrimeExponent(p) {
			result = append(result, MersennePrime{
				P:     p,
				Value: MersenneNumber(p),
			})
		}
	}
	return result
}

func MersennePrimesStream(maxExponent int) <-chan MersennePrime {
	out := make(chan MersennePrime)
	go func() {
		defer close(out)
		if maxExponent < 2 {
			return
		}
		for _, p := range Sieve(maxExponent) {
			if IsMersennePrimeExponent(p) {
				out <- MersennePrime{P: p, Value: MersenneNumber(p)}
			}
		}
	}()
	return out
}

func FirstNMersennePrimesStream(n int) <-chan MersennePrime {
	out := make(chan MersennePrime)
	go func() {
		defer close(out)
		if n <= 0 {
			return
		}
		count := 0
		for p := 2; count < n; p++ {
			if !isPrimeInt(p) {
				continue
			}
			if IsMersennePrimeExponent(p) {
				out <- MersennePrime{P: p, Value: MersenneNumber(p)}
				count++
			}
		}
	}()
	return out
}

func MersenneNumber(p int) *big.Int {
	if p < 0 {
		return big.NewInt(0)
	}
	return new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(p)), big.NewInt(1))
}

func IsMersennePrimeExponent(p int) bool {
	if p == 2 {
		return true
	}
	if p < 2 || !isPrimeInt(p) {
		return false
	}

	m := MersenneNumber(p)
	s := big.NewInt(4)
	two := big.NewInt(2)

	for i := 0; i < p-2; i++ {
		s.Mul(s, s)
		s.Sub(s, two)
		s.Mod(s, m)
	}
	return s.Sign() == 0
}

func isPrimeInt(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for d := 5; d*d <= n; d += 6 {
		if n%d == 0 || n%(d+2) == 0 {
			return false
		}
	}
	return true
}
