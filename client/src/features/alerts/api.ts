import { AlertWrapResponse, AlertsQuery, AlertsResponse } from "@/entities/alert";
import { http } from "@/shared/api";

const toQueryString = (query: AlertsQuery) => {
  const params = new URLSearchParams();
  if (query.locationId) params.set("locationId", String(query.locationId));
  if (query.from) params.set("from", String(query.from));
  if (query.to) params.set("to", String(query.to));
  const value = params.toString();
  return value ? `?${value}` : "";
};

export const alertsApi = {
  list: async (query: AlertsQuery = {}) => {
    const response = await http.get<AlertsResponse>(`/alerts${toQueryString(query)}`, { authMode: "required" });
    return response.data;
  },
  resolve: async (alertId: number) => {
    const response = await http.put<AlertWrapResponse>(`/alerts/${alertId}/resolve`, {}, { authMode: "required" });
    return response.data;
  },
};
