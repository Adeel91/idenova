
import { NextResponse } from "next/server";

const VERIFICATION_API_URL =
  process.env.VERIFICATION_API_URL ?? "http://127.0.0.1:8080";

export async function POST(request: Request) {
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
    !("customerId" in body) ||
    typeof body.customerId !== "string" ||
    !body.customerId.trim()
  ) {
    return NextResponse.json(
      { error: "Customer ID is required" },
      { status: 400 },
    );
  }

  const customerId = body.customerId.trim();
  const id = `ver_${crypto.randomUUID()}`;

  try {
    const response = await fetch(
      `${VERIFICATION_API_URL}/verifications`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ id, customerId }),
        cache: "no-store",
        signal: AbortSignal.timeout(10000),
      },
    );

    const result = await response.json();

    return NextResponse.json(result, {
      status: response.status,
    });
  } catch {
    return NextResponse.json(
      { error: "Verification Service unavailable" },
      { status: 502 },
    );
  }
}
