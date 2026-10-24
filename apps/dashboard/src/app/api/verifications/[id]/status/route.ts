import { NextResponse } from "next/server";

import { getVerification } from "@/lib/verifications";
import { getVerificationDocuments } from "@/lib/documents";

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

  try {
    const verification = await getVerification(id);

    if (!verification) {
      return NextResponse.json(
        { error: "Verification not found" },
        { status: 404 },
      );
    }

    if (verification.status !== "PENDING") {
      return NextResponse.json(
        { error: "Verification has already been reviewed" },
        { status: 409 },
      );
    }

    const documents = await getVerificationDocuments(id);

    if (body.status === "APPROVED") {
      if (
        documents.length === 0 ||
        documents.some((document) => document.status !== "APPROVED")
      ) {
        return NextResponse.json(
          { error: "All documents must be approved first" },
          { status: 409 },
        );
      }
    }

    if (body.status === "REJECTED") {
      if (!documents.some((document) => document.status === "REJECTED")) {
        return NextResponse.json(
          { error: "At least one document must be rejected first" },
          { status: 409 },
        );
      }
    }

    const response = await fetch(
      `${VERIFICATION_API_URL}/verifications/${encodeURIComponent(id)}/status`,
      {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
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
      // Keep the default message.
    }

    return NextResponse.json(
      { error: message },
      { status: response.status },
    );
  } catch {
    return NextResponse.json(
      { error: "Verification or Document Service unavailable" },
      { status: 502 },
    );
  }
}
