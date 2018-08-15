package phe

import (
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	l   *Server
	s   *Client
	pwd = []byte("Password")
)

func init() {
	l = &Server{RandomZ()}
	s = &Client{Y: RandomZ()}
}

func BenchmarkAddP256(b *testing.B) {
	b.ResetTimer()
	p256 := elliptic.P256()
	_, x, y, _ := elliptic.GenerateKey(p256, rand.Reader)
	_, x1, y1, _ := elliptic.GenerateKey(p256, rand.Reader)

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p256.Add(x, y, x1, y1)
	}

}

func Test_PHE(t *testing.T) {

	//first, ask server for random values & proof
	ns, c0, proof := l.GetEnrollment()

	// Enroll account

	nc, t0, err := s.EnrollAccount(pwd, ns, c0, proof)
	assert.NoError(t, err)

	//Check password request
	c0 = s.CreateVerifyPasswordRequest(nc, pwd, t0)
	//Check password on server
	res, proof := l.VerifyPassword(ns, c0)
	//validate response & decrypt M
	err = s.CheckResponseAndDecrypt(t0, pwd, ns, nc, proof, res)
	assert.NoError(t, err)

}

func Test_PHE_InvalidPassword(t *testing.T) {

	//first, ask server for random values & proof
	ns, c0, proof := l.GetEnrollment()

	// Enroll account
	nc, t0, err := s.EnrollAccount(pwd, ns, c0, proof)
	assert.NoError(t, err)

	//Check password request
	c0 = s.CreateVerifyPasswordRequest(nc, []byte("Password1"), t0)
	//Check password on server
	res, proof := l.VerifyPassword(ns, c0)
	//validate response & decrypt M
	err = s.CheckResponseAndDecrypt(t0, []byte("Password1"), ns, nc, proof, res)
	assert.Nil(t, err)

}

func BenchmarkServer_GetEnrollment(b *testing.B) {

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.GetEnrollment()
	}
}

func BenchmarkClient_EnrollAccount(b *testing.B) {

	ns, c0, proof := l.GetEnrollment()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, err := s.EnrollAccount(pwd, ns, c0, proof)
		assert.NoError(b, err)
	}
}

func BenchmarkClient_CreateVerifyPasswordRequest(b *testing.B) {
	//first, ask server for random values & proof
	ns, c0, proof := l.GetEnrollment()

	// Enroll account

	nc, t0, err := s.EnrollAccount(pwd, ns, c0, proof)
	assert.NoError(b, err)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//Check password request
		s.CreateVerifyPasswordRequest(nc, pwd, t0)
	}
}

func BenchmarkLoginFlow(b *testing.B) {

	//first, ask server for random values & proof
	ns, c0, proof := l.GetEnrollment()

	// Enroll account

	nc, t0, err := s.EnrollAccount(pwd, ns, c0, proof)
	assert.NoError(b, err)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//Check password request
		c0 = s.CreateVerifyPasswordRequest(nc, pwd, t0)
		//Check password on server
		res, proof := l.VerifyPassword(ns, c0)
		//validate response & decrypt M
		err := s.CheckResponseAndDecrypt(t0, pwd, ns, nc, proof, res)
		assert.NoError(b, err)
	}
}
