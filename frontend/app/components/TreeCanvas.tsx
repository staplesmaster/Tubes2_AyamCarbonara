"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { lca } from "@/app/lib/api";
import type { LCARequest, LCAStep, TraversalStep, TreeNode } from "@/app/lib/types";
import {
  NODE_H,
  NODE_W,
  buildTreeLayout,
  emptyTreeLayout,
  getNodeColor,
  nodeLabel,
} from "@/app/lib/treeHelpers";
import styles from "./TreeCanvas.module.css";

interface TreeCanvasProps {
  tree: TreeNode | null;
  steps: TraversalStep[];
  lcaSource: Omit<LCARequest, "nodeA" | "nodeB"> | null;
}

export default function TreeCanvas({ tree, steps, lcaSource }: TreeCanvasProps) {
  const [currentStep, setCurrentStep] = useState(0);
  const [playing, setPlaying] = useState(steps.length > 0);
  const [speed, setSpeed] = useState(400);
  const [zoom, setZoom] = useState(1);
  const [lcaMode, setLcaMode] = useState(false);
  const [selectedNodeIds, setSelectedNodeIds] = useState<number[]>([]);
  const [lcaNodeId, setLcaNodeId] = useState<number | null>(null);
  const [lcaSteps, setLcaSteps] = useState<LCAStep[]>([]);
  const [lcaStepIndex, setLcaStepIndex] = useState(0);
  const [lcaLoading, setLcaLoading] = useState(false);
  const [lcaError, setLcaError] = useState<string | null>(null);
  const [tooltip, setTooltip] = useState<{
    x: number;
    y: number;
    text: string;
  } | null>(null);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const layout = useMemo(() => {
    return tree ? buildTreeLayout(tree) : emptyTreeLayout;
  }, [tree]);

  useEffect(() => {
    if (!playing) {
      if (intervalRef.current) clearInterval(intervalRef.current);
      return;
    }
    intervalRef.current = setInterval(() => {
      setCurrentStep((prev) => {
        if (prev >= steps.length - 1) {
          setPlaying(false);
          return prev;
        }
        return prev + 1;
      });
    }, speed);
    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, [playing, speed, steps.length]);

  useEffect(() => {
    if (!lcaMode || lcaSteps.length === 0 || lcaStepIndex >= lcaSteps.length - 1) {
      return;
    }

    const timer = setTimeout(() => {
      setLcaStepIndex((index) => index + 1);
    }, speed);

    return () => clearTimeout(timer);
  }, [lcaMode, lcaStepIndex, lcaSteps.length, speed]);

  if (!tree) {
    return <div className={styles.emptyState}>No tree loaded. Run a traversal to begin.</div>;
  }

  const { positions, edges, allNodes, parentById, svgW, svgH } = layout;
  const selectedSummary =
    selectedNodeIds.length === 0
      ? "pick 2 nodes"
      : `${selectedNodeIds.length}/2 selected`;
  const lcaNode = lcaNodeId === null ? null : allNodes.find((node) => node.id === lcaNodeId) ?? null;
  const currentLCAStep = lcaSteps[lcaStepIndex];
  const activeLCANodeIds = currentLCAStep?.activeNodeIds ?? [];
  const currentLCANodeId = currentLCAStep?.lcaNodeId ?? null;

  async function handleNodeClick(node: TreeNode) {
    if (!tree || !lcaMode) return;

    const alreadySelected = selectedNodeIds.includes(node.id);
    const nextSelected = alreadySelected
      ? selectedNodeIds.filter((id) => id !== node.id)
      : [...selectedNodeIds, node.id].slice(0, 2);

    setSelectedNodeIds(nextSelected);
    setLcaSteps([]);
    setLcaStepIndex(0);
    setLcaNodeId(null);
    setLcaError(null);

    if (nextSelected.length === 2) {
      if (!lcaSource) {
        setLcaError("Run traversal first before using backend LCA.");
        return;
      }

      setLcaLoading(true);
      try {
        const res = await lca({
          ...lcaSource,
          nodeA: nextSelected[0],
          nodeB: nextSelected[1],
        });

        if (!res.success) {
          setLcaError(res.error ?? "LCA request failed");
          return;
        }

        setLcaNodeId(res.nodeId);
        setLcaSteps(res.steps ?? []);
      } catch (e: unknown) {
        setLcaError(e instanceof Error ? e.message : "LCA request failed");
      } finally {
        setLcaLoading(false);
      }
    }
  }

  return (
    <div className={styles.shell}>
      <div className={styles.toolbar}>
        <button
          onClick={() => setPlaying((p) => !p)}
          disabled={steps.length === 0}
          className={cx(styles.toolBtn, styles.playBtn, playing && styles.playing)}
        >
          {playing ? "Pause" : "Play"}
        </button>

        <button
          onClick={() => setCurrentStep(0)}
          className={styles.toolBtn}
        >
          Reset
        </button>

        <button
          onClick={() => {
            setPlaying(false);
            setCurrentStep(0);
          }}
          disabled={steps.length === 0}
          className={cx(styles.toolBtn, styles.dangerBtn)}
        >
          Stop
        </button>

        <input
          type="range"
          min={0}
          max={Math.max(0, steps.length - 1)}
          value={currentStep}
          onChange={(e) => setCurrentStep(Number(e.target.value))}
          className={styles.stepSlider}
        />

        <span>
          Step {steps.length > 0 ? currentStep + 1 : 0} / {steps.length}
        </span>

        <span className={styles.speedLabel}>Speed</span>
        <input
          type="range"
          min={50}
          max={1000}
          step={50}
          value={1050 - speed}
          onChange={(e) => setSpeed(1050 - Number(e.target.value))}
          className={styles.speedSlider}
        />

        <button
          onClick={() => {
            setLcaMode((mode) => !mode);
            setSelectedNodeIds([]);
            setLcaNodeId(null);
            setLcaSteps([]);
            setLcaStepIndex(0);
            setLcaError(null);
          }}
          disabled={!tree}
          className={cx(styles.lcaBtn, lcaMode && styles.lcaActive)}
        >
          {lcaLoading ? "LCA: loading" : lcaMode ? `LCA: ${selectedSummary}` : "LCA mode"}
        </button>

        <button
          onClick={() => setZoom((value) => Math.max(0.4, Number((value - 0.1).toFixed(2))))}
          className={styles.zoomBtn}
        >
          -
        </button>
        <span className={styles.zoomValue}>{Math.round(zoom * 100)}%</span>
        <button
          onClick={() => setZoom((value) => Math.min(2, Number((value + 0.1).toFixed(2))))}
          className={styles.zoomBtn}
        >
          +
        </button>
      </div>

      <div className={styles.canvas}>
        <svg
          width={svgW * zoom}
          height={svgH * zoom}
          className={styles.svg}
          onMouseLeave={() => setTooltip(null)}
        >
          <g transform={`scale(${zoom})`}>
      
            {edges.map((e, i) => (
              <line
                key={i}
                x1={e.x1}
                y1={e.y1}
                x2={e.x2}
                y2={e.y2}
                stroke="rgba(255,255,255,0.08)"
                strokeWidth={1.5}
              />
            ))}


            {allNodes.map((node) => {
              const pos = positions.get(node.id ?? 0);
              if (!pos) return null;
              const { color, border } = getNodeColor(node, currentStep, steps, parentById);
              const label = nodeLabel(node);
              const shortLabel = label.length > 13 ? `${label.slice(0, 12)}...` : label;
              const isSelected = selectedNodeIds.includes(node.id) || activeLCANodeIds.includes(node.id);
              const isLCA = currentLCANodeId === node.id || (currentLCANodeId === null && lcaNodeId === node.id && lcaStepIndex >= lcaSteps.length - 1);

              return (
                <g
                  key={node.id}
                  className={styles.nodeGroup}
                  onMouseEnter={() => {
                    setTooltip({
                      x: (pos.x + NODE_W / 2) * zoom,
                      y: (pos.y - 8) * zoom,
                      text: label,
                    });
                  }}
                  onMouseLeave={() => setTooltip(null)}
                  onClick={() => handleNodeClick(node)}
                >
                  {(isSelected || isLCA) && (
                    <rect
                      x={pos.x - 5}
                      y={pos.y - 5}
                      width={NODE_W + 10}
                      height={NODE_H + 10}
                      rx={9}
                      fill="none"
                      stroke={isLCA ? "#d4537e" : "#ed93b1"}
                      strokeWidth={1.5}
                      strokeOpacity={isLCA ? 0.9 : 0.55}
                    />
                  )}
                  <rect
                    x={pos.x}
                    y={pos.y}
                    width={NODE_W}
                    height={NODE_H}
                    rx={6}
                    fill={isLCA ? "rgba(212,83,126,0.25)" : isSelected ? "rgba(212,83,126,0.14)" : color}
                    stroke={isLCA ? "#d4537e" : isSelected ? "#ed93b1" : border}
                    strokeWidth={1.5}
                  />
                  <text
                    x={pos.x + NODE_W / 2}
                    y={pos.y + NODE_H / 2 + 4}
                    textAnchor="middle"
                    fontSize={10}
                    fontFamily="var(--font-mono), monospace"
                    fill={isLCA ? "#ed93b1" : "#e8e8f0"}
                  >
                    {shortLabel}
                  </text>
                </g>
              );
            })}
          </g>
        </svg>

        {tooltip && (
          <div
            className={styles.tooltip}
            style={{
              left: tooltip.x,
              top: tooltip.y,
            }}
          >
            {tooltip.text}
            {lcaNode && tooltip.text === nodeLabel(lcaNode) ? " · backend LCA" : ""}
          </div>
        )}
        {lcaMode && (
          <div className={styles.lcaStatus}>
            {lcaError
              ? `LCA error: ${lcaError}`
              : lcaLoading
                ? "LCA backend: calculating..."
                : lcaNodeId === null
                  ? `LCA mode: ${selectedSummary}`
                  : `LCA step ${Math.min(lcaStepIndex + 1, lcaSteps.length)} / ${lcaSteps.length}: ${
                      lcaNode ? nodeLabel(lcaNode) : `node ${lcaNodeId}`
                    }`}
          </div>
        )}
      </div>
    </div>
  );
}

function cx(...classes: Array<string | false | null | undefined>) {
  return classes.filter(Boolean).join(" ");
}
