CREATE KEYSPACE IF NOT EXISTS neohome
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};

USE neohome;

CREATE TABLE IF NOT EXISTS users_by_email (
  email text PRIMARY KEY,
  user_id bigint,
  login text,
  phone text,
  password_hash text,
  role text,
  created_at bigint,
  updated_at bigint
);

CREATE TABLE IF NOT EXISTS users_by_login (
  login text PRIMARY KEY,
  user_id bigint,
  email text,
  phone text,
  password_hash text,
  role text,
  created_at bigint,
  updated_at bigint
);

CREATE TABLE IF NOT EXISTS users_by_id (
  user_id bigint PRIMARY KEY,
  email text,
  login text,
  phone text,
  role text,
  created_at bigint,
  updated_at bigint
);

CREATE TABLE IF NOT EXISTS devices_by_user (
  user_id bigint,
  device_id bigint,
  device_name text,
  device_type text,
  room_name text,
  location_id bigint,
  location_name text,
  status text,
  last_seen_at bigint,
  last_metric_at bigint,
  battery_level bigint,
  signal_strength bigint,
  added_at bigint,
  updated_at bigint,
  PRIMARY KEY ((user_id), device_id)
);

CREATE TABLE IF NOT EXISTS devices_by_id (
  device_id bigint PRIMARY KEY,
  user_id bigint,
  device_name text,
  device_type text,
  room_name text,
  location_id bigint,
  location_name text,
  status text,
  last_seen_at bigint,
  last_metric_at bigint,
  battery_level bigint,
  signal_strength bigint,
  added_at bigint,
  updated_at bigint
);

CREATE TABLE IF NOT EXISTS telemetry_by_device (
  device_id bigint,
  recorded_at bigint,
  metric_type text,
  telemetry_id bigint,
  metric_value double,
  unit text,
  room_name text,
  location_name text,
  battery_level bigint,
  signal_strength bigint,
  PRIMARY KEY ((device_id), recorded_at, metric_type, telemetry_id)
) WITH CLUSTERING ORDER BY (recorded_at DESC, metric_type ASC, telemetry_id ASC);

CREATE TABLE IF NOT EXISTS telemetry_by_device_metric (
  device_id bigint,
  metric_type text,
  recorded_at bigint,
  telemetry_id bigint,
  metric_value double,
  unit text,
  room_name text,
  location_name text,
  battery_level bigint,
  signal_strength bigint,
  PRIMARY KEY ((device_id, metric_type), recorded_at, telemetry_id)
) WITH CLUSTERING ORDER BY (recorded_at DESC, telemetry_id ASC);

CREATE TABLE IF NOT EXISTS device_thresholds_by_device (
  device_id bigint,
  metric_type text,
  min_value double,
  max_value double,
  severity text,
  updated_at bigint,
  PRIMARY KEY ((device_id), metric_type)
);

CREATE TABLE IF NOT EXISTS alerts_by_location (
  location_id bigint,
  created_at bigint,
  alert_id bigint,
  device_id bigint,
  severity text,
  message text,
  is_resolved boolean,
  resolved_at bigint,
  PRIMARY KEY ((location_id), created_at, alert_id)
) WITH CLUSTERING ORDER BY (created_at DESC, alert_id ASC);

CREATE TABLE IF NOT EXISTS alerts_by_device (
  device_id bigint,
  created_at bigint,
  alert_id bigint,
  location_id bigint,
  severity text,
  message text,
  is_resolved boolean,
  resolved_at bigint,
  PRIMARY KEY ((device_id), created_at, alert_id)
) WITH CLUSTERING ORDER BY (created_at DESC, alert_id ASC);

CREATE TABLE IF NOT EXISTS alerts_by_id (
  alert_id bigint PRIMARY KEY,
  location_id bigint,
  device_id bigint,
  created_at bigint,
  severity text,
  message text,
  is_resolved boolean,
  resolved_at bigint
);

