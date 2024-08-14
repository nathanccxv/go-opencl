// Package opencl provides a Go interface for interacting with OpenCL devices.
package opencl

import (
	"slices"
	"testing"
	"unsafe"
)

// TestRunner is a test function that tests the functionality of the OpenCL runner.
func TestRunner(t *testing.T) {
	// Step 1: Get OpenCL device information
	info, _ := Info()

	if len(info.Platforms) < 1 {
		t.Skipf("No OpenCL Devices")
	}
	if len(info.Platforms[0].Devices) < 1 {
		t.Skipf("No OpenCL Devices")
	}

	// Step 2: Initialize the OpenCL runner
	device := info.Platforms[0].Devices[0]
	runner, err := device.InitRunner()
	if err != nil {
		t.Fatal("InitRunner err:", err)
	}
	defer runner.Free()

	// Step 3: Compile the OpenCL kernels
	code := `__kernel void helloworld(__global int* in, __global int* out)
		 {
			 int num = get_global_id(0);
			 out[num] = in[num] * in[num];
		 }`
	codes := []string{code}
	kernelNameList := []string{"helloworld"}
	err = runner.CompileKernels(codes, kernelNameList, "")
	if err != nil {
		t.Fatal("CompileKernels err:", err)
	}

	// Step 4: Create kernel parameters
	/* buffer 1 param */
	input := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	itemSize := int(unsafe.Sizeof(input[0]))
	itemCount := len(input)

	input_buf, err := CreateBuffer(runner, READ_ONLY|COPY_HOST_PTR, input)
	if err != nil {
		t.Fatal("CreateBuffer err:", err)
	}
	/* buffer 2 param */
	output_buf, err := runner.CreateEmptyBuffer(WRITE_ONLY, itemCount*itemSize)
	if err != nil {
		t.Fatal("CreateEmptyBuffer err:", err)
	}

	// Step 5: Run the OpenCL kernel
	err = runner.RunKernel("helloworld", 1, nil, []uint64{uint64(itemCount)}, nil, []KernelParam{
		BufferParam(input_buf),
		BufferParam(output_buf),
	}, true)
	if err != nil {
		t.Fatal("RunKernel err:", err)
	}

	// Step 6: Read the output buffer
	result := make([]int32, itemCount)
	err = ReadBuffer(runner, 0, output_buf, result)
	if err != nil {
		t.Fatal("ReadBuffer err:", err)
	}

	// Step 7: Check the result
	expected_result := make([]int32, itemCount, itemCount)
	for i, v := range input {
		expected_result[i] = v * v
	}
	if !slices.Equal(result, expected_result) {
		t.Fatal("result error:", result)
	}

}
