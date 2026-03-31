"use client";

import { TanStackDevtools } from "@tanstack/react-devtools";
import { formDevtoolsPlugin } from "@tanstack/react-form-devtools";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

export function Devtools() {
  return (
    <>
      <ReactQueryDevtools initialIsOpen={false} />
      <TanStackDevtools plugins={[formDevtoolsPlugin()]} />;
    </>
  );
}
