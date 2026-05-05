package primes

import (
	cryptorand "crypto/rand"
	"errors"
	"math"
	"math/big"
	"runtime"
	"sort"
	"sync"
)

func SieveSegmented(limit int) []int {
	if limit < 2 {
		return []int{}
	}

	root := int(math.Sqrt(float64(limit)))
	base := Sieve(root)
	result := make([]int, 0, estimatePrimeCount(limit))
	result = append(result, base...)
	segmentSize := max(root, 32768)

	for low := root + 1; low <= limit; low += segmentSize {
		high := min(low+segmentSize-1, limit)
		composite := make([]bool, high-low+1)

		for _, p := range base {
			start := max(p*p, ((low+p-1)/p)*p)
			for n := start; n <= high; n += p {
				composite[n-low] = true
			}
		}

		for n := low; n <= high; n++ {
			if n > 1 && !composite[n-low] {
				result = append(result, n)
			}
		}
	}

	return result
}

func SieveParallel(limit int, workers int) []int {
	if limit < 2 {
		return []int{}
	}
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	root := int(math.Sqrt(float64(limit)))
	base := Sieve(root)
	segmentSize := max(root, 32768)
	segments := make([]segment, 0)

	for low := root + 1; low <= limit; low += segmentSize {
		segments = append(segments, segment{low: low, high: min(low+segmentSize-1, limit)})
	}

	jobs := make(chan segment)
	results := make(chan segmentResult, len(segments))
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for s := range jobs {
				results <- segmentResult{low: s.low, values: sieveSegment(s.low, s.high, base)}
			}
		}()
	}

	go func() {
		for _, s := range segments {
			jobs <- s
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	chunks := make([]segmentResult, 0, len(segments))
	for r := range results {
		chunks = append(chunks, r)
	}

	sort.Slice(chunks, func(i, j int) bool { return chunks[i].low < chunks[j].low })
	result := make([]int, 0, estimatePrimeCount(limit))
	result = append(result, base...)
	for _, chunk := range chunks {
		result = append(result, chunk.values...)
	}
	return result
}

func SieveWheel(limit int) []int {
	if limit < 2 {
		return []int{}
	}

	result := make([]int, 0, estimatePrimeCount(limit))
	for _, p := range []int{2, 3, 5} {
		if p <= limit {
			result = append(result, p)
		}
	}

	residues := []int{1, 7, 11, 13, 17, 19, 23, 29}
	for base := 0; base <= limit; base += 30 {
		for _, residue := range residues {
			candidate := base + residue
			if candidate <= 5 || candidate > limit {
				continue
			}
			prime := true
			for _, p := range result[3:] {
				if p*p > candidate {
					break
				}
				if candidate%p == 0 {
					prime = false
					break
				}
			}
			if prime {
				result = append(result, candidate)
			}
		}
	}
	return result
}

func RandomPrimeMillerRabin(bits int, rounds int) (*big.Int, error) {
	if bits < 2 {
		return nil, errors.New("bits debe ser mayor o igual a 2")
	}
	if rounds < 1 {
		rounds = 20
	}

	one := big.NewInt(1)
	for {
		candidate, err := cryptorand.Int(cryptorand.Reader, new(big.Int).Lsh(one, uint(bits)))
		if err != nil {
			return nil, err
		}
		candidate.SetBit(candidate, bits-1, 1)
		candidate.SetBit(candidate, 0, 1)
		if IsProbablePrimeMillerRabin(candidate, rounds) {
			return candidate, nil
		}
	}
}

func IsProbablePrimeMillerRabin(n *big.Int, rounds int) bool {
	if n.Cmp(big.NewInt(2)) == 0 || n.Cmp(big.NewInt(3)) == 0 {
		return true
	}
	if n.Sign() <= 0 || n.Bit(0) == 0 || n.Cmp(big.NewInt(2)) < 0 {
		return false
	}

	smallPrimes := []int64{3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}
	for _, p := range smallPrimes {
		bp := big.NewInt(p)
		if n.Cmp(bp) == 0 {
			return true
		}
		if new(big.Int).Mod(n, bp).Sign() == 0 {
			return false
		}
	}

	d := new(big.Int).Sub(n, big.NewInt(1))
	s := 0
	for d.Bit(0) == 0 {
		d.Rsh(d, 1)
		s++
	}

	nMinusOne := new(big.Int).Sub(n, big.NewInt(1))
	nMinusThree := new(big.Int).Sub(n, big.NewInt(3))
	two := big.NewInt(2)

	for i := 0; i < rounds; i++ {
		a, err := cryptorand.Int(cryptorand.Reader, nMinusThree)
		if err != nil {
			return false
		}
		a.Add(a, two)
		x := new(big.Int).Exp(a, d, n)
		if x.Cmp(big.NewInt(1)) == 0 || x.Cmp(nMinusOne) == 0 {
			continue
		}
		witness := true
		for r := 1; r < s; r++ {
			x.Mul(x, x)
			x.Mod(x, n)
			if x.Cmp(nMinusOne) == 0 {
				witness = false
				break
			}
		}
		if witness {
			return false
		}
	}
	return true
}

type segment struct {
	low  int
	high int
}

type segmentResult struct {
	low    int
	values []int
}

func sieveSegment(low int, high int, base []int) []int {
	composite := make([]bool, high-low+1)
	for _, p := range base {
		start := max(p*p, ((low+p-1)/p)*p)
		for n := start; n <= high; n += p {
			composite[n-low] = true
		}
	}
	values := make([]int, 0)
	for n := low; n <= high; n++ {
		if n > 1 && !composite[n-low] {
			values = append(values, n)
		}
	}
	return values
}
