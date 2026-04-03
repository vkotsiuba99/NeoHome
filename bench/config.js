import { fileURLToPath } from "url";
import path from "path";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

export const config = {
  BASE_URL: process.env.BENCH_API_URL || "http://localhost:3434",
  DURATION_SECONDS: Number(process.env.BENCH_DURATION) || 30,
  WARMUP_SECONDS: Number(process.env.BENCH_WARMUP) || 5,
  CONNECTIONS: Number(process.env.BENCH_CONNECTIONS) || 10,
  PIPELINING: Number(process.env.BENCH_PIPELINING) || 1,
  DEVICE_COUNT: Number(process.env.BENCH_DEVICES) || 5,
  METRIC_TYPES: ["temperature", "humidity", "pressure", "co2", "light"],
  RESULTS_DIR: path.join(__dirname, "results"),
};
