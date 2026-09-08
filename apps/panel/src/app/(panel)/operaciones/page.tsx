import { Suspense } from "react";
import { OperacionesPageContent } from "./components/operaciones-page-content";

export default function OperacionesPage() {
  return (
    <Suspense>
      <OperacionesPageContent />
    </Suspense>
  );
}
