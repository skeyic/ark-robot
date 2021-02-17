package statistics

/* statistics/wvariance.go
 * 
 * Copyright (C) 1996, 1997, 1998, 1999, 2000 Jim Davies, Brian Gough
 * 
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or (at
 * your option) any later version.
 * 
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * General Public License for more details.
 * 
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.
 */

import (
	"math"
)

func wvariance(w, data Interface, wmean float64) (wvariance float64) {
	var W float64
	Len := data.Len()

	// the sum of the squares
	for i := 0; i < Len; i++ {
		wi := w.Value(i)

		if wi > 0 {
			delta := data.Value(i) - wmean
			W += wi
			wvariance += (delta*delta - wvariance) * (wi / W)
		}
	}

	return
}

func factor(w Interface) (factor float64) {
	var a, b float64
	Len := w.Len()

	// the sum of the squares
	for i := 0; i < Len; i++ {
		wi := w.Value(i)

		if wi > 0 {
			a += wi
			b += wi * wi
		}
	}

	factor = (a * a) / ((a * a) - b)

	return
}

func WVarianceWithFixedMean(w, data Interface, wmean float64) float64 {
	return wvariance(w, data, wmean)
}

func WsdWithFixedMean(w, data Interface, wmean float64) float64 {
	wvariance := wvariance(w, data, wmean)
	return math.Sqrt(wvariance)
}

func WVarianceMean(w, data Interface, wmean float64) float64 {
	variance := wvariance(w, data, wmean)
	scale := factor(w)

	return scale * variance
}

func WSdMean(w, data Interface, wmean float64) float64 {
	variance := wvariance(w, data, wmean)
	scale := factor(w)
	return math.Sqrt(scale * variance)
}

func WSd(w, data Interface) float64 {
	wmean := WMean(w, data)
	return WSdMean(w, data, wmean)
}

func WVariance(w, data Interface) float64 {
	wmean := WMean(w, data)
	return WVarianceMean(w, data, wmean)
}
