export type Device = {
  deviceId: number;
  userId: number;
  deviceName: string;
  deviceType: string;
  roomName: string;
  locationId: number;
  locationName: string;
  status: string;
  lastSeenAt: number;
  lastMetricAt: number;
  batteryLevel: number;
  signalStrength: number;
  addedAt: number;
  updatedAt: number;
};

export type DevicesResponse = {
  devices: Device[];
};

export type DeviceWrapResponse = {
  device: Device;
};

export type CreateDeviceRequest = {
  deviceName: string;
  deviceType: string;
  roomName: string;
  locationId: number;
  locationName: string;
  status: string;
};

export type UpdateDeviceRequest = Partial<CreateDeviceRequest>;

export type Threshold = {
  metricType: string;
  minValue: number;
  hasMinValue: boolean;
  maxValue: number;
  hasMaxValue: boolean;
  severity: string;
  updatedAt: number;
};

export type ThresholdsResponse = {
  thresholds: Threshold[];
};

export type ThresholdPatch = {
  metricType: string;
  minValue?: number;
  maxValue?: number;
  severity: string;
};

export type PutThresholdsRequest = {
  thresholds: ThresholdPatch[];
};
