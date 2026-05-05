package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"prime/internal/primes"
)

func main() {
	n := flag.Int("n", 100, "cantidad de primos a generar")
	mersenneN := flag.Int("mersenne-n", 15, "cantidad de primos de Mersenne a generar")
	mode := flag.String("mode", "all", "modo: all, sieve, normal, mersenne, segmented, parallel, wheel, random")
	countOnly := flag.Bool("count-only", false, "mostrar solo la cantidad encontrada")
	workers := flag.Int("workers", 0, "cantidad de workers para modo parallel")
	bits := flag.Int("bits", 128, "bits para generar un primo aleatorio con Miller-Rabin")
	rounds := flag.Int("rounds", 20, "rondas de Miller-Rabin")
	flag.Parse()

	handlers := map[string]func(){
		"all": func() {
			printPrimeList("Primeros %d primos (criba de Eratóstenes optimizada):\n", *n, timedSieve(*n), *countOnly)
			printPrimeList("\nPrimeros %d primos (normal optimizado por división):\n", *n, timedNormal(*n), *countOnly)
			printPrimeList("\nPrimeros %d primos (criba segmentada):\n", *n, timedSegmented(*n), *countOnly)
			printPrimeList("\nPrimeros %d primos (criba paralela):\n", *n, timedParallel(*n, *workers), *countOnly)
			printPrimeList("\nPrimeros %d primos (wheel factorization):\n", *n, timedWheel(*n), *countOnly)
			printRandomPrime(*bits, *rounds)
			printMersenneLive(*mersenneN, *countOnly)
		},
		"sieve": func() {
			printPrimeList("Primeros %d primos (criba de Eratóstenes optimizada):\n", *n, timedSieve(*n), *countOnly)
		},
		"normal": func() {
			printPrimeList("Primeros %d primos (normal optimizado por división):\n", *n, timedNormal(*n), *countOnly)
		},
		"segmented": func() {
			printPrimeList("Primeros %d primos (criba segmentada):\n", *n, timedSegmented(*n), *countOnly)
		},
		"parallel": func() {
			printPrimeList("Primeros %d primos (criba paralela):\n", *n, timedParallel(*n, *workers), *countOnly)
		},
		"wheel": func() {
			printPrimeList("Primeros %d primos (wheel factorization):\n", *n, timedWheel(*n), *countOnly)
		},
		"random": func() {
			printRandomPrime(*bits, *rounds)
		},
		"mersenne": func() {
			printMersenneLive(*mersenneN, *countOnly)
		},
	}

	handler, ok := handlers[*mode]
	if !ok {
		fmt.Fprintf(os.Stderr, "modo inválido %q. Usa: all, sieve, normal, mersenne, segmented, parallel, wheel, random\n", *mode)
		os.Exit(2)
	}
	handler()
}

type timedInts struct {
	values   []int
	duration time.Duration
}

type timedBig struct {
	value    string
	duration time.Duration
	err      error
}

func timedSieve(n int) timedInts {
	start := time.Now()
	values := primes.FirstNBySieve(n)
	return timedInts{values: values, duration: time.Since(start)}
}

func timedNormal(n int) timedInts {
	start := time.Now()
	values := primes.FirstNOptimizedNormal(n)
	return timedInts{values: values, duration: time.Since(start)}
}

func timedSegmented(n int) timedInts {
	start := time.Now()
	values := primes.FirstNBySegmented(n)
	return timedInts{values: values, duration: time.Since(start)}
}

func timedParallel(n int, workers int) timedInts {
	start := time.Now()
	values := primes.FirstNByParallel(n, workers)
	return timedInts{values: values, duration: time.Since(start)}
}

func timedWheel(n int) timedInts {
	start := time.Now()
	values := primes.FirstNByWheel(n)
	return timedInts{values: values, duration: time.Since(start)}
}

func timedRandomPrime(bits int, rounds int) timedBig {
	start := time.Now()
	value, err := primes.RandomPrimeMillerRabin(bits, rounds)
	if err != nil {
		return timedBig{duration: time.Since(start), err: err}
	}
	return timedBig{value: value.String(), duration: time.Since(start)}
}

func printPrimeList(title string, n int, result timedInts, countOnly bool) {
	fmt.Printf(title, n)
	if countOnly {
		fmt.Printf("cantidad: %d\n", len(result.values))
	} else {
		fmt.Println(joinInts(result.values))
	}
	fmt.Printf("tiempo: %.6f ms\n", ms(result.duration))
}

func printRandomPrime(bits int, rounds int) {
	result := timedRandomPrime(bits, rounds)
	fmt.Printf("\nPrimo aleatorio Miller-Rabin de %d bits:\n", bits)
	if result.err != nil {
		fmt.Printf("error: %v\n", result.err)
	} else {
		fmt.Println(result.value)
	}
	fmt.Printf("tiempo: %.6f ms\n", ms(result.duration))
}

func printMersenneLive(n int, countOnly bool) {
	fmt.Printf("\nPrimeros %d primos de Mersenne:\n", n)
	start := time.Now()
	count := 0

	for mp := range primes.FirstNMersennePrimesStream(n) {
		count++
		if !countOnly {
			fmt.Printf("%d) p=%d -> 2^%d - 1 = %s | tiempo acumulado: %.6f ms\n", count, mp.P, mp.P, mp.Value.String(), ms(time.Since(start)))
		}
	}

	if countOnly {
		fmt.Printf("cantidad: %d\n", count)
	}
	fmt.Printf("tiempo total: %.6f ms\n", ms(time.Since(start)))
}

func ms(duration time.Duration) float64 {
	return float64(duration.Nanoseconds()) / 1e6
}

func joinInts(values []int) string {
	if len(values) == 0 {
		return "(ninguno)"
	}
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = fmt.Sprint(value)
	}
	return strings.Join(parts, ", ")
}
