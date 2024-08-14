[![Test](https://github.com/nathanccxv/go-opencl/actions/workflows/test.yml/badge.svg)](https://github.com/nathanccxv/go-opencl/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/nathanccxv/go-opencl.svg)](https://pkg.go.dev/github.com/nathanccxv/go-opencl)

# go-opencl

go-opencl provides a high-level interface for OpenCL devices to run OpenCL programs in Go programs conveniently without delving into the annoying details of OpenCL.

## Development Status

**WARNING**: This project is currently under development and has not been fully tested. Use it at your own risk. We welcome any feedback and contributions.

## Requirements

**linux**
```bash
sudo apt install ocl-icd-opencl-dev opencl-headers
```

**windows**

This project incorporates OpenCL-Headers and OpenCL-ICD-Loader, which are included in the `include-3.0.13` and `lib-windows-3.0.13-x64` directories respectively for Windows.

The sources for these components are as follows:

- OpenCL-Headers: [KhronosGroup/OpenCL-Headers v2023.02.06](https://github.com/KhronosGroup/OpenCL-Headers/releases/tag/v2023.02.06)

- OpenCL-ICD-Loader: [KhronosGroup/OpenCL-SDK v2023.02.06](https://github.com/KhronosGroup/OpenCL-SDK/releases/tag/v2023.02.06)


## cl-info command

The cl-info command provides information about the OpenCL platforms and devices on your system.

To install cl-info, run the following command:

```bash
go install github.com/nathanccxv/go-opencl/cmd/cl-info@latest
```

## OpenCL runner

```go
import cl "github.com/nathanccxv/go-opencl"
```

Refer to the [runner_test.go](./runner_test.go) file or [examples](./examples/) for usage examples of the OpenCL runner.

## Other resources

OPENCL 3.0 Reference: https://registry.khronos.org/OpenCL/sdk/3.0/docs/man/html/
