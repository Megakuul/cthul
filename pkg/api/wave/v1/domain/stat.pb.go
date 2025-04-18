// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: wave/v1/domain/stat.proto

package domain

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DomainPowerState int32

const (
	DomainPowerState_DOMAIN_NOSTATE     DomainPowerState = 0
	DomainPowerState_DOMAIN_RUNNING     DomainPowerState = 1
	DomainPowerState_DOMAIN_BLOCKED     DomainPowerState = 2
	DomainPowerState_DOMAIN_PAUSED      DomainPowerState = 3
	DomainPowerState_DOMAIN_SHUTDOWN    DomainPowerState = 4
	DomainPowerState_DOMAIN_SHUTOFF     DomainPowerState = 5
	DomainPowerState_DOMAIN_CRASHED     DomainPowerState = 6
	DomainPowerState_DOMAIN_PMSUSPENDED DomainPowerState = 7
)

// Enum value maps for DomainPowerState.
var (
	DomainPowerState_name = map[int32]string{
		0: "DOMAIN_NOSTATE",
		1: "DOMAIN_RUNNING",
		2: "DOMAIN_BLOCKED",
		3: "DOMAIN_PAUSED",
		4: "DOMAIN_SHUTDOWN",
		5: "DOMAIN_SHUTOFF",
		6: "DOMAIN_CRASHED",
		7: "DOMAIN_PMSUSPENDED",
	}
	DomainPowerState_value = map[string]int32{
		"DOMAIN_NOSTATE":     0,
		"DOMAIN_RUNNING":     1,
		"DOMAIN_BLOCKED":     2,
		"DOMAIN_PAUSED":      3,
		"DOMAIN_SHUTDOWN":    4,
		"DOMAIN_SHUTOFF":     5,
		"DOMAIN_CRASHED":     6,
		"DOMAIN_PMSUSPENDED": 7,
	}
)

func (x DomainPowerState) Enum() *DomainPowerState {
	p := new(DomainPowerState)
	*p = x
	return p
}

func (x DomainPowerState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DomainPowerState) Descriptor() protoreflect.EnumDescriptor {
	return file_wave_v1_domain_stat_proto_enumTypes[0].Descriptor()
}

func (DomainPowerState) Type() protoreflect.EnumType {
	return &file_wave_v1_domain_stat_proto_enumTypes[0]
}

func (x DomainPowerState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DomainPowerState.Descriptor instead.
func (DomainPowerState) EnumDescriptor() ([]byte, []int) {
	return file_wave_v1_domain_stat_proto_rawDescGZIP(), []int{0}
}

type VCpuStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// full vcpu time since machine startup in nanoseconds.
	CpuTime int64 `protobuf:"varint,1,opt,name=cpu_time,json=cpuTime,proto3" json:"cpu_time,omitempty"`
	// wait time (time vcpu wants to run but is not scheduled) since machine startup in nanoseconds.
	WaitTime int64 `protobuf:"varint,2,opt,name=wait_time,json=waitTime,proto3" json:"wait_time,omitempty"`
	// halt/idle (time the vcpu is disabled or halted) time since machine startup in nanoseconds.
	HaltTime int64 `protobuf:"varint,3,opt,name=halt_time,json=haltTime,proto3" json:"halt_time,omitempty"`
	// delay time (time the vcpu is enqueued) since machine startup in nanoseconds.
	DelayTime int64 `protobuf:"varint,4,opt,name=delay_time,json=delayTime,proto3" json:"delay_time,omitempty"`
}

func (x *VCpuStats) Reset() {
	*x = VCpuStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wave_v1_domain_stat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VCpuStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VCpuStats) ProtoMessage() {}

