package imageutil

import (
	"math"
	"testing"
)

func srgbToLinear(srgb float64) float64 {
	if srgb <= 0.04045 {
		return srgb / 12.92
	}
	return math.Pow((srgb+0.055)/1.055, 2.4)
}

func linearToSRGB(lin float64) float64 {
	if lin <= 0.0031308 {
		return lin * 12.92
	}
	return 1.055*math.Pow(lin, 1.0/2.4) - 0.055
}

func TestSRGBToLinear8(t *testing.T) {
	for i := 0; i < 256; i++ {
		exp := uint16(math.RoundToEven(srgbToLinear(float64(i)/255) * 65535))
		res := SRGB8ToLinear(uint8(i))
		if exp != res {
			t.Errorf("at: %d, expected: %d, got: %d", i, exp, res)
		}
	}
}

func TestSRGBToLinear16(t *testing.T) {
	var cnt, abs, sum int
	var prv uint16
	for i := 0; i < 65536; i++ {
		exp := uint16(math.RoundToEven(srgbToLinear(float64(i)/65535) * 65535))
		res := SRGB16ToLinear(uint16(i))
		err := int(res) - int(exp)
		if prv > res {
			t.Errorf("at %d, non-monotonic", i)
		}
		if err < -1 || err > +1 {
			t.Errorf("at: %d, expected: %d, got: %d", i, exp, res)
		}
		switch {
		case err < 0:
			abs -= err
		case err > 0:
			abs += err
		default:
			cnt++
		}
		sum += err
		prv = res
	}
	t.Logf("correct %d/65536, abs error: %d, error bias: %d", cnt, abs, sum)
}

func TestLinearToSRGB8(t *testing.T) {
	var cnt, abs, sum int
	var prv uint8
	for i := 0; i < 65536; i++ {
		exp := uint16(math.RoundToEven(linearToSRGB(float64(i)/65535) * 255))
		res := LinearToSRGB8(uint16(i))
		err := int(res) - int(exp)
		if prv > res {
			t.Errorf("at %d, non-monotonic", i)
		}
		if err < -1 || err > +1 {
			t.Errorf("at: %d, expected: %d, got: %d", i, exp, res)
		}
		switch {
		case err < 0:
			abs -= err
		case err > 0:
			abs += err
		default:
			cnt++
		}
		sum += err
		prv = res
	}
	t.Logf("correct %d/65536, abs error: %d, error bias: %d", cnt, abs, sum)
}

func TestLinearToSRGB16(t *testing.T) {
	var cnt, abs, sum int
	var prv uint16
	for i := 0; i < 65536; i++ {
		exp := uint16(math.RoundToEven(linearToSRGB(float64(i)/65535) * 65535))
		res := LinearToSRGB16(uint16(i))
		err := int(res) - int(exp)
		if prv > res {
			t.Errorf("at %d, non-monotonic", i)
		}
		if i < 8192 {
			if err < -58 || err > +58 {
				t.Errorf("at: %d, expected: %d, got: %d", i, exp, res)
			}
		} else {
			if err < -1 || err > +1 {
				t.Errorf("at: %d, expected: %d, got: %d", i, exp, res)
			}
		}
		switch {
		case err < 0:
			abs -= err
		case err > 0:
			abs += err
		default:
			cnt++
		}
		sum += err
		prv = res
	}
	t.Logf("correct %d/65536, abs error: %d, error bias: %d", cnt, abs, sum)
}

func TestReverseSRGB8(t *testing.T) {
	for i := 0; i < 256; i++ {
		exp := uint8(i)
		res := LinearToSRGB8(SRGB8ToLinear(uint8(i)))
		if exp != res {
			t.Errorf("at: %d, expected: %d, got: %d", i, exp, res)
		}
	}
}

func TestReverseSRGB16(t *testing.T) {
	for i := 0; i < 65536; i++ {
		exp := uint16(i)
		res := SRGB16ToLinear(LinearToSRGB16(uint16(i)))
		err := int(res) - int(exp)
		if err < -8 || err > +8 {
			t.Errorf("at: %d, expected: %d, got: %d, error: %d", i, exp, res, err)
		}
	}
}

func BenchmarkBaseline(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < 65536; i++ {
			math.RoundToEven(linearToSRGB(float64(i)/65535) * 255)
		}
	}
}

func BenchmarkFast(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < 65536; i++ {
			LinearToSRGB8(uint16(i))
		}
	}
}
