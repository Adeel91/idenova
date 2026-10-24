"use client";

import { useRef, useState, type FormEvent } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Loader2, Plus } from "lucide-react";
import { createIdempotencyKey } from "@/lib/idempotency";

const DOCUMENT_TYPES = [
  { value: "SYNTHETIC_ID", label: "Synthetic identity document" },
  { value: "SYNTHETIC_PASSPORT", label: "Synthetic passport" },
  {
    value: "SYNTHETIC_DRIVING_LICENSE",
    label: "Synthetic driving licence",
  },
];

export default function AddDocumentPage() {
  const router = useRouter();
  const params = useParams<{ id: string }>();

  const verificationId = params.id;
  const detailsUrl = `/verifications/${encodeURIComponent(verificationId)}`;

  const [type, setType] = useState("SYNTHETIC_ID");
  const idempotencyKey = useRef<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (loading) return;

    if (!idempotencyKey.current) {
      idempotencyKey.current = createIdempotencyKey();
    }

    setLoading(true);
    setError("");

    try {
      const response = await fetch(
        `/api/verifications/${encodeURIComponent(verificationId)}/documents`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "Idempotency-Key": idempotencyKey.current,
          },
          body: JSON.stringify({ type }),
        },
      );

      if (!response.ok) {
        let message = "Failed to create document";

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

      router.push(detailsUrl);
      router.refresh();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Something went wrong",
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="mx-auto max-w-3xl space-y-8">
      <header>
        <Link
          href={detailsUrl}
          className="inline-flex items-center gap-2 text-sm font-medium text-slate-500 hover:text-indigo-600"
        >
          <ArrowLeft size={16} />
          Back to verification
        </Link>

        <h1 className="mt-5 text-3xl font-bold tracking-tight text-slate-900">
          Add Document
        </h1>

        <p className="mt-2 text-sm text-slate-500">
          Create synthetic document metadata for this verification.
        </p>
      </header>

      <form
        onSubmit={handleSubmit}
        className="space-y-6 rounded-xl border border-slate-200 bg-white p-6"
      >
        <div className="space-y-2">
          <label className="block text-sm font-semibold text-slate-700">
            Verification ID
          </label>

          <div className="break-all rounded-lg border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-600">
            {verificationId}
          </div>
        </div>

        <div className="space-y-2">
          <label
            htmlFor="documentType"
            className="block text-sm font-semibold text-slate-700"
          >
            Document type
          </label>

          <select
            id="documentType"
            value={type}
            onChange={(event) => {
              setType(event.target.value);
              idempotencyKey.current = null;
              setError("");
            }}
            disabled={loading}
            className="w-full rounded-lg border border-slate-300 bg-white px-4 py-3 text-sm text-slate-900 outline-none focus:border-indigo-500"
          >
            {DOCUMENT_TYPES.map((documentType) => (
              <option
                key={documentType.value}
                value={documentType.value}
              >
                {documentType.label}
              </option>
            ))}
          </select>
        </div>

        <p className="text-sm text-slate-500">
          This creates metadata only. No actual document is uploaded.
        </p>

        {error && (
          <p
            role="alert"
            className="rounded-lg bg-rose-50 px-4 py-3 text-sm text-rose-700"
          >
            {error}
          </p>
        )}

        <div className="flex items-center justify-end gap-3">
          <Link
            href={detailsUrl}
            className="rounded-lg border border-slate-200 px-4 py-2.5 text-sm font-semibold text-slate-700 hover:bg-slate-50"
          >
            Cancel
          </Link>

          <button
            type="submit"
            disabled={loading}
            className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-indigo-700 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {loading ? (
              <Loader2 size={18} className="animate-spin" />
            ) : (
              <Plus size={18} />
            )}
            {loading ? "Creating..." : "Add Document"}
          </button>
        </div>
      </form>
    </div>
  );
}
