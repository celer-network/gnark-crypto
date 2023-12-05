// Copyright 2020 Consensys Software Inc.
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

// Code generated by consensys/gnark-crypto DO NOT EDIT

package bw6633

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bw6-633/fp"
	"github.com/consensys/gnark-crypto/ecc/bw6-633/fr"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

// ------------------------------------------------------------
// tests

func TestPairing(t *testing.T) {

	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := GenE6()

	genR1 := GenFr()
	genR2 := GenFr()

	properties.Property("[BW6-633] Having the receiver as operand (final expo) should output the same result", prop.ForAll(
		func(a GT) bool {
			b := FinalExponentiation(&a)
			a = FinalExponentiation(&a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("[BW6-633] Exponentiating FinalExpo(a) to r should output 1", prop.ForAll(
		func(a GT) bool {
			b := FinalExponentiation(&a)
			return !a.IsInSubGroup() && b.IsInSubGroup()
		},
		genA,
	))

	properties.Property("[BW6-633] Exp, CyclotomicExp and ExpGLV results must be the same in GT (small and big exponents)", prop.ForAll(
		func(a GT, e fr.Element) bool {

			var res bool

			// exponent > r
			{
				a = FinalExponentiation(&a)
				var _e big.Int
				_e.SetString("169893631828481842931290008859743243489098146141979830311893424751855271950692001433356165550548410610101138388623573573742608490725625288296502860183437011025036209791574001140592327223981416956942076610555083128655330944007957223952510233203018053264066056080064687038560794652180979019775788172491868553073169893631828481842931290008859743243489098146141979830311893424751855271950692001433356165550548410610101138388623573573742608490725625288296502860183437011025036209791574001140592327223981416956942076610555083128655330944007957223952510233203018053264066056080064687038560794652180979019775788172491868553073", 10)
				var b, c, d GT
				b.Exp(a, &_e)
				c.ExpGLV(a, &_e)
				d.CyclotomicExp(a, &_e)
				res = b.Equal(&c) && c.Equal(&d)
			}

			// exponent < r
			{
				a = FinalExponentiation(&a)
				var _e big.Int
				e.BigInt(&_e)
				var b, c, d GT
				b.Exp(a, &_e)
				c.ExpGLV(a, &_e)
				d.CyclotomicExp(a, &_e)
				res = res && b.Equal(&c) && c.Equal(&d)
			}

			return res
		},
		genA,
		genR1,
	))

	properties.Property("[BW6-633] Expt(Expt) and Exp(t^2) should output the same result in the cyclotomic subgroup", prop.ForAll(
		func(a GT) bool {
			var b, c, d GT
			b.Conjugate(&a)
			a.Inverse(&a)
			b.Mul(&b, &a)

			a.Frobenius(&b).
				Mul(&a, &b)

			c.Expt(&a).Expt(&c)
			d.Exp(a, &xGen).Exp(d, &xGen)
			return c.Equal(&d)
		},
		genA,
	))

	properties.Property("[BW6-633] bilinearity", prop.ForAll(
		func(a, b fr.Element) bool {

			var res, resa, resb, resab, zero GT

			var ag1 G1Affine
			var bg2 G2Affine

			var abigint, bbigint, ab big.Int

			a.BigInt(&abigint)
			b.BigInt(&bbigint)
			ab.Mul(&abigint, &bbigint)

			ag1.ScalarMultiplication(&g1GenAff, &abigint)
			bg2.ScalarMultiplication(&g2GenAff, &bbigint)

			res, _ = Pair([]G1Affine{g1GenAff}, []G2Affine{g2GenAff})
			resa, _ = Pair([]G1Affine{ag1}, []G2Affine{g2GenAff})
			resb, _ = Pair([]G1Affine{g1GenAff}, []G2Affine{bg2})

			resab.Exp(res, &ab)
			resa.Exp(resa, &bbigint)
			resb.Exp(resb, &abigint)

			return resab.Equal(&resa) && resab.Equal(&resb) && !res.Equal(&zero)

		},
		genR1,
		genR2,
	))

	properties.Property("[BW6-633] PairingCheck", prop.ForAll(
		func(a, b fr.Element) bool {

			var g1GenAffNeg G1Affine
			g1GenAffNeg.Neg(&g1GenAff)
			tabP := []G1Affine{g1GenAff, g1GenAffNeg}
			tabQ := []G2Affine{g2GenAff, g2GenAff}

			res, _ := PairingCheck(tabP, tabQ)

			return res
		},
		genR1,
		genR2,
	))

	properties.Property("[BW6-633] Pair should output the same result with MillerLoop or MillerLoopFixedQ", prop.ForAll(
		func(a, b fr.Element) bool {

			var ag1 G1Affine
			var bg2 G2Affine

			var abigint, bbigint big.Int

			a.BigInt(&abigint)
			b.BigInt(&bbigint)

			ag1.ScalarMultiplication(&g1GenAff, &abigint)
			bg2.ScalarMultiplication(&g2GenAff, &bbigint)

			P := []G1Affine{g1GenAff, ag1}
			Q := []G2Affine{g2GenAff, bg2}

			ml1, _ := MillerLoop(P, Q)
			ml2, _ := MillerLoopFixedQ(
				P,
				[][2][len(LoopCounter) - 1]LineEvaluationAff{
					PrecomputeLines(Q[0]),
					PrecomputeLines(Q[1]),
				})

			res1 := FinalExponentiation(&ml1)
			res2 := FinalExponentiation(&ml2)

			return res1.Equal(&res2)
		},
		genR1,
		genR2,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestMillerLoop(t *testing.T) {

	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genR1 := GenFr()
	genR2 := GenFr()

	properties.Property("[BW6-633] MillerLoop of pairs should be equal to the product of MillerLoops", prop.ForAll(
		func(a, b fr.Element) bool {

			var simpleProd, factorizedProd GT

			var ag1 G1Affine
			var bg2 G2Affine

			var abigint, bbigint big.Int

			a.BigInt(&abigint)
			b.BigInt(&bbigint)

			ag1.ScalarMultiplication(&g1GenAff, &abigint)
			bg2.ScalarMultiplication(&g2GenAff, &bbigint)

			P0 := []G1Affine{g1GenAff}
			P1 := []G1Affine{ag1}
			Q0 := []G2Affine{g2GenAff}
			Q1 := []G2Affine{bg2}

			// FE( ML(a,b) * ML(c,d) * ML(e,f) * ML(g,h) )
			M1, _ := MillerLoop(P0, Q0)
			M2, _ := MillerLoop(P1, Q0)
			M3, _ := MillerLoop(P0, Q1)
			M4, _ := MillerLoop(P1, Q1)
			simpleProd.Mul(&M1, &M2).Mul(&simpleProd, &M3).Mul(&simpleProd, &M4)
			simpleProd = FinalExponentiation(&simpleProd)

			tabP := []G1Affine{g1GenAff, ag1, g1GenAff, ag1}
			tabQ := []G2Affine{g2GenAff, g2GenAff, bg2, bg2}

			// FE( ML([a,c,e,g] ; [b,d,f,h]) ) -> saves 3 squares in Fqk
			factorizedProd, _ = Pair(tabP, tabQ)

			return simpleProd.Equal(&factorizedProd)
		},
		genR1,
		genR2,
	))

	properties.Property("[BW6-633] MillerLoop and MillerLoopFixedQ should skip pairs with a point at infinity", prop.ForAll(
		func(a, b fr.Element) bool {

			var one GT

			var ag1, g1Inf G1Affine
			var bg2, g2Inf G2Affine

			var abigint, bbigint big.Int

			one.SetOne()

			a.BigInt(&abigint)
			b.BigInt(&bbigint)

			ag1.ScalarMultiplication(&g1GenAff, &abigint)
			bg2.ScalarMultiplication(&g2GenAff, &bbigint)

			g1Inf.FromJacobian(&g1Infinity)
			g2Inf.FromJacobian(&g2Infinity)

			// e([0,c] ; [b,d])
			// -> should be equal to e(c,d)
			tabP := []G1Affine{g1Inf, ag1}
			tabQ := []G2Affine{g2GenAff, bg2}
			res1, _ := Pair(tabP, tabQ)

			// e([a,c] ; [0,d])
			// -> should be equal to e(c,d)
			tabP = []G1Affine{g1GenAff, ag1}
			tabQ = []G2Affine{g2Inf, bg2}
			res2, _ := Pair(tabP, tabQ)

			// e([0,c] ; [b,d]) with fixed points b and d
			// -> should be equal to e(c,d)
			tabP = []G1Affine{g1Inf, ag1}
			linesQ := [][2][len(LoopCounter) - 1]LineEvaluationAff{
				PrecomputeLines(g2GenAff),
				PrecomputeLines(bg2),
			}
			res3, _ := PairFixedQ(tabP, linesQ)

			// e([a,c] ; [0,d]) with fixed points 0 and d
			// -> should be equal to e(c,d)
			tabP = []G1Affine{g1GenAff, ag1}
			linesQ = [][2][len(LoopCounter) - 1]LineEvaluationAff{
				PrecomputeLines(g2Inf),
				PrecomputeLines(bg2),
			}
			res4, _ := PairFixedQ(tabP, linesQ)

			// e([0,c] ; [d,0])
			// -> should be equal to 1
			tabP = []G1Affine{g1Inf, ag1}
			tabQ = []G2Affine{bg2, g2Inf}
			res5, _ := Pair(tabP, tabQ)

			// e([0,c] ; [d,0]) with fixed points d and 0
			// -> should be equal to 1
			tabP = []G1Affine{g1Inf, ag1}
			linesQ = [][2][len(LoopCounter) - 1]LineEvaluationAff{
				PrecomputeLines(bg2),
				PrecomputeLines(g2Inf),
			}
			res6, _ := PairFixedQ(tabP, linesQ)

			// e([0,0])
			// -> should be equal to 1
			tabP = []G1Affine{g1Inf}
			tabQ = []G2Affine{g2Inf}
			res7, _ := Pair(tabP, tabQ)

			// e([0,0]) with fixed point 0
			// -> should be equal to 1
			tabP = []G1Affine{g1Inf}
			linesQ = [][2][len(LoopCounter) - 1]LineEvaluationAff{
				PrecomputeLines(g2Inf),
			}
			res8, _ := PairFixedQ(tabP, linesQ)

			return res1.Equal(&res2) && res2.Equal(&res3) && res3.Equal(&res4) &&
				res5.Equal(&one) && res6.Equal(&one) && res7.Equal(&one) && res8.Equal(&one)
		},
		genR1,
		genR2,
	))

	properties.Property("[BW6-633] compressed pairing", prop.ForAll(
		func(a, b fr.Element) bool {

			var ag1 G1Affine
			var bg2 G2Affine

			var abigint, bbigint big.Int

			a.BigInt(&abigint)
			b.BigInt(&bbigint)

			ag1.ScalarMultiplication(&g1GenAff, &abigint)
			bg2.ScalarMultiplication(&g2GenAff, &bbigint)

			res, _ := Pair([]G1Affine{ag1}, []G2Affine{bg2})

			compressed, _ := res.CompressTorus()
			decompressed := compressed.DecompressTorus()

			return decompressed.Equal(&res)

		},
		genR1,
		genR2,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// ------------------------------------------------------------
// benches

func BenchmarkPairing(b *testing.B) {

	var g1GenAff G1Affine
	var g2GenAff G2Affine

	g1GenAff.FromJacobian(&g1Gen)
	g2GenAff.FromJacobian(&g2Gen)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Pair([]G1Affine{g1GenAff}, []G2Affine{g2GenAff})
	}
}

func BenchmarkMillerLoop(b *testing.B) {

	var g1GenAff G1Affine
	var g2GenAff G2Affine

	g1GenAff.FromJacobian(&g1Gen)
	g2GenAff.FromJacobian(&g2Gen)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MillerLoop([]G1Affine{g1GenAff}, []G2Affine{g2GenAff})
	}
}

func BenchmarkFinalExponentiation(b *testing.B) {

	var a GT
	a.SetRandom()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FinalExponentiation(&a)
	}

}

func BenchmarkMultiMiller(b *testing.B) {

	var g1GenAff G1Affine
	var g2GenAff G2Affine

	g1GenAff.FromJacobian(&g1Gen)
	g2GenAff.FromJacobian(&g2Gen)

	n := 10
	P := make([]G1Affine, n)
	Q := make([]G2Affine, n)

	for i := 2; i <= n; i++ {
		for j := 0; j < i; j++ {
			P[j].Set(&g1GenAff)
			Q[j].Set(&g2GenAff)
		}
		b.Run(fmt.Sprintf("%d pairs", i), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				MillerLoop(P, Q)
			}
		})
	}
}

func BenchmarkMultiPair(b *testing.B) {

	var g1GenAff G1Affine
	var g2GenAff G2Affine

	g1GenAff.FromJacobian(&g1Gen)
	g2GenAff.FromJacobian(&g2Gen)

	n := 10
	P := make([]G1Affine, n)
	Q := make([]G2Affine, n)

	for i := 2; i <= n; i++ {
		for j := 0; j < i; j++ {
			P[j].Set(&g1GenAff)
			Q[j].Set(&g2GenAff)
		}
		b.Run(fmt.Sprintf("%d pairs", i), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Pair(P, Q)
			}
		})
	}
}

func BenchmarkExpGT(b *testing.B) {

	var a GT
	a.SetRandom()
	a = FinalExponentiation(&a)

	var e fp.Element
	e.SetRandom()

	k := new(big.Int).SetUint64(6)

	e.Exp(e, k)
	var _e big.Int
	e.BigInt(&_e)

	b.Run("Naive windowed Exp", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			a.Exp(a, &_e)
		}
	})

	b.Run("2-NAF cyclotomic Exp", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			a.CyclotomicExp(a, &_e)
		}
	})

	b.Run("windowed 2-dim GLV Exp", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			a.ExpGLV(a, &_e)
		}
	})
}
