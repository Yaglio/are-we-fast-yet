// The Computer Language Benchmarks Game
// http://shootout.alioth.debian.org/
//
// Based on nbody.java and adapted based on the SOM version.

package main

import (
	"fmt"
	"math"
)

const (
	pi          float64 = 3.141592653589793
	solarMass   float64 = 4 * pi * pi
	daysPerYear float64 = 365.24
)

type body struct {
	x    float64
	y    float64
	z    float64
	vx   float64
	vy   float64
	vz   float64
	mass float64
}

func newBody(x, y, z, vx, vy, vz, mass float64) *body {
	return &body{
		x:    x,
		y:    y,
		z:    z,
		vx:   vx * daysPerYear,
		vy:   vy * daysPerYear,
		vz:   vz * daysPerYear,
		mass: mass * solarMass,
	}
}

func (b *body) offsetMomentum(px, py, pz float64) {
	b.vx = 0.0 - (px / solarMass)
	b.vy = 0.0 - (py / solarMass)
	b.vz = 0.0 - (pz / solarMass)
}

func bodyJupiter() *body {
	return newBody(
		4.84143144246472090e+00,
		-1.16032004402742839e+00,
		-1.03622044471123109e-01,
		1.66007664274403694e-03,
		7.69901118419740425e-03,
		-6.90460016972063023e-05,
		9.54791938424326609e-04)
}

func bodySaturn() *body {
	return newBody(
		8.34336671824457987e+00,
		4.12479856412430479e+00,
		-4.03523417114321381e-01,
		-2.76742510726862411e-03,
		4.99852801234917238e-03,
		2.30417297573763929e-05,
		2.85885980666130812e-04)
}

func bodyUranus() *body {
	return newBody(
		1.28943695621391310e+01,
		-1.51111514016986312e+01,
		-2.23307578892655734e-01,
		2.96460137564761618e-03,
		2.37847173959480950e-03,
		-2.96589568540237556e-05,
		4.36624404335156298e-05)
}

func bodyNeptune() *body {
	return newBody(
		1.53796971148509165e+01,
		-2.59193146099879641e+01,
		1.79258772950371181e-01,
		2.68067772490389322e-03,
		1.62824170038242295e-03,
		-9.51592254519715870e-05,
		5.15138902046611451e-05)
}

func bodySun() *body {
	return newBody(0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0)
}

type nBodySystem struct {
	bodies []*body
}

func newNBodySystem() *nBodySystem {
	s := &nBodySystem{}
	s.bodies = s.createBodies()
	return s
}

func (*nBodySystem) createBodies() []*body {
	bodies := []*body{bodySun(),
		bodyJupiter(),
		bodySaturn(),
		bodyUranus(),
		bodyNeptune()}

	px := 0.0
	py := 0.0
	pz := 0.0

	for _, b := range bodies {
		px += b.vx * b.mass
		py += b.vy * b.mass
		pz += b.vz * b.mass
	}

	bodies[0].offsetMomentum(px, py, pz)

	return bodies
}

func (s *nBodySystem) advance(dt float64) {
	for i := 0; i < len(s.bodies); i++ {
		iBody := s.bodies[i]

		for j := i + 1; j < len(s.bodies); j++ {
			jBody := s.bodies[j]
			dx := iBody.x - jBody.x
			dy := iBody.y - jBody.y
			dz := iBody.z - jBody.z

			dSquared := dx*dx + dy*dy + dz*dz
			distance := math.Sqrt(dSquared)
			mag := dt / (dSquared * distance)

			iBody.vx -= dx * jBody.mass * mag
			iBody.vy -= dy * jBody.mass * mag
			iBody.vz -= dz * jBody.mass * mag

			jBody.vx += dx * iBody.mass * mag
			jBody.vy += dy * iBody.mass * mag
			jBody.vz += dz * iBody.mass * mag
		}
	}

	for _, b := range s.bodies {
		b.x += dt * b.vx
		b.y += dt * b.vy
		b.z += dt * b.vz
	}
}

func (s *nBodySystem) energy() float64 {
	e := 0.0

	for i := 0; i < len(s.bodies); i++ {
		iBody := s.bodies[i]
		e += 0.5 * iBody.mass *
			(iBody.vx*iBody.vx +
				iBody.vy*iBody.vy +
				iBody.vz*iBody.vz)

		for j := i + 1; j < len(s.bodies); j++ {
			jBody := s.bodies[j]
			dx := iBody.x - jBody.x
			dy := iBody.y - jBody.y
			dz := iBody.z - jBody.z

			distance := math.Sqrt(dx*dx + dy*dy + dz*dz)
			e -= (iBody.mass * jBody.mass) / distance
		}
	}
	return e
}

// NBody simulates the Jovian planets; the inner iteration count is the
// number of simulation steps.
type NBody struct{}

// InnerBenchmarkLoop overrides the default inner loop (ADR-003).
func (n *NBody) InnerBenchmarkLoop(innerIterations int) bool {
	system := newNBodySystem()
	for i := 0; i < innerIterations; i++ {
		system.advance(0.01)
	}

	return n.verify(system.energy(), innerIterations)
}

func (*NBody) verify(result float64, innerIterations int) bool {
	if innerIterations == 250000 {
		return result == -0.1690859889909308
	}
	if innerIterations == 1 {
		return result == -0.16907495402506745
	}

	fmt.Printf("No verification result for %d found\n", innerIterations)
	fmt.Printf("Result is: %v\n", result)
	return false
}

// Benchmark is never called, InnerBenchmarkLoop is used instead.
func (*NBody) Benchmark() any {
	panic("Should never be reached")
}

// VerifyResult is never called, InnerBenchmarkLoop is used instead.
func (*NBody) VerifyResult(any) bool {
	panic("Should never be reached")
}
