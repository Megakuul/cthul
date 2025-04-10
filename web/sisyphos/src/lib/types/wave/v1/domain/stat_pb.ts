// @generated by protoc-gen-es v2.2.5 with parameter "target=ts"
// @generated from file wave/v1/domain/stat.proto (package wave.v1.domain, syntax proto3)
/* eslint-disable */

import type { GenEnum, GenFile, GenMessage } from "@bufbuild/protobuf/codegenv1";
import { enumDesc, fileDesc, messageDesc } from "@bufbuild/protobuf/codegenv1";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file wave/v1/domain/stat.proto.
 */
export const file_wave_v1_domain_stat: GenFile = /*@__PURE__*/
  fileDesc("Chl3YXZlL3YxL2RvbWFpbi9zdGF0LnByb3RvEg53YXZlLnYxLmRvbWFpbiJXCglWQ3B1U3RhdHMSEAoIY3B1X3RpbWUYASABKAMSEQoJd2FpdF90aW1lGAIgASgDEhEKCWhhbHRfdGltZRgDIAEoAxISCgpkZWxheV90aW1lGAQgASgDIoEBCghDcHVTdGF0cxIRCgl0aW1lc3RhbXAYASABKAMSEAoIY3B1X3RpbWUYAiABKAMSEQoJdXNlcl90aW1lGAMgASgDEhMKC2tlcm5lbF90aW1lGAQgASgDEigKBXZjcHVzGAUgAygLMhkud2F2ZS52MS5kb21haW4uVkNwdVN0YXRzIv8BCgtNZW1vcnlTdGF0cxIRCgl0aW1lc3RhbXAYASABKAMSDwoHc3dhcF9pbhgCIAEoAxIQCghzd2FwX291dBgDIAEoAxIUCgxtaW5vcl9mYXVsdHMYBCABKAMSFAoMbWFqb3JfZmF1bHRzGAUgASgDEhwKFGh1Z2VwYWdlX2FsbG9jYXRpb25zGAYgASgDEhkKEWh1Z2VwYWdlX2ZhaWx1cmVzGAcgASgDEhAKCGJhbGxvbmVkGAggASgDEhEKCWF2YWlsYWJsZRgJIAEoAxIOCgZ1c2FibGUYCiABKAMSDgoGdW51c2VkGAsgASgDEhAKCGhvc3RfcnNzGAwgASgDIrsBCgpJbnRlclN0YXRzEhEKCXRpbWVzdGFtcBgBIAEoAxISCgpyZWN2X2J5dGVzGAIgASgDEhEKCXJlY3ZfcGt0cxgDIAEoAxIRCglyZWN2X2VycnMYBCABKAMSEgoKcmVjdl9kcm9wcxgFIAEoAxISCgpzZW5kX2J5dGVzGAYgASgDEhEKCXNlbmRfcGt0cxgHIAEoAxIRCglzZW5kX2VycnMYCCABKAMSEgoKc2VuZF9kcm9wcxgJIAEoAyK9AQoJRGlza1N0YXRzEhEKCXRpbWVzdGFtcBgBIAEoAxIRCglyZWFkX3JlcXMYAiABKAMSEgoKcmVhZF9ieXRlcxgDIAEoAxIRCglyZWFkX3RpbWUYBCABKAMSEgoKd3JpdGVfcmVxcxgFIAEoAxITCgt3cml0ZV9ieXRlcxgGIAEoAxISCgp3cml0ZV90aW1lGAcgASgDEhIKCmZsdXNoX3JlcXMYCCABKAMSEgoKZmx1c2hfdGltZRgJIAEoAyLoAQoLRG9tYWluU3RhdHMSLwoFc3RhdGUYASABKA4yIC53YXZlLnYxLmRvbWFpbi5Eb21haW5Qb3dlclN0YXRlEiUKA2NwdRgCIAEoCzIYLndhdmUudjEuZG9tYWluLkNwdVN0YXRzEisKBm1lbW9yeRgDIAEoCzIbLndhdmUudjEuZG9tYWluLk1lbW9yeVN0YXRzEioKBmludGVycxgEIAMoCzIaLndhdmUudjEuZG9tYWluLkludGVyU3RhdHMSKAoFZGlza3MYBSADKAsyGS53YXZlLnYxLmRvbWFpbi5EaXNrU3RhdHMqtgEKEERvbWFpblBvd2VyU3RhdGUSEgoORE9NQUlOX05PU1RBVEUQABISCg5ET01BSU5fUlVOTklORxABEhIKDkRPTUFJTl9CTE9DS0VEEAISEQoNRE9NQUlOX1BBVVNFRBADEhMKD0RPTUFJTl9TSFVURE9XThAEEhIKDkRPTUFJTl9TSFVUT0ZGEAUSEgoORE9NQUlOX0NSQVNIRUQQBhIWChJET01BSU5fUE1TVVNQRU5ERUQQB0InWiVjdGh1bC5pby9jdGh1bC9wa2cvYXBpL3dhdmUvdjEvZG9tYWluYgZwcm90bzM");

