package opencl

// #include "cl.h"
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)

func getOneDevie(platform_id C.cl_platform_id, device_id C.cl_device_id) (*OpenCLDevice, error) {
	var device = OpenCLDevice{Platform_id: platform_id, Device_id: device_id}

	// max_clock_frequency
	var err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_MAX_CLOCK_FREQUENCY, C.sizeof_cl_uint,
		unsafe.Pointer(&device.Max_clock_frequency), nil)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}

	// max_mem_alloc_size
	err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_MAX_MEM_ALLOC_SIZE, C.sizeof_cl_ulong,
		unsafe.Pointer(&device.Max_mem_alloc_size), nil)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}

	// global_mem_size
	err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_GLOBAL_MEM_SIZE, C.sizeof_cl_ulong,
		unsafe.Pointer(&device.Global_mem_size), nil)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}

	// max_compute_units
	err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_MAX_COMPUTE_UNITS, C.sizeof_cl_uint,
		unsafe.Pointer(&device.Max_compute_units), nil)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}

	// max_work_group_size
	err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_MAX_WORK_GROUP_SIZE, C.sizeof_size_t,
		unsafe.Pointer(&device.Max_work_group_size), nil)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}

	// device_type
	err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_TYPE, C.sizeof_cl_device_type,
		unsafe.Pointer(&device.Device_type), nil)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}

	// Max_work_item_dimensions
	err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_MAX_WORK_ITEM_DIMENSIONS, C.sizeof_cl_uint,
		unsafe.Pointer(&device.Max_work_item_dimensions), nil)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}

	// Max_work_item_sizes
	device.Max_work_item_sizes = make([]C.size_t, device.Max_work_item_dimensions, device.Max_work_item_dimensions)
	err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_MAX_WORK_ITEM_SIZES, C.sizeof_size_t*C.size_t(device.Max_work_item_dimensions),
		unsafe.Pointer(&device.Max_work_item_sizes[0]), nil)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}

	var infoSize C.size_t
	// name
	err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_NAME, 0, nil, &infoSize)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}
	var info = make([]byte, infoSize, infoSize)
	err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_NAME, infoSize, unsafe.Pointer(&info[0]), nil)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}
	device.Name = string(info[:len(info)-1])

	// profile
	err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_PROFILE, 0, nil, &infoSize)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}
	info = make([]byte, infoSize, infoSize)
	err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_PROFILE, infoSize, unsafe.Pointer(&info[0]), nil)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}
	device.Profile = string(info[:len(info)-1])

	// version
	err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_VERSION, 0, nil, &infoSize)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}
	info = make([]byte, infoSize, infoSize)
	err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_VERSION, infoSize, unsafe.Pointer(&info[0]), nil)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}
	device.Version = strings.Trim(string(info[:len(info)-1]), " ")

	// vendor
	err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_VENDOR, 0, nil, &infoSize)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}
	info = make([]byte, infoSize, infoSize)
	err = C.clGetDeviceInfo(device_id, C.CL_DEVICE_VENDOR, infoSize, unsafe.Pointer(&info[0]), nil)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}
	device.Vendor = string(info[:len(info)-1])

	// driver_version
	err = C.clGetDeviceInfo(device_id, C.CL_DRIVER_VERSION, 0, nil, &infoSize)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}
	info = make([]byte, infoSize, infoSize)
	err = C.clGetDeviceInfo(device_id, C.CL_DRIVER_VERSION, infoSize, unsafe.Pointer(&info[0]), nil)
	if err != C.CL_SUCCESS {
		return &device, fmt.Errorf("clGetDeviceInfo Err: %v", err)
	}
	device.Driver_version = string(info[:len(info)-1])

	return &device, nil
}

