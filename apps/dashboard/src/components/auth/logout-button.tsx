"use client";

import { signOut } from "next-auth/react";
import { LogOut } from "lucide-react";

export function LogoutButton() {
  return (
    <button
      type="button"
      onClick={() =>
        signOut({ callbackUrl: "/" })
      }
      className="flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium text-slate-600 transition hover:bg-slate-100 hover:text-red-600"
    >
      <LogOut size={18} />
      Sign out
    </button>
  );
}
