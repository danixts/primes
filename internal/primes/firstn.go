package primes

import "math"

func FirstNBySieve(n int) []int {
	return firstNFromGenerator(n, Sieve)
}

func FirstNBySegmented(n int) []int {
	return firstNFromGenerator(n, SieveSegmented)
}

func FirstNByParallel(n int, workers int) []int {
	return firstNFromGenerator(n, func(limit int) []int { return SieveParallel(limit, workers) })
}

func FirstNByWheel(n int) []int {
	return firstNFromGenerator(n, SieveWheel)
}

func FirstNOptimizedNormal(n int) []int {
	if n <= 0 {
		return []int{}
	}
	if n == 1 {
		return []int{2}
	}

	values := make([]int, 0, n)
	values = append(values, 2)
	for candidate := 3; len(values) < n; candidate += 2 {
		prime := true
		for _, p := range values[1:] {
			if p*p > candidate {
				break
			}
			if candidate%p == 0 {
				prime = false
				break
			}
		}
		if prime {
			values = append(values, candidate)
		}
	}
	return values
}

func firstNFromGenerator(n int, generator func(int) []int) []int {
	if n <= 0 {
		return []int{}
	}
	limit := nthPrimeUpperBound(n)
	for {
		values := generator(limit)
		if len(values) >= n {
			return values[:n]
		}
		limit *= 2
	}
}

func nthPrimeUpperBound(n int) int {
	if n < 6 {
		return 15
	}
	x := float64(n)
	return int(math.Ceil(x*(math.Log(x)+math.Log(math.Log(x))))) + 16
}
