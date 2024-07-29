package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

//
// NOTE
// The program can be improved signifigantly.
//
// I could think of a few ways to really speed this up but requires more time
// and since this was a challenge I wanted to get it submitted as soon as I could.
//
// For one the points could be ordered by their N value and large portions of them
// could be occluded if their values would be too large or something.

type InputPt struct {
	X, Y, Z        float64
	N              int
	InputFileIndex int
}

// subtract p1 from p2 and return the resulting vector
func Subtract(p1, p2 *InputPt) *InputPt {
	return &InputPt{
		X: p2.X - p1.X,
		Y: p2.Y - p1.Y,
		Z: p2.Z - p1.Z,
	}
}

// calculate the cross product of two vectors and returrn the resulting vector
func CrossProduct(p1, p2 *InputPt) *InputPt {
	return &InputPt{
		X: p1.Y*p2.Z - p1.Z*p2.Y,
		Y: p1.Z*p2.X - p1.X*p2.Z,
		Z: p1.X*p2.Y - p1.Y*p2.X,
	}
}

// calculate the dot product of two vectors
func DotProduct(p1, p2 *InputPt) float64 {
	return p1.X*p2.X + p1.Y*p2.Y + p1.Z*p2.Z
}

// calculate the volume of four points
func Volume(p1, p2, p3, p4 *InputPt) float64 {
	a := Subtract(p1, p2)
	b := Subtract(p1, p3)
	c := Subtract(p1, p4)
	cross := CrossProduct(b, c)
	volume := math.Abs(DotProduct(a, cross)) / 6.0
	return volume
}

// this type I ended up not really needing
type Tetrahedron struct {
	Points [4]*InputPt
	Volume float64
	Wt     int
}

// Parse an input file of 3d points in the form of (x, y, z, n) into a slice and returns it
func ReadPtsFromFile(path string) ([]*InputPt, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pts := []*InputPt{}

	scanner := bufio.NewScanner(file)
	ix := 0
	for scanner.Scan() {
		item := &InputPt{InputFileIndex: ix}
		ix++
		line := scanner.Text()

		_, err := fmt.Sscanf(line, "(%f, %f, %f, %d)", &item.X, &item.Y, &item.Z, &item.N)
		if err != nil {
			return nil, err
		}

		pts = append(pts, item)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return pts, nil
}

// analyze the points to determine the smallest with a combine N value of 100
func analyzePoints(pts []*InputPt) ([]*InputPt, float64) {
	bestPoints := []*InputPt{}
	bestVol := math.MaxFloat64

	ptCount := len(pts)
	for i := 0; i < ptCount; i++ {
		for j := i + 1; j < ptCount; j++ {
			for k := j + 1; k < ptCount; k++ {
				for l := k + 1; l < ptCount; l++ {
					if pts[i].N+pts[j].N+pts[k].N+pts[l].N == 100 {
						tetra := &Tetrahedron{
							Points: [4]*InputPt{
								pts[i],
								pts[j],
								pts[k],
								pts[l],
							},
							Volume: Volume(pts[i], pts[j], pts[k], pts[l]),
							Wt:     100,
						}
						if tetra.Volume > 0 && tetra.Volume < bestVol {
							bestVol = tetra.Volume
							bestPoints = tetra.Points[:]
						}

					}
				}
			}
		}
	}

	return bestPoints, bestVol
}

func main() {

	// points_small.txt
	// Small - 0, 5, 11, 76 - 12.432096
	ptsSm, err := ReadPtsFromFile("input/points_small.txt")
	if err != nil {
		panic(err)
	}

	bestSmPts, bestSmVol := analyzePoints(ptsSm)

	fmt.Printf("[small file] Best Points:\n")
	for _, p := range bestSmPts {
		fmt.Printf("  [%f, %f, %f, %d] - %d\n", p.X, p.Y, p.Z, p.N, p.InputFileIndex)
	}
	fmt.Printf(" with a volume of %f.\n", bestSmVol)

	// points_large.txt
	// Large - 70, 386, 493, 1429 - 0.000375
	fmt.Println("Please wait, this could take a while...")

	ptsLg, err := ReadPtsFromFile("input/points_large.txt")
	if err != nil {
		panic(err)
	}

	bestLgPts, bestLgVol := analyzePoints(ptsLg)

	fmt.Printf("[large file] Best Points:\n")
	for _, p := range bestLgPts {
		fmt.Printf("  [%f, %f, %f, %d] - %d\n", p.X, p.Y, p.Z, p.N, p.InputFileIndex)
	}
	fmt.Printf(" with a volume of %f.\n", bestLgVol)
}
