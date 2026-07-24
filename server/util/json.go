package util

import (
	"encoding/json"
	"project/zj"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	jsonMarshaler = protojson.MarshalOptions{
		AllowPartial:    true,
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}

	jsonMarshalerPretty = protojson.MarshalOptions{
		AllowPartial:    true,
		UseProtoNames:   true,
		EmitUnpopulated: true,
		Indent:          "\t",
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

func JSONBin(m any) []byte {

	var ab []byte
	var err error

	if t, ok := m.(proto.Message); ok {
		ab, err = jsonMarshaler.Marshal(t)
	} else {
		ab, err = json.Marshal(m)
	}
	if err != nil {
		zj.W(err)
		return nil
	}
	return ab
}

func JSONPretty(m any) string {

	var ab []byte
	var err error

	if t, ok := m.(proto.Message); ok {
		ab, err = jsonMarshalerPretty.Marshal(t)
	} else {
		ab, err = json.MarshalIndent(m, ``, "\t")
	}
	if err != nil {
		return err.Error()
	}
	return string(ab)
}
