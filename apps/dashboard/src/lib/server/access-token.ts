import "server-only";

import { getToken } from "next-auth/jwt";

export async function getAccessToken(
  request: Request,
): Promise<string | null> {
  const session = await getToken({
    req: request,
    secret: process.env.AUTH_SECRET!,
    secureCookie: process.env.NODE_ENV === "production",
  });

  if (
    typeof session?.accessToken !== "string" ||
    typeof session.accessTokenExpiresAt !== "number"
  ) {
    return null;
  }

  if (Date.now() >= session.accessTokenExpiresAt * 1000) {
    return null;
  }

  return session.accessToken;
}
