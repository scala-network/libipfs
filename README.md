# libznipfs

## V5 Rework

A lot has changed in IPFS since the V4 alpha test was done. We're doing a complete reword to align with the
latest version of IPFS.



## About

A C-style library implemented in Go to retrieve the seedlist from ZeroNet and IPFS

## Overview

libznipfs is used by the Stellite daemon to retrieve information from ZeroNet and
IPFS. It starts a full IPFS node and exposes some basic functionality to the
daemon.

When libznipfs is started from the daemon it runs a full IPFS node including the
HTTP API. This means standard IPFS commands can be used while the Stellite daemon
is running.

To use the API, you need to set the IPFS path to your Stellite data directory

Example:
`IPFS_PATH=~/.stellite/ipfs ipfs cat /ipfs/QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv/readme`

This will print the default IPFS readme

### Why Go and not C or C++

Currently, no simple implementation or API exists for ZeroNet and IPFS in C or C++. Instead of writing, or re-writing, large parts of ZeroNet and IPFS in C or C++ we rather use Go and compile it to a C or C++ compatible library. IPFS is implemented in Go already and a Go library for ZeroNet already exist.
