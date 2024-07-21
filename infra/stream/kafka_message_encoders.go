package stream

import "encoding/json"

type MessageEncoder interface {
	Encode(s interface{}) (string, error)
	Decode(s string) (map[string]interface{}, error)
}

type jsonEncoder struct{}

func NewMessageJsonEncoder() MessageEncoder {
	return &jsonEncoder{}
}

func (e *jsonEncoder) Encode(s interface{}) (string, error) {
	str, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	return string(str), nil
}

func (e *jsonEncoder) Decode(s string) (map[string]interface{}, error) {
	var d map[string]interface{}

	err := json.Unmarshal([]byte(s), &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}
