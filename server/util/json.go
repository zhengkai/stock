package util

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	jsonMarshaler = protojson.MarshalOptions{
		AllowPartial:    true,
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}

	jsonUnmarshaler = protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
)

func (f *File) ReadJSON(m proto.Message) error {
	ab, err := f.Read()
	if err != nil {
		return err
	}
	return jsonUnmarshaler.Unmarshal(ab, m)
}

func (f *File) WriteJSON(m proto.Message) error {
	ab, err := jsonMarshaler.Marshal(m)
	if err != nil {
		return err
	}
	return f.Write(ab)
}

func JSON(m any) string {

	var ab []byte
	var err error

	if t, ok := m.(proto.Message); ok {
		ab, err = jsonMarshaler.Marshal(t)
	} else {
		ab, err = json.Marshal(m)
	}
	if err != nil {
		return err.Error()
	}
	return string(ab)
}
