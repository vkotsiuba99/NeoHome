import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { devicesApi } from "./api";

export const useDevices = () => {
  return useQuery({
    queryKey: ["devices"],
    queryFn: devicesApi.list,
  });
};

export const useCreateDevice = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: devicesApi.create,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["devices"] });
    },
  });
};

export const useDeviceThresholds = (deviceId?: number) => {
  return useQuery({
    queryKey: ["devices", deviceId, "thresholds"],
    queryFn: () => devicesApi.thresholds(deviceId!),
    enabled: Boolean(deviceId),
  });
};

export const usePutThresholds = (deviceId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: devicesApi.putThresholds.bind(null, deviceId),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["devices", deviceId, "thresholds"] });
      void queryClient.invalidateQueries({ queryKey: ["alerts"] });
    },
  });
};
