import autocannon from "autocannon";

export const name = "mixed-workload";
export const description = "80% writes / 20% reads simulating realistic IoT traffic";
export const cassandraOps = "Writes: 4 round-trips each | Reads: 2 round-trips each | ~3.6 avg CQL ops per request";

export function run(config, seedCtx, duration) {
  let counter = 0;

  const writeRequest = {
    method: "POST",
    path: "/api/v1/telemetry",
    headers: { "Content-Type": "application/json" },
    setupRequest(req) {
      const dev = seedCtx.devices[counter % seedCtx.devices.length];
      const metric = config.METRIC_TYPES[counter % config.METRIC_TYPES.length];
      counter++;
      req.body = JSON.stringify({
        deviceId: dev.deviceId,
        recordedAt: Date.now(),
        metricType: metric,
        metricValue: 20 + Math.random() * 5,
        unit: "C",
        roomName: "BenchRoom",
        locationName: "BenchLocation",
        batteryLevel: 95,
        signalStrength: -45,
      });
      return req;
    },
  };

  const dev = seedCtx.devices[0];
  const readRequest = {
    method: "GET",
    path: `/api/v1/devices/${dev.deviceId}/telemetry?limit=100`,
    headers: { Authorization: `Bearer ${seedCtx.token}` },
  };

  return autocannon({
    url: config.BASE_URL,
    connections: config.CONNECTIONS,
    pipelining: config.PIPELINING,
    duration,
    requests: [writeRequest, writeRequest, writeRequest, writeRequest, readRequest],
  });
}
