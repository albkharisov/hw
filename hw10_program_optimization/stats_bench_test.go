//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkSimple(b *testing.B) {
	b.StopTimer()

	r, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		panic(err)
	}
	if 1 != len(r.File) {
		panic("1 != len(r.File)")
	}
	_, err = r.File[0].Open()
	if err != nil {
		panic(err)
	}
	r.Close()

	for i := 0; i < b.N; i++ {
		r, _ := zip.OpenReader("testdata/users.dat.zip")
		data, err := r.File[0].Open()

		b.StartTimer()
		_, err = GetDomainStat(data, "biz")
		b.StopTimer()

		if err != nil {
			panic(err)
		}
		r.Close()
	}
}
