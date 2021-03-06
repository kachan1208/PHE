package phe

import (
	"math/big"

	"github.com/pkg/errors"
)

//EnrollmentRecord stores all necessary password protection info
type EnrollmentRecord struct {
	NS []byte `json:"ns"`
	NC []byte `json:"nc"`
	T0 []byte `json:"t_0"`
	T1 []byte `json:"t_1"`
}

func (c *EnrollmentRecord) parse() (t0, t1 *Point, err error) {

	if c == nil ||
		len(c.NC) == 0 || len(c.NS) == 0 ||
		len(c.NC) > 32 || len(c.NS) > 32 {
		err = errors.New("invalid record")
		return
	}

	if t0, err = PointUnmarshal(c.T0); err != nil {
		return
	}

	t1, err = PointUnmarshal(c.T1)
	return
}

// ProofOfSuccess contains data for client to validate
type ProofOfSuccess struct {
	Term1  []byte `json:"term_1"`
	Term2  []byte `json:"term_2"`
	Term3  []byte `json:"term_3"`
	BlindX []byte `json:"blind_x"`
}

func (p *ProofOfSuccess) parse() (term1, term2, term3 *Point, blindX *big.Int, err error) {
	if p == nil {
		err = errors.New("invalid proof")
		return
	}

	if term1, err = PointUnmarshal(p.Term1); err != nil {
		return
	}

	if term2, err = PointUnmarshal(p.Term2); err != nil {
		return
	}

	if term3, err = PointUnmarshal(p.Term3); err != nil {
		return
	}

	if len(p.BlindX) == 0 || len(p.BlindX) > 32 {
		err = errors.New("invalid proof")
		return
	}
	blindX = new(big.Int).SetBytes(p.BlindX)

	return
}

// ProofOfFail contains data for client to validate
type ProofOfFail struct {
	Term1  []byte `json:"term_1"`
	Term2  []byte `json:"term_2"`
	Term3  []byte `json:"term_3"`
	Term4  []byte `json:"term_4"`
	BlindA []byte `json:"blind_a"`
	BlindB []byte `json:"blind_b"`
}

func (p *ProofOfFail) parse() (term1, term2, term3, term4 *Point, blindA, blindB *big.Int, err error) {
	if p == nil {
		err = errors.New("invalid proof")
		return
	}

	if term1, err = PointUnmarshal(p.Term1); err != nil {
		return
	}

	if term2, err = PointUnmarshal(p.Term2); err != nil {
		return
	}

	if term3, err = PointUnmarshal(p.Term3); err != nil {
		return
	}

	if term4, err = PointUnmarshal(p.Term4); err != nil {
		return
	}

	if len(p.BlindA) == 0 || len(p.BlindA) > 32 {
		err = errors.New("invalid proof")
		return
	}

	if len(p.BlindB) == 0 || len(p.BlindB) > 32 {
		err = errors.New("invalid proof")
		return
	}

	blindA = new(big.Int).SetBytes(p.BlindA)
	blindB = new(big.Int).SetBytes(p.BlindB)

	return
}

// UpdateToken contains values needed for value rotation
type UpdateToken struct {
	A []byte `json:"a"`
	B []byte `json:"b"`
}

func (t *UpdateToken) parse() (a, b *big.Int, err error) {
	if t == nil {
		return nil, nil, errors.New("invalid token")
	}
	if len(t.A) == 0 || len(t.A) > 32 {
		return nil, nil, errors.New("invalid update token")
	}
	if len(t.B) == 0 || len(t.B) > 32 {
		return nil, nil, errors.New("invalid update token")
	}

	a = new(big.Int).SetBytes(t.A)
	b = new(big.Int).SetBytes(t.B)
	return
}

// EnrollmentResponse contains two pseudo-random points and seed which server used to generate them
type EnrollmentResponse struct {
	NS    []byte          `json:"ns"`
	C0    []byte          `json:"c_0"`
	C1    []byte          `json:"c_1"`
	Proof *ProofOfSuccess `json:"proof"`
}

// VerifyPasswordRequest contains server's nonce and an attempt to verify a password in form of an elliptic curve point
type VerifyPasswordRequest struct {
	NS []byte `json:"ns"`
	C0 []byte `json:"c_0"`
}

//VerifyPasswordResponse returns the result of evaluating an entered password along with the zero knowledge proof
type VerifyPasswordResponse struct {
	Res          bool            `json:"res"`
	C1           []byte          `json:"c_1"`
	ProofSuccess *ProofOfSuccess `json:"proof_success,omitempty"`
	ProofFail    *ProofOfFail    `json:"proof_fail,omitempty"`
}

type keypair struct {
	PublicKey  []byte
	PrivateKey []byte
}
