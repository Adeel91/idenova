"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Check, Loader2, X } from "lucide-react";

type ReviewStatus = "APPROVED" | "REJECTED";

type Props = {
  verificationId: string;
  canApprove: boolean;
  canReject: boolean;
};

export function VerificationActions({
  verificationId,
  canApprove,
  canReject,
}: Props) {
  const router = useRouter();
  const [loading, setLoading] = useState<ReviewStatus | null>(null);
  const [error, setError] = useState("");

  async function updateStatus(status: ReviewStatus) {
    if (loading) return;

    const action = status === "APPROVED" ? "approve" : "reject";

    if (!window.confirm(`Are you sure you want to ${action} this verification?`)) {
      return;
    }

    setLoading(status);
    setError("");

    try {
      const response = await fetch(
        `/api/verifications/${encodeURIComponent(verificationId)}/status`,
        {
          method: "PATCH",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ status }),
        },
      );

      if (!response.ok) {
        let message = "Failed to update verification";

        try {
          const result: { error?: string } = await response.json();
          if (typeof result.error === "string") {
            message = result.error;
          }
        } catch {
          // Keep the default message.
        }

        throw new Error(message);
      }

      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Something went wrong");
    } finally {
      setLoading(null);
    }
  }

  return (
    <div className="space-y-2">
      <div className="flex flex-wrap gap-3">
        <button
          type="button"
          disabled={!canApprove || loading !== null}
          onClick={() => updateStatus("APPROVED")}
          className="inline-flex items-center gap-2 rounded-lg bg-emerald-600 px-4 py-2 text-sm font-semibold text-white hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-40"
        >
          {loading === "APPROVED" ? (
            <Loader2 size={16} className="animate-spin" />
          ) : (
            <Check size={16} />
          )}
          Approve Verification
        </button>

        <button
          type="button"
          disabled={!canReject || loading !== null}
          onClick={() => updateStatus("REJECTED")}
          className="inline-flex items-center gap-2 rounded-lg bg-rose-600 px-4 py-2 text-sm font-semibold text-white hover:bg-rose-700 disabled:cursor-not-allowed disabled:opacity-40"
        >
          {loading === "REJECTED" ? (
            <Loader2 size={16} className="animate-spin" />
          ) : (
            <X size={16} />
          )}
          Reject Verification
        </button>
      </div>

      {error && (
        <p role="alert" className="text-sm text-rose-600">
          {error}
        </p>
      )}
    </div>
  );
}
