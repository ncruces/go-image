// Package rotateflip rotates, flips or rotates and flips images.
//
// The package works with the Image interface described in the image package.
//
// A fast path is used for most of the in-memory image types defined in that package.
// An image of the same type is returned.
//
// A lazy, slow path, is used for other image types,
// as well as for some subsampled YCbCr images.
//
// Example:
//    exf := rotateflip.Orientation(exifOrientation)
//    img := rotateflip.Image(srcImage, exf.Op())
package rotateflip

import (
	"image"
	"image/color"
)

// Operation specifies a clockwise rotation and flip operation to apply to an image.
type Operation int

const (
	None Operation = iota
	Rotate90
	Rotate180
	Rotate270

	FlipX
	Transpose
	FlipY
	Transverse

	FlipXY = Rotate180

	Rotate90FlipX   = Transpose
	Rotate180FlipX  = FlipY
	Rotate270FlipX  = Transverse
	Rotate90FlipY   = Transverse
	Rotate180FlipY  = FlipX
	Rotate270FlipY  = Transpose
	Rotate90FlipXY  = Rotate270
	Rotate180FlipXY = None
	Rotate270FlipXY = Rotate90
)

// Orientation is an image orientation as specified by EXIF 2.2 and TIFF 6.0.
type Orientation int

const (
	TopLeft Orientation = iota + 1
	TopRight
	BottomRight
	BottomLeft
	LeftTop
	RightTop
	RightBottom
	LeftBottom
)

// Op gets the Operation that restores an image with this Orientation to TopLeft Orientation.
func (or Orientation) Op() Operation {
	switch or {
	default:
		return None
	case TopRight:
		return FlipX
	case BottomRight:
		return FlipXY
	case BottomLeft:
		return FlipY
	case LeftTop:
		return Transpose
	case RightTop:
		return Rotate90
	case RightBottom:
		return Transverse
	case LeftBottom:
		return Rotate270
	}
}

