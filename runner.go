package opencl

// #include "cl.h"
import "C"

import (
	"fmt"
	"unsafe"
)

// Buffer represents an OpenCL buffer.
type Buffer struct {
	buffer C.cl_mem
}

// OpenCLRunner represents an OpenCL runner.
type OpenCLRunner struct {
	Device       *OpenCLDevice
	Context      C.cl_context
	CommandQueue C.cl_command_queue

	Program C.cl_program
	Kernels map[string]C.cl_kernel
	Buffers []*Buffer
}

// InitRunner initializes an OpenCLRunner for the given OpenCLDevice.
// It creates a context and a command queue.
func (device *OpenCLDevice) InitRunner() (*OpenCLRunner, error) {
	var runner = OpenCLRunner{Device: device}

	// clCreateContext
	var context_properties = [3]C.cl_context_properties{
		C.CL_CONTEXT_PLATFORM,
		C.cl_context_properties(uintptr(unsafe.Pointer(device.Platform_id))),
		0,
	}
	var err C.cl_int
	var context = C.clCreateContext(&context_properties[0], 1, &device.Device_id, nil, nil, &err)
	if err != C.CL_SUCCESS {
		return nil, fmt.Errorf("clCreateContext Err: %v", err)
	}

	// clCreateCommandQueue
	var commandQueueProperties C.cl_command_queue_properties = 0
	var commandQueue = C.clCreateCommandQueue(context, device.Device_id, commandQueueProperties, &err)
	if err != C.CL_SUCCESS {
		C.clReleaseContext(context)
		return nil, fmt.Errorf("clCreateCommandQueueErr: %v", err)
	}

	runner.Context = context
	runner.CommandQueue = commandQueue

	return &runner, nil
}

// Free releases all resources associated with the OpenCLRunner.
func (runner *OpenCLRunner) Free() error {
	var err C.cl_int

	if len(runner.Kernels) > 0 {
		for _, kernel := range runner.Kernels {
			err = C.clReleaseKernel(kernel)
		}
	}

	if runner.Program != nil {
		err = C.clReleaseProgram(runner.Program)
	}

	if len(runner.Buffers) > 0 {
		for _, buffer := range runner.Buffers {
			err = C.clReleaseMemObject(buffer.buffer)
		}
	}

	err = C.clReleaseCommandQueue(runner.CommandQueue)
	err = C.clReleaseContext(runner.Context)

	if err != C.CL_SUCCESS {
		return fmt.Errorf("OpenCLRunner.Free cl_err: %v", err)
	}

	return nil
}

// CompileKernels compiles OpenCL kernels from the provided source code.
func (runner *OpenCLRunner) CompileKernels(codeSourceList []string, kernelNameList []string, options string) error {
	var codes [](*C.char)
	for _, codeSource := range codeSourceList {
		code_src := C.CString(codeSource)
		defer C.free(unsafe.Pointer(code_src))
		codes = append(codes, code_src)
	}

	var err C.cl_int

	// clCreateProgramWithSource
	var program = C.clCreateProgramWithSource(runner.Context, C.cl_uint(len(codes)), &codes[0], nil, &err)
	if err != C.CL_SUCCESS {
		return fmt.Errorf("clCreateProgramWithSource Err: %v", err)
	}

	cl_options := C.CString(options)
	defer C.free(unsafe.Pointer(cl_options))

	// clBuildProgram
	err = C.clBuildProgram(program, 1, &runner.Device.Device_id, cl_options, nil, nil)
	if err != C.CL_SUCCESS {
		var logSize C.size_t
		var err2 = C.clGetProgramBuildInfo(program, runner.Device.Device_id, C.CL_PROGRAM_BUILD_LOG, 0, nil, &logSize)
		if err2 != C.CL_SUCCESS {
			C.clReleaseProgram(program)
			return fmt.Errorf("clGetProgramBuildInfo Err: %v", err2)
		}

		var log_buf = make([]byte, logSize, logSize)
		err2 = C.clGetProgramBuildInfo(program, runner.Device.Device_id, C.CL_PROGRAM_BUILD_LOG, logSize, unsafe.Pointer(&log_buf[0]), nil)

		if err2 != C.CL_SUCCESS {
			C.clReleaseProgram(program)
			return fmt.Errorf("clGetProgramBuildInfo Err: %v", err2)
		}

		fmt.Printf("clBuildProgram Err log: %s\n", string(log_buf))

		C.clReleaseProgram(program)
		return fmt.Errorf("clBuildProgram Err: %v", err)
	}

	// clCreateKernel
	runner.Kernels = make(map[string]C.cl_kernel)
	for _, kernelName := range kernelNameList {
		var kernel_name = C.CString(kernelName)
		defer C.free(unsafe.Pointer(kernel_name))

		var kernel = C.clCreateKernel(program, kernel_name, &err)
		if err != C.CL_SUCCESS {
			C.clReleaseProgram(program)
			return fmt.Errorf("clCreateKernel Err: %v", err)
		}
		runner.Kernels[kernelName] = kernel
	}

	runner.Program = program

	return nil
}

