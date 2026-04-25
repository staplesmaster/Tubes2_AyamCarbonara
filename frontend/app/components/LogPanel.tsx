import type { TraversalStep } from "@/app/lib/types";
import styles from "./LogPanel.module.css";

export default function LogPanel({ steps }: { steps: TraversalStep[] }) {
  if (!steps || steps.length === 0) {
    return (
      <div className={styles.emptyPanel}>
        <div className={styles.label}>Traversal Log</div>
        No traversal log yet.
      </div>
    );
  }

  return (
    <div className={styles.panel}>
      <div className={styles.label}>Traversal Log</div>
      {steps.map((step, index) => (
        <div
          key={`${step.step}-${step.nodeId ?? index}-${index}`}
          className={cx(styles.entry, step.isMatch ? styles.matched : styles.visited)}
        >
          <span className={styles.step}>{step.step}</span>
          <span className={styles.tag}>&lt;{step.tag}&gt;</span>
          <span className={styles.status}>{step.isMatch ? "✓ matched" : "visited"}</span>
        </div>
      ))}
    </div>
  );
}

function cx(...classes: Array<string | false | null | undefined>) {
  return classes.filter(Boolean).join(" ");
}
