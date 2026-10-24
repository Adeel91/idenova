
export type VerificationStatus = "PENDING" | "APPROVED" | "REJECTED";

export type Verification = {
  id: string;
  customerId: string;
  status: VerificationStatus;
};

type BackendVerification = {
  ID: string;
  CustomerID: string;
  Status: VerificationStatus;
};

const API_URL =
  process.env.VERIFICATION_API_URL ?? "http://127.0.0.1:8080";

export async function getVerifications(): Promise<Verification[]> {
  const response = await fetch(
    `${API_URL}/verifications?limit=100&offset=0`,
    { cache: "no-store" },
  );

  if (!response.ok) {
    throw new Error(`Failed to fetch verifications: ${response.status}`);
  }

  const records: BackendVerification[] = await response.json();

  return records.map((record) => ({
    id: record.ID,
    customerId: record.CustomerID,
    status: record.Status,
  }));
}

export async function getVerification(
  id: string,
): Promise<Verification | null> {
  const response = await fetch(
    `${API_URL}/verifications/${encodeURIComponent(id)}`,
    { cache: "no-store" },
  );

  if (response.status === 404) {
    return null;
  }

  if (!response.ok) {
    throw new Error(`Failed to fetch verification: ${response.status}`);
  }

  const record: BackendVerification = await response.json();

  return {
    id: record.ID,
    customerId: record.CustomerID,
    status: record.Status,
  };
}