func (x *VCpuStats) ProtoReflect() protoreflect.Message {
	mi := &file_wave_v1_domain_stat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VCpuStats.ProtoReflect.Descriptor instead.
func (*VCpuStats) Descriptor() ([]byte, []int) {
	return file_wave_v1_domain_stat_proto_rawDescGZIP(), []int{0}
}

func (x *VCpuStats) GetCpuTime() int64 {
	if x != nil {
		return x.CpuTime
	}
	return 0
}

func (x *VCpuStats) GetWaitTime() int64 {
	if x != nil {
		return x.WaitTime
	}
	return 0
}

func (x *VCpuStats) GetHaltTime() int64 {
	if x != nil {
		return x.HaltTime
	}
	return 0
}

func (x *VCpuStats) GetDelayTime() int64 {
	if x != nil {
		return x.DelayTime
	}
	return 0
}

type CpuStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp int64 `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// The total domain cpu time since machine startup in nanoseconds.
	CpuTime int64 `protobuf:"varint,2,opt,name=cpu_time,json=cpuTime,proto3" json:"cpu_time,omitempty"`
	// The total domain cpu time spent in user space context in nanoseconds.
	UserTime int64 `protobuf:"varint,3,opt,name=user_time,json=userTime,proto3" json:"user_time,omitempty"`
	// The total domain cpu time spent in kernel space context in nanoseconds.
	KernelTime int64        `protobuf:"varint,4,opt,name=kernel_time,json=kernelTime,proto3" json:"kernel_time,omitempty"`
	Vcpus      []*VCpuStats `protobuf:"bytes,5,rep,name=vcpus,proto3" json:"vcpus,omitempty"`
}

func (x *CpuStats) Reset() {
	*x = CpuStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wave_v1_domain_stat_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CpuStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CpuStats) ProtoMessage() {}

func (x *CpuStats) ProtoReflect() protoreflect.Message {
	mi := &file_wave_v1_domain_stat_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CpuStats.ProtoReflect.Descriptor instead.
func (*CpuStats) Descriptor() ([]byte, []int) {
	return file_wave_v1_domain_stat_proto_rawDescGZIP(), []int{1}
}

func (x *CpuStats) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *CpuStats) GetCpuTime() int64 {
	if x != nil {
		return x.CpuTime
	}
	return 0
}

func (x *CpuStats) GetUserTime() int64 {
	if x != nil {
		return x.UserTime
	}
	return 0
}

func (x *CpuStats) GetKernelTime() int64 {
	if x != nil {
		return x.KernelTime
	}
	return 0
}

func (x *CpuStats) GetVcpus() []*VCpuStats {
	if x != nil {
		return x.Vcpus
	}
	return nil
}

type MemoryStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp int64 `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// total amount of bytes moved from main memory to swap since machine startup.
	SwapIn int64 `protobuf:"varint,2,opt,name=swap_in,json=swapIn,proto3" json:"swap_in,omitempty"`
	// total amount of bytes moved from swap to main memory since machine startup.
	SwapOut int64 `protobuf:"varint,3,opt,name=swap_out,json=swapOut,proto3" json:"swap_out,omitempty"`
	// total amount of handled page fault that loaded pages from memory (e.g. shared libs).
	MinorFaults int64 `protobuf:"varint,4,opt,name=minor_faults,json=minorFaults,proto3" json:"minor_faults,omitempty"`
	// total amount of handled page fault that loaded pages from diskio
	MajorFaults int64 `protobuf:"varint,5,opt,name=major_faults,json=majorFaults,proto3" json:"major_faults,omitempty"`
	// total amount of successful hugetable allocations since machine startup.
	HugepageAllocations int64 `protobuf:"varint,6,opt,name=hugepage_allocations,json=hugepageAllocations,proto3" json:"hugepage_allocations,omitempty"`
	// total amount of failed hugetable allocations since machine startup.
	HugepageFailures int64 `protobuf:"varint,7,opt,name=hugepage_failures,json=hugepageFailures,proto3" json:"hugepage_failures,omitempty"`
	// total amount of bytes currently allocated for the domain by the ballon driver.
	Balloned int64 `protobuf:"varint,8,opt,name=balloned,proto3" json:"balloned,omitempty"`
	// total amount of bytes from guestos perspective (may be less then the assigned memory).
	Available int64 `protobuf:"varint,9,opt,name=available,proto3" json:"available,omitempty"`
	// total amount of bytes the guest os can use if all caches are reclaimed (deflatable memory).
	Usable int64 `protobuf:"varint,10,opt,name=usable,proto3" json:"usable,omitempty"`
	// total amount of bytes that are completely unused by the guest os (part of deflatable memory).
	Unused int64 `protobuf:"varint,11,opt,name=unused,proto3" json:"unused,omitempty"`
	// resident set size (used main memory) of the domains process on the host os (typically qemu).
	HostRss int64 `protobuf:"varint,12,opt,name=host_rss,json=hostRss,proto3" json:"host_rss,omitempty"`
}