/**
 * @generated from message wave.v1.domain.VCpuStats
 */
export type VCpuStats = Message<"wave.v1.domain.VCpuStats"> & {
  /**
   * full vcpu time since machine startup in nanoseconds.
   *
   * @generated from field: int64 cpu_time = 1;
   */
  cpuTime: bigint;

  /**
   * wait time (time vcpu wants to run but is not scheduled) since machine startup in nanoseconds.
   *
   * @generated from field: int64 wait_time = 2;
   */
  waitTime: bigint;

  /**
   * halt/idle (time the vcpu is disabled or halted) time since machine startup in nanoseconds.
   *
   * @generated from field: int64 halt_time = 3;
   */
  haltTime: bigint;

  /**
   * delay time (time the vcpu is enqueued) since machine startup in nanoseconds.
   *
   * @generated from field: int64 delay_time = 4;
   */
  delayTime: bigint;
};

/**
 * Describes the message wave.v1.domain.VCpuStats.
 * Use `create(VCpuStatsSchema)` to create a new message.
 */
export const VCpuStatsSchema: GenMessage<VCpuStats> = /*@__PURE__*/
  messageDesc(file_wave_v1_domain_stat, 0);

/**
 * @generated from message wave.v1.domain.CpuStats
 */
export type CpuStats = Message<"wave.v1.domain.CpuStats"> & {
  /**
   * @generated from field: int64 timestamp = 1;
   */
  timestamp: bigint;

  /**
   * The total domain cpu time since machine startup in nanoseconds.
   *
   * @generated from field: int64 cpu_time = 2;
   */
  cpuTime: bigint;

  /**
   * The total domain cpu time spent in user space context in nanoseconds.
   *
   * @generated from field: int64 user_time = 3;
   */
  userTime: bigint;

  /**
   * The total domain cpu time spent in kernel space context in nanoseconds.
   *
   * @generated from field: int64 kernel_time = 4;
   */
  kernelTime: bigint;

  /**
   * @generated from field: repeated wave.v1.domain.VCpuStats vcpus = 5;
   */
  vcpus: VCpuStats[];
};

/**
 * Describes the message wave.v1.domain.CpuStats.
 * Use `create(CpuStatsSchema)` to create a new message.
 */
export const CpuStatsSchema: GenMessage<CpuStats> = /*@__PURE__*/
  messageDesc(file_wave_v1_domain_stat, 1);

/**
 * @generated from message wave.v1.domain.MemoryStats
 */
