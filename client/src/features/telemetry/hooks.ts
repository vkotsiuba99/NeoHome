import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { telemetryApi } from "./api";
import { TelemetryQuery } from "@/entities/telemetry";

export const useLatestTelemetry = (deviceId?: number) => {
  return useQuery({
    queryKey: ["devices", deviceId, "latest"],
    queryFn: () => telemetryApi.latestByDevice(deviceId!),
    enabled: Boolean(deviceId),
    refetchInterval: 30_000,
  });
};

export const useTelemetryHistory = (deviceId?: number, query?: TelemetryQuery) => {
  return useQuery({
    queryKey: ["devices", deviceId, "telemetry", query],
    queryFn: () => telemetryApi.listByDevice(deviceId!, query ?? {}),
    enabled: Boolean(deviceId),
  });
};

export const useIngestTelemetryRest = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: telemetryApi.ingestByRest,
    onSuccess: (_, payload) => {
      void queryClient.invalidateQueries({ queryKey: ["devices", payload.deviceId, "latest"] });
      void queryClient.invalidateQueries({ queryKey: ["devices", payload.deviceId, "telemetry"] });
      void queryClient.invalidateQueries({ queryKey: ["alerts"] });
    },
  });
};

export const useIngestTelemetryMqtt = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: telemetryApi.ingestByMqttProxy,
    onSuccess: (_, payload) => {
      const deviceId = payload.payload.deviceId;
      if (deviceId) {
        void queryClient.invalidateQueries({ queryKey: ["devices", deviceId, "latest"] });
        void queryClient.invalidateQueries({ queryKey: ["devices", deviceId, "telemetry"] });
      }
      void queryClient.invalidateQueries({ queryKey: ["alerts"] });
    },
  });
};
