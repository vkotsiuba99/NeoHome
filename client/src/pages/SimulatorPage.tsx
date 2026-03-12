import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Card } from "@/shared/ui/Card/Card";
import { Input } from "@/shared/ui/Input/Input";
import { Button } from "@/shared/ui/Button/Button";
import { Textarea } from "@/shared/ui/Textarea/Textarea";
import { useIngestTelemetryMqtt, useIngestTelemetryRest } from "@/features/telemetry/hooks";
import { SimulatorMqttValues, SimulatorRestValues, simulatorMqttSchema, simulatorRestSchema } from "@/features/simulator/schema";
import { getErrorMessage } from "@/shared/utils/error";
import styles from "./SimulatorPage.module.scss";

const defaultPayload: SimulatorRestValues = {
  deviceId: 1,
  recordedAt: Date.now(),
  metricType: "temperature",
  metricValue: 24.2,
  unit: "C",
  roomName: "Living Room",
  locationName: "Home",
  batteryLevel: 87,
  signalStrength: -62,
};

const TEMPLATES: Array<{ label: string; value: SimulatorRestValues }> = [
  { label: "Temp normal", value: { ...defaultPayload, metricType: "temperature", metricValue: 23.4, unit: "C" } },
  { label: "Temp critical", value: { ...defaultPayload, metricType: "temperature", metricValue: 34.8, unit: "C" } },
  { label: "Humidity high", value: { ...defaultPayload, metricType: "humidity", metricValue: 78, unit: "%" } },
];

export const SimulatorPage = () => {
  const [lastResponse, setLastResponse] = useState<string>("");

  const restMutation = useIngestTelemetryRest();
  const mqttMutation = useIngestTelemetryMqtt();

  const restForm = useForm<SimulatorRestValues>({
    resolver: zodResolver(simulatorRestSchema),
    defaultValues: defaultPayload,
  });

  const mqttForm = useForm<SimulatorMqttValues>({
    resolver: zodResolver(simulatorMqttSchema),
    defaultValues: {
      topic: "neohome/devices/1/telemetry",
      payload: defaultPayload,
    },
  });

  const fillTemplate = (values: SimulatorRestValues) => {
    const withTime = { ...values, recordedAt: Date.now() };
    restForm.reset(withTime);
    mqttForm.reset({
      topic: `neohome/devices/${withTime.deviceId}/telemetry`,
      payload: withTime,
    });
  };

  const onSendRest = async (values: SimulatorRestValues) => {
    try {
      const response = await restMutation.mutateAsync({
        ...values,
        recordedAt: values.recordedAt || Date.now(),
      });
      setLastResponse(JSON.stringify(response, null, 2));
    } catch (error) {
      setLastResponse(getErrorMessage(error));
    }
  };

  const onSendMqtt = async (values: SimulatorMqttValues) => {
    try {
      const response = await mqttMutation.mutateAsync({
        ...values,
        payload: {
          ...values.payload,
          recordedAt: values.payload.recordedAt || Date.now(),
        },
      });
      setLastResponse(JSON.stringify(response, null, 2));
    } catch (error) {
      setLastResponse(getErrorMessage(error));
    }
  };

  const restCurl = `curl -X POST http://localhost:3434/api/v1/telemetry \\
  -H "Content-Type: application/json" \\
  -d '${JSON.stringify(restForm.getValues())}'`;

  const mqttCurl = `curl -X POST http://localhost:3434/api/v1/telemetry/mqtt \\
  -H "Content-Type: application/json" \\
  -d '${JSON.stringify(mqttForm.getValues())}'`;

  return (
    <div className="container page">
      <h2>Telemetry simulator</h2>
      <p className="muted">
        This page demonstrates how telemetry is sent without a real MQTT device. You can use REST or MQTT proxy.
      </p>

      <Card>
        <h3>Quick templates</h3>
        <div className={styles.templates}>
          {TEMPLATES.map((template) => (
            <Button key={template.label} variant="secondary" onClick={() => fillTemplate(template.value)}>
              {template.label}
            </Button>
          ))}
        </div>
      </Card>

      <div className="grid-2">
        <Card>
          <h3>Send via REST</h3>
          <form className={styles.form} onSubmit={restForm.handleSubmit(onSendRest)}>
            <Input label="Device ID" type="number" {...restForm.register("deviceId")} />
            <Input label="Recorded at (ms)" type="number" {...restForm.register("recordedAt")} />
            <Input label="Metric type" {...restForm.register("metricType")} />
            <Input label="Metric value" type="number" step="0.01" {...restForm.register("metricValue")} />
            <Input label="Unit" {...restForm.register("unit")} />
            <Input label="Room" {...restForm.register("roomName")} />
            <Input label="Location" {...restForm.register("locationName")} />
            <Input label="Battery level" type="number" {...restForm.register("batteryLevel")} />
            <Input label="Signal strength" type="number" {...restForm.register("signalStrength")} />
            <Button type="submit" loading={restMutation.isPending}>
              Send REST payload
            </Button>
          </form>
        </Card>

        <Card>
          <h3>Send via MQTT proxy</h3>
          <form className={styles.form} onSubmit={mqttForm.handleSubmit(onSendMqtt)}>
            <Input label="Topic" {...mqttForm.register("topic")} />
            <Input label="Device ID" type="number" {...mqttForm.register("payload.deviceId")} />
            <Input label="Recorded at (ms)" type="number" {...mqttForm.register("payload.recordedAt")} />
            <Input label="Metric type" {...mqttForm.register("payload.metricType")} />
            <Input label="Metric value" type="number" step="0.01" {...mqttForm.register("payload.metricValue")} />
            <Input label="Unit" {...mqttForm.register("payload.unit")} />
            <Input label="Room" {...mqttForm.register("payload.roomName")} />
            <Input label="Location" {...mqttForm.register("payload.locationName")} />
            <Input label="Battery level" type="number" {...mqttForm.register("payload.batteryLevel")} />
            <Input label="Signal strength" type="number" {...mqttForm.register("payload.signalStrength")} />
            <Button type="submit" loading={mqttMutation.isPending}>
              Send MQTT payload
            </Button>
          </form>
        </Card>
      </div>

      <div className="grid-2">
        <Card>
          <h3>REST example</h3>
          <Textarea value={restCurl} readOnly />
        </Card>
        <Card>
          <h3>MQTT proxy example</h3>
          <Textarea value={mqttCurl} readOnly />
        </Card>
      </div>

      <Card>
        <h3>Last response</h3>
        <Textarea value={lastResponse || "No request yet"} readOnly />
      </Card>
    </div>
  );
};
