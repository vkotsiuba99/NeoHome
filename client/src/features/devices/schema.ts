import { z } from "zod";

export const createDeviceSchema = z.object({
  deviceName: z.string().min(2, "Device name is required"),
  deviceType: z.string().min(2, "Device type is required"),
  roomName: z.string().min(1, "Room name is required"),
  locationId: z.coerce.number().int().positive("Location ID is required"),
  locationName: z.string().min(1, "Location name is required"),
  status: z.string().min(1, "Status is required"),
});

export type CreateDeviceValues = z.infer<typeof createDeviceSchema>;

export const thresholdsSchema = z.object({
  thresholds: z
    .array(
      z.object({
        metricType: z.string().min(1, "Metric type is required"),
        minValue: z
          .union([z.coerce.number(), z.nan()])
          .optional()
          .transform((value) => (Number.isNaN(value) ? undefined : value)),
        maxValue: z
          .union([z.coerce.number(), z.nan()])
          .optional()
          .transform((value) => (Number.isNaN(value) ? undefined : value)),
        severity: z.string().min(1, "Severity is required"),
      }),
    )
    .min(1, "At least one threshold is required"),
});

export type ThresholdsFormValues = z.infer<typeof thresholdsSchema>;