func (x *MemoryStats) Reset() {
	*x = MemoryStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wave_v1_domain_stat_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemoryStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemoryStats) ProtoMessage() {}

func (x *MemoryStats) ProtoReflect() protoreflect.Message {
	mi := &file_wave_v1_domain_stat_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemoryStats.ProtoReflect.Descriptor instead.
func (*MemoryStats) Descriptor() ([]byte, []int) {
	return file_wave_v1_domain_stat_proto_rawDescGZIP(), []int{2}
}

func (x *MemoryStats) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *MemoryStats) GetSwapIn() int64 {
	if x != nil {
		return x.SwapIn
	}
	return 0
}

func (x *MemoryStats) GetSwapOut() int64 {
	if x != nil {
		return x.SwapOut
	}
	return 0
}

func (x *MemoryStats) GetMinorFaults() int64 {
	if x != nil {
		return x.MinorFaults
	}
	return 0
}

func (x *MemoryStats) GetMajorFaults() int64 {
	if x != nil {
		return x.MajorFaults
	}
	return 0
}

func (x *MemoryStats) GetHugepageAllocations() int64 {
	if x != nil {
		return x.HugepageAllocations
	}
	return 0
}

func (x *MemoryStats) GetHugepageFailures() int64 {
	if x != nil {
		return x.HugepageFailures
	}
	return 0
}

func (x *MemoryStats) GetBalloned() int64 {
	if x != nil {
		return x.Balloned
	}
	return 0
}

func (x *MemoryStats) GetAvailable() int64 {
	if x != nil {
		return x.Available
	}
	return 0
}

func (x *MemoryStats) GetUsable() int64 {
	if x != nil {
		return x.Usable
	}
	return 0
}

func (x *MemoryStats) GetUnused() int64 {
	if x != nil {
		return x.Unused
	}
	return 0
}

func (x *MemoryStats) GetHostRss() int64 {
	if x != nil {
		return x.HostRss
	}
	return 0
}

type InterStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp int64 `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// total amount of bytes received since machine startup.
	RecvBytes int64 `protobuf:"varint,2,opt,name=recv_bytes,json=recvBytes,proto3" json:"recv_bytes,omitempty"`
	// total number of ethernet packets received since machine startup.
	RecvPkts int64 `protobuf:"varint,3,opt,name=recv_pkts,json=recvPkts,proto3" json:"recv_pkts,omitempty"`
	// total number of receive-errors since machine startup.
	RecvErrs int64 `protobuf:"varint,4,opt,name=recv_errs,json=recvErrs,proto3" json:"recv_errs,omitempty"`
	// total number of dropped receive-packets since machine startup.
	RecvDrops int64 `protobuf:"varint,5,opt,name=recv_drops,json=recvDrops,proto3" json:"recv_drops,omitempty"`
	// total amount of bytes sent since machine startup.
	SendBytes int64 `protobuf:"varint,6,opt,name=send_bytes,json=sendBytes,proto3" json:"send_bytes,omitempty"`
	// total number of ethernet packets sent since machine startup.
	SendPkts int64 `protobuf:"varint,7,opt,name=send_pkts,json=sendPkts,proto3" json:"send_pkts,omitempty"`
	// total number of send-errors since machine startup.
	SendErrs int64 `protobuf:"varint,8,opt,name=send_errs,json=sendErrs,proto3" json:"send_errs,omitempty"`
	// total number of dropped send-packets since machine startup.
	SendDrops int64 `protobuf:"varint,9,opt,name=send_drops,json=sendDrops,proto3" json:"send_drops,omitempty"`
}

func (x *InterStats) Reset() {
	*x = InterStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wave_v1_domain_stat_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InterStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InterStats) ProtoMessage() {}

func (x *InterStats) ProtoReflect() protoreflect.Message {
	mi := &file_wave_v1_domain_stat_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InterStats.ProtoReflect.Descriptor instead.
func (*InterStats) Descriptor() ([]byte, []int) {
	return file_wave_v1_domain_stat_proto_rawDescGZIP(), []int{3}
}

func (x *InterStats) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *InterStats) GetRecvBytes() int64 {
	if x != nil {
		return x.RecvBytes
	}
	return 0
}

func (x *InterStats) GetRecvPkts() int64 {
	if x != nil {
		return x.RecvPkts
	}
	return 0
}

func (x *InterStats) GetRecvErrs() int64 {
	if x != nil {
		return x.RecvErrs
	}
	return 0
}

func (x *InterStats) GetRecvDrops() int64 {
	if x != nil {
		return x.RecvDrops
	}
	return 0
}

func (x *InterStats) GetSendBytes() int64 {
	if x != nil {
		return x.SendBytes
	}
	return 0
}

func (x *InterStats) GetSendPkts() int64 {
	if x != nil {
		return x.SendPkts
	}
	return 0
}

func (x *InterStats) GetSendErrs() int64 {
	if x != nil {
		return x.SendErrs
	}
	return 0
}

func (x *InterStats) GetSendDrops() int64 {
	if x != nil {
		return x.SendDrops
	}
	return 0
}

type DiskStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp int64 `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// total number of read requests since machine startup.
	ReadReqs int64 `protobuf:"varint,2,opt,name=read_reqs,json=readReqs,proto3" json:"read_reqs,omitempty"`
	// total amount of bytes read since machine startup.
	ReadBytes int64 `protobuf:"varint,3,opt,name=read_bytes,json=readBytes,proto3" json:"read_bytes,omitempty"`
	// total read time since machine startup in nanoseconds.
	ReadTime int64 `protobuf:"varint,4,opt,name=read_time,json=readTime,proto3" json:"read_time,omitempty"`
	// total number of write requests since machine startup.
	WriteReqs int64 `protobuf:"varint,5,opt,name=write_reqs,json=writeReqs,proto3" json:"write_reqs,omitempty"`
	// total amount of bytes write since machine startup.
	WriteBytes int64 `protobuf:"varint,6,opt,name=write_bytes,json=writeBytes,proto3" json:"write_bytes,omitempty"`
	// total write time since machine startup in nanoseconds.
	WriteTime int64 `protobuf:"varint,7,opt,name=write_time,json=writeTime,proto3" json:"write_time,omitempty"`
	// total number of flush requests since machine startup.
	FlushReqs int64 `protobuf:"varint,8,opt,name=flush_reqs,json=flushReqs,proto3" json:"flush_reqs,omitempty"`
	// total flush time since machine startup in nanoseconds.
	FlushTime int64 `protobuf:"varint,9,opt,name=flush_time,json=flushTime,proto3" json:"flush_time,omitempty"`
}

func (x *DiskStats) Reset() {
	*x = DiskStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wave_v1_domain_stat_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiskStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiskStats) ProtoMessage() {}

func (x *DiskStats) ProtoReflect() protoreflect.Message {
	mi := &file_wave_v1_domain_stat_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiskStats.ProtoReflect.Descriptor instead.
func (*DiskStats) Descriptor() ([]byte, []int) {
	return file_wave_v1_domain_stat_proto_rawDescGZIP(), []int{4}
}

func (x *DiskStats) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *DiskStats) GetReadReqs() int64 {
	if x != nil {
		return x.ReadReqs
	}
	return 0
}

func (x *DiskStats) GetReadBytes() int64 {
	if x != nil {
		return x.ReadBytes
	}
	return 0
}

func (x *DiskStats) GetReadTime() int64 {
	if x != nil {
		return x.ReadTime
	}
	return 0
}

func (x *DiskStats) GetWriteReqs() int64 {
	if x != nil {
		return x.WriteReqs
	}
	return 0
}

func (x *DiskStats) GetWriteBytes() int64 {
	if x != nil {
		return x.WriteBytes
	}
	return 0
}

func (x *DiskStats) GetWriteTime() int64 {
	if x != nil {
		return x.WriteTime
	}
	return 0
}

func (x *DiskStats) GetFlushReqs() int64 {
	if x != nil {
		return x.FlushReqs
	}
	return 0
}

func (x *DiskStats) GetFlushTime() int64 {
	if x != nil {
		return x.FlushTime
	}
	return 0
}

type DomainStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State  DomainPowerState `protobuf:"varint,1,opt,name=state,proto3,enum=wave.v1.domain.DomainPowerState" json:"state,omitempty"`
	Cpu    *CpuStats        `protobuf:"bytes,2,opt,name=cpu,proto3" json:"cpu,omitempty"`
	Memory *MemoryStats     `protobuf:"bytes,3,opt,name=memory,proto3" json:"memory,omitempty"`
	Inters []*InterStats    `protobuf:"bytes,4,rep,name=inters,proto3" json:"inters,omitempty"`
	Disks  []*DiskStats     `protobuf:"bytes,5,rep,name=disks,proto3" json:"disks,omitempty"`
}

func (x *DomainStats) Reset() {
	*x = DomainStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wave_v1_domain_stat_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DomainStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DomainStats) ProtoMessage() {}

func (x *DomainStats) ProtoReflect() protoreflect.Message {
	mi := &file_wave_v1_domain_stat_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DomainStats.ProtoReflect.Descriptor instead.
func (*DomainStats) Descriptor() ([]byte, []int) {
	return file_wave_v1_domain_stat_proto_rawDescGZIP(), []int{5}
}

func (x *DomainStats) GetState() DomainPowerState {
	if x != nil {
		return x.State
	}
	return DomainPowerState_DOMAIN_NOSTATE
}

func (x *DomainStats) GetCpu() *CpuStats {
	if x != nil {
		return x.Cpu
	}
	return nil
}

func (x *DomainStats) GetMemory() *MemoryStats {
	if x != nil {
		return x.Memory
	}
	return nil
}

func (x *DomainStats) GetInters() []*InterStats {
	if x != nil {
		return x.Inters
	}
	return nil
}

func (x *DomainStats) GetDisks() []*DiskStats {
	if x != nil {
		return x.Disks
	}
	return nil
}

var File_wave_v1_domain_stat_proto protoreflect.FileDescriptor

var file_wave_v1_domain_stat_proto_rawDesc = []byte{
	0x0a, 0x19, 0x77, 0x61, 0x76, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x2f, 0x73, 0x74, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x77, 0x61, 0x76,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x7f, 0x0a, 0x09, 0x56,
	0x43, 0x70, 0x75, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x63, 0x70, 0x75, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x63, 0x70, 0x75, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x61, 0x69, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x77, 0x61, 0x69, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x68, 0x61, 0x6c, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x08, 0x68, 0x61, 0x6c, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1d, 0x0a,
	0x0a, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x22, 0xb2, 0x01, 0x0a,
	0x08, 0x43, 0x70, 0x75, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x19, 0x0a, 0x08, 0x63, 0x70, 0x75, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x63, 0x70, 0x75, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x1f, 0x0a, 0x0b, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x2f, 0x0a, 0x05, 0x76, 0x63, 0x70, 0x75, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x2e, 0x56, 0x43, 0x70, 0x75, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x05, 0x76, 0x63, 0x70, 0x75,
	0x73, 0x22, 0x8a, 0x03, 0x0a, 0x0b, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x53, 0x74, 0x61, 0x74,
	0x73, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12,
	0x17, 0x0a, 0x07, 0x73, 0x77, 0x61, 0x70, 0x5f, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x73, 0x77, 0x61, 0x70, 0x49, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x77, 0x61, 0x70,
	0x5f, 0x6f, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x73, 0x77, 0x61, 0x70,
	0x4f, 0x75, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x5f, 0x66, 0x61, 0x75,
	0x6c, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x6d, 0x69, 0x6e, 0x6f, 0x72,
	0x46, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x61, 0x6a, 0x6f, 0x72, 0x5f,
	0x66, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x6d, 0x61,
	0x6a, 0x6f, 0x72, 0x46, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x31, 0x0a, 0x14, 0x68, 0x75, 0x67,
	0x65, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x61, 0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x13, 0x68, 0x75, 0x67, 0x65, 0x70, 0x61, 0x67,
	0x65, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x2b, 0x0a, 0x11,
	0x68, 0x75, 0x67, 0x65, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65,
	0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x68, 0x75, 0x67, 0x65, 0x70, 0x61, 0x67,
	0x65, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x61, 0x6c,
	0x6c, 0x6f, 0x6e, 0x65, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x62, 0x61, 0x6c,
	0x6c, 0x6f, 0x6e, 0x65, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62,
	0x6c, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61,
	0x62, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x75,
	0x6e, 0x75, 0x73, 0x65, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x6e, 0x75,
	0x73, 0x65, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x72, 0x73, 0x73, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x68, 0x6f, 0x73, 0x74, 0x52, 0x73, 0x73, 0x22, 0x9a,
	0x02, 0x0a, 0x0a, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x1c, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1d, 0x0a, 0x0a, 0x72,
	0x65, 0x63, 0x76, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x72, 0x65, 0x63, 0x76, 0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65,
	0x63, 0x76, 0x5f, 0x70, 0x6b, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x72,
	0x65, 0x63, 0x76, 0x50, 0x6b, 0x74, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x63, 0x76, 0x5f,
	0x65, 0x72, 0x72, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x72, 0x65, 0x63, 0x76,
	0x45, 0x72, 0x72, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x63, 0x76, 0x5f, 0x64, 0x72, 0x6f,
	0x70, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x72, 0x65, 0x63, 0x76, 0x44, 0x72,
	0x6f, 0x70, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x62, 0x79, 0x74, 0x65,
	0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x42, 0x79, 0x74,
	0x65, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x70, 0x6b, 0x74, 0x73, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x50, 0x6b, 0x74, 0x73, 0x12,
	0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x65, 0x72, 0x72, 0x73, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x45, 0x72, 0x72, 0x73, 0x12, 0x1d, 0x0a, 0x0a,
	0x73, 0x65, 0x6e, 0x64, 0x5f, 0x64, 0x72, 0x6f, 0x70, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x44, 0x72, 0x6f, 0x70, 0x73, 0x22, 0x9f, 0x02, 0x0a, 0x09,
	0x44, 0x69, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x61, 0x64, 0x5f,
	0x72, 0x65, 0x71, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x72, 0x65, 0x61, 0x64,
	0x52, 0x65, 0x71, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x62, 0x79, 0x74,
	0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x72, 0x65, 0x61, 0x64, 0x42, 0x79,
	0x74, 0x65, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x72, 0x65, 0x61, 0x64, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x1d, 0x0a, 0x0a, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x71, 0x73, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x77, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x73, 0x12,
	0x1f, 0x0a, 0x0b, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x77, 0x72, 0x69, 0x74, 0x65, 0x42, 0x79, 0x74, 0x65, 0x73,
	0x12, 0x1d, 0x0a, 0x0a, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x77, 0x72, 0x69, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x1d, 0x0a, 0x0a, 0x66, 0x6c, 0x75, 0x73, 0x68, 0x5f, 0x72, 0x65, 0x71, 0x73, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x66, 0x6c, 0x75, 0x73, 0x68, 0x52, 0x65, 0x71, 0x73, 0x12, 0x1d,
	0x0a, 0x0a, 0x66, 0x6c, 0x75, 0x73, 0x68, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x66, 0x6c, 0x75, 0x73, 0x68, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x8b, 0x02,
	0x0a, 0x0b, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x36, 0x0a,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x77,
	0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x44, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x2a, 0x0a, 0x03, 0x63, 0x70, 0x75, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x18, 0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x2e, 0x43, 0x70, 0x75, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x03, 0x63, 0x70,
	0x75, 0x12, 0x33, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1b, 0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x2e, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x06,
	0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x12, 0x32, 0x0a, 0x06, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x73,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x77, 0x61, 0x76, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61,
	0x74, 0x73, 0x52, 0x06, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x73, 0x12, 0x2f, 0x0a, 0x05, 0x64, 0x69,
	0x73, 0x6b, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x77, 0x61, 0x76, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x44, 0x69, 0x73, 0x6b, 0x53,
	0x74, 0x61, 0x74, 0x73, 0x52, 0x05, 0x64, 0x69, 0x73, 0x6b, 0x73, 0x2a, 0xb6, 0x01, 0x0a, 0x10,
	0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x12, 0x0a, 0x0e, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x4e, 0x4f, 0x53, 0x54, 0x41,
	0x54, 0x45, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x52,
	0x55, 0x4e, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x44, 0x4f, 0x4d, 0x41,
	0x49, 0x4e, 0x5f, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x45, 0x44, 0x10, 0x02, 0x12, 0x11, 0x0a, 0x0d,
	0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x50, 0x41, 0x55, 0x53, 0x45, 0x44, 0x10, 0x03, 0x12,
	0x13, 0x0a, 0x0f, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x53, 0x48, 0x55, 0x54, 0x44, 0x4f,
	0x57, 0x4e, 0x10, 0x04, 0x12, 0x12, 0x0a, 0x0e, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x53,
	0x48, 0x55, 0x54, 0x4f, 0x46, 0x46, 0x10, 0x05, 0x12, 0x12, 0x0a, 0x0e, 0x44, 0x4f, 0x4d, 0x41,
	0x49, 0x4e, 0x5f, 0x43, 0x52, 0x41, 0x53, 0x48, 0x45, 0x44, 0x10, 0x06, 0x12, 0x16, 0x0a, 0x12,
	0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x50, 0x4d, 0x53, 0x55, 0x53, 0x50, 0x45, 0x4e, 0x44,
	0x45, 0x44, 0x10, 0x07, 0x42, 0x27, 0x5a, 0x25, 0x63, 0x74, 0x68, 0x75, 0x6c, 0x2e, 0x69, 0x6f,
	0x2f, 0x63, 0x74, 0x68, 0x75, 0x6c, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x77,
	0x61, 0x76, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_wave_v1_domain_stat_proto_rawDescOnce sync.Once
	file_wave_v1_domain_stat_proto_rawDescData = file_wave_v1_domain_stat_proto_rawDesc
)

func file_wave_v1_domain_stat_proto_rawDescGZIP() []byte {
	file_wave_v1_domain_stat_proto_rawDescOnce.Do(func() {
		file_wave_v1_domain_stat_proto_rawDescData = protoimpl.X.CompressGZIP(file_wave_v1_domain_stat_proto_rawDescData)
	})
	return file_wave_v1_domain_stat_proto_rawDescData
}

var file_wave_v1_domain_stat_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_wave_v1_domain_stat_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_wave_v1_domain_stat_proto_goTypes = []any{
	(DomainPowerState)(0), // 0: wave.v1.domain.DomainPowerState
	(*VCpuStats)(nil),     // 1: wave.v1.domain.VCpuStats
	(*CpuStats)(nil),      // 2: wave.v1.domain.CpuStats
	(*MemoryStats)(nil),   // 3: wave.v1.domain.MemoryStats
	(*InterStats)(nil),    // 4: wave.v1.domain.InterStats
	(*DiskStats)(nil),     // 5: wave.v1.domain.DiskStats
	(*DomainStats)(nil),   // 6: wave.v1.domain.DomainStats
}
var file_wave_v1_domain_stat_proto_depIdxs = []int32{
	1, // 0: wave.v1.domain.CpuStats.vcpus:type_name -> wave.v1.domain.VCpuStats
	0, // 1: wave.v1.domain.DomainStats.state:type_name -> wave.v1.domain.DomainPowerState
	2, // 2: wave.v1.domain.DomainStats.cpu:type_name -> wave.v1.domain.CpuStats
	3, // 3: wave.v1.domain.DomainStats.memory:type_name -> wave.v1.domain.MemoryStats
	4, // 4: wave.v1.domain.DomainStats.inters:type_name -> wave.v1.domain.InterStats
	5, // 5: wave.v1.domain.DomainStats.disks:type_name -> wave.v1.domain.DiskStats
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_wave_v1_domain_stat_proto_init() }
func file_wave_v1_domain_stat_proto_init() {
	if File_wave_v1_domain_stat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_wave_v1_domain_stat_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*VCpuStats); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wave_v1_domain_stat_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CpuStats); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wave_v1_domain_stat_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*MemoryStats); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wave_v1_domain_stat_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*InterStats); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wave_v1_domain_stat_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*DiskStats); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wave_v1_domain_stat_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*DomainStats); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_wave_v1_domain_stat_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_wave_v1_domain_stat_proto_goTypes,
		DependencyIndexes: file_wave_v1_domain_stat_proto_depIdxs,
		EnumInfos:         file_wave_v1_domain_stat_proto_enumTypes,
		MessageInfos:      file_wave_v1_domain_stat_proto_msgTypes,
	}.Build()
	File_wave_v1_domain_stat_proto = out.File
	file_wave_v1_domain_stat_proto_rawDesc = nil
	file_wave_v1_domain_stat_proto_goTypes = nil
	file_wave_v1_domain_stat_proto_depIdxs = nil
}
