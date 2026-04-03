const api = async (url, opts = {}) => {
  const res = await fetch(url, {
    headers: { "Content-Type": "application/json", ...opts.headers },
    ...opts,
  });
  const text = await res.text();
  const body = text ? JSON.parse(text) : null;
  if (!res.ok) throw new Error(`Seed failed: ${opts.method || "GET"} ${url} → ${res.status} ${text}`);
  return body;
};

export async function seed(config) {
  const ts = Date.now();
  const email = `bench-${ts}@test.com`;
  const login = `bench-${ts}`;
  const base = config.BASE_URL;

  console.log("  Registering user...");
  await api(`${base}/api/v1/auth/register`, {
    method: "POST",
    body: JSON.stringify({ email, password: "benchbench1", login, phone: "00000000" }),
  });

  console.log("  Logging in...");
  const auth = await api(`${base}/api/v1/auth/login/email`, {
    method: "POST",
    body: JSON.stringify({ email, password: "benchbench1" }),
  });
  const token = auth.accessToken;

  const devices = [];
  for (let i = 0; i < config.DEVICE_COUNT; i++) {
    console.log(`  Creating device ${i + 1}/${config.DEVICE_COUNT}...`);
    const dev = await api(`${base}/api/v1/devices`, {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify({
        deviceName: `Bench Sensor ${i + 1}`,
        deviceType: "sensor",
        roomName: `Room-${i + 1}`,
        locationId: i + 1,
        locationName: `Location-${i + 1}`,
        status: "online",
      }),
    });
    devices.push({ deviceId: dev.device.deviceId, locationId: i + 1 });

    console.log(`  Setting thresholds for device ${dev.device.deviceId}...`);
    await api(`${base}/api/v1/devices/${dev.device.deviceId}/thresholds`, {
      method: "PUT",
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify({
        thresholds: [
          { metricType: "temperature", minValue: 10, maxValue: 35, severity: "critical" },
          { metricType: "humidity", minValue: 15, maxValue: 80, severity: "warning" },
        ],
      }),
    });
  }

  const now = Date.now();
  const hourAgo = now - 60 * 60 * 1000;
  const rowsPerDevice = 50;
  console.log(`  Seeding ${rowsPerDevice} telemetry rows per device...`);

  for (const dev of devices) {
    for (let j = 0; j < rowsPerDevice; j++) {
      const recordedAt = hourAgo + Math.floor((j / rowsPerDevice) * 60 * 60 * 1000);
      const metricType = config.METRIC_TYPES[j % config.METRIC_TYPES.length];
      await api(`${base}/api/v1/telemetry`, {
        method: "POST",
        body: JSON.stringify({
          deviceId: dev.deviceId,
          recordedAt,
          metricType,
          metricValue: 20 + Math.random() * 5,
          unit: metricType === "temperature" ? "C" : "%",
          roomName: "BenchRoom",
          locationName: "BenchLocation",
          batteryLevel: 95,
          signalStrength: -45,
        }),
      });
    }
  }

  console.log(`  Seed complete: ${devices.length} devices, ${devices.length * rowsPerDevice} telemetry rows`);
  return { token, email, devices };
}