const (
	READ_WRITE     C.cl_mem_flags = C.CL_MEM_READ_WRITE
	WRITE_ONLY                    = C.CL_MEM_WRITE_ONLY
	READ_ONLY                     = C.CL_MEM_READ_ONLY
	USE_HOST_PTR                  = C.CL_MEM_USE_HOST_PTR
	ALLOC_HOST_PTR                = C.CL_MEM_ALLOC_HOST_PTR
	COPY_HOST_PTR                 = C.CL_MEM_COPY_HOST_PTR
)

// CreateBuffer creates an OpenCL buffer with the specified flags and source data.
func CreateBuffer[E any](runner *OpenCLRunner, flags C.cl_mem_flags, source []E) (*Buffer, error) {
	if len(source) == 0 {
		return nil, fmt.Errorf("clCreateBuffer Err: source is empty")
	}
	var err C.cl_int
	size := C.size_t(int(unsafe.Sizeof(source[0])) * len(source))
	host_ptr := unsafe.Pointer(&source[0])
	cl_mem := C.clCreateBuffer(runner.Context, flags, size, host_ptr, &err)
	if err != C.CL_SUCCESS {
		return nil, fmt.Errorf("clCreateBuffer Err: %v", err)
	}

	buffer := &Buffer{cl_mem}
	runner.Buffers = append(runner.Buffers, buffer)
	return buffer, nil
}

// CreateEmptyBuffer creates an empty OpenCL buffer with the specified flags and size.
func (runner *OpenCLRunner) CreateEmptyBuffer(flags C.cl_mem_flags, size int) (*Buffer, error) {
	var err C.cl_int
	cl_mem := C.clCreateBuffer(runner.Context, flags, C.size_t(size), nil, &err)
	if err != C.CL_SUCCESS {
		return nil, fmt.Errorf("clCreateBuffer Err: %v", err)
	}
	buffer := &Buffer{cl_mem}
	runner.Buffers = append(runner.Buffers, buffer)
	return buffer, nil
}

// ReadBuffer reads data from an OpenCL buffer into the target slice.
func ReadBuffer[E any](runner *OpenCLRunner, offset int, buffer *Buffer, target []E) error {
	if len(target) == 0 {
		return fmt.Errorf("clEnqueueReadBuffer Err: target is nil")
	}
	err := C.clEnqueueReadBuffer(runner.CommandQueue, buffer.buffer, C.CL_TRUE, C.size_t(offset),
		C.size_t(int(unsafe.Sizeof(target[0]))*len(target)),
		unsafe.Pointer(&target[0]), 0, nil, nil)
	if err != C.CL_SUCCESS {
		return fmt.Errorf("clEnqueueReadBuffer Err: %v", err)
	}
	return nil
}

// WriteBuffer writes data from the source slice to an OpenCL buffer.
func WriteBuffer[E any](runner *OpenCLRunner, offset int, buffer *Buffer, source []E, blocking bool) error {
	if len(source) == 0 {
		return fmt.Errorf("clEnqueueWriteBuffer Err: source is empty")
	}
	var _blocking C.cl_bool = C.CL_FALSE
	if blocking {
		_blocking = C.CL_TRUE
	}
	err := C.clEnqueueWriteBuffer(runner.CommandQueue, buffer.buffer, _blocking, C.size_t(offset),
		C.size_t(int(unsafe.Sizeof(source[0]))*len(source)), unsafe.Pointer(&source[0]), 0, nil, nil)
	if err != C.CL_SUCCESS {
		return fmt.Errorf("clEnqueueWriteBuffer Err: %v", err)
	}
	return nil
}

