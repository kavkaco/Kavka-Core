package utils

import (
	"testing"
)

type TestData struct {
	Field1 string
	Field2 int
}

func BenchmarkTypeConverter(b *testing.B) {
	data := TestData{"hello", 123}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := TypeConverter[TestData](data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
