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
	ep := vips.NewDefaultPNGExportParams()
	ep.Quality = 100
	ep.Lossless = true

	transColor := vips.ColorRGBA{R: 0, G: 0, B: 0, A: 0}
	img, err := vips.NewImageFromFile(fn)
	checkError(err)

	imgbg, err := img.Copy()
	checkError(err)
	vs := 1.0
	if img.Width() > screenW {
		vs = float64(img.Width()) / float64(screenW)
	} else {
		vs = float64(screenW) / float64(img.Width())
	}
	imgbg.ResizeWithVScale(vs, float64(screenH)/float64(imgbg.Height()), vips.KernelAuto)

	if img.Height() < screenH {
		imgbg.GaussianBlur(100)
	}

	// resize image to screen heigh
	if img.Height() > screenH {
		err = img.Resize(float64(screenH)/float64(img.Height()), vips.KernelAuto)
		checkError(err)
	}

	// tile image from right
	for x := screenW - img.Width(); x >= -img.Width(); x -= img.Width() {
		err = imgbg.Insert(img, x, (screenH-img.Height())/2, false, &transColor)
		checkError(err)
	}

	// crop to screen size
	imgbg.ExtractArea(imgbg.Width()-screenW, 0, screenW, screenH)

	epbytes, _, err := imgbg.Export(ep)
	checkError(err)
	err = os.WriteFile(fmt.Sprintf("%v.png", time.Now().UnixNano()), epbytes, 0644)
	checkError(err)
}
