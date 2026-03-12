import { TelemetryIngestRequest, TelemetryMqttRequest } from "@/entities/telemetry";
import { telemetryApi } from "@/features/telemetry/api";

export const simulatorApi = {
  sendRest: (payload: TelemetryIngestRequest) => telemetryApi.ingestByRest(payload),
  sendMqtt: (payload: TelemetryMqttRequest) => telemetryApi.ingestByMqttProxy(payload),
};
