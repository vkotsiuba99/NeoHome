import { z } from "zod";
import { telemetryIngestSchema } from "@/features/telemetry/schema";

export const simulatorRestSchema = telemetryIngestSchema;
export type SimulatorRestValues = z.infer<typeof simulatorRestSchema>;

export const simulatorMqttSchema = z.object({
  topic: z.string().min(1),
  payload: telemetryIngestSchema,
});

export type SimulatorMqttValues = z.infer<typeof simulatorMqttSchema>;
