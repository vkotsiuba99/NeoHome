import styles from "./LineChart.module.scss";

type ChartPoint = {
  xLabel: string;
  value: number;
};

type Props = {
  points: ChartPoint[];
  color?: string;
  height?: number;
};

const buildPath = (points: ChartPoint[], width: number, height: number, padding: number) => {
  if (points.length === 0) {
    return "";
  }
  const values = points.map((point) => point.value);
  const min = Math.min(...values);
  const max = Math.max(...values);
  const safeRange = max - min || 1;

  return points
    .map((point, index) => {
      const x = padding + (index * (width - padding * 2)) / Math.max(points.length - 1, 1);
      const y = height - padding - ((point.value - min) * (height - padding * 2)) / safeRange;
      return `${index === 0 ? "M" : "L"} ${x} ${y}`;
    })
    .join(" ");
};

export const LineChart = ({ points, color = "#1f6feb", height = 180 }: Props) => {
  const width = 640;
  const padding = 18;
  const path = buildPath(points, width, height, padding);

  return (
    <div className={styles.wrap}>
      <svg viewBox={`0 0 ${width} ${height}`} className={styles.chart} role="img" aria-label="line chart">
        <rect x={0} y={0} width={width} height={height} className={styles.bg} />
        {path ? <path d={path} stroke={color} strokeWidth={2.5} fill="none" /> : null}
      </svg>
      <div className={styles.caption}>
        {points.length > 0 ? `${points[0].xLabel} ... ${points[points.length - 1].xLabel}` : "No data"}
      </div>
    </div>
  );
};
