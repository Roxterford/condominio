import type { CodegenConfig } from "@graphql-codegen/cli";

const config: CodegenConfig = {
  schema: "http://localhost:8081/query",
  documents: ["src/**/*.{ts,tsx}"],
  ignoreNoDocuments: true,
  verbose: true,
  generates: {
    "./src/providers/graphql/": {
      preset: "client",
      config: {
        documentMode: "string",
        scalars: {
          DateTime: {
            input: "Date",
            output: "Date",
          },
        },
      },
    },
    "./schema.graphql": {
      plugins: ["schema-ast"],
      config: {
        includeDirectives: true,
      },
    },
  },
};

export default config;
