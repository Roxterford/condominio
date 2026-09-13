import { OverlayProps } from "@/components/overlay/overlay";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Check, Plus, X } from "lucide-react";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useDebounce } from "@/hooks/useDebounce";
import { useOverlay } from "@/hooks/useOverlay";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { useQuery } from "@tanstack/react-query";
import { SearchIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { GastoAProveedor } from "@/providers/graphql/graphql";
import { toast } from "sonner";

export interface SeleccionarGastosOverlayProps
  extends Omit<OverlayProps, "onDone"> {
  omitIDs?: string[];
  onDone?(values: Array<GastoAProveedor>): void;
}

export function SeleccionarGastosOverlay(props: SeleccionarGastosOverlayProps) {
  const [gastos_selectos, setGastosSelectos] = useState<Array<GastoAProveedor>>(
    [],
  );

  return (
    <Dialog open={props.open} onOpenChange={props.onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Agregar Gastos</DialogTitle>
          <DialogDescription>
            Se encontraron gastos no asociados a ninguna cuota.
          </DialogDescription>
        </DialogHeader>
        <form>
          <Busqueda
            onAdd={(it) => setGastosSelectos((s) => [...s, it])}
            omitIDs={gastos_selectos
              .map((it) => it.operacion ?? "")
              .concat(props.omitIDs ?? [])}
          />
        </form>
        <section>
          <p className="font-semibold">
            Gastos
            {gastos_selectos.length > 0 && (
              <>
                {" "}
                <span>({gastos_selectos.length})</span>
              </>
            )}
          </p>
          <hr />
          <ul className="grid gap-5">
            {gastos_selectos.map((g) => (
              <li
                className="flex justify-between items-center py-2"
                key={g.operacion}
              >
                <div>
                  <p className="font-semibold">{g.concepto}</p>
                  <p className="text-sm text-gray-500">
                    {g.fecha?.toLocaleString()}
                  </p>
                </div>
                <div></div>
                <div className="flex gap-2">
                  <Button variant={"outline"}>Ver Detalles</Button>
                  <Button
                    variant={"outline"}
                    onClick={() =>
                      setGastosSelectos((s) =>
                        s.filter((it) => it.operacion != g.operacion),
                      )
                    }
                  >
                    Remover <X />
                  </Button>
                </div>
              </li>
            ))}
          </ul>
        </section>
        <section className="flex gap-2 justify-end items-center">
          <Button variant={"outline"}>Atras</Button>
          <Button
            onClick={() => {
              setGastosSelectos([]);
              props.onDone?.(gastos_selectos);
            }}
          >
            Añadir <Check />
          </Button>
        </section>
      </DialogContent>
    </Dialog>
  );
}

// TODO: Añade un filtro para omitir los gastos omitidos (gastos ya selectos)
const BusquedaQuery = graphql(/* GraphQL */ `
  query BuscarGastosHuerfanos($busqueda: String) {
    gastos: obtenerGastos(
      filter: { concepto: { like: $busqueda }, cuota: { eq: null } }
    ) {
      data {
        __typename
        ... on Gasto {
          monto
          total
          operacion
          concepto
          metodo
          moneda
          tasa
          fecha
          metodo
          registrado_por
          registro
        }
        ... on GastoAProveedor {
          proveedor {
            id
            nombre
            rif
            telefono
            email
            actualizado_en
            creado_en
            direccion
          }
        }
      }
    }
  }
`);

interface BusquedaProps {
  onAdd(transaccion: GastoAProveedor): void;

  omitIDs: string[];
}

function Busqueda({ onAdd, omitIDs }: BusquedaProps) {
  const state = useOverlay();

  const [busqueda, setBusqueda] = useState("");
  const buscarOperaciones = useQuery({
    queryKey: ["operaciones.search", busqueda],
    queryFn: ({ queryKey: [, busqueda] }) =>
      execute(BusquedaQuery, { busqueda }),
  });

  const handleDebouceChange = useDebounce((v: string) => {
    setBusqueda(`%${v}%`);
  });

  useEffect(() => {
    if (busqueda === "") return;
    state.open();
  }, [busqueda]);

  useEffect(() => {
    if (
      buscarOperaciones.isSuccess &&
      buscarOperaciones.data &&
      buscarOperaciones.data.errors
    ) {
      console.error(buscarOperaciones.data.errors);
      toast.error("error", {
        description: () => (
          <pre>{JSON.stringify(buscarOperaciones.data?.errors, null, 2)}</pre>
        ),
      });
    }
  }, [buscarOperaciones.isSuccess, buscarOperaciones.data]);

  return (
    <Popover open={state.isOpen}>
      <PopoverTrigger
        nativeButton={false}
        render={
          <InputGroup>
            <InputGroupAddon align="inline-start">
              <SearchIcon className="text-muted-foreground" />
            </InputGroupAddon>
            <InputGroupInput
              onChange={(e) => handleDebouceChange(e.target.value)}
              onBlur={state.close}
              onFocus={state.open}
              placeholder="Ej. Reparación de Bomba de Agua"
            ></InputGroupInput>
          </InputGroup>
        }
      />
      <PopoverContent className="w-(--anchor-width)" initialFocus={false}>
        <ul className="grid gap-5">
          {buscarOperaciones.data?.data?.gastos.data
            .filter((it) => !omitIDs.includes(it.operacion))
            .map((gasto) => (
              <li
                className="flex justify-between items-center py-2"
                key={gasto.operacion}
              >
                <div>
                  <p className="font-semibold">{gasto.concepto}</p>
                  <p className="text-sm text-gray-500">
                    {gasto.fecha.toLocaleDateString()}
                  </p>
                </div>
                <div></div>
                <div className="flex gap-2">
                  <Button variant={"outline"}>Ver Detalles</Button>
                  <Button
                    onClick={() => {
                      console.log({ tx: gasto });

                      if (gasto.__typename !== "GastoAProveedor") return;

                      onAdd(gasto);
                    }}
                  >
                    Agregar <Plus />
                  </Button>
                </div>
              </li>
            ))}
        </ul>
      </PopoverContent>
    </Popover>
  );
}
