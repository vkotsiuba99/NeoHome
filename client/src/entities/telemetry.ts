export type Telemetry = {
  telemetryId: number;
  deviceId: number;
  recordedAt: number;
  metricType: string;
  metricValue: number;
  unit: string;
  roomName: string;
  locationName: string;
  batteryLevel: number;
  signalStrength: number;
};

export type TelemetryListResponse = {
  telemetry: Telemetry[];
};

export type TelemetryIngestRequest = {
  deviceId: number;
  recordedAt?: number;
  metricType: string;
  metricValue: number;
  unit: string;
  roomName: string;
  locationName: string;
  batteryLevel?: number;
  signalStrength?: number;
};

export type TelemetryMqttRequest = {
  topic: string;
  payload: TelemetryIngestRequest;
};

export type TelemetryIngestResponse = {
  telemetry: Telemetry;
  alerts: {
    alertId: number;
    locationId: number;
    deviceId: number;
    createdAt: number;
    severity: string;
    message: string;
    isResolved: boolean;
    resolvedAt: number;
  }[];
};

export type TelemetryQuery = {
  metricType?: string;
  from?: number;
  to?: number;
  limit?: number;
};
