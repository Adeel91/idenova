import Link from "next/link";
import { Plus } from "lucide-react";

import { getVerifications } from "@/lib/verifications";
import { statusStyles } from "@/lib/mock-verifications";

export const instant = false;

export default async function VerificationsPage() {
  const verifications = await getVerifications();

  return (
    <div className="mx-auto max-w-7xl space-y-8">
      <header className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-slate-900">
            Verifications
          </h1>

          <p className="mt-2 text-sm text-slate-500">
            View and manage identity verification requests.
          </p>
        </div>

        <Link
          href="/verifications/new"
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-indigo-700"
        >
          <Plus size={18} />
          Create Verification
        </Link>
      </header>

      <section className="overflow-hidden rounded-xl border border-slate-200 bg-white">
        <div className="border-b border-slate-200 px-6 py-5">
          <h2 className="text-lg font-semibold text-slate-900">
            All verifications
          </h2>

          <p className="mt-1 text-sm text-slate-500">
            Showing {verifications.length} verification requests
          </p>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
              <tr>
                <th className="px-6 py-4 font-semibold">
                  Verification ID
                </th>
                <th className="px-6 py-4 font-semibold">
                  Customer ID
                </th>
                <th className="px-6 py-4 font-semibold">
                  Status
                </th>
              </tr>
            </thead>

            <tbody className="divide-y divide-slate-100">
              {verifications.map((verification) => (
                <tr
                  key={verification.id}
                  className="hover:bg-slate-50"
                >
                  <td className="px-6 py-4 font-medium text-slate-900">
                    {verification.id}
                  </td>

                  <td className="px-6 py-4 text-slate-600">
                    {verification.customerId}
                  </td>

                  <td className="px-6 py-4">
                    <span
                      className={`inline-flex rounded-full px-3 py-1 text-xs font-semibold ${statusStyles[verification.status]}`}
                    >
                      {verification.status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          {verifications.length === 0 && (
            <div className="px-6 py-12 text-center">
              <p className="text-sm text-slate-500">
                No verification requests found.
              </p>

              <Link
                href="/verifications/new"
                className="mt-3 inline-block text-sm font-semibold text-indigo-600 hover:text-indigo-700"
              >
                Create your first verification
              </Link>
            </div>
          )}
        </div>
      </section>
    </div>
  );
}
