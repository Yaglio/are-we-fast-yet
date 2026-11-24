package main

import "fmt"

// Mandelbrot computes a bitmap of the Mandelbrot set; the inner iteration
// count is the image size.
type Mandelbrot struct{}

// InnerBenchmarkLoop overrides the default inner loop (ADR-003).
func (m *Mandelbrot) InnerBenchmarkLoop(innerIterations int) bool {
	return m.verify(m.mandelbrot(innerIterations), innerIterations)
}

// Benchmark is never called, InnerBenchmarkLoop is used instead.
func (*Mandelbrot) Benchmark() any {
	panic("Should never be reached")
}

// VerifyResult is never called, InnerBenchmarkLoop is used instead.
func (*Mandelbrot) VerifyResult(any) bool {
	panic("Should never be reached")
}

func (*Mandelbrot) verify(result, innerIterations int) bool {
	if innerIterations == 500 {
		return result == 191
	}
	if innerIterations == 750 {
		return result == 50
	}
	if innerIterations == 1 {
		return result == 128
	}

	fmt.Printf("No verification result for %d found\n", innerIterations)
	fmt.Printf("Result is: %d\n", result)
	return false
}

func (*Mandelbrot) mandelbrot(size int) int {
	sum := 0
	byteAcc := 0
	bitNum := 0

	y := 0

	for y < size {
		ci := (2.0 * float64(y) / float64(size)) - 1.0
		x := 0

		for x < size {
			zrzr := 0.0
			zi := 0.0
			zizi := 0.0
			cr := (2.0 * float64(x) / float64(size)) - 1.5

			z := 0
			notDone := true
			escape := 0
			for notDone && z < 50 {
				zr := zrzr - zizi + cr
				zi = 2.0*zr*zi + ci

				// preserve recalculation
				zrzr = zr * zr
				zizi = zi * zi

				if zrzr+zizi > 4.0 {
					notDone = false
					escape = 1
				}
				z += 1
			}

			byteAcc = (byteAcc << 1) + escape
			bitNum += 1

			// Code is very similar for these cases, but using separate blocks
			// ensures we skip the shifting when it's unnecessary, which is most cases.
			if bitNum == 8 {
				sum ^= byteAcc
				byteAcc = 0
				bitNum = 0
			} else if x == size-1 {
				byteAcc <<= 8 - bitNum
				sum ^= byteAcc
				byteAcc = 0
				bitNum = 0
			}
			x += 1
		}
		y += 1
	}
	return sum
}
