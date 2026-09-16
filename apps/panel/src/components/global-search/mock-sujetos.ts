export type MockSujeto = {
  id: string;
  tipo: "PERSONA_NATURAL" | "ENTE_JURIDICO";
  display_name: string;
  documento_identidad: string;
  unidad_codigo: string;
  rol: string;
};

export const MOCK_SUJETOS: MockSujeto[] = [
  {
    id: "mock-sujeto-juan",
    tipo: "PERSONA_NATURAL",
    display_name: "Juan Andrés Ramírez",
    documento_identidad: "V-12345678",
    unidad_codigo: "villa-1",
    rol: "Propietario",
  },
  {
    id: "mock-sujeto-santiago",
    tipo: "PERSONA_NATURAL",
    display_name: "Santiago Mariño",
    documento_identidad: "V-12345679",
    unidad_codigo: "villa-2",
    rol: "Representante",
  },
  {
    id: "mock-sujeto-roxterford",
    tipo: "ENTE_JURIDICO",
    display_name: "Roxterford",
    documento_identidad: "J-12345678",
    unidad_codigo: "villa-2",
    rol: "Titular",
  },
  {
    id: "mock-sujeto-ana",
    tipo: "PERSONA_NATURAL",
    display_name: "Ana Pérez",
    documento_identidad: "V-20987123",
    unidad_codigo: "villa-10",
    rol: "Propietario",
  },
  {
    id: "mock-sujeto-carlos",
    tipo: "PERSONA_NATURAL",
    display_name: "Carlos Ramírez",
    documento_identidad: "V-15023456",
    unidad_codigo: "villa-5",
    rol: "Propietario",
  },
];

export function buscarMockSujetos(query: string, limit = 5): MockSujeto[] {
  const q = query.trim().toLocaleLowerCase("es");
  if (!q) return [];
  return MOCK_SUJETOS.filter((s) =>
    s.display_name.toLocaleLowerCase("es").includes(q),
  ).slice(0, limit);
}
