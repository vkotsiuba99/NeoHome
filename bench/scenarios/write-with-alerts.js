import autocannon from "autocannon";

export const name = "write-with-alerts";
export const description = "Telemetry writes that trigger alerts (3 extra table writes per request)";
export const cassandraOps = "5 round-trips: GetDevice + AddTelemetry(batch 2) + UpdateDevice(batch 2) + GetThresholds + CreateAlert(batch 3)";

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
          counter++;
          req.body = JSON.stringify({
            deviceId: dev.deviceId,
            recordedAt: Date.now(),
            metricType: "temperature",
            metricValue: 42.0,
            unit: "C",
            roomName: "BenchRoom",
            locationName: "BenchLocation",
            batteryLevel: 85,
            signalStrength: -40,
          });
          return req;
        },
      },
    ],
  });
}
