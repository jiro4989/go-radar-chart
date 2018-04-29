package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	x      = 0
	y      = 0
	width  = 255
	height = 255
)

var (
	whiteCol = color.RGBA{255, 255, 255, 255}
	blackCol = color.RGBA{80, 80, 80, 255}
	redCol   = color.RGBA{255, 0, 0, 100}
	r        = float64(width) * 0.8 * 0.5
	center   = image.Point{width / 2, width / 2}
)

func createWhiteImage() *image.RGBA {
	// 画像の生成
	img := image.NewRGBA(image.Rect(x, y, width, height))
	// 背景色の描画
	rect := img.Rect
	for h := rect.Min.Y; h < rect.Max.Y; h++ {
		for w := rect.Min.X; w < rect.Max.X; w++ {
			img.Set(w, h, whiteCol)
		}
	}
	return img
}

func addLabel(img *image.RGBA, x, y int, label string) {
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(blackCol),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func TestDraw(t *testing.T) {
	for rcnt := 3.0; rcnt <= 6; rcnt++ {
		img := createWhiteImage()

		// 背景色の塗りつぶし
		cx := float64(center.X)
		cy := float64(center.Y)
		for rr := 0.0; rr < r; rr += 0.1 {
			for rad := 0.0; rad < width*rr; rad++ {
				var x int = int(cx + math.Cos(rad)*rr)
				var y int = int(cy + math.Sin(rad)*rr)
				img.Set(x, y, redCol)
			}
		}

		// メモリ線の描画
		rad := 2 * math.Pi / rcnt
		drad := 2*math.Pi - rad - 2*math.Pi/4
		for j := 1.0; j <= rcnt; j++ {
			y1, x1 := math.Sincos(rad*j + drad)
			for i := 0.0; i < r; i += 0.1 {
				x2 := int(float64(center.X) + x1*i)
				y2 := int(float64(center.Y) + y1*i)
				img.Set(x2, y2, blackCol)
			}
			// ラベル描画
			x2 := int(float64(center.X) + x1*r)
			y2 := int(float64(center.Y) + y1*r)
			addLabel(img, int(x2), int(y2), fmt.Sprintf("%f", j))
		}

		// 円の枠線描画
		div := 5
		m := int(r) / div
		for rr := 0; rr < int(math.Ceil(r)); rr++ {
			if rr%m == 0 {
				rrf := float64(rr)
				for rad := 0.0; rad < float64(width)*rrf; rad += 0.1 {
					var x int = int(cx + math.Cos(rad)*rrf)
					var y int = int(cy + math.Sin(rad)*rrf)
					img.Set(x, y, blackCol)
				}
			}
		}

		// 画像ファイルの生成
		w, err := os.Create(fmt.Sprintf("img/line%d.png", int(rcnt)))
		if err != nil {
			log.Fatal(err)
		}
		defer w.Close()

		// 画像データの書き込み
		if err := png.Encode(w, img); err != nil {
			log.Fatal(err)
		}
	}
}

type data struct {
	label string
	value float64
}

func TestDrawChart(t *testing.T) {
	datas := []data{
		data{label: "ATK", value: 100.0},
		data{label: "DEF", value: 60.0},
		data{label: "MAT", value: 5.0},
	}
	rcnt := float64(len(datas))

	img := createWhiteImage()

	// 背景色の塗りつぶし
	cx := float64(center.X)
	cy := float64(center.Y)

	// メモリ線の描画
	rad := 2 * math.Pi / rcnt
	drad := 2*math.Pi - rad - 2*math.Pi/4
	for j, v := range datas {
		y1, x1 := math.Sincos(rad*float64(j) + drad)
		for i := 0.0; i < r; i += 0.1 {
			x2 := int(float64(center.X) + x1*i)
			y2 := int(float64(center.Y) + y1*i)
			img.Set(x2, y2, blackCol)
		}
		// ラベル描画
		x2 := int(float64(center.X) + x1*r)
		y2 := int(float64(center.Y) + y1*r)
		addLabel(img, int(x2), int(y2), fmt.Sprintf("%s", v.label))
	}

	// 円の枠線描画
	div := 5
	m := int(r) / div
	for rr := 0; rr < int(math.Ceil(r)); rr++ {
		if rr%m == 0 {
			rrf := float64(rr)
			for rad := 0.0; rad < float64(width)*rrf; rad += 0.1 {
				var x int = int(cx + math.Cos(rad)*rrf)
				var y int = int(cy + math.Sin(rad)*rrf)
				img.Set(x, y, blackCol)
			}
		}
	}

	// 画像ファイルの生成
	w, err := os.Create(fmt.Sprintf("img/chart%d.png", int(rcnt)))
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	// 画像データの書き込み
	if err := png.Encode(w, img); err != nil {
		log.Fatal(err)
	}
}

func TestDrawText(t *testing.T) {
	img := createWhiteImage()
	addLabel(img, 20, 30, "Hello Go")

	f, err := os.Create("img/text.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func TestFillBackground(t *testing.T) {
	img := createWhiteImage()

	// 背景色の塗りつぶし
	for rr := 0.0; rr < r; rr += 0.1 {
		for rad := 0.0; rad < width*rr; rad++ {
			x := int(float64(center.X) + math.Cos(rad)*rr)
			y := int(float64(center.Y) + math.Sin(rad)*rr)
			img.Set(x, y, redCol)
		}
	}

	f, err := os.Create("img/fill_background.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func TestMath(t *testing.T) {
	vs := []map[string]float64{
		map[string]float64{"r": 1.0, "div": 1.0},
		map[string]float64{"r": 2.0, "div": 1.0},
		map[string]float64{"r": 50.0, "div": 1.0},
		map[string]float64{"r": 1.0, "div": 4.0},
		map[string]float64{"r": 2.0, "div": 4.0},
		map[string]float64{"r": 50.0, "div": 4.0},
	}
	for _, v := range vs {
		r := v["r"]
		div := v["div"]

		area := math.Pi * r * r / div
		arc := 2 * math.Pi * r / div
		theta := arc / r
		fmt.Println(fmt.Sprintf("半径:%f,面積:%f,弧の長さ:%f,θ:%f", r, area, arc, theta))
	}
}
