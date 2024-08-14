package opencl

// #include "cl.h"
import "C"

type OpenCLDevice struct {
	Device_id   C.cl_device_id
	Platform_id C.cl_platform_id

	Device_type    C.cl_device_type
	Name           string
	Profile        string
	Version        string
	Vendor         string
	Driver_version string

	Max_clock_frequency C.cl_uint
	Max_mem_alloc_size  C.cl_ulong
	Global_mem_size     C.cl_ulong
	Max_compute_units   C.cl_uint
	Max_work_group_size C.size_t

	Max_work_item_dimensions C.cl_uint
	Max_work_item_sizes      []C.size_t
}

type OpenCLPlatform struct {
	Platform_id  C.cl_platform_id
	Device_count C.cl_uint
	Name         string
	Profile      string
	Version      string
	Vendor       string
	Devices      []*OpenCLDevice
}

type OpenCLInfo struct {
	Platform_count C.cl_uint
	Platforms      []*OpenCLPlatform
}
