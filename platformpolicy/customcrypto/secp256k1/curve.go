package secp256k1

import (
	"crypto/elliptic"
	"math/big"
)

// A Curve represents a short-form Weierstrass curve with a=-3.
// See https://www.hyperelliptic.org/EFD/g1p/auto-shortw.html
type Curve interface {
	// Params returns the parameters for the curve.
	Params() *CurveParams
	// IsOnCurve reports whether the given (x,y) lies on the curve.
	IsOnCurve(x, y *big.Int) bool
	// Add returns the sum of (x1,y1) and (x2,y2)
	Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int)
	// Double returns 2*(x,y)
	Double(x1, y1 *big.Int) (x, y *big.Int)
	// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
	ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int)
	// ScalarBaseMult returns k*G, where G is the base point of the group
	// and k is an integer in big-endian form.
	ScalarBaseMult(k []byte) (x, y *big.Int)
}

type CurveParams struct {
	P       *big.Int
	N       *big.Int
	B       *big.Int
	Gx, Gy  *big.Int
	BitSize int
	Name    string
}

var (
	secp256k1 *CurveParams
	p, _      = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	n, _      = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	b, _      = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
	gx, _     = new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
	gy, _     = new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)
	bitSize   = 256
	name      = "secp256k1"
)

func Secp256k1() elliptic.Curve {
	return &CurveParams{
		P:       p,
		N:       n,
		B:       b,
		Gx:      gx,
		Gy:      gy,
		BitSize: bitSize,
		Name:    name,
	}
}

func (curve *CurveParams) Params() *elliptic.CurveParams {
	return &elliptic.CurveParams{
		P:       p,
		N:       n,
		B:       b,
		Gx:      gx,
		Gy:      gy,
		BitSize: bitSize,
		Name:    name,
	}
}

func (curve *CurveParams) IsOnCurve(x, y *big.Int) bool {
	y2 := new(big.Int).Exp(y, big.NewInt(2), nil)
	x3 := new(big.Int).Exp(x, big.NewInt(3), nil)
	ans := new(big.Int).Mod(y2.Sub(y2, x3.Add(x3, curve.B)), curve.P)
	return ans.Cmp(big.NewInt(0)) == 0
}

func (curve *CurveParams) addJacobian(x1, y1, z1, x2, y2, z2 *big.Int) (*big.Int, *big.Int, *big.Int) {
	x3, y3, z3 := new(big.Int), new(big.Int), new(big.Int)
	if z1.Sign() == 0 {
		x3.Set(x2)
		y3.Set(y2)
		z3.Set(z2)
		return x3, y3, z3
	}
	if z2.Sign() == 0 {
		x3.Set(x1)
		y3.Set(y1)
		z3.Set(z1)
		return x3, y3, z3
	}
	z1z1 := new(big.Int).Mul(z1, z1)
	z1z1.Mod(z1z1, curve.P)
	z2z2 := new(big.Int).Mul(z2, z2)
	z2z2.Mod(z2z2, curve.P)
	u1 := new(big.Int).Mul(x1, z2z2)
	u1.Mod(u1, curve.P)
	u2 := new(big.Int).Mul(x2, z1z1)
	u2.Mod(u2, curve.P)
	h := new(big.Int).Sub(u2, u1)
	xEqual := h.Sign() == 0
	if h.Sign() == -1 {
		h.Add(h, curve.P)
	}
	i := new(big.Int).Lsh(h, 1)
	i.Mul(i, i)
	j := new(big.Int).Mul(h, i)
	s1 := new(big.Int).Mul(y1, z2)
	s1.Mul(s1, z2z2)
	s1.Mod(s1, curve.P)
	s2 := new(big.Int).Mul(y2, z1)
	s2.Mul(s2, z1z1)
	s2.Mod(s2, curve.P)
	r := new(big.Int).Sub(s2, s1)
	if r.Sign() == -1 {
		r.Add(r, curve.P)
	}
	yEqual := r.Sign() == 0
	if xEqual && yEqual {
		return curve.doubleJacobian(x1, y1, z1)
	}
	r.Lsh(r, 1)
	v := new(big.Int).Mul(u1, i)
	x3.Set(r)
	x3.Mul(x3, x3)
	x3.Sub(x3, j)
	x3.Sub(x3, v)
	x3.Sub(x3, v)
	x3.Mod(x3, curve.P)
	y3.Set(r)
	v.Sub(v, x3)
	y3.Mul(y3, v)
	s1.Mul(s1, j)
	s1.Lsh(s1, 1)
	y3.Sub(y3, s1)
	y3.Mod(y3, curve.P)
	z3.Add(z1, z2)
	z3.Mul(z3, z3)
	z3.Sub(z3, z1z1)
	z3.Sub(z3, z2z2)
	z3.Mul(z3, h)
	z3.Mod(z3, curve.P)
	return x3, y3, z3
}

