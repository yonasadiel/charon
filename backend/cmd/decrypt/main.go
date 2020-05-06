package main

import (
	"bufio"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func getNextPrime(startValue *big.Int) *big.Int {
	var z *big.Int = new(big.Int)
	*z = *startValue
	if new(big.Int).And(z, big.NewInt(1)).Cmp(big.NewInt(1)) != 0 {
		z = new(big.Int).Add(z, big.NewInt(1))
	}
	for !z.ProbablyPrime(40) {
		z = z.Add(z, big.NewInt(2))
	}
	return z
}

func generateRsa(pseed, qseed int64) (*big.Int, *big.Int, *big.Int) {
	var primelength uint = 1024
	var twoPower *big.Int = new(big.Int).Lsh(big.NewInt(1), primelength)
	var prandSrc *rand.Rand = rand.New(rand.NewSource(pseed))
	var qrandSrc *rand.Rand = rand.New(rand.NewSource(qseed))
	var five = big.NewInt(5)
	var prand *big.Int = new(big.Int).Rand(prandSrc, twoPower)
	var qrand *big.Int = new(big.Int).Rand(qrandSrc, twoPower)

	var p *big.Int = getNextPrime(new(big.Int).Exp(five, prand, twoPower))
	var q *big.Int = getNextPrime(new(big.Int).Exp(five, qrand, twoPower))

	var n *big.Int = new(big.Int).Mul(p, q)

	return p, q, n
}

func calculateSquaringsPerSecond() int64 {
	var guess int64 = 1000000
	var b *big.Int = big.NewInt(2)
	var t *big.Int = big.NewInt(guess) // do 1000.000 squarings
	_, _, n := generateRsa(123, 456)
	start := time.Now()
	for i := big.NewInt(0); i.Cmp(t) == -1; i = i.Add(i, big.NewInt(1)) {
		b = b.Mul(b, b)
		b = b.Mod(b, n)
	}
	duration := time.Since(start)
	return int64(float64(guess) / duration.Seconds())
}

func createPuzzle(minutes int64) (*big.Int, *big.Int, *big.Int) {
	var sqrPerSecond *big.Int = big.NewInt(calculateSquaringsPerSecond())
	fmt.Printf("Assumed number of sqr/second = %s\n", sqrPerSecond.String())

	var sqrPerMinute *big.Int = new(big.Int).Mul(sqrPerSecond, big.NewInt(60))
	fmt.Printf("Squarings per minute = %s\n", sqrPerMinute)

	var t *big.Int = new(big.Int).Mul(big.NewInt(minutes), sqrPerMinute)
	fmt.Printf("Squarings total = %s\n", t.String())

	p, q, n := generateRsa(123, 456)
	fmt.Printf("p = %s\n", p.String())
	fmt.Printf("q = %s\n", q.String())
	fmt.Printf("n = %s\n", n.String())

	var pm1 *big.Int = new(big.Int).Sub(p, big.NewInt(1))
	var qm1 *big.Int = new(big.Int).Sub(q, big.NewInt(1))
	var phi *big.Int = new(big.Int).Mul(pm1, qm1)
	fmt.Printf("phi = %s\n", phi.String())

	var a *big.Int = big.NewInt(2)
	var eps *big.Int = new(big.Int).Exp(big.NewInt(2), t, phi)
	var b *big.Int = new(big.Int).Exp(a, eps, n)
	fmt.Printf("b = %s\n", b.Text(16))
	return a, n, t
}

func solvePuzzle(a, n, t *big.Int) {
	var b *big.Int = new(big.Int)
	*b = *a
	for i := big.NewInt(0); i.Cmp(t) == -1; i = i.Add(i, big.NewInt(1)) {
		b = b.Mul(b, b) // b = b * b
		b = b.Mod(b, n) // b = b % n
	}
	fmt.Printf("b = %s\n", b.Text(16))
}

func main() {
	fmt.Printf("OS: %s\nArchitecture: %s\n", runtime.GOOS, runtime.GOARCH)
	start := time.Now()
	a, n, t := createPuzzle(10)
	fmt.Printf("createPuzzle took %v\n", time.Since(start))
	solvePuzzle(a, n, t)
	fmt.Printf("solvePuzzle took %v\n", time.Since(start))

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Press enter to exit ")
	reader.ReadString('\n')
}