// ReleaseBuffer releases the specified OpenCL buffer.
func (runner *OpenCLRunner) ReleaseBuffer(buffer *Buffer) error {
	err := C.clReleaseMemObject(buffer.buffer)
	if err != C.CL_SUCCESS {
		return fmt.Errorf("clReleaseMemObject Err: %v", err)
	}
	return nil
}

// map_size_t converts a slice of uint64 to a slice of C.size_t.
func map_size_t[E uint64](slice []E) []C.size_t {
	size_t_slice := make([]C.size_t, len(slice), len(slice))
	for i, v := range slice {
		size_t_slice[i] = C.size_t(v)
	}
	return size_t_slice
}

// KernelParam represents a parameter for an OpenCL kernel.
type KernelParam struct {
	Size    uintptr
	Pointer unsafe.Pointer
}

// BufferParam creates a KernelParam for an OpenCL buffer.
func BufferParam(v *Buffer) KernelParam {
	return KernelParam{Size: unsafe.Sizeof(v.buffer), Pointer: unsafe.Pointer(&v.buffer)}
}

// Param creates a KernelParam for a value.
func Param[E any](v *E) KernelParam {
	return KernelParam{Size: unsafe.Sizeof(*v), Pointer: unsafe.Pointer(v)}
}

// SetKernelArgs sets the arguments for a specific OpenCL kernel.
func (runner *OpenCLRunner) SetKernelArgs(kernelName string, args []KernelParam) error {
	var kernel = runner.Kernels[kernelName]
	var err C.cl_int
	for i, arg := range args {
		err = C.clSetKernelArg(kernel, C.cl_uint(i), C.size_t(arg.Size), arg.Pointer)
		if err != C.CL_SUCCESS {
			return fmt.Errorf("clSetKernelArg Err: %v", err)
		}
	}
	return nil
}

// RunKernel runs an OpenCL kernel with the specified work dimensions, work sizes, and arguments.
func (runner *OpenCLRunner) RunKernel(kernelName string, work_dim int,
	global_work_offset []uint64, global_work_size []uint64, local_work_size []uint64, args []KernelParam, wait bool) error {
	var kernel = runner.Kernels[kernelName]
	var err C.cl_int
	for i, arg := range args {
		err = C.clSetKernelArg(kernel, C.cl_uint(i), C.size_t(arg.Size), arg.Pointer)
		if err != C.CL_SUCCESS {
			return fmt.Errorf("clSetKernelArg Err: %v", err)
		}
	}

	var global_work_offset_ptr, global_work_size_ptr, local_work_size_ptr *C.size_t = nil, nil, nil

	if len(global_work_offset) != 0 {
		_global_work_offset := map_size_t(global_work_offset)
		global_work_offset_ptr = &_global_work_offset[0]
	}

	if len(global_work_size) != 0 {
		_global_work_size := map_size_t(global_work_size)
		global_work_size_ptr = &_global_work_size[0]
	}

	if len(local_work_size) != 0 {
		_local_work_size := map_size_t(local_work_size)
		local_work_size_ptr = &_local_work_size[0]
	}

	var evt *C.cl_event = nil
	if wait {
		var evt_obj C.cl_event
		evt = &evt_obj
		defer C.clReleaseEvent(evt_obj)
	}
	err = C.clEnqueueNDRangeKernel(runner.CommandQueue, kernel, C.cl_uint(work_dim),
		global_work_offset_ptr, global_work_size_ptr, local_work_size_ptr, 0, nil, evt)
	if err != C.CL_SUCCESS {
		return fmt.Errorf("clEnqueueNDRangeKernel Err: %v", err)
	}

	if wait {
		err = C.clWaitForEvents(1, evt)
		if err != C.CL_SUCCESS {
			return fmt.Errorf("clWaitForEvents Err: %v", err)
		}
	}

	return nil
}
