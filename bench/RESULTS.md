# NeoHome Cassandra Benchmark Results

## Environment

| Parameter | Value |
|-----------|-------|
| Cassandra | 5.0, single node, Docker |
| Replication | SimpleStrategy, RF=1 |
| Consistency | QUORUM (= 1 node with RF=1) |
| Go API | go1.24, gorilla/mux, gocql |
| Machine | macOS, Apple Silicon |
| Connections | 10 concurrent |
| Duration | 30s per scenario (+ 5s warm-up) |
| Seeded data | 5 devices, 250 telemetry rows, 2 thresholds/device |

## Results

| Scenario | Req/s | p50 ms | p90 ms | p99 ms | Avg ms | Max ms | CQL ops/req | Errors |
|----------|------:|-------:|-------:|-------:|-------:|-------:|------------:|-------:|
| write-telemetry | 1,441 | 6 | 8 | 10 | 6.39 | 57 | 4 | 0 |
| write-with-alerts | 1,135 | 8 | 10 | 13 | 8.27 | 83 | 5 | 0 |
| read-telemetry | 2,301 | 4 | 5 | 8 | 3.85 | 36 | 2 | 0 |
| read-latest | 1,328 | 7 | 10 | 17 | 7.03 | 30 | 2 | 0 |
| mixed-workload | 1,647 | 6 | 8 | 10 | 5.55 | 62 | ~3.6 | 0 |

## Analysis

### Write Throughput (Telemetry Ingestion)

The primary IoT workload — `POST /api/v1/telemetry` — sustains **~1,441 req/s** with 10 connections. Each request performs 4 Cassandra round-trips:

1. `SELECT` from `devices_by_id` (device lookup)
2. `INSERT` batch into `telemetry_by_device` + `telemetry_by_device_metric`
3. `INSERT` batch into `devices_by_id` + `devices_by_user` (update last_seen, battery, etc.)
4. `SELECT` from `device_thresholds_by_device` (alert threshold check)

This means Cassandra is handling **~5,764 CQL operations/sec** (1,441 x 4) for the write path alone. The p99 latency of 10ms is excellent for a single-node deployment.

### Alert Write Amplification

When every telemetry reading triggers an alert (`write-with-alerts`), throughput drops to **1,135 req/s** — a **21% decrease** compared to the no-alert baseline. Each alert adds a 5th round-trip: a batch of 3 `INSERT`s into `alerts_by_id`, `alerts_by_location`, and `alerts_by_device`.

The p99 latency increases from 10ms to 13ms (+30%). This quantifies the cost of Cassandra's denormalized alert model: 3 extra writes per alert, across 3 tables, to support querying alerts by ID, by location, and by device.

### Read Performance

**Telemetry history** (`read-telemetry`) at **2,301 req/s** is the fastest scenario — expected, since it only requires 2 CQL queries (device ownership check + partition scan with LIMIT 100). Cassandra's clustering key ordering (`recorded_at DESC`) means no sorting is needed.

**Latest telemetry** (`read-latest`) is notably slower at **1,328 req/s** despite also being 2 CQL queries. The difference is the scan width: this endpoint reads up to 1000 rows from the partition and deduplicates by metric type in Go. The p99 of 17ms (vs 8ms for history) reflects this heavier scan.

### Mixed Workload

The 80/20 write/read mix achieves **1,647 req/s** — between the pure write (1,441) and pure read (2,301) numbers, as expected. The p99 of 10ms matches the pure write scenario, suggesting that reads don't create significant contention. This aligns with Cassandra's architecture: reads and writes use separate code paths and memtables.

### Cassandra-Specific Observations

1. **Write-optimized model validated**: 1,400+ writes/sec on a single Docker node confirms Cassandra's log-structured merge tree (LSM) handles IoT ingestion well. Writes go to memtable (RAM) first, making them fast regardless of dataset size.

2. **Denormalization cost is measurable**: Writing to 2 telemetry tables + 2 device tables per ingest is a 4x write amplification. With alerts, it's 7x. This is the explicit tradeoff of Cassandra's query-driven data model — fast reads at the cost of write amplification.

3. **Partition distribution matters**: We distribute writes across 5 devices (5 partitions). A single-device benchmark would show worse numbers due to partition-level contention. In production with thousands of devices, the distribution would be even better.

4. **Latency is consistent**: The gap between p50 and p99 is small (6ms vs 10ms for writes). Cassandra's deterministic write path (memtable → commitlog) avoids the latency spikes common in B-tree databases under write pressure.

5. **Single-node limitation**: These results don't demonstrate Cassandra's horizontal scaling. With a 3-node cluster and RF=3, write throughput would remain similar (writes go to all replicas in parallel) while read throughput could scale linearly by spreading partition ownership.

## How to Reproduce

```bash
docker compose up -d          # start all services
cd bench && npm install
npm run bench                 # run all 5 scenarios

# Individual scenarios:
npm run bench:write           # write-telemetry only
npm run bench:alerts          # write-with-alerts only
npm run bench:read            # read-telemetry only
npm run bench:latest          # read-latest only
npm run bench:mixed           # mixed-workload only

# Tune parameters:
BENCH_DURATION=60 BENCH_CONNECTIONS=20 npm run bench
```
