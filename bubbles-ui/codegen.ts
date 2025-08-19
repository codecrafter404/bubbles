
import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  schema: '../server/gql_schema/*.graphql',
  documents: ['src/**/*.svelte'],
  ignoreNoDocuments: true, // for better experience with the watcher
  generates: {
    './src/gql/': {
      preset: 'client',
      plugins: ['typescript'],
      config: {
        useTypeImports: true,
        scalars: {
          ID: {
            input: 'number',
            output: 'number'
          },
          Time: {
            input: 'string',
            output: 'string'
          }
        }
      }
    },
  },
};

export default config;
