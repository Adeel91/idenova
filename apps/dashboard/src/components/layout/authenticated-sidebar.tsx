import { auth } from "@/auth";
import { Sidebar } from "./sidebar";

export async function AuthenticatedSidebar() {
  const session = await auth();

  return (
    <Sidebar
      user={{
        name: session?.user?.name ?? "Idenova User",
        email: session?.user?.email ?? "",
      }}
    />
  );
}
