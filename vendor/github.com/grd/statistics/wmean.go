package statistics

/* statistics/wmean.go
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

// WMean calculates the weighted arithmetic mean of a dataset
func WMean(w, data Interface) (wmean float64) {
	var W float64
	Len := data.Len()

	for i := 0; i < Len; i++ {
		wi := w.Value(i)

		if wi > 0 {
			W += wi
			wmean += (data.Value(i) - wmean) * (wi / W)
		}
	}

	return
}