export type MemoryStats = Message<"wave.v1.domain.MemoryStats"> & {
  /**
   * @generated from field: int64 timestamp = 1;
   */
  timestamp: bigint;

  /**
   * total amount of bytes moved from main memory to swap since machine startup.
   *
   * @generated from field: int64 swap_in = 2;
   */
  swapIn: bigint;

  /**
   * total amount of bytes moved from swap to main memory since machine startup.
   *
   * @generated from field: int64 swap_out = 3;
   */
  swapOut: bigint;

  /**
   * total amount of handled page fault that loaded pages from memory (e.g. shared libs).
   *
   * @generated from field: int64 minor_faults = 4;
   */
  minorFaults: bigint;

  /**
   * total amount of handled page fault that loaded pages from diskio
   *
   * @generated from field: int64 major_faults = 5;
   */
  majorFaults: bigint;

  /**
   * total amount of successful hugetable allocations since machine startup.
   *
   * @generated from field: int64 hugepage_allocations = 6;
   */
  hugepageAllocations: bigint;

  /**
   * total amount of failed hugetable allocations since machine startup.
   *
   * @generated from field: int64 hugepage_failures = 7;
   */
  hugepageFailures: bigint;

  /**
   * total amount of bytes currently allocated for the domain by the ballon driver.
   *
   * @generated from field: int64 balloned = 8;
   */
  balloned: bigint;

  /**
   * total amount of bytes from guestos perspective (may be less then the assigned memory).
   *
   * @generated from field: int64 available = 9;
   */
  available: bigint;

  /**
   * total amount of bytes the guest os can use if all caches are reclaimed (deflatable memory).
   *
   * @generated from field: int64 usable = 10;
   */
  usable: bigint;

  /**
   * total amount of bytes that are completely unused by the guest os (part of deflatable memory).
   *
   * @generated from field: int64 unused = 11;
   */
  unused: bigint;

  /**
   * resident set size (used main memory) of the domains process on the host os (typically qemu).
   *
   * @generated from field: int64 host_rss = 12;
   */
  hostRss: bigint;
};

/**
 * Describes the message wave.v1.domain.MemoryStats.
 * Use `create(MemoryStatsSchema)` to create a new message.
 */
export const MemoryStatsSchema: GenMessage<MemoryStats> = /*@__PURE__*/
  messageDesc(file_wave_v1_domain_stat, 2);

/**
 * @generated from message wave.v1.domain.InterStats
 */
export type InterStats = Message<"wave.v1.domain.InterStats"> & {
  /**
   * @generated from field: int64 timestamp = 1;
   */
  timestamp: bigint;

  /**
   * total amount of bytes received since machine startup.
   *
   * @generated from field: int64 recv_bytes = 2;
   */
  recvBytes: bigint;

  /**
   * total number of ethernet packets received since machine startup.
   *
   * @generated from field: int64 recv_pkts = 3;
   */
  recvPkts: bigint;

  /**
   * total number of receive-errors since machine startup.
   *
   * @generated from field: int64 recv_errs = 4;
   */
  recvErrs: bigint;

  /**
   * total number of dropped receive-packets since machine startup.
   *
   * @generated from field: int64 recv_drops = 5;
   */
  recvDrops: bigint;

  /**
   * total amount of bytes sent since machine startup.
   *
   * @generated from field: int64 send_bytes = 6;
   */
  sendBytes: bigint;

  /**
   * total number of ethernet packets sent since machine startup.
   *
   * @generated from field: int64 send_pkts = 7;
   */
  sendPkts: bigint;

  /**
   * total number of send-errors since machine startup.
   *
   * @generated from field: int64 send_errs = 8;
   */
  sendErrs: bigint;

  /**
   * total number of dropped send-packets since machine startup.
   *
   * @generated from field: int64 send_drops = 9;
   */
  sendDrops: bigint;
};

/**
 * Describes the message wave.v1.domain.InterStats.
 * Use `create(InterStatsSchema)` to create a new message.
 */
export const InterStatsSchema: GenMessage<InterStats> = /*@__PURE__*/
  messageDesc(file_wave_v1_domain_stat, 3);

/**
 * @generated from message wave.v1.domain.DiskStats
 */
