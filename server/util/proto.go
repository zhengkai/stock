package util

import "google.golang.org/protobuf/proto"

var (
	protoMarshaler = proto.MarshalOptions{
		AllowPartial: true,
	}

	protoUnmarshaler = proto.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	MarshalProto   = protoMarshaler.Marshal
	UnmarshalProto = protoUnmarshaler.Unmarshal
)

func (f *File) ReadProto(m proto.Message) error {
	ab, err := f.Read()
	if err != nil {
		return err
	}
	return UnmarshalProto(ab, m)
}

func (f *File) WriteProto(m proto.Message) error {
	ab, err := MarshalProto(m)
	if err != nil {
		return err
	}
	return f.WriteWithHash(ab)
}
