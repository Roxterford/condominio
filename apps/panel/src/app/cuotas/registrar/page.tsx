"use client";

import { Checkbox } from "@/components/ui/checkbox";
import { DatePickerInput } from "@/components/ui/date-picker-input";
import {
  Field,
  FieldContent,
  FieldDescription,
  FieldLabel,
} from "@/components/ui/field";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import { AgregarGastoOverlay } from "@/components/overlay/agregar-gasto-overlay";
import { RegistrarGastoOverlay } from "@/components/overlay/registrar-gasto-overlay";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useOverlay } from "@/hooks/useOverlay";
import { MoreHorizontalIcon, Plus } from "lucide-react";

export default function RegistrarCuota() {
  const agregarGastoOverlay = useOverlay();
  const registrarGastoOverlay = useOverlay();

  return (
    <>
      <h1>Crear Nueva Cuota</h1>
      <p className="text-gray-500">
        Complete el formulario para crear una nueva cuota en el sistema.
      </p>
      <form>
        <section>
          <h3>Tipo de cuota</h3>
          <RadioGroup defaultValue="regular" className="w-fit">
            <Field orientation="horizontal">
              <RadioGroupItem value="regular" id="regular" />
              <FieldContent>
                <FieldLabel htmlFor="regular">Mensualidad Regular</FieldLabel>
                <FieldDescription>Cuota mensual estándar</FieldDescription>
              </FieldContent>
            </Field>
            <Field orientation="horizontal">
              <RadioGroupItem value="especial" id="especial" />
              <FieldContent>
                <FieldLabel htmlFor="especial">Cuota Especial</FieldLabel>
                <FieldDescription>
                  Cuota destinada a algún tipo de inprevisto
                </FieldDescription>
              </FieldContent>
            </Field>
          </RadioGroup>
        </section>
        <section>
          <h3>Periodo</h3>
          {/* Año */}
          <div className="flex gap-5">
            <Field className="">
              <FieldLabel>Año</FieldLabel>
              <Select>
                <SelectTrigger>
                  <SelectValue placeholder="Seleccione el año" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem value="2023">2023</SelectItem>
                    <SelectItem value="2024">2024</SelectItem>
                    <SelectItem value="2025">2025</SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
            </Field>
            {/* Mes */}
            <Field className="">
              <FieldLabel>Mes</FieldLabel>
              <Select>
                <SelectTrigger>
                  <SelectValue placeholder="Seleccione el mes" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem value="1">Enero</SelectItem>
                    <SelectItem value="2">Febrero</SelectItem>
                    <SelectItem value="3">Marzo</SelectItem>
                    <SelectItem value="4">Abril</SelectItem>
                    <SelectItem value="5">Mayo</SelectItem>
                    <SelectItem value="6">Junio</SelectItem>
                    <SelectItem value="7">Julio</SelectItem>
                    <SelectItem value="8">Agosto</SelectItem>
                    <SelectItem value="9">Septiembre</SelectItem>
                    <SelectItem value="10">Octubre</SelectItem>
                    <SelectItem value="11">Noviembre</SelectItem>
                    <SelectItem value="12">Diciembre</SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
            </Field>
          </div>
          <Field orientation="horizontal">
            <Checkbox id="anio_actual" name="anio_actual" />
            <FieldLabel htmlFor="anio_actual">Usar año actual</FieldLabel>
          </Field>

          <div className="flex gap-5">
            {/* TODO: mejorar el input con su mask y mejor seleccion */}
            <Field className="">
              <FieldLabel htmlFor="fecha_emision">Fecha de emisión</FieldLabel>
              <DatePickerInput
                placeholder="dd/mm/aaaa"
                locale="es-VE"
                id="fecha_emision"
              />
            </Field>
            {/* TODO: mejorar el input con su mask y mejor seleccion */}
            <Field className="">
              <FieldLabel htmlFor="fecha_limite">
                Fecha limite de pago
              </FieldLabel>
              <DatePickerInput
                placeholder="dd/mm/aaaa"
                locale="es-VE"
                id="fecha_limite"
              />
            </Field>
          </div>
        </section>

        <section>
          <div className="flex justify-between items-center">
            <h3>Desglose de gastos</h3>
            <Button
              type="button"
              variant="outline"
              onClick={registrarGastoOverlay.open}
            >
              <Plus /> Agregar Gasto
            </Button>
          </div>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="font-semibold uppercase">ID</TableHead>
                <TableHead className="font-semibold uppercase">
                  Concepto
                </TableHead>
                <TableHead className="font-semibold uppercase">Monto</TableHead>
                <TableHead className="font-semibold uppercase">Fecha</TableHead>
                <TableHead className="font-semibold uppercase">
                  Proveedor
                </TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow>
                <TableCell className="font-medium">bogy5a</TableCell>
                <TableCell>$29.99</TableCell>
                <TableCell className="text-right">
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="icon" className="size-8">
                        <MoreHorizontalIcon />
                        <span className="sr-only">Open menu</span>
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem>Edit</DropdownMenuItem>
                      <DropdownMenuItem>Duplicate</DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem variant="destructive">
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell className="font-medium">
                  Mechanical Keyboard
                </TableCell>
                <TableCell>$129.99</TableCell>
                <TableCell className="text-right">
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="icon" className="size-8">
                        <MoreHorizontalIcon />
                        <span className="sr-only">Open menu</span>
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem>Edit</DropdownMenuItem>
                      <DropdownMenuItem>Duplicate</DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem variant="destructive">
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell className="font-medium">USB-C Hub</TableCell>
                <TableCell>$49.99</TableCell>
                <TableCell className="text-right">
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="icon" className="size-8">
                        <MoreHorizontalIcon />
                        <span className="sr-only">Open menu</span>
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem>Edit</DropdownMenuItem>
                      <DropdownMenuItem>Duplicate</DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem variant="destructive">
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </section>
      </form>

      <AgregarGastoOverlay {...agregarGastoOverlay.overlayProps} />
      <RegistrarGastoOverlay {...registrarGastoOverlay.overlayProps} />
    </>
  );
}
