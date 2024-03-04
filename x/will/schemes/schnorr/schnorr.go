package schnorr

import (
	"math/big"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

// Set parameters
var (
	curve  = edwards25519.NewBlakeSHA256Ed25519()
	sha256 = curve.Hash()
	g      = curve.Point().Base()
)

type Signature struct {
	R kyber.Point
	S kyber.Scalar
}

// Sign using Schnorr EdDSA
// m: Message
// z: Private key
func Sign(m kyber.Scalar, z kyber.Scalar) Signature {
	// Pick a random k from allowed set.
	k := curve.Scalar().Pick(curve.RandomStream())

	// r = k * G (a.k.a the same operation as r = g^k)
	r := curve.Point().Mul(k, g)

	// h := Hash(r.String() + m + publicKey)
	// publicKey := curve.Point().Mul(z, g)
	// h := Hash(publicKey.String() + m)

	// s = k + e * x
	s := curve.Scalar().Add(k, curve.Scalar().Mul(m, z))

	return Signature{R: r, S: s}
}

// Verify Schnorr EdDSA signatures
// m: Message
// s: Signature
// y: Public key
func Verify(m kyber.Scalar, S Signature, y kyber.Point) bool {
	// Attempt to reconstruct 's * G' with a provided signature; s * G = r - h * y
	sGv := curve.Point().Add(S.R, curve.Point().Mul(m, y))

	// Construct the actual 's * G'
	sG := curve.Point().Mul(S.S, g)

	// Equality check; ensure signature and public key outputs to s * G.
	return sG.Equal(sGv)
}

// Return a new random key pair
func RandomKeyPair() (kyber.Scalar, kyber.Point) {
	privateKey := curve.Scalar().Pick(curve.RandomStream())
	publicKey := curve.Point().Mul(privateKey, g)

	return privateKey, publicKey
}

// Given string, return hash Scalar
func Hash(s string) kyber.Scalar {
	sha256.Reset()
	sha256.Write([]byte(s))

	return curve.Scalar().SetBytes(sha256.Sum(nil))
}

// Given private key in big.Int, return its hash as Scalar
func convKey(d *big.Int) kyber.Scalar {
	sha256.Reset()
	sha256.Write(d.Bytes())

	return curve.Scalar().SetBytes(sha256.Sum(nil))
}

// ------------------------------------ //
// Generate a multi-signature given 2 signatures and the corresponding public keys
func mulSig(sigA, sigB Signature, pubKeyA, pubKeyB kyber.Point) (Signature, kyber.Point) {
	newR := curve.Point().Add(sigA.R, sigB.R)
	newS := curve.Scalar().Add(sigA.S, sigB.S)
	newPubKey := curve.Point().Add(pubKeyA, pubKeyB)
	// fmt.Println("Aggregated publicKey in aggSig: ", newPubKey)

	var sigC Signature
	sigC = Signature{
		R: newR,
		S: newS,
	}

	return sigC, newPubKey
}
