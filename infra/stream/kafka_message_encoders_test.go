package stream

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type myModel struct {
	Data string `json:"data"`
}

func Test_MessageJsonEncoder(t *testing.T) {
	encoder := NewMessageJsonEncoder()

	model := myModel{Data: "Sample"}

	raw, err := encoder.Encode(model)
	require.NoError(t, err)
	require.NotEmpty(t, raw)

	decoded, err := encoder.Decode(raw)
	require.NoError(t, err)

	require.Equal(t, decoded["data"], model.Data)
}
