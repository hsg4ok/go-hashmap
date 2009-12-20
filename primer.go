package main

// Compute primes close to (but greater than) powers of two.
//
// Note that there's some evidence that primes used for hash
// tables should be as far as possible between powers of two
// instead.

import "math"
import "fmt"

const maxExponent = 40

func isPrime(n int64) bool {
	bound := int64(math.Sqrt(float64(n)))+1
	for i := int64(2); i <= bound; i++ {
		if n % i == 0 {
			return false
		}
	}
	return true
}

func main() {
	for i := uint64(0); i <= maxExponent; i++ {
		x := (int64(1)<<i)+1; y := int64(1)<<(i+1)
		for x < y {
			if isPrime(x) {
				fmt.Println(x)
				break
			}
			x += 2
		}
	}
}
