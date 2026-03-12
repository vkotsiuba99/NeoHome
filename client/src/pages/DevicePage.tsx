import { useEffect, useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import { useFieldArray, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useDeviceThresholds, useDevices, usePutThresholds } from "@/features/devices/hooks";
import { ThresholdsFormValues, thresholdsSchema } from "@/features/devices/schema";
import { useTelemetryHistory } from "@/features/telemetry/hooks";
import { Card } from "@/shared/ui/Card/Card";
import { Input } from "@/shared/ui/Input/Input";
import { Button } from "@/shared/ui/Button/Button";
import { Spinner } from "@/shared/ui/Spinner/Spinner";
import { LineChart } from "@/shared/ui/LineChart/LineChart";
import { formatDateTime, metricColor } from "@/shared/utils/format";
import { getErrorMessage } from "@/shared/utils/error";
import styles from "./DevicePage.module.scss";

const PERIODS = [
  { key: "24h", label: "24h", ms: 24 * 60 * 60 * 1000 },
  { key: "7d", label: "7d", ms: 7 * 24 * 60 * 60 * 1000 },
] as const;

export const DevicePage = () => {
  const params = useParams<{ deviceId: string }>();
  const deviceId = Number(params.deviceId);
  const [periodKey, setPeriodKey] = useState<(typeof PERIODS)[number]["key"]>("24h");

  const nowMs = Date.now();
  const selectedPeriod = PERIODS.find((item) => item.key === periodKey) ?? PERIODS[0];
  const fromMs = nowMs - selectedPeriod.ms;

  const devicesQuery = useDevices();
  const thresholdsQuery = useDeviceThresholds(deviceId);
  const upsertThresholds = usePutThresholds(deviceId);

  const temperatureTelemetry = useTelemetryHistory(deviceId, {
    metricType: "temperature",
    from: fromMs,
    to: nowMs,
    limit: 500,
  });
  const humidityTelemetry = useTelemetryHistory(deviceId, {
    metricType: "humidity",
    from: fromMs,
    to: nowMs,
    limit: 500,
  });
  const tableTelemetry = useTelemetryHistory(deviceId, {
    from: fromMs,
    to: nowMs,
    limit: 100,
  });

  const form = useForm<ThresholdsFormValues>({
    resolver: zodResolver(thresholdsSchema),
    defaultValues: {
      thresholds: [
        { metricType: "temperature", minValue: undefined, maxValue: 29, severity: "critical" },
        { metricType: "humidity", minValue: 20, maxValue: 70, severity: "warning" },
      ],
    },
  });
  const fieldArray = useFieldArray({
    control: form.control,
    name: "thresholds",
  });

  useEffect(() => {
    const current = thresholdsQuery.data?.thresholds;
    if (!current || current.length === 0) {
      return;
    }

    form.reset({
      thresholds: current.map((item) => ({
        metricType: item.metricType,
        minValue: item.hasMinValue ? item.minValue : undefined,
        maxValue: item.hasMaxValue ? item.maxValue : undefined,
        severity: item.severity,
      })),
    });
  }, [form, thresholdsQuery.data?.thresholds]);

  const device = useMemo(() => {
    return devicesQuery.data?.devices.find((item) => item.deviceId === deviceId);
  }, [deviceId, devicesQuery.data?.devices]);

  const onSubmitThresholds = async (values: ThresholdsFormValues) => {
    try {
      await upsertThresholds.mutateAsync({
        thresholds: values.thresholds.map((item) => ({
          metricType: item.metricType,
          minValue: item.minValue,
          maxValue: item.maxValue,
          severity: item.severity,
        })),
      });
    } catch {
      // message rendered below
    }
  };

  const temperaturePoints = (temperatureTelemetry.data?.telemetry ?? [])
    .slice()
    .reverse()
    .map((item) => ({
      xLabel: new Date(item.recordedAt).toLocaleTimeString(),
      value: item.metricValue,
    }));

  const humidityPoints = (humidityTelemetry.data?.telemetry ?? [])
    .slice()
    .reverse()
    .map((item) => ({
      xLabel: new Date(item.recordedAt).toLocaleTimeString(),
      value: item.metricValue,
    }));

  if (!Number.isFinite(deviceId) || deviceId <= 0) {
    return (
      <div className="container page">
        <Card>
          <p>Invalid device ID.</p>
        </Card>
      </div>
    );
  }

  return (
    <div className="container page">
      <h2>Device details</h2>
      <Card>
        {devicesQuery.isLoading ? (
          <Spinner label="Loading device..." />
        ) : device ? (
          <div className={styles.deviceMeta}>
            <p>
              <strong>{device.deviceName}</strong> ({device.deviceType})
            </p>
            <p className="muted">
              Room: {device.roomName}, location: {device.locationName}, status: {device.status}
            </p>
          </div>
        ) : (
          <p className="muted">Device not found in your account.</p>
        )}
      </Card>

      <div className={styles.periods}>
        {PERIODS.map((period) => (
          <Button
            key={period.key}
            variant={period.key === periodKey ? "primary" : "secondary"}
            onClick={() => setPeriodKey(period.key)}
          >
            {period.label}
          </Button>
        ))}
      </div>

      <div className="grid-2">
        <Card>
          <h3>Temperature graph ({selectedPeriod.label})</h3>
          {temperatureTelemetry.isLoading ? (
            <Spinner label="Loading temperature..." />
          ) : (
            <LineChart points={temperaturePoints} color={metricColor("temperature")} />
          )}
        </Card>

        <Card>
          <h3>Humidity graph ({selectedPeriod.label})</h3>
          {humidityTelemetry.isLoading ? (
            <Spinner label="Loading humidity..." />
          ) : (
            <LineChart points={humidityPoints} color={metricColor("humidity")} />
          )}
        </Card>
      </div>

      <Card>
        <h3>Thresholds</h3>
        <form className={styles.thresholdsForm} onSubmit={form.handleSubmit(onSubmitThresholds)}>
          {fieldArray.fields.map((field, index) => (
            <div key={field.id} className={styles.thresholdRow}>
              <Input
                label="Metric"
                {...form.register(`thresholds.${index}.metricType`)}
                error={form.formState.errors.thresholds?.[index]?.metricType?.message}
              />
              <Input
                label="Min"
                type="number"
                step="0.01"
                {...form.register(`thresholds.${index}.minValue`)}
                error={form.formState.errors.thresholds?.[index]?.minValue?.message}
              />
              <Input
                label="Max"
                type="number"
                step="0.01"
                {...form.register(`thresholds.${index}.maxValue`)}
                error={form.formState.errors.thresholds?.[index]?.maxValue?.message}
              />
              <Input
                label="Severity"
                {...form.register(`thresholds.${index}.severity`)}
                error={form.formState.errors.thresholds?.[index]?.severity?.message}
              />
              <Button variant="ghost" type="button" onClick={() => fieldArray.remove(index)}>
                Remove
              </Button>
            </div>
          ))}
          <div className={styles.thresholdActions}>
            <Button
              variant="secondary"
              type="button"
              onClick={() => fieldArray.append({ metricType: "", minValue: undefined, maxValue: undefined, severity: "warning" })}
            >
              Add threshold
            </Button>
            <Button type="submit" loading={upsertThresholds.isPending}>
              Save thresholds
            </Button>
          </div>
          {upsertThresholds.isError ? <p className={styles.error}>{getErrorMessage(upsertThresholds.error)}</p> : null}
        </form>
      </Card>

      <Card>
        <h3>Telemetry history ({selectedPeriod.label})</h3>
        {tableTelemetry.isLoading ? (
          <Spinner label="Loading history..." />
        ) : (
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>Time</th>
                  <th>Metric</th>
                  <th>Value</th>
                  <th>Room</th>
                </tr>
              </thead>
              <tbody>
                {(tableTelemetry.data?.telemetry ?? []).map((item) => (
                  <tr key={item.telemetryId}>
                    <td>{formatDateTime(item.recordedAt)}</td>
                    <td>{item.metricType}</td>
                    <td>
                      {item.metricValue} {item.unit}
                    </td>
                    <td>{item.roomName}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Card>
    </div>
  );
};
