package main

import (
	closests "L4/pkg/closest"
	"fmt"
	"os"
	"time"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

func main() {
	generators := []func([2]float64, [2]float64, uint, uint64) [][2]float64{
		randXNormalYNormal,
		randXUniformYNormal,
		randXNormalYExponential,
		randXUniformYExponential,
		randXNormalYPoisson,
		randXUniformYPoisson,
		randXNormalYUniform,
		randXUniformYUniform,
		randXNormalYNormalSame,
		randXNormalYNormaUnbounded,
	}

	filenames := []string{
		"data/randXNormalYNormal.csv",
		"data/randXUniformYNormal.csv",
		"data/randXNormalYExponential.csv",
		"data/randXUniformYExponential.csv",
		"data/randXNormalYPoisson.csv",
		"data/randXUniformYPoisson.csv",
		"data/randXNormalYUniform.csv",
		"data/randXUniformYUniform.csv",
		"data/randXNormalYNormalSame.csv",
		"data/randXNormalYNormaUnbounded.csv",
	}

	for i, g := range generators {
		experiment(g, filenames[i])
	}
}

type Res struct {
	n uint
	d float64
	t int64
	c uint64
}

func (r *Res) String() string {
	return fmt.Sprintf("%v;%v;%v;%v", r.n, r.d, r.t, r.c)
}

func experiment(generator func([2]float64, [2]float64, uint, uint64) [][2]float64, filename string) {

	res := make([]Res, 0)
	for n := uint(500); n < 30000; n += 500 {
		for i := 0; i < 300; i++ {
			points := generator([2]float64{0, 1}, [2]float64{0, 1}, n, uint64(time.Now().UnixNano()))
			start := time.Now()
			_, d, c := closests.FindClosests(points)
			elapsed := time.Since(start).Nanoseconds()

			res = append(res, Res{n: n, d: d, t: elapsed, c: c})
		}
		fmt.Println("Done "+filename+" for n = ", n)
	}

	f, e := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)

	if e != nil {
		fmt.Println(e.Error())
		return
	}

	defer f.Close()
	fmt.Fprintln(f, "n;d;t;c")
	for _, r := range res {
		fmt.Fprintln(f, r.String())
	}

}

func randomPoints(xBounds, yBounds [2]float64, n uint) [][2]float64 {
	points := make([][2]float64, n)
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	for i := range points {
		points[i][0] = xBounds[0] + r.Float64()*(xBounds[1]-xBounds[0])
		points[i][1] = yBounds[0] + r.Float64()*(yBounds[1]-yBounds[0])
	}

	return points
}

func printPoints(points [][2]float64, name string) {
	fmt.Println(name, "= [")
	for _, p := range points {
		fmt.Println("\t[", p[0], ", ", p[1], "],")
	}
	fmt.Println("]")
}

func randXNormalYUniform(xBounds, yBounds [2]float64, n uint, seed uint64) [][2]float64 {
	dist1 := distuv.Normal{
		Mu:    0.5,
		Sigma: 0.5 / 3.0,
		Src:   rand.NewSource(seed),
	}
	dist2 := distuv.Uniform{
		Min: 0.0,
		Max: 1.0,
		Src: rand.NewSource(seed),
	}

	points := make([][2]float64, n)

	for i := range points {
		points[i][0] = cutValue(0.0, 1.0, dist1.Rand())
		points[i][1] = cutValue(0.0, 1.0, dist2.Rand())
	}

	return points
}

func randXUniformYNormal(xBounds, yBounds [2]float64, n uint, seed uint64) [][2]float64 {
	dist1 := distuv.Uniform{
		Min: 0.0,
		Max: 1.0,
		Src: rand.NewSource(seed),
	}

	dist2 := distuv.Normal{
		Mu:    0.5,
		Sigma: 0.5 / 3.0,
		Src:   rand.NewSource(seed),
	}

	points := make([][2]float64, n)

	for i := range points {
		points[i][0] = cutValue(0.0, 1.0, dist1.Rand())
		points[i][1] = cutValue(0.0, 1.0, dist2.Rand())
	}

	return points
}

func randXNormalYNormal(xBounds, yBounds [2]float64, n uint, seed uint64) [][2]float64 {
	dist1 := distuv.Normal{
		Mu:    0.5,
		Sigma: 0.5 / 3.0,
		Src:   rand.NewSource(seed),
	}
	dist2 := distuv.Normal{
		Mu:    0.7,
		Sigma: 0.1,
		Src:   rand.NewSource(seed),
	}

	points := make([][2]float64, n)

	for i := range points {
		points[i][0] = cutValue(0.0, 1.0, dist1.Rand())
		points[i][1] = cutValue(0.0, 1.0, dist2.Rand())
	}

	return points
}

