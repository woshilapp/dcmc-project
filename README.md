# dcmc-project
Direct Connection Minecraft Project(unfinished)

## Important
It's a unfinish project. It can't run properly and keep it's security. Code still shit
and I haven't tidy it.

## Introduction

DCMC (Direct Connection Minecraft) is a network communication tool developed in Go language, designed to establish direct connections between clients without relying on central servers. The project implements NAT traversal technology and room management mechanisms to enable peer-to-peer communication in distributed environments.

## Key Features

- Client role selection (Peer/Host)
- Server-based room management
- NAT traversal implementation
- Message passing between clients
- Command-line interface (CLI)

## Requires

- Go 1.23.3 or later

## Protocol

The project uses a custom data transmission protocol that supports integers, strings, booleans, and floating-point numbers. Refer to `protocol_id.txt` for detailed protocol specifications.

## License

This project is licensed under the GNU General Public License v3.0. See the `LICENSE` file for more information.
