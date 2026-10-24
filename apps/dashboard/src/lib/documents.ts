export type DocumentStatus = "PENDING" | "APPROVED" | "REJECTED";

export type VerificationDocument = {
  id: string;
  verificationId: string;
  type: string;
  status: DocumentStatus;
};

type BackendDocument = {
  ID: string;
  VerificationID: string;
  Type: string;
  Status: DocumentStatus;
};

const DOCUMENT_API_URL =
  process.env.DOCUMENT_API_URL ?? "http://127.0.0.1:8083";

export async function getVerificationDocuments(
  verificationId: string,
): Promise<VerificationDocument[]> {
  const response = await fetch(
    `${DOCUMENT_API_URL}/verifications/${encodeURIComponent(verificationId)}/documents`,
    {
      cache: "no-store",
      signal: AbortSignal.timeout(10000),
    },
  );

  if (!response.ok) {
    throw new Error(`Failed to fetch documents: ${response.status}`);
  }

  const records: BackendDocument[] = await response.json();

  if (!Array.isArray(records)) {
    throw new Error("Invalid response from Document Service");
  }

  return records.map((record) => ({
    id: record.ID,
    verificationId: record.VerificationID,
    type: record.Type,
    status: record.Status,
  }));
}
