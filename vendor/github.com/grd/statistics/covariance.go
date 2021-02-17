package statistics

/* statistics/covariance.go
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

// takes a dataset and calculates the covariance
func covariance(data1, data2 Interface, mean1, mean2 float64) (ret float64) {
	// calculate the sum of the squares
	for i := 0; i < data1.Len(); i++ {
		delta1 := (data1.Value(i) - mean1)
		delta2 := (data2.Value(i) - mean2)
		ret += (delta1*delta2 - ret) / float64(i+1)
	}
	return
}

func CovarianceMean(data1, data2 Interface, mean1, mean2 float64) float64 {
	n := data1.Len()
	covariance := covariance(data1, data2, mean1, mean2)
	return covariance * float64(n) / float64(n-1)
}

func Covariance(data1, data2 Interface) float64 {
	mean1 := Mean(data1)
	mean2 := Mean(data2)

	return CovarianceMean(data1, data2, mean1, mean2)
}
