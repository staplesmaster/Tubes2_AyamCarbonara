import type { LCARequest, LCAResponse, TraverseRequest, TraverseResponse, UploadResponse } from "@/app/lib/types";

export const traverse = (body: TraverseRequest): Promise<TraverseResponse> =>
  fetch("/api/traverse", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  }).then((res) => res.json());

export const lca = (body: LCARequest): Promise<LCAResponse> =>
  fetch("/api/lca", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  }).then((res) => res.json());

export const uploadHTML = (file: File): Promise<UploadResponse> => {
  const body = new FormData();
  body.append("file", file);

  return fetch("/api/upload", {
    method: "POST",
    body,
  }).then((res) => res.json());
};
