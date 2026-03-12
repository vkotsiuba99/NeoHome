import { TelemetryIngestRequest, TelemetryIngestResponse, TelemetryListResponse, TelemetryMqttRequest, TelemetryQuery } from "@/entities/telemetry";
import { http } from "@/shared/api";

const toQueryString = (query: TelemetryQuery) => {
  const params = new URLSearchParams();
  if (query.metricType) params.set("metricType", query.metricType);
  if (query.from) params.set("from", String(query.from));
  if (query.to) params.set("to", String(query.to));
  if (query.limit) params.set("limit", String(query.limit));
  const value = params.toString();
  return value ? `?${value}` : "";
};

export const telemetryApi = {
  latestByDevice: async (deviceId: number) => {
    const response = await http.get<TelemetryListResponse>(`/devices/${deviceId}/latest`, { authMode: "required" });
    return response.data;
  },
  listByDevice: async (deviceId: number, query: TelemetryQuery) => {
    const response = await http.get<TelemetryListResponse>(`/devices/${deviceId}/telemetry${toQueryString(query)}`, {
      authMode: "required",
    });
    return response.data;
  },
  ingestByRest: async (payload: TelemetryIngestRequest) => {
    const response = await http.post<TelemetryIngestResponse>("/telemetry", payload, { authMode: "none" });
    return response.data;
  },
  ingestByMqttProxy: async (payload: TelemetryMqttRequest) => {
    const response = await http.post<TelemetryIngestResponse>("/telemetry/mqtt", payload, { authMode: "none" });
    return response.data;
  },
};
