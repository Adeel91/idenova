import { randomUUID } from "node:crypto";
import { NextResponse } from "next/server";

const DOCUMENT_API_URL =
  process.env.DOCUMENT_API_URL ?? "http://127.0.0.1:8083";

const ALLOWED_TYPES = [
  "SYNTHETIC_ID",
  "SYNTHETIC_PASSPORT",
  "SYNTHETIC_DRIVING_LICENSE",
] as const;

type Context = {
  params: Promise<{ id: string }>;
};

export async function POST(request: Request, { params }: Context) {
  const { id: verificationId } = await params;

  if (!verificationId.trim()) {
    return NextResponse.json(
      { error: "Verification ID is required" },
      { status: 400 },
    );
  }

  let body: unknown;

  try {
    body = await request.json();
  } catch {
    return NextResponse.json(
      { error: "Invalid JSON" },
      { status: 400 },
    );
  }

  if (
    typeof body !== "object" ||
    body === null ||
    !("type" in body) ||
    typeof body.type !== "string" ||
    !ALLOWED_TYPES.some((type) => type === body.type)
  ) {
    return NextResponse.json(
      { error: "Invalid document type" },
      { status: 400 },
    );
  }

  try {
    const response = await fetch(`${DOCUMENT_API_URL}/documents`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        id: `doc_${randomUUID()}`,
        verificationId,
        type: body.type,
      }),
      cache: "no-store",
      signal: AbortSignal.timeout(10000),
    });

    const result: unknown = await response.json();

    return NextResponse.json(result, {
      status: response.status,
    });
  } catch {
    return NextResponse.json(
      { error: "Document Service unavailable" },
      { status: 502 },
    );
  }
}
