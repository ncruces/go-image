package imageutil

import (
	"image"
	"math/rand"
	"testing"
)

func Test_YCbCrUpsample(t *testing.T) {
	var subsample string
	rect := image.Rect(0, 0, 16, 16)

	testSub := func(img image.Image) {
		var dst image.Image
		if src, ok := img.(*image.YCbCr); ok {
			dst = YCbCrUpsample(src)
		}
		if src, ok := img.(*image.NYCbCrA); ok {
			dst = NYCbCrAUpsample(src)
		}

		bounds := img.Bounds()
		if bounds != dst.Bounds() {
			t.Errorf("%T%s: bounds don't match", img, subsample)
		}
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				if img.At(x, y) != dst.At(x, y) {
					t.Errorf("%T%s: colors don't match at %2dx%d", img, subsample, x, y)
					t.Fatal(img, dst)
					return
				}
			}
		}
	}

	testImg := func(img imageWithSubImage) {
		testSub(img)

		testSub(img.SubImage(image.Rect(0, 1, 16, 16)))
		testSub(img.SubImage(image.Rect(1, 0, 16, 16)))
		testSub(img.SubImage(image.Rect(1, 1, 16, 16)))

		testSub(img.SubImage(image.Rect(0, 0, 16, 15)))
		testSub(img.SubImage(image.Rect(0, 0, 15, 16)))
		testSub(img.SubImage(image.Rect(0, 0, 15, 15)))

		testSub(img.SubImage(image.Rect(2, 2, 14, 14)))
		testSub(img.SubImage(image.Rect(3, 3, 13, 13)))
	}

	for sr := image.YCbCrSubsampleRatio444; sr <= image.YCbCrSubsampleRatio410; sr++ {
		subsample = "(" + sr.String() + ")"
		{
			img := image.NewYCbCr(rect, sr)
			random(img.Y)
			random(img.Cb)
			random(img.Cr)
			testImg(img)
		}
		{
			img := image.NewNYCbCrA(rect, sr)
			random(img.Y)
			random(img.Cb)
			random(img.Cr)
			random(img.A)
			testImg(img)
		}
	}
}

func random(pix []uint8) {
	for i := range pix {
		pix[i] = uint8(rand.Int63())
	}
}

type imageWithSubImage interface {
	image.Image
	SubImage(image.Rectangle) image.Image
}
