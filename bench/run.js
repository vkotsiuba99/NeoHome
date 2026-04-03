import fs from "fs";
import path from "path";
import autocannon from "autocannon";
import { config } from "./config.js";
import { seed } from "./seed.js";

import * as writeTelemetry from "./scenarios/write-telemetry.js";
import * as writeWithAlerts from "./scenarios/write-with-alerts.js";
import * as readTelemetry from "./scenarios/read-telemetry.js";
import * as readLatest from "./scenarios/read-latest.js";
import * as mixedWorkload from "./scenarios/mixed-workload.js";

const ALL_SCENARIOS = [writeTelemetry, writeWithAlerts, readTelemetry, readLatest, mixedWorkload];

function parseArgs() {
  const args = process.argv.slice(2);
  const idx = args.indexOf("--scenario");
  if (idx !== -1 && args[idx + 1]) {
    return args[idx + 1];
  }
  return null;
}

function extractMetrics(result) {
  return {
    name: result._name,
    description: result._description,
    cassandraOps: result._cassandraOps,
    requests: {
      total: result.requests.total,
      average: Math.round(result.requests.average),
      mean: Math.round(result.requests.mean),
    },
    latency: {
      p50: result.latency.p50,
      p90: result.latency.p90,
      p99: result.latency.p99,
      p999: result.latency.p999,
      avg: Math.round(result.latency.average * 100) / 100,
      max: result.latency.max,
    },
    throughput: {
      avgMBs: Math.round((result.throughput.average / 1024 / 1024) * 100) / 100,
    },
    errors: result.errors,
    timeouts: result.timeouts,
    non2xx: result.non2xx,
    duration: result.duration,
    connections: result.connections,
    pipelining: result.pipelining,
  };
}

function printSummaryTable(allMetrics) {
  console.log("\n" + "=".repeat(100));
  console.log("BENCHMARK SUMMARY");
  console.log("=".repeat(100));
  console.log(
    [
      "Scenario".padEnd(24),
      "Req/s".padStart(8),
      "p50 ms".padStart(9),
      "p90 ms".padStart(9),
      "p99 ms".padStart(9),
      "Avg ms".padStart(9),
      "Max ms".padStart(9),
      "Errors".padStart(8),
      "Non-2xx".padStart(9),
    ].join(" | ")
  );
  console.log("-".repeat(100));

  for (const m of allMetrics) {
    console.log(
      [
        m.name.padEnd(24),
        String(m.requests.average).padStart(8),
        String(m.latency.p50).padStart(9),
        String(m.latency.p90).padStart(9),
        String(m.latency.p99).padStart(9),
        String(m.latency.avg).padStart(9),
        String(m.latency.max).padStart(9),
        String(m.errors).padStart(8),
        String(m.non2xx).padStart(9),
      ].join(" | ")
    );
  }
  console.log("=".repeat(100));
}

async function runScenario(scenario, seedCtx) {
  console.log(`\n--- Warm-up: ${scenario.name} (${config.WARMUP_SECONDS}s) ---`);
  const warmup = scenario.run(config, seedCtx, config.WARMUP_SECONDS);
  await new Promise((resolve) => {
    warmup.on("done", resolve);
    warmup.on("error", resolve);
  });

  console.log(`--- Benchmark: ${scenario.name} (${config.DURATION_SECONDS}s, ${config.CONNECTIONS} connections) ---`);
  console.log(`    ${scenario.description}`);
  console.log(`    Cassandra: ${scenario.cassandraOps}`);

  const instance = scenario.run(config, seedCtx, config.DURATION_SECONDS);
  const result = await new Promise((resolve, reject) => {
    instance.on("done", resolve);
    instance.on("error", reject);
  });

  result._name = scenario.name;
  result._description = scenario.description;
  result._cassandraOps = scenario.cassandraOps;

  const metrics = extractMetrics(result);

  autocannon.printResult(result);

  const ts = new Date().toISOString().replace(/[:.]/g, "-");
  const outPath = path.join(config.RESULTS_DIR, `${scenario.name}-${ts}.json`);
  fs.mkdirSync(config.RESULTS_DIR, { recursive: true });
  fs.writeFileSync(outPath, JSON.stringify(metrics, null, 2));
  console.log(`    Results saved: ${outPath}`);

  return metrics;
}

async function main() {
  const scenarioFilter = parseArgs();

  let scenarios = ALL_SCENARIOS;
  if (scenarioFilter) {
    scenarios = ALL_SCENARIOS.filter((s) => s.name === scenarioFilter);
    if (scenarios.length === 0) {
      console.error(`Unknown scenario: ${scenarioFilter}`);
      console.error(`Available: ${ALL_SCENARIOS.map((s) => s.name).join(", ")}`);
      process.exit(1);
    }
  }

  console.log("=".repeat(60));
  console.log("NeoHome Cassandra Benchmark Suite");
  console.log("=".repeat(60));
  console.log(`Target:      ${config.BASE_URL}`);
  console.log(`Duration:    ${config.DURATION_SECONDS}s per scenario`);
  console.log(`Warm-up:     ${config.WARMUP_SECONDS}s per scenario`);
  console.log(`Connections: ${config.CONNECTIONS}`);
  console.log(`Devices:     ${config.DEVICE_COUNT}`);
  console.log(`Scenarios:   ${scenarios.map((s) => s.name).join(", ")}`);
  console.log();

  console.log("--- Seeding test data ---");
  const seedCtx = await seed(config);
  console.log(`  Token:   ${seedCtx.token.slice(0, 20)}...`);
  console.log(`  Devices: ${seedCtx.devices.map((d) => d.deviceId).join(", ")}`);

  const allMetrics = [];
  for (const scenario of scenarios) {
    const metrics = await runScenario(scenario, seedCtx);
    allMetrics.push(metrics);
  }

  printSummaryTable(allMetrics);

  const hasErrors = allMetrics.some((m) => m.errors > 0 || m.non2xx > 0);
  if (hasErrors) {
    console.log("\nWARNING: Some scenarios had errors or non-2xx responses!");
  }
  process.exit(hasErrors ? 1 : 0);
}

main().catch((err) => {
  console.error("Benchmark failed:", err.message);
  process.exit(1);
});
