
import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  overwrite: true,
  schema: "../server/graph/schema.graphqls",
  documents: "./src/gql-queries/**/*.{gql,graphql}",
  generates: {
    "src/generated/graphql.ts": {
      plugins: ['typescript', 'typescript-operations', 'typed-document-node'],
      config: {
        useTypeImports: true
      }
    }
  }
};

export default config;