func getOnePlatform(platform_id C.cl_platform_id) (*OpenCLPlatform, error) {
	var platform = OpenCLPlatform{Platform_id: platform_id}

	var infoSize C.size_t
	// name
	var err = C.clGetPlatformInfo(platform_id, C.CL_PLATFORM_NAME, 0, nil, &infoSize)
	if err != C.CL_SUCCESS {
		return &platform, fmt.Errorf("clGetPlatformInfo Err: %v", err)
	}
	var info = make([]byte, infoSize, infoSize)
	err = C.clGetPlatformInfo(platform_id, C.CL_PLATFORM_NAME, infoSize, unsafe.Pointer(&info[0]), nil)
	if err != C.CL_SUCCESS {
		return &platform, fmt.Errorf("clGetPlatformInfo Err: %v", err)
	}
	platform.Name = string(info[:len(info)-1])

	// profile
	err = C.clGetPlatformInfo(platform_id, C.CL_PLATFORM_PROFILE, 0, nil, &infoSize)
	if err != C.CL_SUCCESS {
		return &platform, fmt.Errorf("clGetPlatformInfo Err: %v", err)
	}
	info = make([]byte, infoSize, infoSize)
	err = C.clGetPlatformInfo(platform_id, C.CL_PLATFORM_PROFILE, infoSize, unsafe.Pointer(&info[0]), nil)
	if err != C.CL_SUCCESS {
		return &platform, fmt.Errorf("clGetPlatformInfo Err: %v", err)
	}
	platform.Profile = string(info[:len(info)-1])

	// version
	err = C.clGetPlatformInfo(platform_id, C.CL_PLATFORM_VERSION, 0, nil, &infoSize)
	if err != C.CL_SUCCESS {
		return &platform, fmt.Errorf("clGetPlatformInfo Err: %v", err)
	}
	info = make([]byte, infoSize, infoSize)
	err = C.clGetPlatformInfo(platform_id, C.CL_PLATFORM_VERSION, infoSize, unsafe.Pointer(&info[0]), nil)
	if err != C.CL_SUCCESS {
		return &platform, fmt.Errorf("clGetPlatformInfo Err: %v", err)
	}
	platform.Version = strings.Trim(string(info[:len(info)-1]), " ")

	// vendor
	err = C.clGetPlatformInfo(platform_id, C.CL_PLATFORM_VENDOR, 0, nil, &infoSize)
	if err != C.CL_SUCCESS {
		return &platform, fmt.Errorf("clGetPlatformInfo Err: %v", err)
	}
	info = make([]byte, infoSize, infoSize)
	err = C.clGetPlatformInfo(platform_id, C.CL_PLATFORM_VENDOR, infoSize, unsafe.Pointer(&info[0]), nil)
	if err != C.CL_SUCCESS {
		return &platform, fmt.Errorf("clGetPlatformInfo Err: %v", err)
	}
	platform.Vendor = string(info[:len(info)-1])

	// devices
	err = C.clGetDeviceIDs(platform_id, C.CL_DEVICE_TYPE_ALL, 0, nil, &platform.Device_count)
	if err != C.CL_SUCCESS {
		return &platform, fmt.Errorf("clGetDeviceIDs Err: %v", err)
	}
	var device_ids = make([]C.cl_device_id, platform.Device_count, platform.Device_count)
	err = C.clGetDeviceIDs(platform_id, C.CL_DEVICE_TYPE_ALL, platform.Device_count, &device_ids[0], nil)
	if err != C.CL_SUCCESS {
		return &platform, fmt.Errorf("clGetDeviceIDs Err: %v", err)
	}

	for _, device_id := range device_ids {
		device, _ := getOneDevie(platform.Platform_id, device_id)
		platform.Devices = append(platform.Devices, device)
	}

	return &platform, nil
}

func Info() (*OpenCLInfo, error) {
	var info OpenCLInfo

	var err = C.clGetPlatformIDs(0, nil, &info.Platform_count)
	if err != C.CL_SUCCESS {
		return &info, fmt.Errorf("clGetPlatformIDs Err: %v", err)
	}

	if info.Platform_count == 0 {
		return &info, nil
	}

	var platform_ids = make([]C.cl_platform_id, info.Platform_count, info.Platform_count)
	err = C.clGetPlatformIDs(info.Platform_count, &platform_ids[0], nil)
	if err != C.CL_SUCCESS {
		return &info, fmt.Errorf("clGetPlatformIDs Err: %v", err)
	}

	for _, platform_id := range platform_ids {
		platform, _ := getOnePlatform(platform_id)
		info.Platforms = append(info.Platforms, platform)
	}

	return &info, nil
}
