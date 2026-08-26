"use client";

import { Button } from "@/components/ui/button";
import { graphql } from "@/providers/graphql";
import { execute } from "@/providers/graphql/execute";
import { LoginMutationVariables } from "@/providers/graphql/graphql";
import { useMutation } from "@tanstack/react-query";

const LoginMutation = graphql(/* GraphQL */ `
  mutation Login($email: String!, $pass: String!) {
    login(email: $email, password: $pass) {
      token
    }
  }
`);

export default function Home() {
  const login = useMutation({
    // mutationKey: ["login"],
    mutationFn: (credentials: LoginMutationVariables) =>
      execute(LoginMutation, credentials),
    onSuccess: (data) => {
      const token = data.data?.login?.token;
      if (token) {
        document.cookie = `api_token=${token}; path=/; max-age=604800; SameSite=Lax`;
      }
    },
  });

  function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = e.currentTarget;
    const email = (form.elements.namedItem("email") as HTMLInputElement).value;
    const pass = (form.elements.namedItem("password") as HTMLInputElement)
      .value;
    login.mutate({ email, pass });
  }

  return (
    <>
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="max-w-md w-full space-y-8 p-8 bg-white rounded-xl shadow-lg">
          <div className="text-center">
            <h2 className="text-3xl font-bold text-gray-900">Iniciar sesión</h2>
            <p className="mt-2 text-sm text-gray-600">
              Ingresa tus credenciales para continuar
            </p>
          </div>
          <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
            <div className="space-y-4">
              <div>
                <label
                  htmlFor="email"
                  className="block text-sm font-medium text-gray-700"
                >
                  Correo electrónico
                </label>
                <input
                  id="email"
                  name="email"
                  type="email"
                  autoComplete="email"
                  required
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="tu@email.com"
                />
              </div>
              <div>
                <label
                  htmlFor="password"
                  className="block text-sm font-medium text-gray-700"
                >
                  Contraseña
                </label>
                <input
                  id="password"
                  name="password"
                  type="password"
                  autoComplete="current-password"
                  required
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="••••••••"
                />
              </div>
            </div>
            {login.isError && (
              <p className="text-sm text-red-600 text-center">
                Credenciales incorrectas. Intenta de nuevo.
              </p>
            )}
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <input
                  id="remember-me"
                  name="remember-me"
                  type="checkbox"
                  className="h-4 w-4 text-blue-600 border-gray-300 rounded"
                />
                <label
                  htmlFor="remember-me"
                  className="ml-2 block text-sm text-gray-700"
                >
                  Recordarme
                </label>
              </div>
              <a href="#" className="text-sm text-blue-600 hover:text-blue-500">
                ¿Olvidaste tu contraseña?
              </a>
            </div>
            <Button type="submit" className="w-full" disabled={login.isPending}>
              {login.isPending ? "Iniciando sesión..." : "Iniciar sesión"}
            </Button>
          </form>
          {login.isSuccess && login.data && (
            <div className="mt-6 p-4 bg-green-50 border border-green-200 rounded-md">
              <h3 className="text-sm font-medium text-green-800 mb-2">
                Sesión iniciada correctamente
              </h3>
              <pre className="text-xs text-green-700 overflow-auto whitespace-pre-wrap break-all">
                {JSON.stringify(login.data, null, 2)}
              </pre>
            </div>
          )}
          {login.isError && (
            <div className="mt-6 p-4 bg-red-50 border border-red-200 rounded-md">
              <h3 className="text-sm font-medium text-red-800 mb-2">
                Error al iniciar sesión
              </h3>
              <pre className="text-xs text-red-700 overflow-auto whitespace-pre-wrap break-all">
                {JSON.stringify(login.error, null, 2)}
              </pre>
            </div>
          )}
        </div>
      </div>
    </>
  );
}
