# SparkCell Roadmap

SparkCell is a lightweight database cell runtime for Linux hosts. The goal is to run database workloads in isolated cells with their own filesystem, process tree, network namespace, and resource limits, while keeping the control plane simple: a client talks to a daemon over a Unix socket.

## Current Direction

- Primary target database: PostgreSQL
- Control plane split: `sparkcell` client and persistent daemon
- Base filesystem source: `base_linux_environment.gz`
- Host storage root: `/var/lib/sparkcell/`

## Roadmap

### Phase 1: Engine Foundation and IPC

- [ ] Build the CLI command structure
  - `sparkcell run postgres --port <port> --memory <limit> --cpu <limit>`
  - `sparkcell list`
  - `sparkcell stop <cell_id>`
- [ ] Build the persistent daemon service
  - Listen on `/var/run/sparkcell.sock`
  - Run as a root-managed service
- [ ] Add daemon state storage
  - Track `id`, `host_pid`, `port`, `status`, `memory_limit`, `cpu_limit`, `created_at`
  - Use SQLite inside the daemon
- [ ] Define the JSON message contract over Unix sockets
  - Client sends request payloads
  - Daemon returns structured success or error responses

### Phase 2: Filesystem Provisioning

- [ ] Create the host storage layout
  - `/var/lib/sparkcell/images/`
  - `/var/lib/sparkcell/cells/`
- [ ] Add image provisioning
  - Use bundled base image when present
  - Download release asset when missing
- [ ] Add cell root creation
  - Generate unique cell IDs
  - Extract the base rootfs into each cell directory
- [ ] Add cell deletion and cleanup
  - Stop the running process
  - Remove cell storage when requested

### Phase 3: Runtime Isolation

- [ ] Add process isolation entry flow
  - Re-entrant fork pattern
  - Child execution via `/proc/self/exe`
- [ ] Add namespace isolation
  - UTS
  - PID
  - Mount
  - Network
- [ ] Add root filesystem anchoring
  - `chroot`
  - `chdir("/")`
  - Mount required system paths
  - Mount `/proc`

### Phase 4: Resource Control

- [ ] Add cgroups v2 support
  - Create a per-cell cgroup leaf
  - Track processes in `cgroup.procs`
- [ ] Enforce memory limits
  - Parse memory flags
  - Write `memory.max`
  - Disable swap for cells
- [ ] Enforce CPU limits
  - Parse CPU flags
  - Write `cpu.max`

### Phase 5: Network Connectivity

- [ ] Create bridge networking on the host
  - `sparkbr0`
  - Gateway subnet assignment
- [ ] Create per-cell veth pairs
- [ ] Attach host-side interfaces to the bridge
- [ ] Move cell-side interfaces into the namespace
- [ ] Configure cell networking
  - `lo`
  - `eth0`
  - IP assignment
  - Default route
- [ ] Add NAT and port forwarding
  - Map host ports to database ports

### Phase 6: PostgreSQL Lifecycle

- [ ] Bootstrap PostgreSQL inside the cell
  - Run `initdb`
  - Prepare data directories
- [ ] Configure PostgreSQL for cell access
  - `listen_addresses = '*'`
- [ ] Start PostgreSQL as the payload process
- [ ] Validate end-to-end execution
  - `sparkcell run postgres --port 8001 --memory 512M --cpu 0.5`
  - Confirm isolation, limits, and external connectivity

## Working Rule

As implementation lands, move completed items from unchecked to checked so this file stays a live roadmap rather than a static note.
