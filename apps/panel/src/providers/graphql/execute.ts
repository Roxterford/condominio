import type { TypedDocumentString } from "./graphql";

interface Error {
  message: string;
  locations?: Array<{ line: number; column: number }>;
  path?: string[];
  extensions?: Record<string, unknown>;
}

interface GraphqlResponse<TResult> {
  data: TResult;
  errors?: Error[];
}

export async function execute<TResult, TVariables>(
  query: TypedDocumentString<TResult, TVariables>,
  ...[variables]: TVariables extends Record<string, never> ? [] : [TVariables]
): Promise<GraphqlResponse<TResult>>;

export async function execute<TResult, TVariables>(
  query: TypedDocumentString<TResult, TVariables>,
  options: RequestInit,
  ...[variables]: TVariables extends Record<string, never> ? [] : [TVariables]
): Promise<GraphqlResponse<TResult>>;

export async function execute<TResult, TVariables>(
  query: TypedDocumentString<TResult, TVariables>,
  optionsOrVariables?:
    | RequestInit
    | (TVariables extends Record<string, never> ? never : TVariables),
  ...[variables]: TVariables extends Record<string, never> ? [] : [TVariables]
): Promise<GraphqlResponse<TResult>> {
  const isOptions =
    optionsOrVariables &&
    typeof optionsOrVariables === "object" &&
    "headers" in optionsOrVariables;
  const options = isOptions ? (optionsOrVariables as RequestInit) : undefined;
  const vars = isOptions
    ? (variables as TVariables extends Record<string, never>
        ? never
        : TVariables)
    : (optionsOrVariables as TVariables);

  const endpoint =
    process.env.GRAPHQL_ENDPOINT || "http://localhost:8081/query";

  let token: string | undefined;

  // 1. Lógica para obtener el token según el entorno
  if (typeof window === "undefined") {
    // ESTAMOS EN EL SERVIDOR (Server Components / Actions)
    // Usamos importación dinámica para evitar errores en el cliente
    const { cookies } = await import("next/headers");
    const cookieStore = await cookies();
    token = cookieStore.get("api_token")?.value;
  } else {
    // ESTAMOS EN EL CLIENTE (Navegador)
    token = document.cookie
      .split("; ")
      .find((row) => row.startsWith("api_token="))
      ?.split("=")[1];
  }

  const response = await fetch(endpoint, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/graphql-response+json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options?.headers,
    },
    body: JSON.stringify({
      query,
      variables: vars,
    }),
  });

  if (!response.ok) {
    throw new Error("Network response was not ok");
  }

  const result = await response.json();
  return result as GraphqlResponse<TResult>;
}
