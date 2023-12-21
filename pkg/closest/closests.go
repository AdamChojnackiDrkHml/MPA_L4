package closests

import (
	"math"
	"sort"
)

func FindClosests(points [][2]float64) ([][2]float64, float64, uint64) {
	sort.Slice(points, func(i int, j int) bool {
		return points[i][1] < points[j][1]
	})

	return findRecursive(points)
}

func findRecursive(points [][2]float64) (closestsPoints [][2]float64, d float64, c uint64) {
	if len(points) == 3 {
		closestsPoints, d = dealWithThree(points)
		return closestsPoints, d, 2
	}
	if len(points) == 2 {
		return [][2]float64(points), distance(points[0], points[1]), 1
	}

	mid := uint(len(points) / 2)

	l, ld, cl := findRecursive(points[mid:])
	r, rd, cr := findRecursive(points[:mid])

	minPoints, d := min3(l, r, ld, rd)

	Sy := findMiddle(points, mid, d)

	comps := uint64(0)
	for i := range Sy {
		for j := i + 1; j < len(Sy) && j < 15; j++ {
			minPoints, d = min2([][2]float64{Sy[i], Sy[j]}, minPoints, d)
			comps++
		}
	}
	return minPoints, d, comps + cl + cr
}

func distance(p1, p2 [2]float64) float64 {
	return math.Sqrt(math.Pow(p1[0]-p2[0], 2.0) + math.Pow(p1[1]-p2[1], 2.0))
}

func dealWithThree(threePoints [][2]float64) ([][2]float64, float64) {
	dL := distance(threePoints[0], threePoints[1])
	dR := distance(threePoints[1], threePoints[2])
	if dL < dR {
		return threePoints[:2], dL
	}
	return threePoints[1:], dR
}

func findMiddle(points [][2]float64, mid uint, d float64) [][2]float64 {
	var high, low uint
	for i := mid; i > 0; i-- {
		low = i
		if points[mid][1]-points[i][1] > d {
			break
		}
	}

	for i := mid; i < uint(len(points)-1); i++ {
		high = i
		if points[i][1]-points[mid][1] > d {
			break
		}
	}

	return points[low : high+1]
}

func min(l, r [][2]float64) ([][2]float64, float64) {
	ld := distance(l[0], l[1])
	rd := distance(r[0], r[1])
	if ld < rd {
		return l, ld
	}

	return r, rd
}

func min2(l, r [][2]float64, d float64) ([][2]float64, float64) {
	ld := distance(l[0], l[1])
	if ld < d {
		return l, ld
	}

	return r, d
}

func min3(l, r [][2]float64, ld, rd float64) ([][2]float64, float64) {
	if ld < rd {
		return l, ld
	}

	return r, rd
}