// Image applies an Operation to an image.
func Image(src image.Image, op Operation) image.Image {
	op &= 7 // sanitize

	if op == 0 {
		return src // nop
	}

	bounds := rotateBounds(src.Bounds(), op)

	// fast path, eager
	switch src := src.(type) {
	case *image.Alpha:
		dst := image.NewAlpha(bounds)
		rotateFlip(dst.Pix, dst.Stride, dst.Bounds().Dx(), dst.Bounds().Dy(), src.Pix, src.Stride, src.Bounds().Dx(), src.Bounds().Dy(), op, 1)
		return dst

	case *image.Alpha16:
		dst := image.NewAlpha16(bounds)
		rotateFlip(dst.Pix, dst.Stride, dst.Bounds().Dx(), dst.Bounds().Dy(), src.Pix, src.Stride, src.Bounds().Dx(), src.Bounds().Dy(), op, 2)
		return dst

	case *image.CMYK:
		dst := image.NewCMYK(bounds)
		rotateFlip(dst.Pix, dst.Stride, dst.Bounds().Dx(), dst.Bounds().Dy(), src.Pix, src.Stride, src.Bounds().Dx(), src.Bounds().Dy(), op, 4)
		return dst

	case *image.Gray:
		dst := image.NewGray(bounds)
		rotateFlip(dst.Pix, dst.Stride, dst.Bounds().Dx(), dst.Bounds().Dy(), src.Pix, src.Stride, src.Bounds().Dx(), src.Bounds().Dy(), op, 1)
		return dst

	case *image.Gray16:
		dst := image.NewGray16(bounds)
		rotateFlip(dst.Pix, dst.Stride, dst.Bounds().Dx(), dst.Bounds().Dy(), src.Pix, src.Stride, src.Bounds().Dx(), src.Bounds().Dy(), op, 2)
		return dst

	case *image.NRGBA:
		dst := image.NewNRGBA(bounds)
		rotateFlip(dst.Pix, dst.Stride, dst.Bounds().Dx(), dst.Bounds().Dy(), src.Pix, src.Stride, src.Bounds().Dx(), src.Bounds().Dy(), op, 4)
		return dst

	case *image.NRGBA64:
		dst := image.NewNRGBA64(bounds)
		rotateFlip(dst.Pix, dst.Stride, dst.Bounds().Dx(), dst.Bounds().Dy(), src.Pix, src.Stride, src.Bounds().Dx(), src.Bounds().Dy(), op, 8)
		return dst

	case *image.RGBA:
		dst := image.NewRGBA(bounds)
		rotateFlip(dst.Pix, dst.Stride, dst.Bounds().Dx(), dst.Bounds().Dy(), src.Pix, src.Stride, src.Bounds().Dx(), src.Bounds().Dy(), op, 4)
		return dst

	case *image.RGBA64:
		dst := image.NewRGBA64(bounds)
		rotateFlip(dst.Pix, dst.Stride, dst.Bounds().Dx(), dst.Bounds().Dy(), src.Pix, src.Stride, src.Bounds().Dx(), src.Bounds().Dy(), op, 8)
		return dst

	case *image.Paletted:
		dst := image.NewPaletted(bounds, src.Palette)
		rotateFlip(dst.Pix, dst.Stride, dst.Bounds().Dx(), dst.Bounds().Dy(), src.Pix, src.Stride, src.Bounds().Dx(), src.Bounds().Dy(), op, 1)
		return dst

	case *image.YCbCr:
		if sr, ok := rotateYCbCrSubsampleRatio(src.SubsampleRatio, src.Bounds(), op); ok {
			dst := image.NewYCbCr(bounds, sr)
			srcCBounds := subsampledBounds(src.Bounds(), src.SubsampleRatio)
			dstCBounds := subsampledBounds(dst.Bounds(), dst.SubsampleRatio)
			rotateFlip(dst.Y, dst.YStride, dst.Bounds().Dx(), dst.Bounds().Dy(), src.Y, src.YStride, src.Bounds().Dx(), src.Bounds().Dy(), op, 1)
			rotateFlip(dst.Cb, dst.CStride, dstCBounds.Dx(), dstCBounds.Dy(), src.Cb, src.CStride, srcCBounds.Dx(), srcCBounds.Dy(), op, 1)
			rotateFlip(dst.Cr, dst.CStride, dstCBounds.Dx(), dstCBounds.Dy(), src.Cr, src.CStride, srcCBounds.Dx(), srcCBounds.Dy(), op, 1)
			return dst
		}

	case *image.NYCbCrA:
		if sr, ok := rotateYCbCrSubsampleRatio(src.SubsampleRatio, src.Bounds(), op); ok {
			dst := image.NewNYCbCrA(bounds, sr)
			srcCBounds := subsampledBounds(src.Bounds(), src.SubsampleRatio)
			dstCBounds := subsampledBounds(dst.Bounds(), dst.SubsampleRatio)
			rotateFlip(dst.Y, dst.YStride, dst.Bounds().Dx(), dst.Bounds().Dy(), src.Y, src.YStride, src.Bounds().Dx(), src.Bounds().Dy(), op, 1)
			rotateFlip(dst.A, dst.AStride, dst.Bounds().Dx(), dst.Bounds().Dy(), src.A, src.AStride, src.Bounds().Dx(), src.Bounds().Dy(), op, 1)
			rotateFlip(dst.Cb, dst.CStride, dstCBounds.Dx(), dstCBounds.Dy(), src.Cb, src.CStride, srcCBounds.Dx(), srcCBounds.Dy(), op, 1)
			rotateFlip(dst.Cr, dst.CStride, dstCBounds.Dx(), dstCBounds.Dy(), src.Cr, src.CStride, srcCBounds.Dx(), srcCBounds.Dy(), op, 1)
			return dst
		}
	}

	// slow path, lazy
	return &rotateFlipImage{src, op}
}

type rotateFlipImage struct {
	src image.Image
	op  Operation
}

func (rft *rotateFlipImage) ColorModel() color.Model {
	return rft.src.ColorModel()
}

func (rft *rotateFlipImage) Bounds() image.Rectangle {
	return rotateBounds(rft.src.Bounds(), rft.op)
}

func (rft *rotateFlipImage) At(x, y int) color.Color {
	bounds := rft.src.Bounds()
	switch rft.op {
	default:
		return rft.src.At(bounds.Min.X+x, bounds.Min.Y+y)
	case FlipX:
		return rft.src.At(bounds.Max.X-x-1, bounds.Min.Y+y)
	case FlipXY:
		return rft.src.At(bounds.Max.X-x-1, bounds.Max.Y-y-1)
	case FlipY:
		return rft.src.At(bounds.Min.X+x, bounds.Max.Y-y-1)
	case Transpose:
		return rft.src.At(bounds.Min.X+y, bounds.Min.Y+x)
	case Rotate90:
		return rft.src.At(bounds.Min.X+y, bounds.Max.Y-x-1)
	case Transverse:
		return rft.src.At(bounds.Max.X-y-1, bounds.Max.Y-x-1)
	case Rotate270:
		return rft.src.At(bounds.Max.X-y-1, bounds.Min.Y+x)
	}
}

