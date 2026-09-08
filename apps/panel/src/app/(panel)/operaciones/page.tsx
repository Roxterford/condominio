import { Suspense } from "react";
import { OperacionesContent } from "./components/operaciones-content";

export default function OperacionesPage() {
  return (
    <Suspense>
      <OperacionesContent />
    </Suspense>
  );
}
