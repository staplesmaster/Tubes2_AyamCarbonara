"use client";

import { useState } from "react";
import type { ReactNode } from "react";
import { traverse, uploadHTML } from "@/app/lib/api";
import type {
  Algorithm,
  InputMode,
  LCARequest,
  TraversalStats,
  TraversalStep,
  TreeNode,
} from "@/app/lib/types";
import LogPanel from "@/app/components/LogPanel";
import StatsPanel from "@/app/components/StatsPanel";
import styles from "./Sidebar.module.css";

interface SidebarProps {
  stats: TraversalStats | null;
  steps: TraversalStep[];
  setTree: (tree: TreeNode) => void;
  setSteps: (steps: TraversalStep[]) => void;
  setStats: (stats: TraversalStats) => void;
  setLcaSource: (source: Omit<LCARequest, "nodeA" | "nodeB">) => void;
}

export default function Sidebar({ stats, steps, setTree, setSteps, setStats, setLcaSource }: SidebarProps) {
  const [mode, setMode] = useState<InputMode>("url");
  const [url, setUrl] = useState("");
  const [html, setHTML] = useState("");
  const [uploadedFilename, setUploadedFilename] = useState("");
  const [selector, setSelector] = useState("div > p");
  const [algorithm, setAlgorithm] = useState<Algorithm>("bfs");
  const [parallel, setParallel] = useState(false);
  const [limit, setLimit] = useState(5);
  const [allResult, setAllResult] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const canRun = mode === "url" ? url.trim().length > 0 : html.trim().length > 0;

  async function run() {
    setError(null);
    setLoading(true);

    try {
      const res = await traverse({
        inputMode: mode,
        url,
        html,
        selector,
        algorithm,
        parallel,
        limit,
        allResult,
      });

      if (!res.success) {
        setError(res.error ?? "Unknown error");
        return;
      }

      setTree(res.tree);
      setSteps(res.steps);
      setStats(res.stats);
      setLcaSource({ inputMode: mode, url, html });
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to connect to server");
    } finally {
      setLoading(false);
    }
  }

  async function handleUpload(file: File | undefined) {
    if (!file) return;

    setError(null);
    if (!file.name.toLowerCase().endsWith(".html")) {
      setError("File harus berekstensi .html");
      return;
    }

    try {
      const res = await uploadHTML(file);
      if (!res.success || !res.html) {
        setError(res.error ?? "Gagal upload file");
        return;
      }

      setMode("html");
      setHTML(res.html);
      setUploadedFilename(res.filename ?? file.name);
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Gagal upload file");
    }
  }

  return (
    <aside className={styles.sidebar}>
      <header className={styles.header}>
        <h1 className={styles.title}>DOM Traverser</h1>
        <div className={styles.subtitle}>IF2211 - BFS / DFS Visualizer</div>
      </header>

      <Section label="Input Mode">
        <ButtonGroup>
          {(["url", "html"] as const).map((nextMode) => (
            <button
              key={nextMode}
              onClick={() => setMode(nextMode)}
              className={toggleClass(mode === nextMode, "bfs")}
            >
              {nextMode.toUpperCase()}
            </button>
          ))}
        </ButtonGroup>
      </Section>

      <Section label={mode === "url" ? "Target URL" : "HTML Source"}>
        {mode === "url" ? (
          <input
            className={styles.input}
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="https://example.com"
          />
        ) : (
          <textarea
            className={styles.textarea}
            value={html}
            onChange={(e) => setHTML(e.target.value)}
            placeholder="<div><p>Hello</p></div>"
            rows={5}
          />
        )}
        <input
          className={styles.fileInput}
          type="file"
          accept=".html,text/html"
          onChange={(e) => handleUpload(e.target.files?.[0])}
        />
        {uploadedFilename && <div className={styles.fileHint}>loaded: {uploadedFilename}</div>}
      </Section>

      <Section label="CSS Selector">
        <input
          className={styles.input}
          value={selector}
          onChange={(e) => setSelector(e.target.value)}
          placeholder="div > p"
        />
      </Section>

      <Section label="Algorithm">
        <ButtonGroup>
          <button onClick={() => setAlgorithm("bfs")} className={toggleClass(algorithm === "bfs", "bfs")}>
            BFS
          </button>
          <button onClick={() => setAlgorithm("dfs")} className={toggleClass(algorithm === "dfs", "dfs")}>
            DFS
          </button>
        </ButtonGroup>
      </Section>

      <Section label="Execution Mode">
        <ButtonGroup>
          <button onClick={() => setParallel(false)} className={toggleClass(!parallel, "bfs")}>
            Normal
          </button>
          <button onClick={() => setParallel(true)} className={toggleClass(parallel, "dfs")}>
            Threaded
          </button>
        </ButtonGroup>
      </Section>

      <Section label="Result Limit">
        <div className={styles.limitRow}>
          <span className={styles.smallText}>top</span>
          <input
            className={styles.numberInput}
            type="number"
            min={1}
            max={999}
            value={limit}
            disabled={allResult}
            onChange={(e) => setLimit(Math.max(1, Number(e.target.value) || 1))}
          />
          <button
            onClick={() => setAllResult((value) => !value)}
            className={cx(styles.allBtn, allResult && styles.allActive)}
          >
            all
          </button>
        </div>
      </Section>

      <section className={styles.runSection}>
        <button onClick={run} disabled={loading || !canRun} className={styles.runBtn}>
          {loading ? "Running..." : "Run Traversal"}
        </button>
        {error && <div className={styles.error}>x {error}</div>}
      </section>

      <StatsPanel stats={stats} steps={steps} />
      <LogPanel steps={steps} />
    </aside>
  );
}

function Section({ label, children }: { label: string; children: ReactNode }) {
  return (
    <section className={styles.section}>
      <div className={styles.label}>{label}</div>
      {children}
    </section>
  );
}

function ButtonGroup({ children }: { children: ReactNode }) {
  return <div className={styles.buttonGroup}>{children}</div>;
}

function toggleClass(active: boolean, variant: "bfs" | "dfs") {
  return cx(styles.toggleBtn, active && (variant === "bfs" ? styles.activeBfs : styles.activeDfs));
}

function cx(...classes: Array<string | false | null | undefined>) {
  return classes.filter(Boolean).join(" ");
}
