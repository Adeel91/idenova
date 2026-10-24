"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Check, Loader2, X } from "lucide-react";

type DocumentStatus = "APPROVED" | "REJECTED";

type Props = {
  documentId: string;
};

export function DocumentActions({ documentId }: Props) {
  const router = useRouter();

  const [loading, setLoading] = useState<DocumentStatus | null>(null);
  const [error, setError] = useState("");

  async function updateStatus(status: DocumentStatus) {
    if (loading) return;

    const action = status === "APPROVED" ? "approve" : "reject";

    if (!window.confirm(`Are you sure you want to ${action} this document?`)) {
      return;
    }

    setLoading(status);
    setError("");

    try {
      const response = await fetch(
        `/api/documents/${encodeURIComponent(documentId)}/status`,
        {
          method: "PATCH",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ status }),
        },
      );

      if (!response.ok) {
        let message = "Failed to update document";

        try {
          const result: { error?: string } = await response.json();

          if (typeof result.error === "string") {
            message = result.error;
          }
        } catch {
          // Keep the default error message.
        }

        throw new Error(message);
      }

      router.refresh();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Something went wrong",
      );
    } finally {
      setLoading(null);
    }
  }

  return (
    <div className="space-y-2">
      <div className="flex flex-wrap items-center gap-2">
        <button
          type="button"
          disabled={loading !== null}
          onClick={() => updateStatus("APPROVED")}
          className="inline-flex items-center gap-1 rounded-lg bg-emerald-50 px-3 py-2 text-xs font-semibold text-emerald-700 hover:bg-emerald-100 disabled:opacity-50"
        >
          {loading === "APPROVED" ? (
            <Loader2 size={14} className="animate-spin" />
          ) : (
            <Check size={14} />
          )}
          Approve
        </button>

        <button
          type="button"
          disabled={loading !== null}
          onClick={() => updateStatus("REJECTED")}
          className="inline-flex items-center gap-1 rounded-lg bg-rose-50 px-3 py-2 text-xs font-semibold text-rose-700 hover:bg-rose-100 disabled:opacity-50"
        >
          {loading === "REJECTED" ? (
            <Loader2 size={14} className="animate-spin" />
          ) : (
            <X size={14} />
          )}
          Reject
        </button>
      </div>

      {error && (
        <p role="alert" className="text-xs text-rose-600">
          {error}
        </p>
      )}
    </div>
  );
}
