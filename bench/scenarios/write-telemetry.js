import autocannon from "autocannon";

export const name = "write-telemetry";
export const description = "Pure telemetry write throughput (no alerts triggered)";
export const cassandraOps = "4 round-trips: GetDevice + AddTelemetry(batch 2) + UpdateDevice(batch 2) + GetThresholds";

export function run(config, seedCtx, duration) {
  let counter = 0;

  return autocannon({
    url: `${config.BASE_URL}/api/v1/telemetry`,
    method: "POST",
    connections: config.CONNECTIONS,
    pipelining: config.PIPELINING,
    duration,
    headers: { "Content-Type": "application/json" },
    requests: [
      {
        method: "POST",
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
      },
    ],
  });
}
