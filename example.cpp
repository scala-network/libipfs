#include "libipfs-linux-amd64.h"
#include <unistd.h>
#include <iostream>
#include <nlohmann/json.hpp>

using json = nlohmann::json;

bool isStatusSuccess(const char* response) {
    try {
        json parsed = json::parse(response);
        return (parsed["Status"] == "success");
    } catch (const json::exception& e) {
        std::cerr << "JSON Parsing Error: " << e.what() << std::endl;
        return false;
    }
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

    try {
        json parsedStart = json::parse(startResponse);
        std::string peerId = parsedStart["Data"]["peerId"];
        std::cout << "Started with Peer ID: " << peerId << std::endl;

        const char* addResponse = Add((char*)"/home/hayzam/Projects/libipfs-new/bin/example.cpp");
        json parsedAdd = json::parse(addResponse);
        std::string cid = parsedAdd["Data"]["cid"];
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
    } catch (const json::exception& e) {
        std::cerr << "JSON Parsing Error: " << e.what() << std::endl;
        return 1;
    }

    return 0;
}