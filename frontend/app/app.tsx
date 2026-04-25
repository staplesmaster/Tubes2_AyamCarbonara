"use client";

import { useState } from "react";
import type { LCARequest, TraversalStats, TraversalStep, TreeNode } from "@/app/lib/types";
import Sidebar from "@/app/components/Sidebar";
import TreeCanvas from "@/app/components/TreeCanvas";
import styles from "./page.module.css";

export default function Page() {
  const [tree, setTree] = useState<TreeNode | null>(null);
  const [steps, setSteps] = useState<TraversalStep[]>([]);
  const [stats, setStats] = useState<TraversalStats | null>(null);
  const [lcaSource, setLcaSource] = useState<Omit<LCARequest, "nodeA" | "nodeB"> | null>(null);
  const [runId, setRunId] = useState(0);

  function handleSteps(nextSteps: TraversalStep[]) {
    setSteps(nextSteps);
    setRunId((id) => id + 1);
  }

  return (
    <div className={styles.app}>
      <Sidebar
        stats={stats}
        steps={steps}
        setTree={setTree}
        setSteps={handleSteps}
        setStats={setStats}
        setLcaSource={setLcaSource}
      />
      <main className={styles.mainPanel}>
        <TreeCanvas key={runId} tree={tree} steps={steps} lcaSource={lcaSource} />
      </main>
    </div>
  );
}
