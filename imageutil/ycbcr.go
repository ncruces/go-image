// Package imageutil contains code shared by image-related packages.
package imageutil

import (
	"image"
)

// YCbCrUpsample upsamples a chroma subsampled YCbCr image.
// The returned image has YCbCrSubsampleRatio444.
func YCbCrUpsample(img *image.YCbCr) *image.YCbCr {
	if img.SubsampleRatio == image.YCbCrSubsampleRatio444 {
		return img
	}

	dst := image.NewYCbCr(img.Rect, image.YCbCrSubsampleRatio444)
	resample(dst.Y, dst.YStride, img.Y, img.YStride, img.Rect.Dy())
	upsample(img, dst)
	return dst
}

// NYCbCrAUpsample upsamples a chroma subsampled NYCbCrA image.
// The returned image has YCbCrSubsampleRatio444.
func NYCbCrAUpsample(img *image.NYCbCrA) *image.NYCbCrA {
	if img.SubsampleRatio == image.YCbCrSubsampleRatio444 {
		return img
	}

	dst := image.NewNYCbCrA(img.Rect, image.YCbCrSubsampleRatio444)
	resample(dst.Y, dst.YStride, img.Y, img.YStride, img.Rect.Dy())
	resample(dst.A, dst.AStride, img.A, img.AStride, img.Rect.Dy())
	upsample(&img.YCbCr, &dst.YCbCr)
	return dst
}

func resample(dst []uint8, dstStride int, src []uint8, srcStride int, count int) {
	var srcOffset, dstOffset int
	for i := 0; i < count; i++ {
		copy(dst[dstOffset:dstOffset+dstStride], src[srcOffset:srcOffset+dstStride])
		dstOffset += dstStride
		srcOffset += srcStride
	}
}

func upsample(src, dst *image.YCbCr) {
	sx, sy := subsampleRatios(src.SubsampleRatio)

	var i int
	for y := src.Rect.Min.Y; y < src.Rect.Max.Y; y++ {
		cy := (y/sy - src.Rect.Min.Y/sy) * src.CStride
		for x := src.Rect.Min.X; x < src.Rect.Max.X; x++ {
			j := cy + (x/sx - src.Rect.Min.X/sx)
			dst.Cb[i] = src.Cb[j]
			dst.Cr[i] = src.Cr[j]
			i++
		}
	}
}

func subsampleRatios(subsampleRatio image.YCbCrSubsampleRatio) (sx, sy int) {
	switch subsampleRatio {
	case image.YCbCrSubsampleRatio444:
		return 1, 1
	case image.YCbCrSubsampleRatio422:
		return 2, 1
	case image.YCbCrSubsampleRatio420:
		return 2, 2
	case image.YCbCrSubsampleRatio440:
		return 1, 2
	case image.YCbCrSubsampleRatio411:
		return 4, 1
	case image.YCbCrSubsampleRatio410:
		return 4, 2
	}
	panic("Unknown YCbCrSubsampleRatio")
}
