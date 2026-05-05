package primes

func Sieve(limit int) []int {
	if limit < 2 {
		return []int{}
	}
	if limit == 2 {
		return []int{2}
	}

	oddCount := (limit - 1) / 2
	composite := make([]bool, oddCount)

	for i := 0; ; i++ {
		p := 2*i + 3
		if p*p > limit {
			break
		}
		if composite[i] {
			continue
		}

		for j := (p*p - 3) / 2; j < oddCount; j += p {
			composite[j] = true
		}
	}

	primes := make([]int, 0, estimatePrimeCount(limit))
	primes = append(primes, 2)
	for i, isComposite := range composite {
		if !isComposite {
			primes = append(primes, 2*i+3)
		}
	}
	return primes
}

func GeneratePrimesOptimizedNormal(limit int) []int {
	if limit < 2 {
		return []int{}
	}

	primes := make([]int, 0, estimatePrimeCount(limit))
	primes = append(primes, 2)

	for candidate := 3; candidate <= limit; candidate += 2 {
		isPrime := true
		for _, p := range primes[1:] {
			if p*p > candidate {
				break
			}
			if candidate%p == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			primes = append(primes, candidate)
		}
	}
	return primes
}

func estimatePrimeCount(limit int) int {
	if limit < 17 {
		return 6
	}
	return limit/10 + 16
}
