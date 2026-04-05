import { Button } from "@/components/ui/button";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { CuotasTable } from "@/features/administracion/components/cuotas_table/cuotas_table";
import { Plus } from "lucide-react";

export default function CuotasPage() {
  return (
    <>
      <header className="flex items-center justify-between">
        <div>
          <h1 className="page-title">Cuotas</h1>
          <p className="page-description">
            Gestión de cuotas mensuales y espaciales
          </p>
        </div>

        <Button>
          <Plus /> Nueva cuota
        </Button>
      </header>
      <Tabs defaultValue="regulares">
        <TabsList variant="line">
          <TabsTrigger value="regulares">Mensualidades</TabsTrigger>
          <TabsTrigger value="especiales">Cuotas Especiales</TabsTrigger>
        </TabsList>
      </Tabs>
      <CuotasTable data={[]} />
    </>
  );
}
