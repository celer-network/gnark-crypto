// Copyright 2020 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by gurvy/internal/generators DO NOT EDIT

package bn256

// E12 elmt in degree 12 extension
type E12 struct {
	C0, C1 E6
}

// Equal returns true if z equals x, fasle otherwise
func (z *E12) Equal(x *E12) bool {
	return z.C0.Equal(&x.C0) && z.C1.Equal(&x.C1)
}

// String puts E12 in string form
func (z *E12) String() string {
	return (z.C0.String() + "+(" + z.C1.String() + ")*w")
}

// SetString sets a E12 from string
func (z *E12) SetString(s0, s1, s2, s3, s4, s5, s6, s7, s8, s9, s10, s11 string) *E12 {
	z.C0.SetString(s0, s1, s2, s3, s4, s5)
	z.C1.SetString(s6, s7, s8, s9, s10, s11)
	return z
}

// Set copies x into z and returns z
func (z *E12) Set(x *E12) *E12 {
	z.C0 = x.C0
	z.C1 = x.C1
	return z
}

// SetOne sets z to 1 in Montgomery form and returns z
func (z *E12) SetOne() *E12 {
	z.C0.B0.A0.SetOne()
	z.C0.B0.A1.SetZero()
	z.C0.B1.A0.SetZero()
	z.C0.B1.A1.SetZero()
	z.C0.B2.A0.SetZero()
	z.C0.B2.A1.SetZero()
	z.C1.B0.A0.SetZero()
	z.C1.B0.A1.SetZero()
	z.C1.B1.A0.SetZero()
	z.C1.B1.A1.SetZero()
	z.C1.B2.A0.SetZero()
	z.C1.B2.A1.SetZero()
	return z
}

// ToMont converts to Mont form
// TODO can this be deleted?
func (z *E12) ToMont() *E12 {
	z.C0.ToMont()
	z.C1.ToMont()
	return z
}

// FromMont converts from Mont form
// TODO can this be deleted?
func (z *E12) FromMont() *E12 {
	z.C0.FromMont()
	z.C1.FromMont()
	return z
}

// Add set z=x+y in E12 and return z
func (z *E12) Add(x, y *E12) *E12 {
	z.C0.Add(&x.C0, &y.C0)
	z.C1.Add(&x.C1, &y.C1)
	return z
}

// Sub set z=x-y in E12 and return z
func (z *E12) Sub(x, y *E12) *E12 {
	z.C0.Sub(&x.C0, &y.C0)
	z.C1.Sub(&x.C1, &y.C1)
	return z
}

// Double sets z=2*x and returns z
func (z *E12) Double(x *E12) *E12 {
	z.C0.Double(&x.C0)
	z.C1.Double(&x.C1)
	return z
}

// SetRandom used only in tests
// TODO eliminate this method!
func (z *E12) SetRandom() *E12 {
	z.C0.B0.A0.SetRandom()
	z.C0.B0.A1.SetRandom()
	z.C0.B1.A0.SetRandom()
	z.C0.B1.A1.SetRandom()
	z.C0.B2.A0.SetRandom()
	z.C0.B2.A1.SetRandom()
	z.C1.B0.A0.SetRandom()
	z.C1.B0.A1.SetRandom()
	z.C1.B1.A0.SetRandom()
	z.C1.B1.A1.SetRandom()
	z.C1.B2.A0.SetRandom()
	z.C1.B2.A1.SetRandom()
	return z
}

// Mul set z=x*y in E12 and return z
func (z *E12) Mul(x, y *E12) *E12 {
	// Algorithm 20 from https://eprint.iacr.org/2010/354.pdf

	var t0, t1, xSum, ySum E6

	t0.Mul(&x.C0, &y.C0) // step 1
	t1.Mul(&x.C1, &y.C1) // step 2

	// finish processing input in case z==x or y
	xSum.Add(&x.C0, &x.C1)
	ySum.Add(&y.C0, &y.C1)

	// step 3
	{ // begin: inline z.C0.MulByNonResidue(&t1)
		var result E6
		result.B1.Set(&(&t1).B0)
		result.B2.Set(&(&t1).B1)
		{ // begin: inline result.B0.MulByNonResidue(&(&t1).B2)
			var buf, buf9 E2
			buf.Set(&(&t1).B2)
			buf9.Double(&buf).
				Double(&buf9).
				Double(&buf9).
				Add(&buf9, &buf)
			result.B0.A1.Add(&buf.A0, &buf9.A1)
			{ // begin: inline MulByNonResidue(&(result.B0).A0, &buf.A1)
				(&(result.B0).A0).Neg(&buf.A1)
			} // end: inline MulByNonResidue(&(result.B0).A0, &buf.A1)
			result.B0.A0.AddAssign(&buf9.A0)
		} // end: inline result.B0.MulByNonResidue(&(&t1).B2)
		z.C0.Set(&result)
	} // end: inline z.C0.MulByNonResidue(&t1)
	z.C0.Add(&z.C0, &t0)

	// step 4
	z.C1.Mul(&xSum, &ySum).
		Sub(&z.C1, &t0).
		Sub(&z.C1, &t1)

	return z
}

