package structs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Sample struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func Test_SetFieldByBSON(t *testing.T) {
	samples := []*Sample{{Name: "Sample", Age: 17}}
	newAge := 25

	for _, sample := range samples {
		err := SetFieldByBSON(sample, "age", newAge)

		require.NoError(t, err)
	}

	require.Equal(t, samples[0].Age, newAge)
}
