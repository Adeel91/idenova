import { NextResponse } from "next/server";
import { getAccessToken } from "@/lib/server/access-token";

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

  const idempotencyKey = request.headers.get("Idempotency-Key");

  if (!idempotencyKey || !idempotencyKey.trim()) {
    return NextResponse.json(
      { error: "Idempotency-Key header is required" },
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

  const accessToken = await getAccessToken(request);

  if (!accessToken) {
    return NextResponse.json(
      { error: "Unauthorized. Please sign in again." },
      { status: 401 },
    );
  }

  try {
    const response = await fetch(`${DOCUMENT_API_URL}/documents`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${accessToken}`,
        "Idempotency-Key": idempotencyKey,
      },
      body: JSON.stringify({
        verificationId,
        type: body.type,
      }),
      cache: "no-store",
      signal: AbortSignal.timeout(10000),
    });

    const result: unknown = await response.json();

    const headers = new Headers();

    if (response.headers.get("Idempotency-Replayed") === "true") {
      headers.set("Idempotency-Replayed", "true");
    }

    return NextResponse.json(result, {
      status: response.status,
      headers,
    });
  } catch {
    return NextResponse.json(
      { error: "Document Service unavailable" },
      { status: 502 },
    );
  }
}
