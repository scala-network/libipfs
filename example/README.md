
# Building the usage sample

```
# Copy the library and header
cp ../bin/libznipfs-linux.* ./
# Compile
g++ -pthread -o usage usage.cpp libznipfs-linux.a
# Run
./usage
```

## Example output

```
This is a test C application to call libznipfs.
{"Status":"ok","Message":"Seedlist retrieved from ZeroNet and IPFS","Seedlist":["139.180.198.68:20188","89.158.106.154:20188"]}⏎   
```
