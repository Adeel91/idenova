import { auth } from "@/auth";
import { NextResponse } from "next/server";

export const proxy = auth((request) => {
  if (request.auth?.user) {
    return NextResponse.next();
  }

  if (request.nextUrl.pathname.startsWith("/api/")) {
    return NextResponse.json(
      { error: "Unauthorized" },
      { status: 401 },
    );
  }

  return NextResponse.redirect(
    new URL("/api/auth/signin", request.url),
  );
});

export const config = {
  matcher: [
    "/",
    "/verifications/:path*",
    "/api/verifications/:path*",
    "/api/documents/:path*",
  ],
};
