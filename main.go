package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/color"
	"io"
	"log"
	"os"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	epochs := flag.Int("epochs", 10000, "number of epochs for gradient decent")
	alpha := flag.Float64("alpha", 0.001, "learning rate")
	flag.Parse()

	f, err := os.Open("data.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	X, Y, err := readData(f)
	if err != nil {
		fmt.Println(err)
	}

	_, _ = performLinreg(X, Y, *epochs, *alpha)

	// out, err := os.Create("out.png")
	// if err != nil {
	// 	fmt.Printf("failed to create out.png: %v\n", err)
	// }
	//
	// if err := plotData(out, X, Y, m, c, 1, 30); err != nil {
	// 	fmt.Println(err)
	// }
}

func readData(data io.Reader) (xs, ys []float64, err error) {
	s := bufio.NewScanner(data)
	for s.Scan() {
		var x, y float64
		_, err := fmt.Sscanf(s.Text(), "%f,%f", &x, &y)
		if err != nil {
			log.Printf("discarding bad data point %q: %v", s.Text(), err)
			continue
		}
		xs = append(xs, x)
		ys = append(ys, y)
	}
	if err := s.Err(); err != nil {
		return nil, nil, fmt.Errorf("could not scan: %v", err)
	}
	return xs, ys, nil
}

func performLinreg(X, Y []float64, epochs int, alpha float64) (m, c float64) {

	for i := 0; i < epochs; i++ {
		cost, dm, dc := gradient(X, Y, m, c)
		fmt.Printf("cost(%f,%f) = %f\n", m, c, cost)
		m += -dm * alpha
		c += -dc * alpha

		// --Visualise the line after each epoch-- //
		time.Sleep(time.Millisecond * 50)
		out, err := os.Create("out.png")
		if err != nil {
			fmt.Printf("failed to create out.png: %v\n", err)
		}

		if err := plotData(out, X, Y, m, c, 1, 30); err != nil {
			fmt.Println(err)
		}
		// ---------------------------------------- //
	}
	return m, c
}

func gradient(X, Y []float64, m, c float64) (cost, dm, dc float64) {

	n := float64(len(X))
	for i := 0; i < len(X); i++ {
		partial := Y[i] - (X[i]*m + c)
		dm += -X[i] * partial
		dc += -partial
		cost += partial * partial
	}

	return cost / n, 2 / n * dm, 2 / n * dc
}

type xyer struct{ xs, ys []float64 }

func (x xyer) Len() int                    { return len(x.xs) }
func (x xyer) XY(i int) (float64, float64) { return x.xs[i], x.ys[i] }

func plotData(out io.Writer, X, Y []float64, m, c float64, min, max float64) error {
	p, err := plot.New()
	if err != nil {
		return fmt.Errorf("failed to create plot: %v", err)
	}

	// create scatter with all data points
	s, err := plotter.NewScatter(xyer{X, Y})
	if err != nil {
		return fmt.Errorf("failed to create scatter: %v", err)
	}
	s.GlyphStyle.Shape = draw.CrossGlyph{}
	s.Color = color.RGBA{R: 255, A: 255}
	p.Add(s)

	l, err := plotter.NewLine(plotter.XYs{
		{min, min*m + c}, {max, max*m + c},
	})
	if err != nil {
		return fmt.Errorf("failed to create line: %v", err)
	}
	p.Add(l)

	wt, err := p.WriterTo(256, 256, "png")
	if err != nil {
		return fmt.Errorf("failed to create writer: %v", err)
	}

	_, err = wt.WriteTo(out)
	if err != nil {
		return fmt.Errorf("failed to write: %v", err)
	}

	return nil
}
