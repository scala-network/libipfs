cp -rf ../bin/* .
g++ -pthread -o usage usage.cpp libipfs-linux.a
