import Link from "next/link";
import { notFound } from "next/navigation";
import {
  ArrowLeft,
  FileText,
  Fingerprint,
  Plus,
  ShieldCheck,
  UserRound,
} from "lucide-react";

import { getVerification } from "@/lib/verifications";
import { getVerificationDocuments } from "@/lib/documents";
import { statusStyles } from "@/lib/mock-verifications";
import { DocumentActions } from "@/components/verifications/document-actions";
import { VerificationActions } from "@/components/verifications/verification-actions";

export const instant = false;

type Props = {
  params: Promise<{ id: string }>;
};

export default async function VerificationDetailsPage({ params }: Props) {
  const { id } = await params;

  const verification = await getVerification(id);

  if (!verification) {
    notFound();
  }

  const documents = await getVerificationDocuments(verification.id);
  const canApprove =
  documents.length > 0 &&
  documents.every((document) => document.status === "APPROVED");

const canReject = documents.some((document) => document.status === "REJECTED");

  return (
    <div className="mx-auto max-w-5xl space-y-8">
      <header className="space-y-4">
        <Link
          href="/verifications"
          className="inline-flex items-center gap-2 text-sm font-medium text-slate-500 hover:text-indigo-600"
        >
          <ArrowLeft size={16} />
          Back to verifications
        </Link>

        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p className="mb-1 text-sm font-medium text-indigo-600">
              Verification details
            </p>
            <h1 className="break-all text-3xl font-bold tracking-tight text-slate-900">
              {verification.id}
            </h1>
          </div>

          <span
            className={`inline-flex rounded-full px-4 py-2 text-sm font-semibold ${statusStyles[verification.status]}`}
          >
            {verification.status}
          </span>
        </div>
      </header>

      <section className="overflow-hidden rounded-xl border border-slate-200 bg-white">
        <div className="border-b border-slate-200 px-6 py-5">
          <h2 className="text-lg font-semibold text-slate-900">
            Verification information
          </h2>
        </div>

        <div className="grid gap-8 p-6 sm:grid-cols-2">
          <div className="space-y-2">
            <div className="flex items-center gap-2 text-sm text-slate-500">
              <Fingerprint size={18} />
              Verification ID
            </div>
            <p className="break-all font-medium text-slate-900">
              {verification.id}
            </p>
          </div>

          <div className="space-y-2">
            <div className="flex items-center gap-2 text-sm text-slate-500">
              <UserRound size={18} />
              Customer ID
            </div>
            <p className="break-all font-medium text-slate-900">
              {verification.customerId}
            </p>
          </div>

          <div className="space-y-2">
            <div className="flex items-center gap-2 text-sm text-slate-500">
              <ShieldCheck size={18} />
              Verification status
            </div>
            <p className="font-medium text-slate-900">
              {verification.status}
            </p>
          </div>
        </div>
      </section>

      {verification.status === "PENDING" && (
        <section className="space-y-4 rounded-xl border border-slate-200 bg-white p-6">
          <h2 className="text-lg font-semibold text-slate-900">
            Verification review
          </h2>

          <p className="text-sm text-slate-500">
            Approve when all documents are approved, or reject when at least one
            document is rejected.
          </p>

          <VerificationActions
            verificationId={verification.id}
            canApprove={canApprove}
            canReject={canReject}
          />
        </section>
      )}

      <section className="overflow-hidden rounded-xl border border-slate-200 bg-white">
        <div className="border-b border-slate-200 px-6 py-5">
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <FileText size={20} className="text-indigo-600" />
              <h2 className="text-lg font-semibold text-slate-900">
                Documents
              </h2>
            </div>

            <Link
              href={`/verifications/${encodeURIComponent(verification.id)}/documents/new`}
              className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white hover:bg-indigo-700"
            >
              <Plus size={16} />
              Add Document
            </Link>
          </div>
          <p className="mt-1 text-sm text-slate-500">
            {documents.length} associated documents
          </p>
        </div>

        {documents.length === 0 ? (
          <div className="px-6 py-12 text-center">
            <FileText size={32} className="mx-auto text-slate-300" />
            <p className="mt-3 text-sm font-medium text-slate-700">
              No documents added yet
            </p>
            <p className="mt-1 text-sm text-slate-500">
              Documents associated with this verification will appear here.
            </p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm">
              <thead className="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
                <tr>
                  <th className="px-6 py-4 font-semibold">
                    Document ID
                  </th>
                  <th className="px-6 py-4 font-semibold">
                    Type
                  </th>
                  <th className="px-6 py-4 font-semibold">
                    Status
                  </th>
                  <th className="px-6 py-4 font-semibold">
                    Actions
                  </th>
                </tr>
              </thead>

              <tbody className="divide-y divide-slate-100">
                {documents.map((document) => (
                  <tr key={document.id} className="hover:bg-slate-50">
                    <td className="px-6 py-4 font-medium text-slate-900">
                      {document.id}
                    </td>
                    <td className="px-6 py-4 text-slate-600">
                      {document.type}
                    </td>
                    <td className="px-6 py-4">
                      <span
                        className={`inline-flex rounded-full px-3 py-1 text-xs font-semibold ${statusStyles[document.status]}`}
                      >
                        {document.status}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      {document.status === "PENDING" ? (
                        <DocumentActions documentId={document.id} />
                      ) : (
                        <span className="text-xs text-slate-400">
                          Reviewed
                        </span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </div>
  );
}