export type DiskStats = Message<"wave.v1.domain.DiskStats"> & {
  /**
   * @generated from field: int64 timestamp = 1;
   */
  timestamp: bigint;

  /**
   * total number of read requests since machine startup.
   *
   * @generated from field: int64 read_reqs = 2;
   */
  readReqs: bigint;

  /**
   * total amount of bytes read since machine startup.
   *
   * @generated from field: int64 read_bytes = 3;
   */
  readBytes: bigint;

  /**
   * total read time since machine startup in nanoseconds.
   *
   * @generated from field: int64 read_time = 4;
   */
  readTime: bigint;

  /**
   * total number of write requests since machine startup.
   *
   * @generated from field: int64 write_reqs = 5;
   */
  writeReqs: bigint;

  /**
   * total amount of bytes write since machine startup.
   *
   * @generated from field: int64 write_bytes = 6;
   */
  writeBytes: bigint;

  /**
   * total write time since machine startup in nanoseconds.
   *
   * @generated from field: int64 write_time = 7;
   */
  writeTime: bigint;

  /**
   * total number of flush requests since machine startup.
   *
   * @generated from field: int64 flush_reqs = 8;
   */
  flushReqs: bigint;

  /**
   * total flush time since machine startup in nanoseconds.
   *
   * @generated from field: int64 flush_time = 9;
   */
  flushTime: bigint;
};

/**
 * Describes the message wave.v1.domain.DiskStats.
 * Use `create(DiskStatsSchema)` to create a new message.
 */
export const DiskStatsSchema: GenMessage<DiskStats> = /*@__PURE__*/
  messageDesc(file_wave_v1_domain_stat, 4);

/**
 * @generated from message wave.v1.domain.DomainStats
 */
export type DomainStats = Message<"wave.v1.domain.DomainStats"> & {
  /**
   * @generated from field: wave.v1.domain.DomainPowerState state = 1;
   */
  state: DomainPowerState;

  /**
   * @generated from field: wave.v1.domain.CpuStats cpu = 2;
   */
  cpu?: CpuStats;

  /**
   * @generated from field: wave.v1.domain.MemoryStats memory = 3;
   */
  memory?: MemoryStats;

  /**
   * @generated from field: repeated wave.v1.domain.InterStats inters = 4;
   */
  inters: InterStats[];

  /**
   * @generated from field: repeated wave.v1.domain.DiskStats disks = 5;
   */
  disks: DiskStats[];
};

/**
 * Describes the message wave.v1.domain.DomainStats.
 * Use `create(DomainStatsSchema)` to create a new message.
 */
export const DomainStatsSchema: GenMessage<DomainStats> = /*@__PURE__*/
  messageDesc(file_wave_v1_domain_stat, 5);

/**
 * @generated from enum wave.v1.domain.DomainPowerState
 */
export enum DomainPowerState {
  /**
   * @generated from enum value: DOMAIN_NOSTATE = 0;
   */
  DOMAIN_NOSTATE = 0,

  /**
   * @generated from enum value: DOMAIN_RUNNING = 1;
   */
  DOMAIN_RUNNING = 1,

  /**
   * @generated from enum value: DOMAIN_BLOCKED = 2;
   */
  DOMAIN_BLOCKED = 2,

  /**
   * @generated from enum value: DOMAIN_PAUSED = 3;
   */
  DOMAIN_PAUSED = 3,

  /**
   * @generated from enum value: DOMAIN_SHUTDOWN = 4;
   */
  DOMAIN_SHUTDOWN = 4,

  /**
   * @generated from enum value: DOMAIN_SHUTOFF = 5;
   */
  DOMAIN_SHUTOFF = 5,

  /**
   * @generated from enum value: DOMAIN_CRASHED = 6;
   */
  DOMAIN_CRASHED = 6,

  /**
   * @generated from enum value: DOMAIN_PMSUSPENDED = 7;
   */
  DOMAIN_PMSUSPENDED = 7,
}

/**
 * Describes the enum wave.v1.domain.DomainPowerState.
 */
export const DomainPowerStateSchema: GenEnum<DomainPowerState> = /*@__PURE__*/
  enumDesc(file_wave_v1_domain_stat, 0);

