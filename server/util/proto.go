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
)

func (f *File) ReadProto(m proto.Message) error {
	ab, err := f.Read()
	if err != nil {
		return err
	}
	return protoUnmarshaler.Unmarshal(ab, m)
}

func (f *File) WriteProto(m proto.Message) error {
	ab, err := protoMarshaler.Marshal(m)
	if err != nil {
		return err
	}
	return f.Write(ab)
}
