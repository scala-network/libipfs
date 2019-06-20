
# Building the usage sample

```
# Copy the library and header
cp ../bin/libznipfs-linux.* ./
# Compile
g++ -pthread -o usage usage.cpp libznipfs-linux.a
# Run
./usage
```
