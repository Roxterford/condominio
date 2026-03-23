import { ApolloClient, HttpLink, InMemoryCache } from "@apollo/client";
import { ApolloProvider } from "@apollo/client/react";
  // import { setClient } from "svelte-apollo";

const client = new ApolloClient({
  // TODO: env uri
  link: new HttpLink({ uri: "http://localhost:8081/query" }),
  cache: new InMemoryCache(),
});

// setClient(client)

export { ApolloProvider, client };

export default client;