func randXNormalYNormalSame(xBounds, yBounds [2]float64, n uint, seed uint64) [][2]float64 {
	dist1 := distuv.Normal{
		Mu:    0.5,
		Sigma: 0.5 / 3.0,
		Src:   rand.NewSource(seed),
	}
	dist2 := distuv.Normal{
		Mu:    0.5,
		Sigma: 0.5 / 3.0,
		Src:   rand.NewSource(seed),
	}

	points := make([][2]float64, n)

	for i := range points {
		points[i][0] = cutValue(0.0, 1.0, dist1.Rand())
		points[i][1] = cutValue(0.0, 1.0, dist2.Rand())
	}

	return points
}

func randXNormalYNormaUnbounded(xBounds, yBounds [2]float64, n uint, seed uint64) [][2]float64 {
	dist1 := distuv.Normal{
		Mu:    0.5,
		Sigma: 0.5 / 3.0,
		Src:   rand.NewSource(seed),
	}
	dist2 := distuv.Normal{
		Mu:    0.5,
		Sigma: 0.5 / 3.0,
		Src:   rand.NewSource(seed),
	}

	points := make([][2]float64, n)

	for i := range points {
		points[i][0] = dist1.Rand()
		points[i][1] = dist2.Rand()
	}

	return points
}

func randXUniformYUniform(xBounds, yBounds [2]float64, n uint, seed uint64) [][2]float64 {
	dist1 := distuv.Uniform{
		Min: 0.0,
		Max: 1.0,
		Src: rand.NewSource(seed),
	}

	dist2 := distuv.Normal{
		Mu:    0.5,
		Sigma: 0.5 / 3.0,
		Src:   rand.NewSource(seed),
	}

	points := make([][2]float64, n)

	for i := range points {
		points[i][0] = cutValue(0.0, 1.0, dist1.Rand())
		points[i][1] = cutValue(0.0, 1.0, dist2.Rand())
	}

	return points
}

func randXNormalYExponential(xBounds, yBounds [2]float64, n uint, seed uint64) [][2]float64 {
	dist1 := distuv.Normal{
		Mu:    0.5,
		Sigma: 0.5 / 3.0,
		Src:   rand.NewSource(seed),
	}
	dist2 := distuv.Exponential{
		Rate: 0.2,
		Src:  rand.NewSource(seed),
	}

	points := make([][2]float64, n)

	for i := range points {
		points[i][0] = cutValue(0.0, 1.0, dist1.Rand())
		points[i][1] = cutValue(0.0, 1.0, dist2.Rand())
	}

	return points
}

func randXUniformYExponential(xBounds, yBounds [2]float64, n uint, seed uint64) [][2]float64 {
	dist1 := distuv.Uniform{
		Min: 0.0,
		Max: 1.0,
		Src: rand.NewSource(seed),
	}
	dist2 := distuv.Exponential{
		Rate: 0.2,
		Src:  rand.NewSource(seed),
	}

	points := make([][2]float64, n)

	for i := range points {
		points[i][0] = cutValue(0.0, 1.0, dist1.Rand())
		points[i][1] = cutValue(0.0, 1.0, dist2.Rand())
	}

	return points
}

func randXNormalYPoisson(xBounds, yBounds [2]float64, n uint, seed uint64) [][2]float64 {
	dist1 := distuv.Normal{
		Mu:    0.5,
		Sigma: 0.5 / 3.0,
		Src:   rand.NewSource(seed),
	}
	dist2 := distuv.Poisson{
		Lambda: 0.2,
		Src:    rand.NewSource(seed),
	}

	points := make([][2]float64, n)

	for i := range points {
		points[i][0] = cutValue(0.0, 1.0, dist1.Rand())
		points[i][1] = cutValue(0.0, 1.0, dist2.Rand())
	}

	return points
}

func randXUniformYPoisson(xBounds, yBounds [2]float64, n uint, seed uint64) [][2]float64 {
	dist1 := distuv.Uniform{
		Min: 0.0,
		Max: 1.0,
		Src: rand.NewSource(seed),
	}
	dist2 := distuv.Poisson{
		Lambda: 0.2,
		Src:    rand.NewSource(seed),
	}

	points := make([][2]float64, n)

	for i := range points {
		points[i][0] = cutValue(0.0, 1.0, dist1.Rand())
		points[i][1] = cutValue(0.0, 1.0, dist2.Rand())
	}

	return points
}

func cutValue(min, max, val float64) float64 {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}
