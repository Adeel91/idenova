import { NextResponse } from "next/server";

const DOCUMENT_API_URL =
  process.env.DOCUMENT_API_URL ?? "http://127.0.0.1:8083";

type Context = {
  params: Promise<{ id: string }>;
};

export async function PATCH(request: Request, { params }: Context) {
  const { id } = await params;

  if (!id.trim()) {
    return NextResponse.json(
      { error: "Document ID is required" },
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
    !("status" in body) ||
    (body.status !== "APPROVED" && body.status !== "REJECTED")
  ) {
    return NextResponse.json(
      { error: "Invalid document status" },
      { status: 400 },
    );
  }

  try {
    const response = await fetch(
      `${DOCUMENT_API_URL}/documents/${encodeURIComponent(id)}/status`,
      {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ status: body.status }),
        cache: "no-store",
        signal: AbortSignal.timeout(10000),
      },
    );

    if (response.ok) {
      return new Response(null, { status: 204 });
    }

    let message = "Failed to update document status";

    try {
      const result: unknown = await response.json();

      if (
        typeof result === "object" &&
        result !== null &&
        "error" in result &&
        typeof result.error === "string"
      ) {
        message = result.error;
      }
    } catch {
      // Keep the default error message.
    }

    return NextResponse.json(
      { error: message },
      { status: response.status },
    );
  } catch {
    return NextResponse.json(
      { error: "Document Service unavailable" },
      { status: 502 },
    );
  }
}