func rotateFlip(dst []uint8, dst_stride, dst_width, dst_height int, src []uint8, src_stride, src_width, src_height int, op Operation, bpp int) {
	rotate := op&1 != 0
	flip_y := op&2 != 0
	flip_x := parity(op)

	var dst_row, src_row int

	if flip_x {
		dst_row += bpp * (dst_width - 1)
	}
	if flip_y {
		dst_row += dst_stride * (dst_height - 1)
	}

	var dst_x_offset, dst_y_offset int

	if rotate {
		if flip_x {
			dst_y_offset = -bpp
		} else {
			dst_y_offset = +bpp
		}
		if flip_y {
			dst_x_offset = -dst_stride
		} else {
			dst_x_offset = +dst_stride
		}
	} else {
		if flip_x {
			dst_x_offset = -bpp
		} else {
			dst_x_offset = +bpp
		}
		if flip_y {
			dst_y_offset = -dst_stride
		} else {
			dst_y_offset = +dst_stride
		}
	}

	if dst_x_offset == bpp {
		for y := 0; y < src_height; y++ {
			copy(dst[dst_row:], src[src_row:src_row+src_width*bpp])
			dst_row += dst_y_offset
			src_row += src_stride
		}
	} else {
		for y := 0; y < src_height; y++ {
			dst_pixel := dst_row
			src_pixel := src_row

			for x := 0; x < src_width; x++ {
				copy(dst[dst_pixel:], src[src_pixel:src_pixel+bpp])
				dst_pixel += dst_x_offset
				src_pixel += bpp
			}

			dst_row += dst_y_offset
			src_row += src_stride
		}
	}
}

func rotateBounds(bounds image.Rectangle, op Operation) image.Rectangle {
	var dx, dy int
	if rotate := op&1 != 0; rotate {
		dx = bounds.Dy()
		dy = bounds.Dx()
	} else {
		dx = bounds.Dx()
		dy = bounds.Dy()
	}
	return image.Rectangle{image.ZP, image.Point{dx, dy}}

}

func rotateYCbCrSubsampleRatio(subsampleRatio image.YCbCrSubsampleRatio, bounds image.Rectangle, op Operation) (image.YCbCrSubsampleRatio, bool) {
	switch subsampleRatio {
	case image.YCbCrSubsampleRatio444:
		return subsampleRatio, true

	case image.YCbCrSubsampleRatio420:
		if (bounds.Min.X|bounds.Max.X|bounds.Min.Y|bounds.Max.Y)&1 != 0 {
			break
		}
		return subsampleRatio, true

	case image.YCbCrSubsampleRatio422:
		if (bounds.Min.X|bounds.Max.X)&1 != 0 {
			break
		}
		if rotate := op&1 != 0; rotate {
			return image.YCbCrSubsampleRatio440, true
		}
		return subsampleRatio, true

	case image.YCbCrSubsampleRatio440:
		if (bounds.Min.Y|bounds.Max.Y)&1 != 0 {
			break
		}
		if rotate := op&1 != 0; rotate {
			return image.YCbCrSubsampleRatio422, true
		}
		return subsampleRatio, true

	case image.YCbCrSubsampleRatio411, image.YCbCrSubsampleRatio410:
		if (bounds.Min.X|bounds.Max.X)&3 != 0 {
			break
		}
		if (bounds.Min.Y|bounds.Max.Y)&1 != 0 {
			break
		}
		if rotate := op&1 != 0; rotate {
			break
		}
		return subsampleRatio, true
	}

	return -1, false
}

func subsampledBounds(bounds image.Rectangle, subsampleRatio image.YCbCrSubsampleRatio) image.Rectangle {
	switch subsampleRatio {
	default:
		panic("Unknown YCbCrSubsampleRatio")
	case image.YCbCrSubsampleRatio440:
		bounds.Min.Y = (bounds.Min.Y + 0) / 2
		bounds.Max.Y = (bounds.Max.Y + 1) / 2
		fallthrough
	case image.YCbCrSubsampleRatio444:

	case image.YCbCrSubsampleRatio420:
		bounds.Min.Y = (bounds.Min.Y + 0) / 2
		bounds.Max.Y = (bounds.Max.Y + 1) / 2
		fallthrough
	case image.YCbCrSubsampleRatio422:
		bounds.Min.X = (bounds.Min.X + 0) / 2
		bounds.Max.X = (bounds.Max.X + 1) / 2

	case image.YCbCrSubsampleRatio410:
		bounds.Min.Y = (bounds.Min.Y + 0) / 2
		bounds.Max.Y = (bounds.Max.Y + 1) / 2
		fallthrough
	case image.YCbCrSubsampleRatio411:
		bounds.Min.X = (bounds.Min.X + 0) / 4
		bounds.Max.X = (bounds.Max.X + 3) / 4
	}
	return bounds
}

func parity(op Operation) bool {
	op = 0226 >> uint8(op)
	return op&1 != 0
}
