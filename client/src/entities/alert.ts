export type Alert = {
  alertId: number;
  locationId: number;
  deviceId: number;
  createdAt: number;
  severity: string;
  message: string;
  isResolved: boolean;
  resolvedAt: number;
};

export type AlertsResponse = {
  alerts: Alert[];
};

export type AlertWrapResponse = {
  alert: Alert;
};

export type AlertsQuery = {
  locationId?: number;
  from?: number;
  to?: number;
};
