// @generated by protoc-gen-es v2.2.5 with parameter "target=ts"
// @generated from file wave/v1/serial/config.proto (package wave.v1.serial, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc } from "@bufbuild/protobuf/codegenv1";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file wave/v1/serial/config.proto.
 */
export const file_wave_v1_serial_config: GenFile = /*@__PURE__*/
  fileDesc("Cht3YXZlL3YxL3NlcmlhbC9jb25maWcucHJvdG8SDndhdmUudjEuc2VyaWFsIhwKDFNlcmlhbENvbmZpZxIMCgRwYXRoGAEgASgJQidaJWN0aHVsLmlvL2N0aHVsL3BrZy9hcGkvd2F2ZS92MS9zZXJpYWxiBnByb3RvMw");

/**
 * @generated from message wave.v1.serial.SerialConfig
 */
export type SerialConfig = Message<"wave.v1.serial.SerialConfig"> & {
  /**
   * @generated from field: string path = 1;
   */
  path: string;
};

/**
 * Describes the message wave.v1.serial.SerialConfig.
 * Use `create(SerialConfigSchema)` to create a new message.
 */
export const SerialConfigSchema: GenMessage<SerialConfig> = /*@__PURE__*/
  messageDesc(file_wave_v1_serial_config, 0);

