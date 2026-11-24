package main

import "awfy/som"

type ball struct {
	x    int
	y    int
	xVel int
	yVel int
}

func newBall(random *som.Random) *ball {
	return &ball{
		x:    random.Next() % 500,
		y:    random.Next() % 500,
		xVel: (random.Next() % 300) - 150,
		yVel: (random.Next() % 300) - 150,
	}
}

func (b *ball) bounce() bool {
	xLimit := 500
	yLimit := 500
	bounced := false

	b.x += b.xVel
	b.y += b.yVel
	if b.x > xLimit {
		b.x = xLimit
		b.xVel = 0 - abs(b.xVel)
		bounced = true
	}
	if b.x < 0 {
		b.x = 0
		b.xVel = abs(b.xVel)
		bounced = true
	}
	if b.y > yLimit {
		b.y = yLimit
		b.yVel = 0 - abs(b.yVel)
		bounced = true
	}
	if b.y < 0 {
		b.y = 0
		b.yVel = abs(b.yVel)
		bounced = true
	}
	return bounced
}

// Bounce simulates balls bouncing in a box.
type Bounce struct{}

// Benchmark runs one iteration of the benchmark.
func (*Bounce) Benchmark() any {
	random := som.NewRandom()

	ballCount := 100
	bounces := 0
	balls := make([]*ball, ballCount)

	for i := range balls {
		balls[i] = newBall(random)
	}

	for i := 0; i < 50; i++ {
		for _, ball := range balls {
			if ball.bounce() {
				bounces += 1
			}
		}
	}
	return bounces
}

// VerifyResult checks the benchmark result.
func (*Bounce) VerifyResult(result any) bool {
	return 1331 == result.(int)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
