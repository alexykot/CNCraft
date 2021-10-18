# CNCraft
_cloud native Minecraft server_

This is a stub of a cloud native Minecraft server implementation in Golang.

Uses embedded NATS as async messaging backbone and Postgres for persistence. Designed to be a multi-process horizontally scalable system, but the implementation is  far from complete. Implements Minecraft network protocol v754 as per [wiki.vg spec](https://wiki.vg/index.php?title=Protocol&oldid=16676), supports v1.16.5 vanilla client.

**Feature completeness:** allows to login and run around a hardcoded limited size flatworld. Mining implementation unfinished. World generation is not intended by design (plan is to import worlds generated elsewhere).

**Development status:** shelved project until it will have a use case again.
