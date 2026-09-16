import { Suspense } from "react";
import { VillasPageContent } from "./components/villas-page-content";

export default function VillasPage() {
  return (
    <Suspense>
      <VillasPageContent />
    </Suspense>
  );
}