package random

import (
	"testing"
)

func BenchmarkRandomUserID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateUserID()
	}
}

func BenchmarkRandomAckID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateAckID()
	}
}
