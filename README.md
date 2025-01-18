# libIPFS

`libIPFS` is a C-style library that wraps around `go-ipfs` as a library (not an embedded executable) and provides a simple and intuitive API.

## Building

Use the provided `Makefile` to build the library. For example, to build for Linux on `amd64`:

```bash
make build_linux_amd64
```

You can also build for other platforms by specifying the target:

```bash
make build_windows_amd64
make build_darwin_amd64
make build_darwin_arm64
make build_linux_riscv64
make build_linux_arm64
make build_freebsd_amd64
```

> **Note:**  
> - Cross-compilation has been tested on Debian 12.  
> - For FreeBSD builds, you may need a proper sysroot configuration.  
> - For macOS builds, `osxcross` is required.  

## Usage

An example program, `example.cpp`, is provided in the root of the repository. After building the library, copy the example file to the `bin/` folder and compile it with `g++`:

```bash
g++ -pthread -o libipfs-example example.cpp libipfs-linux-amd64.a -Wl,--no-as-needed -ldl -lresolv
```

Run the compiled example:

```bash
./libipfs-example
```

## License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.