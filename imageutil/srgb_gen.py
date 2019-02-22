import numpy as np  # https://www.numpy.org/
import pwlf         # https://github.com/cjekel/piecewise_linear_fit_py
import sys

def srgb(lin):
    if lin <= 0.0031308:
        return lin * 12.92
    return 1.055 * lin**(1.0/2.4) - 0.055

x = np.array(range(0, 65536, 5))        # data points
x0 = np.array(range(0, 65536, 257))     # break points

y = []
for i in x:
    y.append(srgb(i/65535.0)*65535)
y = np.array(y)

my_pwlf = pwlf.PiecewiseLinFit(x, y)

# force both end-points and first two points, where absolute error is high
my_pwlf.fit_with_breaks_force_points(x0, [0,257,514,65535], [0,3324,5625,65535])

for i, yy in enumerate(my_pwlf.predict(x0)):
    sys.stdout.write("0x%04x, " % round(yy))
    if i % 8 == 7: print
