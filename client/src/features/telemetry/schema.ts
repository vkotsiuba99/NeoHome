import { z } from "zod";

export const telemetryIngestSchema = z.object({
  deviceId: z.coerce.number().int().positive(),
  recordedAt: z.coerce.number().int().positive().optional(),
  metricType: z.string().min(1),
  metricValue: z.coerce.number(),
  unit: z.string().min(1),
  roomName: z.string().min(1),
  locationName: z.string().min(1),
  batteryLevel: z.coerce.number().int().optional(),
  signalStrength: z.coerce.number().int().optional(),
});

export type TelemetryIngestFormValues = z.infer<typeof telemetryIngestSchema>;

export const telemetryMqttSchema = z.object({
  topic: z.string().min(1),
  payload: telemetryIngestSchema,
});

export type TelemetryMqttFormValues = z.infer<typeof telemetryMqttSchema>;
