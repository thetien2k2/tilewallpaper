package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
)

const (
	screenW = 3840
	screenH = 2160
)

func checkError(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	fn := flag.Arg(0)
	vips.LoggingSettings(nil, vips.LogLevelCritical)
	vips.Startup(nil)
	defer vips.Shutdown()
	wp, err := vips.Black(screenW, screenH)
	checkError(err)
	transColor := vips.ColorRGBA{R: 0, G: 0, B: 0, A: 0}
	img, err := vips.NewImageFromFile(fn)
	checkError(err)

	// resize image to screen heigh
	if img.Height() > screenH {
		img.Resize(float64(screenH)/float64(img.Height()), vips.KernelAuto)
	}
	// create gaussian blur background
	if img.Height() < screenH {
		imgbg, err := img.Copy()
		checkError(err)
		imgbg.ResizeWithVScale(1, float64(screenH)/float64(imgbg.Height()), vips.KernelAuto)
		imgbg.GaussianBlur(40)
		for x := screenW - imgbg.Width(); x >= -imgbg.Width(); x -= imgbg.Width() {
			wp.Insert(imgbg, x, (screenH-imgbg.Height())/2, true, &transColor)
		}
		wp.ExtractArea(wp.Width()-screenW, 0, screenW, screenH)
	}
	// tile image from right
	for x := screenW - img.Width(); x >= -img.Width(); x -= img.Width() {
		wp.Insert(img, x, (screenH-img.Height())/2, true, &transColor)
	}
	// crop to screen size
	wp.ExtractArea(wp.Width()-screenW, 0, screenW, screenH)

	ep := vips.NewDefaultPNGExportParams()
	ep.Quality = 100
	ep.Lossless = true
	epbytes, _, err := wp.Export(ep)
	checkError(err)
	err = os.WriteFile(fmt.Sprintf("%v.png", time.Now().UnixNano()), epbytes, 0644)
	checkError(err)
}
