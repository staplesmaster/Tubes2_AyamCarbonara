import { NextRequest, NextResponse } from "next/server";

const BACKEND = process.env.BACKEND_URL || "http://localhost:8080";

export async function POST(req: NextRequest) {
  const body = await req.formData();

  const res = await fetch(`${BACKEND}/api/upload`, {
    method: "POST",
    body,
  });

  const data = await res.json();
  return NextResponse.json(data, { status: res.status });
}
