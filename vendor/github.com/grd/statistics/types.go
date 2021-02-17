package statistics

/* statistics/types.go
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

//
// Interface is used throughout the package.
//
type Interface interface {
	Value(int) float64
	SetValue(int, float64)
	Len() int
	Less(int, int) bool
	Swap(int, int)
}

type Float64 []float64

func (f *Float64) Value(i int) float64         { return (*f)[i] }
func (f *Float64) SetValue(i int, val float64) { (*f)[i] = val }
func (f *Float64) Len() int                    { return len(*f) }
func (f *Float64) Less(i, j int) bool          { return (*f)[i] < (*f)[j] }
func (f *Float64) Swap(i, j int)               { (*f)[i], (*f)[j] = (*f)[j], (*f)[i] }

type Int64 []int64

func (f *Int64) Value(i int) float64           { return float64((*f)[i]) }
func (f *Int64) SetValue(i int, value float64) { (*f)[i] = int64(value) }
func (f *Int64) Len() int                      { return len(*f) }
func (f *Int64) Less(i, j int) bool            { return (*f)[i] < (*f)[j] }
func (f *Int64) Swap(i, j int)                 { (*f)[i], (*f)[j] = (*f)[j], (*f)[i] }

//
// Strider strides over the data, for sampling purposes.
//
type Strider struct {
	Interface
	stride int
}

func NewStrider(data Interface, stride int) *Strider {
	return &Strider{data, stride}
}

func (p *Strider) Value(i int) float64 {
	return p.Interface.Value(i * p.stride)
}

func (p *Strider) SetValue(i int, value float64) {
	p.Interface.SetValue(i*p.stride, value)
}

func (p *Strider) Len() int {
	return p.Interface.Len() / p.stride
}

func (p *Strider) Less(i, j int) bool {
	return p.Interface.Less(i*p.stride, j*p.stride)
}

func (p *Strider) Swap(i, j int) {
	p.Interface.Swap(i*p.stride, j*p.stride)
}
