#if defined(__linux__) && defined(__x86_64__)
#  include "libipfs-linux-amd64.h"
#elif defined(__linux__) && defined(__aarch64__)
#  include "libipfs-linux-arm64.h"
#elif defined(__linux__) && defined(__riscv) && (__riscv_xlen == 64)
#  include "libipfs-linux-riscv64.h"
#elif defined(__APPLE__) && defined(__x86_64__)
#  include "libipfs-darwin-amd64.h"
#elif defined(__APPLE__) && defined(__aarch64__)
#  include "libipfs-darwin-arm64.h"
#elif defined(_WIN32) && defined(_M_X64)
#  include "libipfs-windows-amd64.h"
#elif defined(__FreeBSD__) && defined(__x86_64__)
#  include "libipfs-freebsd-amd64.h"
#else
#  error "Unsupported platform/architecture"
#endif

#include <unistd.h>
#include <iostream>
#include <string>

bool isStatusSuccess(const char* response) {
    if (!response) return false;
    return std::string(response).find("\"status\":\"success\"") != std::string::npos;
}

std::string extractField(const std::string& json, const std::string& key) {
    const std::string needle = "\"" + key + "\":\"";
    auto start = json.find(needle);
    if (start == std::string::npos) return "";
    start += needle.size();
    auto end = json.find('"', start);
    if (end == std::string::npos) return "";
    return json.substr(start, end - start);
}

void handleResponse(bool success, const std::string& successMsg, const std::string& errorMsg, const char* response = nullptr) {
    if (success) {
        std::cout << successMsg << std::endl;
    } else {
        std::cerr << errorMsg << std::endl;
        if (response) {
            std::cerr << response << std::endl;
        }
    }
}

void performAction(char* (*actionFunc)(char*), const std::string& arg, const std::string& actionName) {
    const char* response = actionFunc((char*)arg.c_str());
    handleResponse(isStatusSuccess(response), actionName + " succeeded", actionName + " failed", response);
}

void performAction(char* (*actionFunc)(), const std::string& actionName) {
    const char* response = actionFunc();
    handleResponse(isStatusSuccess(response), actionName + " succeeded", actionName + " failed", response);
}

int main() {
    std::cout << "Starting IPFS from C++ Land ..." << std::endl;
    const char* startResponse = Start((char*)"./ipfs", 16969);

    std::cout << "Add Peer" << std::endl;
    performAction(AddPeer, "/ip4/104.131.131.82/udp/4001/quic-v1/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ", "Add Peer");

    sleep(5);

    std::cout << "Remove Peer" << std::endl;
    performAction(RemovePeer, "/ip4/104.131.131.82/udp/4001/quic-v1/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ", "Remove Peer");

    if (!isStatusSuccess(startResponse)) {
        std::cerr << "Failed to start IPFS node: " << startResponse << std::endl;
        return 1;
    }

    std::string peerId = extractField(startResponse, "peerId");
    std::cout << "Started with Peer ID: " << peerId << std::endl;

    const char* addResponse = Add((char*)"./example.cpp");
    if (!isStatusSuccess(addResponse)) {
        std::cerr << "Failed to add file: " << addResponse << std::endl;
        return 1;
    }
    std::string cid = extractField(addResponse, "cid");
    std::cout << "Added example.cpp file -> CID: " << cid << std::endl;

    performAction(Pin, cid, "Pin CID");

    std::cout << "Sleeping for a bit after pinning..." << std::endl;
    sleep(5);

    performAction(Unpin, cid, "Unpin CID");
    performAction(GarbageCollect, "Garbage Collect");

    const std::string targetPath = "./downloaded.txt";
    const std::string getCid = "/ipfs/bafybeifx7yeb55armcsxwwitkymga5xf53dxiarykms3ygqic223w5sk3m";
    if (!isStatusSuccess(Get((char*)getCid.c_str(), (char*)targetPath.c_str(), true))) {
        std::cerr << "Failed to get CID: " << getCid << std::endl;
    } else {
        std::cout << "Got CID: " << getCid << " and saved to " << targetPath << std::endl;
    }

    sleep(60);

    std::cout << Stop() << std::endl;
    return 0;
}