import autocannon from "autocannon";

export const name = "read-telemetry";
export const description = "Time-series telemetry history reads (LIMIT 100)";
export const cassandraOps = "2 round-trips: GetDevice (ownership check) + ListTelemetry (partition scan with LIMIT)";

export function run(config, seedCtx, duration) {
  const requests = seedCtx.devices.map((dev) => ({
    method: "GET",
    path: `/api/v1/devices/${dev.deviceId}/telemetry?limit=100`,
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
