
"use client";

import { useState, type FormEvent } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { ArrowLeft, Loader2 } from "lucide-react";

export default function NewVerificationPage() {
  const router = useRouter();

  const [customerId, setCustomerId] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const value = customerId.trim();

    if (!value) {
      setError("Customer ID is required");
      return;
    }

    setLoading(true);
    setError("");

    try {
      const response = await fetch("/api/verifications", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ customerId: value }),
      });

      if (!response.ok) {
        let message = "Failed to create verification";

        try {
          const result = await response.json();

          if (typeof result.error === "string") {
            message = result.error;
          }
        } catch {
          // Keep the default error message.
        }

        throw new Error(message);
      }

      router.push("/verifications");
      router.refresh();
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Something went wrong",
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="mx-auto max-w-2xl space-y-8">
      <div>
        <Link
          href="/verifications"
          className="inline-flex items-center gap-2 text-sm text-slate-500 hover:text-slate-900"
        >
          <ArrowLeft size={16} />
          Back to verifications
        </Link>

        <h1 className="mt-6 text-3xl font-bold text-slate-900">
          Create Verification
        </h1>

        <p className="mt-2 text-sm text-slate-500">
          Start a new identity verification request.
        </p>
      </div>

      <form
        onSubmit={handleSubmit}
        className="rounded-xl border border-slate-200 bg-white p-6"
      >
        <h2 className="text-lg font-semibold text-slate-900">
          Customer details
        </h2>

        <p className="mt-2 text-sm text-slate-500">
          Enter the customer ID to create a verification.
        </p>

        <div className="mt-6">
          <label
            htmlFor="customerId"
            className="mb-2 block text-sm font-medium text-slate-700"
          >
            Customer ID
          </label>

          <input
            id="customerId"
            name="customerId"
            type="text"
            value={customerId}
            onChange={(event) => {
              setCustomerId(event.target.value);
              setError("");
            }}
            placeholder="e.g. customer_123"
            required
            disabled={loading}
            className="w-full rounded-lg border border-slate-300 bg-white px-4 py-3 text-sm text-slate-900 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-100 disabled:opacity-60"
          />
        </div>

        {error && (
          <div
            role="alert"
            className="mt-4 rounded-lg bg-red-50 p-4 text-sm text-red-700"
          >
            {error}
          </div>
        )}

        <div className="mt-8 flex justify-end gap-3">
          <Link
            href="/verifications"
            className="rounded-lg border border-slate-200 px-5 py-2.5 text-sm font-medium text-slate-700 hover:bg-slate-50"
          >
            Cancel
          </Link>

          <button
            type="submit"
            disabled={loading}
            className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white hover:bg-indigo-700 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {loading && <Loader2 size={16} className="animate-spin" />}
            {loading ? "Creating..." : "Create Verification"}
          </button>
        </div>
      </form>
    </div>
  );
}
