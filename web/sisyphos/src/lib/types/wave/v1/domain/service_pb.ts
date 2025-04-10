// @generated by protoc-gen-es v2.2.5 with parameter "target=ts"
// @generated from file wave/v1/domain/service.proto (package wave.v1.domain, syntax proto3)
/* eslint-disable */

import type { GenFile, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { AttachRequestSchema, AttachResponseSchema, CreateRequestSchema, CreateResponseSchema, DeleteRequestSchema, DeleteResponseSchema, DetachRequestSchema, DetachResponseSchema, GetRequestSchema, GetResponseSchema, ListRequestSchema, ListResponseSchema, StatRequestSchema, StatResponseSchema, UpdateRequestSchema, UpdateResponseSchema } from "./message_pb";
import { file_wave_v1_domain_message } from "./message_pb";

/**
 * Describes the file wave/v1/domain/service.proto.
 */
export const file_wave_v1_domain_service: GenFile = /*@__PURE__*/
  fileDesc("Chx3YXZlL3YxL2RvbWFpbi9zZXJ2aWNlLnByb3RvEg53YXZlLnYxLmRvbWFpbjLSBAoNRG9tYWluU2VydmljZRJACgNHZXQSGi53YXZlLnYxLmRvbWFpbi5HZXRSZXF1ZXN0Ghsud2F2ZS52MS5kb21haW4uR2V0UmVzcG9uc2UiABJDCgRTdGF0Ehsud2F2ZS52MS5kb21haW4uU3RhdFJlcXVlc3QaHC53YXZlLnYxLmRvbWFpbi5TdGF0UmVzcG9uc2UiABJDCgRMaXN0Ehsud2F2ZS52MS5kb21haW4uTGlzdFJlcXVlc3QaHC53YXZlLnYxLmRvbWFpbi5MaXN0UmVzcG9uc2UiABJJCgZDcmVhdGUSHS53YXZlLnYxLmRvbWFpbi5DcmVhdGVSZXF1ZXN0Gh4ud2F2ZS52MS5kb21haW4uQ3JlYXRlUmVzcG9uc2UiABJJCgZVcGRhdGUSHS53YXZlLnYxLmRvbWFpbi5VcGRhdGVSZXF1ZXN0Gh4ud2F2ZS52MS5kb21haW4uVXBkYXRlUmVzcG9uc2UiABJJCgZBdHRhY2gSHS53YXZlLnYxLmRvbWFpbi5BdHRhY2hSZXF1ZXN0Gh4ud2F2ZS52MS5kb21haW4uQXR0YWNoUmVzcG9uc2UiABJJCgZEZXRhY2gSHS53YXZlLnYxLmRvbWFpbi5EZXRhY2hSZXF1ZXN0Gh4ud2F2ZS52MS5kb21haW4uRGV0YWNoUmVzcG9uc2UiABJJCgZEZWxldGUSHS53YXZlLnYxLmRvbWFpbi5EZWxldGVSZXF1ZXN0Gh4ud2F2ZS52MS5kb21haW4uRGVsZXRlUmVzcG9uc2UiAEInWiVjdGh1bC5pby9jdGh1bC9wa2cvYXBpL3dhdmUvdjEvZG9tYWluYgZwcm90bzM", [file_wave_v1_domain_message]);

/**
 * @generated from service wave.v1.domain.DomainService
 */
export const DomainService: GenService<{
  /**
   * @generated from rpc wave.v1.domain.DomainService.Get
   */
  get: {
    methodKind: "unary";
    input: typeof GetRequestSchema;
    output: typeof GetResponseSchema;
  },
  /**
   * @generated from rpc wave.v1.domain.DomainService.Stat
   */
  stat: {
    methodKind: "unary";
    input: typeof StatRequestSchema;
    output: typeof StatResponseSchema;
  },
  /**
   * @generated from rpc wave.v1.domain.DomainService.List
   */
  list: {
    methodKind: "unary";
    input: typeof ListRequestSchema;
    output: typeof ListResponseSchema;
  },
  /**
   * @generated from rpc wave.v1.domain.DomainService.Create
   */
  create: {
    methodKind: "unary";
    input: typeof CreateRequestSchema;
    output: typeof CreateResponseSchema;
  },
  /**
   * @generated from rpc wave.v1.domain.DomainService.Update
   */
  update: {
    methodKind: "unary";
    input: typeof UpdateRequestSchema;
    output: typeof UpdateResponseSchema;
  },
  /**
   * @generated from rpc wave.v1.domain.DomainService.Attach
   */
  attach: {
    methodKind: "unary";
    input: typeof AttachRequestSchema;
    output: typeof AttachResponseSchema;
  },
  /**
   * @generated from rpc wave.v1.domain.DomainService.Detach
   */
  detach: {
    methodKind: "unary";
    input: typeof DetachRequestSchema;
    output: typeof DetachResponseSchema;
  },
  /**
   * @generated from rpc wave.v1.domain.DomainService.Delete
   */
  delete: {
    methodKind: "unary";
    input: typeof DeleteRequestSchema;
    output: typeof DeleteResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_wave_v1_domain_service, 0);

