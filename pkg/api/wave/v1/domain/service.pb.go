// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: wave/v1/domain/service.proto

package domain

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_wave_v1_domain_service_proto protoreflect.FileDescriptor

var file_wave_v1_domain_service_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x77, 0x61, 0x76, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e,
	0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x1a, 0x1c,
	0x77, 0x61, 0x76, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xd2, 0x04, 0x0a,
	0x0d, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x40,
	0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1a, 0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1b, 0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x43, 0x0a, 0x04, 0x53, 0x74, 0x61, 0x74, 0x12, 0x1b, 0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x43, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1b, 0x2e,
	0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x77, 0x61, 0x76,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x49, 0x0a, 0x06, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x49, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12,
	0x1d, 0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e,
	0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x49, 0x0a, 0x06, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x12, 0x1d, 0x2e, 0x77, 0x61, 0x76,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x41, 0x74, 0x74, 0x61,
	0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x77, 0x61, 0x76, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x41, 0x74, 0x74, 0x61, 0x63,
	0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x49, 0x0a, 0x06, 0x44,
	0x65, 0x74, 0x61, 0x63, 0x68, 0x12, 0x1d, 0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x44, 0x65, 0x74, 0x61, 0x63, 0x68, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x44, 0x65, 0x74, 0x61, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x49, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x12, 0x1d, 0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69,
	0x6e, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1e, 0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x27, 0x5a, 0x25, 0x63, 0x74, 0x68, 0x75, 0x6c, 0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x74,
	0x68, 0x75, 0x6c, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x77, 0x61, 0x76, 0x65,
	0x2f, 0x76, 0x31, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var file_wave_v1_domain_service_proto_goTypes = []any{
	(*GetRequest)(nil),     // 0: wave.v1.domain.GetRequest
	(*StatRequest)(nil),    // 1: wave.v1.domain.StatRequest
	(*ListRequest)(nil),    // 2: wave.v1.domain.ListRequest
	(*CreateRequest)(nil),  // 3: wave.v1.domain.CreateRequest
	(*UpdateRequest)(nil),  // 4: wave.v1.domain.UpdateRequest
	(*AttachRequest)(nil),  // 5: wave.v1.domain.AttachRequest
	(*DetachRequest)(nil),  // 6: wave.v1.domain.DetachRequest
	(*DeleteRequest)(nil),  // 7: wave.v1.domain.DeleteRequest
	(*GetResponse)(nil),    // 8: wave.v1.domain.GetResponse
	(*StatResponse)(nil),   // 9: wave.v1.domain.StatResponse
	(*ListResponse)(nil),   // 10: wave.v1.domain.ListResponse
	(*CreateResponse)(nil), // 11: wave.v1.domain.CreateResponse
	(*UpdateResponse)(nil), // 12: wave.v1.domain.UpdateResponse
	(*AttachResponse)(nil), // 13: wave.v1.domain.AttachResponse
	(*DetachResponse)(nil), // 14: wave.v1.domain.DetachResponse
	(*DeleteResponse)(nil), // 15: wave.v1.domain.DeleteResponse
}
var file_wave_v1_domain_service_proto_depIdxs = []int32{
	0,  // 0: wave.v1.domain.DomainService.Get:input_type -> wave.v1.domain.GetRequest
	1,  // 1: wave.v1.domain.DomainService.Stat:input_type -> wave.v1.domain.StatRequest
	2,  // 2: wave.v1.domain.DomainService.List:input_type -> wave.v1.domain.ListRequest
	3,  // 3: wave.v1.domain.DomainService.Create:input_type -> wave.v1.domain.CreateRequest
	4,  // 4: wave.v1.domain.DomainService.Update:input_type -> wave.v1.domain.UpdateRequest
	5,  // 5: wave.v1.domain.DomainService.Attach:input_type -> wave.v1.domain.AttachRequest
	6,  // 6: wave.v1.domain.DomainService.Detach:input_type -> wave.v1.domain.DetachRequest
	7,  // 7: wave.v1.domain.DomainService.Delete:input_type -> wave.v1.domain.DeleteRequest
	8,  // 8: wave.v1.domain.DomainService.Get:output_type -> wave.v1.domain.GetResponse
	9,  // 9: wave.v1.domain.DomainService.Stat:output_type -> wave.v1.domain.StatResponse
	10, // 10: wave.v1.domain.DomainService.List:output_type -> wave.v1.domain.ListResponse
	11, // 11: wave.v1.domain.DomainService.Create:output_type -> wave.v1.domain.CreateResponse
	12, // 12: wave.v1.domain.DomainService.Update:output_type -> wave.v1.domain.UpdateResponse
	13, // 13: wave.v1.domain.DomainService.Attach:output_type -> wave.v1.domain.AttachResponse
	14, // 14: wave.v1.domain.DomainService.Detach:output_type -> wave.v1.domain.DetachResponse
	15, // 15: wave.v1.domain.DomainService.Delete:output_type -> wave.v1.domain.DeleteResponse
	8,  // [8:16] is the sub-list for method output_type
	0,  // [0:8] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_wave_v1_domain_service_proto_init() }
func file_wave_v1_domain_service_proto_init() {
	if File_wave_v1_domain_service_proto != nil {
		return
	}
	file_wave_v1_domain_message_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_wave_v1_domain_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_wave_v1_domain_service_proto_goTypes,
		DependencyIndexes: file_wave_v1_domain_service_proto_depIdxs,
	}.Build()
	File_wave_v1_domain_service_proto = out.File
	file_wave_v1_domain_service_proto_rawDesc = nil
	file_wave_v1_domain_service_proto_goTypes = nil
	file_wave_v1_domain_service_proto_depIdxs = nil
}
