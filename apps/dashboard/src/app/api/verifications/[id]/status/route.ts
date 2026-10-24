import { NextResponse } from "next/server";
import { getAccessToken } from "@/lib/server/access-token";

const VERIFICATION_API_URL =
  process.env.VERIFICATION_API_URL ?? "http://127.0.0.1:8080";

type Context = {
  params: Promise<{ id: string }>;
};

export async function PATCH(request: Request, { params }: Context) {
  const { id } = await params;

  if (!id.trim()) {
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
    !("status" in body) ||
    (body.status !== "APPROVED" && body.status !== "REJECTED")
  ) {
    return NextResponse.json(
      { error: "Invalid verification status" },
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
    const response = await fetch(
      `${VERIFICATION_API_URL}/verifications/${encodeURIComponent(id)}/status`,
      {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify({ status: body.status }),
        cache: "no-store",
        signal: AbortSignal.timeout(10000),
      },
    );

    if (response.ok) {
      return new Response(null, { status: 204 });
    }

    let message = "Failed to update verification";

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
      { error: "Verification Service unavailable" },
      { status: 502 },
    );
  }
}
