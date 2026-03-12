import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { alertsApi } from "./api";
import { AlertsQuery } from "@/entities/alert";

export const useAlerts = (query: AlertsQuery = {}) => {
  return useQuery({
    queryKey: ["alerts", query],
    queryFn: () => alertsApi.list(query),
    refetchInterval: 30_000,
  });
};

export const useResolveAlert = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: alertsApi.resolve,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["alerts"] });
    },
  });
};
