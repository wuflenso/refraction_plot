package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/wuflenso/refraction_model/refraction"
	"github.com/wuflenso/refraction_model/refraction/utilities"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

var (
	populateData = func(x []float64, y []float64) (plotter.XYs, error) {
		if len(x) != len(y) {
			return nil, errors.New("differing x and y data count")
		}

		pts := make(plotter.XYs, len(x))
		for i, _ := range pts {
			pts[i].X = x[i]
			pts[i].Y = y[i]
		}
		return pts, nil
	}
)

func main() {
	// Inputs
	// Convert angles from Degrees to Radians first
	velocities := []float64{200, 400, 500, 550, 300, 350, 100}
	layerThicknesses := []float64{-500, -300, -500, -200, -500, -1000, -500}
	grids := [][]float64{{0, 0}}

	var initialAngleDeg float64 = -5
	currentAngle := initialAngleDeg
	angleIncrement := -5.00
	iterations := 5

	xArrMaxes := []float64{}

	linePlots := []*plotter.Line{}

	for i := 0; i < iterations; i++ {
		angles := []float64{utilities.DegreeToRadians(currentAngle)}

		// Execute function
		resGrids, resAnglesRad, message := refraction.TraceRayRefraction(layerThicknesses, velocities, grids, angles)

		// Convert to angles to degrees
		var resAnglesDeg []float64
		for _, o := range resAnglesRad {
			resAnglesDeg = append(resAnglesDeg, utilities.RadiansToDegree(o))
		}

		// Print results
		for i, _ := range resGrids {
			s := fmt.Sprintf("Coordinate: %.2f, θ2: %.2f°", resGrids[i], resAnglesDeg[i])
			fmt.Println(s)
		}

		fmt.Println(message)

		// Prepare plot data
		xArr := []float64{}
		yArr := []float64{}
		for _, grid := range resGrids {
			xArr = append(xArr, grid[0])
			yArr = append(yArr, grid[1])
		}

		data, err := populateData(xArr, yArr)
		if err != nil {
			log.Fatal(err.Error())
		}

		xArrMaxes = append(xArrMaxes, xArr...)

		// create line plots of wave ray
		filled, err := plotter.NewLine(data)
		if err != nil {
			log.Panic(err)
		}
		linePlots = append(linePlots, filled)

		currentAngle += angleIncrement
	}

	// Make plot instance
	p := plot.New()
	p.Title.Text = "Refraction Paths"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	// plot layers
	rnd := rand.New(rand.NewSource(1))

	layer := layerThicknesses[1]
	for i, thickness := range layerThicknesses {
		if i == 0 {
			pts := make(plotter.XYs, 2)
			pts[0].X = 0
			pts[0].Y = 0
			pts[1].X = getMax(xArrMaxes)
			pts[1].Y = 0
			layers, err := plotter.NewLine(pts)
			if err != nil {
				log.Panic(err)
			}

			layers.FillColor = color.RGBA{R: uint8(rnd.Intn(255)), G: uint8(rnd.Intn(255)), B: uint8(rnd.Intn(255)), A: uint8(rnd.Intn(255))}
			p.Add(layers)
		}

		pts := make(plotter.XYs, 2)
		pts[0].X = 0
		pts[0].Y = layer
		pts[1].X = getMax(xArrMaxes)
		pts[1].Y = layer
		layers, err := plotter.NewLine(pts)
		if err != nil {
			log.Panic(err)
		}

		layers.FillColor = color.RGBA{R: uint8(rnd.Intn(255)), G: uint8(rnd.Intn(255)), B: uint8(rnd.Intn(255)), A: uint8(rnd.Intn(255))}
		p.Add(layers)
		layer += thickness

	}

	// plot wave rays
	for _, item := range linePlots {
		p.Add(item)
	}

	// V. Plot to a graph
	err := p.Save(200, 200, "./refraction.png")
	if err != nil {
		log.Panic(err)
	}

}

func getMax(slice []float64) float64 {
	max := slice[0]

	// iterate through the slice and update max if a larger value is found
	for _, num := range slice {
		if num > max {
			max = num
		}
	}

	return max
}
