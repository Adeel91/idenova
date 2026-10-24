import { Suspense } from "react";
import AddDocumentForm from "./add-document-form";

export const instant = false;

export default function NewDocumentPage() {
  return (
    <Suspense
      fallback={
        <div className="mx-auto max-w-3xl p-6 text-sm text-slate-500">
          Loading document form...
        </div>
      }
    >
      <AddDocumentForm />
    </Suspense>
  );
}
