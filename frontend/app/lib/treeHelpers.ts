import type { TraversalStep, TreeNode } from "@/app/lib/types";

export const NODE_W = 90;
export const NODE_H = 34;

const X_GAP = 110;
const Y_GAP = 94;

export interface NodePosition {
  x: number;
  y: number;
}

export interface EdgePosition {
  x1: number;
  y1: number;
  x2: number;
  y2: number;
}

export interface TreeLayout {
  positions: Map<number, NodePosition>;
  edges: EdgePosition[];
  allNodes: TreeNode[];
  parentById: Map<number, number>;
  svgW: number;
  svgH: number;
}

export const emptyTreeLayout: TreeLayout = {
  positions: new Map<number, NodePosition>(),
  edges: [],
  allNodes: [],
  parentById: new Map<number, number>(),
  svgW: 900,
  svgH: 420,
};

export function buildTreeLayout(root: TreeNode): TreeLayout {
  const positions = new Map<number, NodePosition>();
  const levels = collectLevels(root);
  const maxLevelWidth = Math.max(1, ...levels.map((level) => level.length));
  const svgW = Math.max(900, maxLevelWidth * X_GAP);

  levels.forEach((level, depth) => {
    const rowWidth = (level.length - 1) * X_GAP;
    const startX = svgW / 2 - rowWidth / 2;

    level.forEach((node, index) => {
      positions.set(node.id, {
        x: startX + index * X_GAP - NODE_W / 2,
        y: 40 + depth * Y_GAP,
      });
    });
  });

  const edges: EdgePosition[] = [];
  const allNodes: TreeNode[] = [];
  const parentById = new Map<number, number>();

  collectEdges(root, positions, edges);
  collectNodes(root, allNodes, parentById);

  return {
    positions,
    edges,
    allNodes,
    parentById,
    svgW,
    svgH: Math.max(420, levels.length * Y_GAP + 80),
  };
}

export function getNodeColor(
  node: TreeNode,
  currentStep: number,
  steps: TraversalStep[],
  parentById: Map<number, number>
) {
  const stepsUpTo = steps.slice(0, currentStep + 1);
  const current = steps[currentStep];
  const parentId = parentById.get(node.id);

  let color = "var(--bg3)";
  let border = "rgba(255,255,255,0.1)";

  for (const step of stepsUpTo) {
    const isAffectedText = node.tag === "#text" && step.nodeId === parentId && step.isMatch;
    if (step.nodeId !== node.id && !isAffectedText) continue;

    if (step.isMatch) {
      color = "rgba(29,158,117,0.25)";
      border = "var(--accent-dfs)";
    } else {
      color = "rgba(55,138,221,0.2)";
      border = "var(--accent-bfs)";
    }
  }

  if (current?.nodeId === node.id) {
    color = "var(--current)";
    border = "var(--current)";
  }

  return { color, border };
}

export function nodeLabel(node: TreeNode) {
  if (node.tag === "#text") {
    const text = node.content ?? "";
    return text.length > 24 ? `#text "${text.slice(0, 21)}..."` : `#text "${text}"`;
  }

  const htmlId = node.attrs?.id ? `#${node.attrs.id}` : "";
  const className = node.attrs?.class ? `.${node.attrs.class.split(" ")[0]}` : "";
  return `<${node.tag}${htmlId}${className}>`;
}

function collectLevels(root: TreeNode) {
  const levels: TreeNode[][] = [];
  const queue: Array<{ node: TreeNode; depth: number }> = [{ node: root, depth: 0 }];

  for (let head = 0; head < queue.length; head++) {
    const { node, depth } = queue[head];
    if (!levels[depth]) levels[depth] = [];
    levels[depth].push(node);

    for (const child of node.children ?? []) {
      queue.push({ node: child, depth: depth + 1 });
    }
  }

  return levels;
}

function collectEdges(
  node: TreeNode,
  positions: Map<number, NodePosition>,
  edges: EdgePosition[]
) {
  const parentPos = positions.get(node.id);
  if (!parentPos) return;

  for (const child of node.children ?? []) {
    const childPos = positions.get(child.id);
    if (childPos) {
      edges.push({
        x1: parentPos.x + NODE_W / 2,
        y1: parentPos.y + NODE_H,
        x2: childPos.x + NODE_W / 2,
        y2: childPos.y,
      });
    }
    collectEdges(child, positions, edges);
  }
}

function collectNodes(
  node: TreeNode,
  allNodes: TreeNode[],
  parentById: Map<number, number>,
  parentId?: number
) {
  allNodes.push(node);
  if (parentId !== undefined) parentById.set(node.id, parentId);

  for (const child of node.children ?? []) {
    collectNodes(child, allNodes, parentById, node.id);
  }
}
