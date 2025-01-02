/**
 * Cthul System
 *
 * Copyright (C) 2024 Linus Ilian Moser <linus.moser@megakuul.ch>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package structure

type DOMAIN_STATE int64

const (
	DOMAIN_NOSTATE DOMAIN_STATE = iota
	DOMAIN_RUNNING
	DOMAIN_BLOCKED
	DOMAIN_PAUSED
	DOMAIN_SHUTDOWN
	DOMAIN_SHUTOFF
	DOMAIN_CRASHED
	DOMAIN_PMSUSPENDED
)

type DomainStats struct {
	State  DOMAIN_STATE `json:"state"`
	Cpu    CpuStats     `json:"cpu"`
	Memory MemoryStats  `json:"memory"`
}

type CpuStats struct {
	Timestamp int64 `json:"timestamp"`

	// The total domain cpu time since machine startup in nanoseconds.
	CpuTime int64 `json:"cpu_time"`
	// The total domain cpu time spent in user space context in nanoseconds.
	UserTime int64 `json:"user_time"`
	// The total domain cpu time spent in kernel space context in nanoseconds.
	KernelTime int64 `json:"kernel_time"`

	VCpus []VCpuStats `json:"vcpus"`
}

type VCpuStats struct {
	// full vcpu time since machine startup in nanoseconds.
	CpuTime int64 `json:"cpu_time"`
	// wait time (time vcpu wants to run but is not scheduled) since machine startup in nanoseconds.
	WaitTime int64 `json:"wait_time"`
	// halt/idle (time the vcpu is disabled or halted) time since machine startup in nanoseconds.
	HaltTime int64 `json:"halt_time"`
	// delay time (time the vcpu is enqueued) since machine startup in nanoseconds.
	DelayTime int64 `json:"delay_time"`
}

type MemoryStats struct {
	Timestamp int64 `json:"timestamp"`

	// total amount of bytes moved from main memory to swap since machine startup.
	SwapIn int64 `json:"swap_in"`
	// total amount of bytes moved from swap to main memory since machine startup.
	SwapOut int64 `json:"swap_out"`
	// total amount of handled page fault that loaded pages from memory (e.g. shared libs).
	MinorFaults int64 `json:"minor_faults"`
	// total amount of handled page fault that loaded pages from diskio
	MajorFaults int64 `json:"major_faults"`
	// total amount of successful hugetable allocations since machine startup.
	HugepageAllocations int64 `json:"hugepage_allocations"`
	// total amount of failed hugetable allocations since machine startup.
	HugepageFailures int64 `json:"hugepage_failures"`
	// total amount of bytes currently allocated for the domain by the ballon driver.
	Balloned int64 `json:"balloned"`
	// total amount of bytes from guestos perspective (may be less then the assigned memory).
	Available int64 `json:"available"`
	// total amount of bytes the guest os can use if all caches are reclaimed (deflatable memory).
	Usable int64 `json:"usable"`
	// total amount of bytes that are completely unused by the guest os (part of deflatable memory).
	Unused int64 `json:"unused"`
	// resident set size (used main memory) of the domains process on the host os (typically qemu).
	HostRSS int64 `json:"host_rss"`
}

type Interface struct {
	Timestamp int64 `json:"timestamp"`

	// total amount of bytes received since machine startup.
	RecvBytes int64 `json:"recv_bytes"`
	// total number of ethernet packets received since machine startup.
	RecvPkts int64 `json:"recv_pkts"`
	// total number of receive-errors since machine startup.
	RecvErrs int64 `json:"recv_errs"`
	// total number of dropped receive-packets since machine startup.
	RecvDrops int64 `json:"recv_drops"`
	// total amount of bytes sent since machine startup.
	SendBytes int64 `json:"send_bytes"`
	// total number of ethernet packets sent since machine startup.
	SendPkts int64 `json:"send_pkts"`
	// total number of send-errors since machine startup.
	SendErrs int64 `json:"send_errs"`
	// total number of dropped send-packets since machine startup.
	SendDrops int64 `json:"send_drops"`
}

type Block struct {
	Timestamp int64 `json:"timestamp"`

	// total number of read requests since machine startup.
	ReadReqs int64 `json:"read_reqs"`
	// total amount of bytes read since machine startup.
	ReadBytes int64 `json:"read_bytes"`
	// total read time since machine startup in nanoseconds.
	ReadTime int64 `json:"read_time"`
	// total number of write requests since machine startup.
	WriteReqs int64 `json:"write_reqs"`
	// total amount of bytes write since machine startup.
	WriteBytes int64 `json:"write_bytes"`
	// total write time since machine startup in nanoseconds.
	WriteTime int64 `json:"write_time"`
	// total number of flush requests since machine startup.
	FlushReqs int64 `json:"flush_reqs"`
	// total flush time since machine startup in nanoseconds.
	FlushTime int64 `json:"flush_time"`
}
