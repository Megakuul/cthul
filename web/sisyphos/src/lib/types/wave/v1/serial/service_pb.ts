// @generated by protoc-gen-es v2.2.5 with parameter "target=ts"
// @generated from file wave/v1/serial/service.proto (package wave.v1.serial, syntax proto3)
/* eslint-disable */

import type { GenFile, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { CreateRequestSchema, CreateResponseSchema, DeleteRequestSchema, DeleteResponseSchema, GetRequestSchema, GetResponseSchema, ListRequestSchema, ListResponseSchema, UpdateRequestSchema, UpdateResponseSchema } from "./message_pb";
import { file_wave_v1_serial_message } from "./message_pb";

/**
 * Describes the file wave/v1/serial/service.proto.
 */
export const file_wave_v1_serial_service: GenFile = /*@__PURE__*/
  fileDesc("Chx3YXZlL3YxL3NlcmlhbC9zZXJ2aWNlLnByb3RvEg53YXZlLnYxLnNlcmlhbDL3AgoNU2VyaWFsU2VydmljZRJACgNHZXQSGi53YXZlLnYxLnNlcmlhbC5HZXRSZXF1ZXN0Ghsud2F2ZS52MS5zZXJpYWwuR2V0UmVzcG9uc2UiABJDCgRMaXN0Ehsud2F2ZS52MS5zZXJpYWwuTGlzdFJlcXVlc3QaHC53YXZlLnYxLnNlcmlhbC5MaXN0UmVzcG9uc2UiABJJCgZDcmVhdGUSHS53YXZlLnYxLnNlcmlhbC5DcmVhdGVSZXF1ZXN0Gh4ud2F2ZS52MS5zZXJpYWwuQ3JlYXRlUmVzcG9uc2UiABJJCgZVcGRhdGUSHS53YXZlLnYxLnNlcmlhbC5VcGRhdGVSZXF1ZXN0Gh4ud2F2ZS52MS5zZXJpYWwuVXBkYXRlUmVzcG9uc2UiABJJCgZEZWxldGUSHS53YXZlLnYxLnNlcmlhbC5EZWxldGVSZXF1ZXN0Gh4ud2F2ZS52MS5zZXJpYWwuRGVsZXRlUmVzcG9uc2UiAEInWiVjdGh1bC5pby9jdGh1bC9wa2cvYXBpL3dhdmUvdjEvc2VyaWFsYgZwcm90bzM", [file_wave_v1_serial_message]);

/**
 * @generated from service wave.v1.serial.SerialService
 */
export const SerialService: GenService<{
  /**
   * @generated from rpc wave.v1.serial.SerialService.Get
   */
  get: {
    methodKind: "unary";
    input: typeof GetRequestSchema;
    output: typeof GetResponseSchema;
  },
  /**
   * @generated from rpc wave.v1.serial.SerialService.List
   */
  list: {
    methodKind: "unary";
    input: typeof ListRequestSchema;
    output: typeof ListResponseSchema;
  },
  /**
   * @generated from rpc wave.v1.serial.SerialService.Create
   */
  create: {
    methodKind: "unary";
    input: typeof CreateRequestSchema;
    output: typeof CreateResponseSchema;
  },
  /**
   * @generated from rpc wave.v1.serial.SerialService.Update
   */
  update: {
    methodKind: "unary";
    input: typeof UpdateRequestSchema;
    output: typeof UpdateResponseSchema;
  },
  /**
   * @generated from rpc wave.v1.serial.SerialService.Delete
   */
  delete: {
    methodKind: "unary";
    input: typeof DeleteRequestSchema;
    output: typeof DeleteResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_wave_v1_serial_service, 0);

