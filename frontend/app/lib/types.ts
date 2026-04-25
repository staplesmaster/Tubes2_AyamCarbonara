export type InputMode = "url" | "html";
export type Algorithm = "bfs" | "dfs";

export interface TreeNode {
  id: number;
  tag: string;
  attrs?: Record<string, string>;
  content?: string;
  children?: TreeNode[];
}

export interface TraversalStep {
  step: number;
  nodeId: number;
  tag: string;
  isMatch: boolean;
  visitedCount: number;
  matchedCount: number;
}

export interface TraversalStats {
  visited: number;
  matched: number;
  maxDepth: number;
  elapsedMs: number;
}

export interface TraverseRequest {
  inputMode: InputMode;
  url: string;
  html: string;
  selector: string;
  algorithm: Algorithm;
  parallel: boolean;
  limit: number;
  allResult: boolean;
}

export interface LCARequest {
  inputMode: InputMode;
  url: string;
  html: string;
  nodeA: number;
  nodeB: number;
}

export interface TraverseResponse {
  success: boolean;
  tree: TreeNode;
  steps: TraversalStep[];
  stats: TraversalStats;
  error?: string;
}

export interface LCAStep {
  step: number;
  nodeA: number;
  nodeB: number;
  activeNodeIds: number[];
  lcaNodeId?: number;
}

export interface LCAResponse {
  success: boolean;
  nodeId: number;
  tag: string;
  label: string;
  steps: LCAStep[];
  error?: string;
}

export interface UploadResponse {
  success: boolean;
  filename?: string;
  html?: string;
  error?: string;
}
