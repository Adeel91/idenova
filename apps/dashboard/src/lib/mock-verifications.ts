
export type VerificationStatus = "PENDING" | "APPROVED" | "REJECTED";

export type Verification = {
  id: string;
  customerId: string;
  status: VerificationStatus;
  createdAt: string;
};

export const verifications: Verification[] = [
  {
    id: "ver_1005",
    customerId: "customer_105",
    status: "PENDING",
    createdAt: "2026-10-24",
  },
  {
    id: "ver_1004",
    customerId: "customer_104",
    status: "APPROVED",
    createdAt: "2026-10-23",
  },
  {
    id: "ver_1003",
    customerId: "customer_103",
    status: "REJECTED",
    createdAt: "2026-10-22",
  },
  {
    id: "ver_1002",
    customerId: "customer_102",
    status: "APPROVED",
    createdAt: "2026-10-21",
  },
  {
    id: "ver_1001",
    customerId: "customer_101",
    status: "PENDING",
    createdAt: "2026-10-20",
  },
];

export const statusStyles: Record<VerificationStatus, string> = {
  PENDING: "bg-amber-50 text-amber-700",
  APPROVED: "bg-emerald-50 text-emerald-700",
  REJECTED: "bg-rose-50 text-rose-700",
};
