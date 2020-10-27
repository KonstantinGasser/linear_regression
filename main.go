package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	xys, err := readFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}

	err = plotData("out.png", xys)
	if err != nil {
		log.Fatal(err)
	}

}

func plotData(out string, xys []XY) error {

	p, err := plot.New()
	if err != nil {
		return fmt.Errorf("Could not create plot: %v", err)
	}

	pxys := make(plotter.XYs, len(xys))
	for i, xy := range xys {
		pxys[i].X = xy.X
		pxys[i].Y = xy.Y
	}
	s, err := plotter.NewScatter(pxys)
	if err != nil {
		return err
	}
	s.GlyphStyle.Shape = draw.CrossGlyph{}
	s.Color = color.RGBA{B: 255, A: 255}
	p.Add(s)
	wt, err := p.WriterTo(521, 521, "png")
	if err != nil {
		return fmt.Errorf("Could not create writer: %v", err)
	}

	f, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("Could not create file: %v", err)
	}
	defer f.Close()
	_, err = wt.WriteTo(f)
	if err != nil {
		return fmt.Errorf("Could not write file: %v", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("Could not close file: %v", err)
	}
	return nil
}

type XY struct {
	X, Y float64
}

func readFile(path string) ([]XY, error) {

	var xys []XY

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		var x, y float64
		_, err := fmt.Sscanf(s.Text(), "%f,%f", &x, &y)
		if err != nil {
			log.Printf("discarding bad data point %q: %v", s.Text(), err)
		}

		if err := s.Err(); err != nil {
			return nil, err
		}
		xys = append(xys, XY{x, y})
	}
	return xys, nil
}
