import type { TraversalStats, TraversalStep } from "@/app/lib/types";
import styles from "./StatsPanel.module.css";

export default function StatsPanel({
  stats,
  steps,
}: {
  stats: TraversalStats | null;
  steps: TraversalStep[];
}) {
  if (!stats) return null;

  const cards = [
    { label: "visited", value: stats.visited },
    { label: "matched", value: stats.matched },
    { label: "time", value: `${formatMs(stats.elapsedMs)}ms` },
    { label: "max depth", value: stats.maxDepth },
    { label: "steps", value: steps.length },
  ];

  return (
    <div className={styles.panel}>
      <div className={styles.title}>Statistics</div>
      {cards.map((card) => (
        <div key={card.label} className={styles.card}>
          <div className={styles.cardLabel}>{card.label}</div>
          <div className={styles.cardValue}>{card.value}</div>
        </div>
      ))}
    </div>
  );
}

function formatMs(value: number) {
  if (value === 0) return "0";
  if (value < 0.01) return "<0.01";
  if (value < 1) return value.toFixed(3);
  if (value < 10) return value.toFixed(2);
  return value.toFixed(1);
}
