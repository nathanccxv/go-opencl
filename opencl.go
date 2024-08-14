package opencl

// #cgo CFLAGS: -DCL_TARGET_OPENCL_VERSION=120
// #cgo linux LDFLAGS: -lOpenCL
// #cgo darwin LDFLAGS: -framework OpenCL
// #cgo windows CFLAGS: -I${SRCDIR}/include-3.0.13
// #cgo windows LDFLAGS: -L${SRCDIR}/lib-windows-3.0.13-x64 -lOpenCL
import "C"