// Square set z=x*x in E12 and return z
func (z *E12) Square(x *E12) *E12 {
	// TODO implement Algorithm 22 from https://eprint.iacr.org/2010/354.pdf
	// or the complex method from fp2
	// for now do it the dumb way
	var b0, b1 E6

	b0.Square(&x.C0)
	b1.Square(&x.C1)
	{ // begin: inline b1.MulByNonResidue(&b1)
		var result E6
		result.B1.Set(&(&b1).B0)
		result.B2.Set(&(&b1).B1)
		{ // begin: inline result.B0.MulByNonResidue(&(&b1).B2)
			var buf, buf9 E2
			buf.Set(&(&b1).B2)
			buf9.Double(&buf).
				Double(&buf9).
				Double(&buf9).
				Add(&buf9, &buf)
			result.B0.A1.Add(&buf.A0, &buf9.A1)
			{ // begin: inline MulByNonResidue(&(result.B0).A0, &buf.A1)
				(&(result.B0).A0).Neg(&buf.A1)
			} // end: inline MulByNonResidue(&(result.B0).A0, &buf.A1)
			result.B0.A0.AddAssign(&buf9.A0)
		} // end: inline result.B0.MulByNonResidue(&(&b1).B2)
		b1.Set(&result)
	} // end: inline b1.MulByNonResidue(&b1)
	b1.Add(&b0, &b1)

	z.C1.Mul(&x.C0, &x.C1).Double(&z.C1)
	z.C0 = b1

	return z
}

// squares an element a+by as an Fp4 elmt, where y**2=1+u
func fp4Square(a, b, c, d *E2) {
	var tmp E2
	c.Square(a)
	tmp.Square(b).MulByNonResidue(&tmp)
	c.Add(c, &tmp)
	d.Mul(a, b).Double(d)
}

// CyclotomicSquare https://eprint.iacr.org/2009/565.pdf, 3.2
func (z *E12) CyclotomicSquare(x *E12) *E12 {

	var res, b, a E12
	var tmp E2

	// A
	fp4Square(&x.C0.B0, &x.C1.B1, &b.C0.B0, &b.C1.B1)
	a.C0.B0.Set(&x.C0.B0)
	a.C1.B1.Neg(&x.C1.B1)

	// B
	tmp.MulByNonResidueInv(&x.C1.B0)
	fp4Square(&x.C0.B2, &tmp, &b.C0.B1, &b.C1.B2)
	b.C0.B1.MulByNonResidue(&b.C0.B1)
	b.C1.B2.MulByNonResidue(&b.C1.B2)
	a.C0.B1.Set(&x.C0.B1)
	a.C1.B2.Neg(&x.C1.B2)

	// C
	fp4Square(&x.C0.B1, &x.C1.B2, &b.C0.B2, &b.C1.B0)
	b.C1.B0.MulByNonResidue(&b.C1.B0)
	a.C0.B2.Set(&x.C0.B2)
	a.C1.B0.Neg(&x.C1.B0)

	res.Set(&b)
	b.Sub(&b, &a).Double(&b)
	z.Add(&res, &b)

	return z
}

// Inverse set z to the inverse of x in E12 and return z
func (z *E12) Inverse(x *E12) *E12 {
	// Algorithm 23 from https://eprint.iacr.org/2010/354.pdf

	var t [2]E6

	t[0].Square(&x.C0)
	t[1].Square(&x.C1)
	{
		var buf E6
		{ // begin: inline buf.MulByNonResidue(&t[1])
			var result E6
			result.B1.Set(&(&t[1]).B0)
			result.B2.Set(&(&t[1]).B1)
			{ // begin: inline result.B0.MulByNonResidue(&(&t[1]).B2)
				var buf, buf9 E2
				buf.Set(&(&t[1]).B2)
				buf9.Double(&buf).
					Double(&buf9).
					Double(&buf9).
					Add(&buf9, &buf)
				result.B0.A1.Add(&buf.A0, &buf9.A1)
				{ // begin: inline MulByNonResidue(&(result.B0).A0, &buf.A1)
					(&(result.B0).A0).Neg(&buf.A1)
				} // end: inline MulByNonResidue(&(result.B0).A0, &buf.A1)
				result.B0.A0.AddAssign(&buf9.A0)
			} // end: inline result.B0.MulByNonResidue(&(&t[1]).B2)
			buf.Set(&result)
		} // end: inline buf.MulByNonResidue(&t[1])
		t[0].Sub(&t[0], &buf)
	}
	t[1].Inverse(&t[0])               // step 4
	z.C0.Mul(&x.C0, &t[1])            // step 5
	z.C1.Mul(&x.C1, &t[1]).Neg(&z.C1) // step 6

	return z
}

