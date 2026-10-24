
"use client";

import {
  LayoutDashboard,
  ClipboardCheck,
  Users,
  Settings,
  ShieldCheck,
} from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";

const navigation = [
  { label: "Dashboard", href: "/", icon: LayoutDashboard },
  { label: "Verifications", href: "/verifications", icon: ClipboardCheck },
  { label: "Customers", href: "/customers", icon: Users },
  { label: "Settings", href: "/settings", icon: Settings },
];

export function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="flex w-64 shrink-0 flex-col border-r border-slate-200 bg-white p-6">
      <div className="mb-10 flex items-center gap-3">
        <div className="rounded-xl bg-indigo-600 p-2 text-white">
          <ShieldCheck size={24} />
        </div>
        <span className="text-xl font-bold text-slate-900">
          Idenova
        </span>
      </div>

      <nav className="space-y-2" aria-label="Main navigation">
        {navigation.map((item) => {
          const Icon = item.icon;
          const active =
            item.href === "/"
              ? pathname === "/"
              : pathname.startsWith(item.href);

          return (
            <Link
              key={item.href}
              href={item.href}
              aria-current={active ? "page" : undefined}
              className={`flex items-center gap-3 rounded-lg px-4 py-3 text-sm font-medium transition-colors ${
                active
                  ? "bg-indigo-50 text-indigo-700"
                  : "text-slate-600 hover:bg-slate-100 hover:text-slate-900"
              }`}
            >
              <Icon size={19} />
              {item.label}
            </Link>
          );
        })}
      </nav>

      <div className="mt-auto border-t border-slate-200 pt-5 text-xs text-slate-400">
        Idenova · Development
      </div>
    </aside>
  );
}
