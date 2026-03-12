import { Link } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useAlerts, useResolveAlert } from "@/features/alerts/hooks";
import { useCreateDevice, useDevices } from "@/features/devices/hooks";
import { useLatestTelemetry } from "@/features/telemetry/hooks";
import { CreateDeviceValues, createDeviceSchema } from "@/features/devices/schema";
import { Card } from "@/shared/ui/Card/Card";
import { Input } from "@/shared/ui/Input/Input";
import { Button } from "@/shared/ui/Button/Button";
import { Badge } from "@/shared/ui/Badge/Badge";
import { Spinner } from "@/shared/ui/Spinner/Spinner";
import { formatDateTime } from "@/shared/utils/format";
import { getErrorMessage } from "@/shared/utils/error";
import styles from "./DashboardPage.module.scss";

type DeviceTileProps = {
  deviceId: number;
  title: string;
  type: string;
  room: string;
  status: string;
};

const DeviceTile = ({ deviceId, title, type, room, status }: DeviceTileProps) => {
  const latest = useLatestTelemetry(deviceId);
  const latestRows = latest.data?.telemetry ?? [];

  return (
    <Card className={styles.deviceCard}>
      <div className={styles.deviceHead}>
        <div>
          <h3>{title}</h3>
          <p className="muted">
            {type} - {room}
          </p>
        </div>
        <Badge tone={status.toLowerCase() === "online" ? "success" : "neutral"}>{status}</Badge>
      </div>
      <div className={styles.metrics}>
        {latest.isLoading ? (
          <Spinner label="Fetching latest..." />
        ) : latestRows.length > 0 ? (
          latestRows.slice(0, 3).map((metric) => (
            <div key={`${metric.telemetryId}-${metric.metricType}`} className={styles.metric}>
              <span>{metric.metricType}</span>
              <strong>
                {metric.metricValue} {metric.unit}
              </strong>
            </div>
          ))
        ) : (
          <p className="muted">No telemetry yet.</p>
        )}
      </div>
      <Link to={`/devices/${deviceId}`} className={styles.deviceLink}>
        Open details
      </Link>
    </Card>
  );
};

export const DashboardPage = () => {
  const devicesQuery = useDevices();
  const alertsQuery = useAlerts();
  const createDevice = useCreateDevice();
  const resolveAlert = useResolveAlert();

  const form = useForm<CreateDeviceValues>({
    resolver: zodResolver(createDeviceSchema),
    defaultValues: {
      deviceName: "",
      deviceType: "sensor",
      roomName: "Living Room",
      locationId: 1,
      locationName: "Home",
      status: "online",
    },
  });

  const onCreateDevice = async (values: CreateDeviceValues) => {
    try {
      await createDevice.mutateAsync(values);
      form.reset({
        deviceName: "",
        deviceType: values.deviceType,
        roomName: values.roomName,
        locationId: values.locationId,
        locationName: values.locationName,
        status: values.status,
      });
    } catch {
      // handled below
    }
  };

  const activeAlerts = (alertsQuery.data?.alerts ?? []).filter((item) => !item.isResolved);

  return (
    <div className="container page">
      <h2>Dashboard</h2>

      <div className="grid-2">
        <Card>
          <h3>Add device</h3>
          <form className={styles.form} onSubmit={form.handleSubmit(onCreateDevice)}>
            <Input
              label="Device name"
              {...form.register("deviceName")}
              error={form.formState.errors.deviceName?.message}
            />
            <Input
              label="Device type"
              {...form.register("deviceType")}
              error={form.formState.errors.deviceType?.message}
            />
            <Input label="Room" {...form.register("roomName")} error={form.formState.errors.roomName?.message} />
            <Input
              label="Location ID"
              type="number"
              {...form.register("locationId")}
              error={form.formState.errors.locationId?.message}
            />
            <Input
              label="Location name"
              {...form.register("locationName")}
              error={form.formState.errors.locationName?.message}
            />
            <Input label="Status" {...form.register("status")} error={form.formState.errors.status?.message} />
            <Button type="submit" loading={createDevice.isPending}>
              Add device
            </Button>
            {createDevice.isError ? <p className={styles.error}>{getErrorMessage(createDevice.error)}</p> : null}
          </form>
        </Card>

        <Card>
          <h3>Active alerts</h3>
          {alertsQuery.isLoading ? (
            <Spinner label="Loading alerts..." />
          ) : activeAlerts.length === 0 ? (
            <p className="muted">No active alerts right now.</p>
          ) : (
            <div className={styles.alerts}>
              {activeAlerts.map((alert) => (
                <div key={alert.alertId} className={styles.alertRow}>
                  <div>
                    <p className={styles.alertMessage}>{alert.message}</p>
                    <small className="muted">
                      Device #{alert.deviceId} - {formatDateTime(alert.createdAt)}
                    </small>
                  </div>
                  <div className={styles.alertActions}>
                    <Badge tone="danger">{alert.severity}</Badge>
                    <Button
                      variant="secondary"
                      onClick={() => resolveAlert.mutate(alert.alertId)}
                      loading={resolveAlert.isPending}
                    >
                      Resolve
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>
      </div>

      <section className={styles.devicesSection}>
        <h3>Devices</h3>
        {devicesQuery.isLoading ? (
          <Spinner label="Loading devices..." />
        ) : (devicesQuery.data?.devices.length ?? 0) === 0 ? (
          <p className="muted">No devices yet. Add your first one above.</p>
        ) : (
          <div className={styles.devicesGrid}>
            {devicesQuery.data?.devices.map((device) => (
              <DeviceTile
                key={device.deviceId}
                deviceId={device.deviceId}
                title={device.deviceName}
                type={device.deviceType}
                room={device.roomName}
                status={device.status}
              />
            ))}
          </div>
        )}
      </section>
    </div>
  );
};
