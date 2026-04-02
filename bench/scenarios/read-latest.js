import autocannon from "autocannon";

export const name = "read-latest";
export const description = "Latest telemetry per device (scans up to 1000 rows, deduplicates in memory)";
export const cassandraOps = "2 round-trips: GetDevice (ownership check) + GetLatestTelemetry (scan LIMIT 1000 + in-memory dedup)";

export function run(config, seedCtx, duration) {
  const requests = seedCtx.devices.map((dev) => ({
    method: "GET",
    path: `/api/v1/devices/${dev.deviceId}/latest`,
    headers: { Authorization: `Bearer ${seedCtx.token}` },
  }));

  return autocannon({
    url: config.BASE_URL,
    connections: config.CONNECTIONS,
    pipelining: config.PIPELINING,
    duration,
    requests,
  });
}