func (curve *CurveParams) affineFromJacobian(x, y, z *big.Int) (xOut, yOut *big.Int) {
	if z.Sign() == 0 {
		return new(big.Int), new(big.Int)
	}
	zinv := new(big.Int).ModInverse(z, curve.P)
	zinvsq := new(big.Int).Mul(zinv, zinv)
	xOut = new(big.Int).Mul(x, zinvsq)
	xOut.Mod(xOut, curve.P)
	zinvsq.Mul(zinvsq, zinv)
	yOut = new(big.Int).Mul(y, zinvsq)
	yOut.Mod(yOut, curve.P)
	return
}

func zForAffine(x, y *big.Int) *big.Int {
	z := new(big.Int)
	if x.Sign() != 0 || y.Sign() != 0 {
		z.SetInt64(1)
	}
	return z
}

func (curve *CurveParams) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	z1 := zForAffine(x1, y1)
	z2 := zForAffine(x2, y2)
	return curve.affineFromJacobian(curve.addJacobian(x1, y1, z1, x2, y2, z2))
}

func (curve *CurveParams) doubleJacobian(x, y, z *big.Int) (*big.Int, *big.Int, *big.Int) {
	a := new(big.Int).Mul(x, x)
	b := new(big.Int).Mul(y, y)
	c := new(big.Int).Mul(b, b)
	d := new(big.Int).Add(x, b)
	d.Mul(d, d)
	d.Sub(d, a)
	d.Sub(d, c)
	d.Mul(d, big.NewInt(2))
	e := new(big.Int).Mul(big.NewInt(3), a)
	f := new(big.Int).Mul(e, e)
	x3 := new(big.Int).Mul(big.NewInt(2), d)
	x3.Sub(f, x3)
	x3.Mod(x3, curve.P)
	y3 := new(big.Int).Sub(d, x3)
	y3.Mul(e, y3)
	y3.Sub(y3, new(big.Int).Mul(big.NewInt(8), c))
	y3.Mod(y3, curve.P)
	z3 := new(big.Int).Mul(y, z)
	z3.Mul(big.NewInt(2), z3)
	z3.Mod(z3, curve.P)
	return x3, y3, z3
}

func (curve *CurveParams) Double(x1, y1 *big.Int) (x, y *big.Int) {
	z1 := zForAffine(x1, y1)
	return curve.affineFromJacobian(curve.doubleJacobian(x1, y1, z1))
}

func (curve *CurveParams) ScalarMult(Bx, By *big.Int, k []byte) (x, y *big.Int) {
	Bz := new(big.Int).SetInt64(1)
	x, y, z := new(big.Int), new(big.Int), new(big.Int)
	for _, byte := range k {
		for bitNum := 0; bitNum < 8; bitNum++ {
			x, y, z = curve.doubleJacobian(x, y, z)
			if byte&0x80 == 0x80 {
				x, y, z = curve.addJacobian(Bx, By, Bz, x, y, z)
			}
			byte <<= 1
		}
	}
	return curve.affineFromJacobian(x, y, z)
}

func (curve *CurveParams) ScalarBaseMult(k []byte) (x, y *big.Int) {
	return curve.ScalarMult(curve.Gx, curve.Gy, k)
}