// InverseUnitary inverse a unitary element
func (z *E12) InverseUnitary(x *E12) *E12 {
	return z.Conjugate(x)
}

// Conjugate set z to (x.C0, -x.C1) and return z
func (z *E12) Conjugate(x *E12) *E12 {
	z.Set(x)
	z.C1.Neg(&z.C1)
	return z
}

// MulByVW set z to x*(y*v*w) and return z
// here y*v*w means the E12 element with C1.B1=y and all other components 0
func (z *E12) MulByVW(x *E12, y *E2) *E12 {
	var result E12
	var yNR E2

	{ // begin: inline yNR.MulByNonResidue(y)
		var buf, buf9 E2
		buf.Set(y)
		buf9.Double(&buf).
			Double(&buf9).
			Double(&buf9).
			Add(&buf9, &buf)
		yNR.A1.Add(&buf.A0, &buf9.A1)
		{ // begin: inline MulByNonResidue(&(yNR).A0, &buf.A1)
			(&(yNR).A0).Neg(&buf.A1)
		} // end: inline MulByNonResidue(&(yNR).A0, &buf.A1)
		yNR.A0.AddAssign(&buf9.A0)
	} // end: inline yNR.MulByNonResidue(y)
	result.C0.B0.Mul(&x.C1.B1, &yNR)
	result.C0.B1.Mul(&x.C1.B2, &yNR)
	result.C0.B2.Mul(&x.C1.B0, y)
	result.C1.B0.Mul(&x.C0.B2, &yNR)
	result.C1.B1.Mul(&x.C0.B0, y)
	result.C1.B2.Mul(&x.C0.B1, y)
	z.Set(&result)
	return z
}

// MulByV set z to x*(y*v) and return z
// here y*v means the E12 element with C0.B1=y and all other components 0
func (z *E12) MulByV(x *E12, y *E2) *E12 {
	var result E12
	var yNR E2

	{ // begin: inline yNR.MulByNonResidue(y)
		var buf, buf9 E2
		buf.Set(y)
		buf9.Double(&buf).
			Double(&buf9).
			Double(&buf9).
			Add(&buf9, &buf)
		yNR.A1.Add(&buf.A0, &buf9.A1)
		{ // begin: inline MulByNonResidue(&(yNR).A0, &buf.A1)
			(&(yNR).A0).Neg(&buf.A1)
		} // end: inline MulByNonResidue(&(yNR).A0, &buf.A1)
		yNR.A0.AddAssign(&buf9.A0)
	} // end: inline yNR.MulByNonResidue(y)
	result.C0.B0.Mul(&x.C0.B2, &yNR)
	result.C0.B1.Mul(&x.C0.B0, y)
	result.C0.B2.Mul(&x.C0.B1, y)
	result.C1.B0.Mul(&x.C1.B2, &yNR)
	result.C1.B1.Mul(&x.C1.B0, y)
	result.C1.B2.Mul(&x.C1.B1, y)
	z.Set(&result)
	return z
}

// MulByV2W set z to x*(y*v^2*w) and return z
// here y*v^2*w means the E12 element with C1.B2=y and all other components 0
func (z *E12) MulByV2W(x *E12, y *E2) *E12 {
	var result E12
	var yNR E2

	{ // begin: inline yNR.MulByNonResidue(y)
		var buf, buf9 E2
		buf.Set(y)
		buf9.Double(&buf).
			Double(&buf9).
			Double(&buf9).
			Add(&buf9, &buf)
		yNR.A1.Add(&buf.A0, &buf9.A1)
		{ // begin: inline MulByNonResidue(&(yNR).A0, &buf.A1)
			(&(yNR).A0).Neg(&buf.A1)
		} // end: inline MulByNonResidue(&(yNR).A0, &buf.A1)
		yNR.A0.AddAssign(&buf9.A0)
	} // end: inline yNR.MulByNonResidue(y)
	result.C0.B0.Mul(&x.C1.B0, &yNR)
	result.C0.B1.Mul(&x.C1.B1, &yNR)
	result.C0.B2.Mul(&x.C1.B2, &yNR)
	result.C1.B0.Mul(&x.C0.B1, &yNR)
	result.C1.B1.Mul(&x.C0.B2, &yNR)
	result.C1.B2.Mul(&x.C0.B0, y)
	z.Set(&result)
	return z
}
