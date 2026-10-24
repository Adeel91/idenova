import { connection } from "next/server";
import Link from "next/link";
import {
  ArrowRight,
  CheckCircle2,
  Clock3,
  FileCheck2,
  Plus,
  XCircle,
} from "lucide-react";

import { getVerifications } from "@/lib/verifications";
import { statusStyles } from "@/lib/mock-verifications";

export const instant = false;

export default async function DashboardPage() {
  await connection();

  const verifications = await getVerifications();

  const pending = verifications.filter(
    (verification) => verification.status === "PENDING",
  ).length;

  const approved = verifications.filter(
    (verification) => verification.status === "APPROVED",
  ).length;

  const rejected = verifications.filter(
    (verification) => verification.status === "REJECTED",
  ).length;

  const stats = [
    {
      label: "Total verifications",
      value: verifications.length,
      icon: FileCheck2,
      color: "text-indigo-600",
    },
    {
      label: "Pending",
      value: pending,
      icon: Clock3,
      color: "text-amber-600",
    },
    {
      label: "Approved",
      value: approved,
      icon: CheckCircle2,
      color: "text-emerald-600",
    },
    {
      label: "Rejected",
      value: rejected,
      icon: XCircle,
      color: "text-rose-600",
    },
  ];

  return (
    <div className="mx-auto max-w-7xl space-y-8">
      <header className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <p className="mb-1 text-sm font-medium text-indigo-600">
            Overview
          </p>

          <h1 className="text-3xl font-bold tracking-tight text-slate-900">
            Dashboard
          </h1>

          <p className="mt-2 text-sm text-slate-500">
            Monitor and manage your identity verifications.
          </p>
        </div>

        <Link
          href="/verifications/new"
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-indigo-700"
        >
          <Plus size={18} />
          Create Verification
        </Link>
      </header>

      <section
        className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4"
        aria-label="Verification statistics"
      >
        {stats.map((stat) => {
          const Icon = stat.icon;

          return (
            <div
              key={stat.label}
              className="rounded-xl border border-slate-200 bg-white p-5"
            >
              <div className="flex items-center justify-between">
                <p className="text-sm font-medium text-slate-500">
                  {stat.label}
                </p>

                <Icon size={20} className={stat.color} />
              </div>

              <p className="mt-4 text-3xl font-bold text-slate-900">
                {stat.value}
              </p>
            </div>
          );
        })}
      </section>

      <section className="overflow-hidden rounded-xl border border-slate-200 bg-white">
        <div className="flex items-center justify-between gap-4 border-b border-slate-200 px-6 py-5">
          <div>
            <h2 className="text-lg font-semibold text-slate-900">
              Recent verifications
            </h2>

            <p className="mt-1 text-sm text-slate-500">
              Latest verification requests and their statuses
            </p>
          </div>

          <Link
            href="/verifications"
            className="inline-flex items-center gap-2 text-sm font-semibold text-indigo-600 hover:text-indigo-700"
          >
            View all
            <ArrowRight size={16} />
          </Link>
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
              {verifications.slice(0, 5).map((verification) => (
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
            <p className="px-6 py-12 text-center text-sm text-slate-500">
              No verifications yet. Create your first verification.
            </p>
          )}
        </div>
      </section>
    </div>
  );
}
