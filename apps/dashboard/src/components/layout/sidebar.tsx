
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
import { LogoutButton } from "@/components/auth/logout-button";

const navigation = [
  { label: "Dashboard", href: "/", icon: LayoutDashboard },
  { label: "Verifications", href: "/verifications", icon: ClipboardCheck },
  { label: "Customers", href: "/customers", icon: Users },
  { label: "Settings", href: "/settings", icon: Settings },
];

type SidebarProps = {
  user: {
    name: string;
    email: string;
  };
};

export function Sidebar({ user }: SidebarProps) {
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

      <div className="mt-auto space-y-3 border-t border-slate-200 pt-5">
        <div className="px-3">
          <p className="truncate text-sm font-semibold text-slate-900">
            {user.name}
          </p>
          <p className="truncate text-xs text-slate-500">
            {user.email}
          </p>
        </div>

        <LogoutButton />

        <p className="px-3 text-xs text-slate-400">
          Idenova · Development
        </p>
      </div>
    </aside>
  );
}
