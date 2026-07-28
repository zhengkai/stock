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

	MarshalJSON       = jsonMarshaler.Marshal
	MarshalJSONPretty = jsonMarshalerPretty.Marshal
	UnmarshalJSON     = jsonUnmarshaler.Unmarshal
)

func (f *File) ReadJSON(m proto.Message) error {
	ab, err := f.Read()
	if err != nil {
		return err
	}
	return UnmarshalJSON(ab, m)
}

func (f *File) WriteJSON(m proto.Message) error {
	ab, err := MarshalJSON(m)
	if err != nil {
		return err
	}
	return f.WriteWithHash(ab)
}

func JSON(m any) string {

	var ab []byte
	var err error

	if t, ok := m.(proto.Message); ok {
		ab, err = MarshalJSON(t)
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
		ab, err = MarshalJSON(t)
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
		ab, err = MarshalJSONPretty(t)
	} else {
		ab, err = json.MarshalIndent(m, ``, "\t")
	}
	if err != nil {
		return err.Error()
	}
	return string(ab)
}
