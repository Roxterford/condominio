import * as v from "valibot";
import { MetodoDePago } from "@/features/pagos/shemas/pago.schema";
import { Moneda } from "@/providers/graphql/graphql";

export enum DestinoDePago {
	DeudaEspecifica = "DEUDA_ESPECIFICA",
	MasAntigua = "MAS_ANTIGUA",
	Recargo = "RECARGO",
	MayorSaldo = "MAYOR_SALDO",
	Prorrateo = "PRORRATEO",
	AbonoCuenta = "ABONO_CUENTA",
}

export const RegistrarPagoFormSchema = v.object({
	unidad: v.pipe(v.string(), v.nonEmpty("Seleccione la unidad")),
	monto: v.pipe(v.number(), v.minValue(0.01, "El monto debe ser mayor a 0")),
	metodo: v.enum(MetodoDePago),
	referencia: v.pipe(v.string(), v.nonEmpty("Ingrese la referencia")),
	moneda: v.enum(Moneda),
	tasa: v.pipe(v.number(), v.minValue(0)),
	concepto: v.pipe(v.string(), v.nonEmpty("Ingrese el concepto")),
	destino: v.optional(v.enum(DestinoDePago)),
	deuda_especifica: v.string(),
});

export type RegistrarPagoForm = v.InferOutput<typeof RegistrarPagoFormSchema>;

export const registrarPagoDefaultValues: RegistrarPagoForm = {
	unidad: "",
	monto: 0,
	metodo: null as any,
	referencia: "",
	moneda: Moneda.Usd,
	tasa: 36.5,
	concepto: "",
	destino: DestinoDePago.AbonoCuenta,
	deuda_especifica: "",
};
