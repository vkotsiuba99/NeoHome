import { http } from "@/shared/api";
import {
  CreateDeviceRequest,
  DeviceWrapResponse,
  DevicesResponse,
  PutThresholdsRequest,
  ThresholdsResponse,
  UpdateDeviceRequest,
} from "@/entities/device";

export const devicesApi = {
  list: async () => {
    const response = await http.get<DevicesResponse>("/devices", { authMode: "required" });
    return response.data;
  },
  create: async (payload: CreateDeviceRequest) => {
    const response = await http.post<DeviceWrapResponse>("/devices", payload, { authMode: "required" });
    return response.data;
  },
  update: async (deviceId: number, payload: UpdateDeviceRequest) => {
    const response = await http.put<DeviceWrapResponse>(`/devices/${deviceId}`, payload, { authMode: "required" });
    return response.data;
  },
  thresholds: async (deviceId: number) => {
    const response = await http.get<ThresholdsResponse>(`/devices/${deviceId}/thresholds`, { authMode: "required" });
    return response.data;
  },
  putThresholds: async (deviceId: number, payload: PutThresholdsRequest) => {
    const response = await http.put<ThresholdsResponse>(`/devices/${deviceId}/thresholds`, payload, {
      authMode: "required",
    });
    return response.data;
  },
};
