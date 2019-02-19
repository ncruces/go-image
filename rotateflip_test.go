package rotateflip

import (
	"image"
	"image/color"
	"image/color/palette"
	"math/rand"
	"testing"
)

func Test_Image(t *testing.T) {
	var subsample string
	rect := image.Rect(0, 0, 16, 16)

	init := func(pix []uint8) {
		for i := range pix {
			pix[i] = uint8(rand.Int63())
		}
	}

	testSub := func(img image.Image, op Operation) {
		rf1 := Image(img, op)
		rf2 := Image(&wrapper{img}, op)
		bounds := rf1.Bounds()

		if bounds != rf2.Bounds() {
			t.Errorf("%T%s/%d: bounds don't match", img, subsample, op)
		}
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				if rf1.At(x, y) != rf2.At(x, y) {
					t.Errorf("%T%s/%d: colors don't match at %2dx%d", img, subsample, op, x, y)
					return
				}
			}
		}
	}

	testImg := func(img imageWithSubImage) {
		for op := None; op <= Transverse; op++ {
			testSub(img, op)

			testSub(img.SubImage(image.Rect(0, 1, 16, 16)), op)
			testSub(img.SubImage(image.Rect(1, 0, 16, 16)), op)
			testSub(img.SubImage(image.Rect(1, 1, 16, 16)), op)

			testSub(img.SubImage(image.Rect(0, 0, 16, 15)), op)
			testSub(img.SubImage(image.Rect(0, 0, 15, 16)), op)
			testSub(img.SubImage(image.Rect(0, 0, 15, 15)), op)

			testSub(img.SubImage(image.Rect(2, 2, 14, 14)), op)
			testSub(img.SubImage(image.Rect(3, 3, 13, 13)), op)
		}
	}

	{
		img := image.NewAlpha(rect)
		init(img.Pix)
		testImg(img)
	}
	{
		img := image.NewAlpha16(rect)
		init(img.Pix)
		testImg(img)
	}
	{
		img := image.NewCMYK(rect)
		init(img.Pix)
		testImg(img)
	}
	{
		img := image.NewGray(rect)
		init(img.Pix)
		testImg(img)
	}
	{
		img := image.NewGray16(rect)
		init(img.Pix)
		testImg(img)
	}
	{
		img := image.NewNRGBA(rect)
		init(img.Pix)
		testImg(img)
	}
	{
		img := image.NewNRGBA64(rect)
		init(img.Pix)
		testImg(img)
	}
	{
		img := image.NewRGBA(rect)
		init(img.Pix)
		testImg(img)
	}
	{
		img := image.NewRGBA64(rect)
		init(img.Pix)
		testImg(img)
	}
	{
		img := image.NewPaletted(rect, palette.Plan9)
		init(img.Pix)
		testImg(img)
	}

	for sr := image.YCbCrSubsampleRatio444; sr <= image.YCbCrSubsampleRatio410; sr++ {
		subsample = "(" + sr.String() + ")"
		{
			img := image.NewYCbCr(rect, sr)
			init(img.Y)
			init(img.Cb)
			init(img.Cr)
			testImg(img)
		}
		{
			img := image.NewNYCbCrA(rect, sr)
			init(img.Y)
			init(img.Cb)
			init(img.Cr)
			init(img.A)
			testImg(img)
		}
	}
}

type wrapper struct {
	i image.Image
}

func (w *wrapper) ColorModel() color.Model {
	return w.i.ColorModel()
}

func (w *wrapper) Bounds() image.Rectangle {
	return w.i.Bounds()
}

func (w *wrapper) At(x, y int) color.Color {
	return w.i.At(x, y)
}

type imageWithSubImage interface {
	image.Image
	SubImage(image.Rectangle) image.Image
}
