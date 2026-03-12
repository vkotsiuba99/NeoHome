export const formatDateTime = (unixMilliseconds: number) => {
  if (!unixMilliseconds) {
    return "-";
  }
  return new Date(unixMilliseconds).toLocaleString();
};

export const toIsoFromUnix = (unixMilliseconds: number) => new Date(unixMilliseconds).toISOString();

export const metricColor = (metricType: string) => {
  switch (metricType.toLowerCase()) {
    case "temperature":
      return "#ff7b00";
    case "humidity":
      return "#0077cc";
    case "co2":
      return "#6f42c1";
    default:
      return "#1f6feb";
  }
};
