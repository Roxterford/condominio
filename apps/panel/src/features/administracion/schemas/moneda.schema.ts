import * as v from "valibot";

export enum Moneda {
  USD = "USD",
  VED = "VED",
}

export const MonedaSchema = v.enum(Moneda);
