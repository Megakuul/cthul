/**
 * Cthul System
 *
 * Copyright (C) 2024 Linus Ilian Moser <linus.moser@megakuul.ch>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package structure

type DomainStats struct {
	
}

type CpuStats struct {
	Timestamp int64
	CpuTime int64 // The total domain cpu time since machine startup in nanoseconds.
	UserTime int64 // The total domain cpu time spent in user space context in nanoseconds.
 	KernelTime int64 // The total domain cpu time spent in kernel space context in nanoseconds.
}

type MemoryStats struct {
	Timestamp int64
	SwapIn int64 // total amount of bytes moved from main memory to swap since machine startup.
	SwapOut int64 // total amount of bytes moved from swap to main memory since machine startup.
	MinorFaults int64 // total amount of handled page fault that loaded pages from memory (e.g. shared libs).
	MajorFaults int64 // total amount of handled page fault that loaded pages from diskio
	HugepageAllocations int64 // total amount of successful hugetable allocations since machine startup.
	HugepageFailures int64 // total amount of failed hugetable allocations since machine startup.
	Balloned int64 // total amount of bytes currently allocated for the domain by the ballon driver.
	Available int64 // total amount of bytes from guestos perspective (may be less then the assigned memory).
	Usable int64 // total amount of bytes the guest os can use if all caches are reclaimed (deflatable memory).
	Unused int64 // total amount of bytes that are completely unused by the guest os (part of deflatable memory).
	HostRSS int64 // resident set size (used main memory) of the domains process on the host os (typically qemu).
}
