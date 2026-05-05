# Prime

Proyecto en Go para generar y comparar números primos con diferentes algoritmos.

## Qué hace

Este CLI permite generar los primeros `n` números primos usando varias estrategias:

- Criba de Eratóstenes optimizada.
- Generación normal optimizada por división.
- Criba segmentada.
- Criba paralela con goroutines.
- Wheel factorization.
- Primos de Mersenne, pintados uno a uno mientras se encuentran.
- Generación de un primo aleatorio grande usando Miller-Rabin.

Cada ejecución muestra el tiempo en milisegundos.

## Requisitos

- Go instalado.

## Uso rápido

```bash
go run ./cmd/prime
```

## Modos disponibles

```bash
go run ./cmd/prime -mode all
go run ./cmd/prime -mode sieve
go run ./cmd/prime -mode normal
go run ./cmd/prime -mode segmented
go run ./cmd/prime -mode parallel
go run ./cmd/prime -mode wheel
go run ./cmd/prime -mode mersenne
go run ./cmd/prime -mode random
```

## Ejemplos

Generar los primeros 10 primos con criba:

```bash
go run ./cmd/prime -mode sieve -n 10
```

Generar los primeros 10 primos con todos los algoritmos y mostrar solo cantidades:

```bash
go run ./cmd/prime -mode all -n 10 -count-only
```

Generar 15 primos de Mersenne por defecto:

```bash
go run ./cmd/prime -mode mersenne
```

Generar otra cantidad de primos de Mersenne:

```bash
go run ./cmd/prime -mode mersenne -mersenne-n 5
```

Generar un primo aleatorio de 256 bits:

```bash
go run ./cmd/prime -mode random -bits 256 -rounds 20
```

Usar criba paralela con 8 workers:

```bash
go run ./cmd/prime -mode parallel -n 100000 -workers 8 -count-only
```

## Flags principales

| Flag | Descripción | Valor por defecto |
|---|---|---|
| `-mode` | Algoritmo a ejecutar | `all` |
| `-n` | Cantidad de primos normales a generar | `100` |
| `-mersenne-n` | Cantidad de primos de Mersenne | `15` |
| `-count-only` | Muestra solo cantidad y tiempo | `false` |
| `-workers` | Workers para criba paralela | CPUs disponibles |
| `-bits` | Bits del primo aleatorio Miller-Rabin | `128` |
| `-rounds` | Rondas de Miller-Rabin | `20` |

## Estructura

```txt
cmd/prime/main.go          CLI
internal/primes/sieve.go   Criba y método normal optimizado
internal/primes/advanced.go Algoritmos segmentado, paralelo, wheel y Miller-Rabin
internal/primes/firstn.go  Adaptadores para primeros n primos
internal/primes/mersenne.go Primos de Mersenne
```
