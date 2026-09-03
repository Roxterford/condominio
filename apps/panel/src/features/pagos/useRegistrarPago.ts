import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { RegistrarPagoDto } from "@/providers/graphql/graphql";
import { useMutation } from "@tanstack/react-query";

const Mutation = graphql(/* GraphQL */`
    mutation RegistrarPago($input: RegistrarPagoDTO!){
        registrarPago(input: $input) {
            id
        }
    }
`)

export function useRegistrarPago() {
    return useMutation({
        mutationKey: ["pagos.registrar"],
        mutationFn: (input: RegistrarPagoDto) => execute(Mutation, {input: {
            concepto: input.concepto,
            metodo: input.metodo,
            moneda: input.moneda,
            monto: input.monto,
            tasa: input.tasa,
            unidad: input.unidad,
        }})
    })

}