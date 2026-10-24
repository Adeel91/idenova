
import type { Metadata } from "next";
import { Suspense } from "react";
import { Sidebar } from "@/components/layout/sidebar";
import "./globals.css";

export const metadata: Metadata = {
  title: "Idenova Dashboard",
  description: "Identity verification management platform",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="bg-slate-50 text-slate-900 antialiased">
        <div className="flex min-h-screen">
          <Suspense
            fallback={
              <aside className="w-64 shrink-0 border-r border-slate-200 bg-white" />
            }
          >
            <Sidebar />
          </Suspense>

          <main className="min-w-0 flex-1 bg-[#f8fafc] p-6 lg:p-10">
            {children}
          </main>
        </div>
      </body>
    </html>
  );
}
