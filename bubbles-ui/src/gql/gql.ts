/* eslint-disable */
import * as types from './graphql';
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 * Learn more about it here: https://the-guild.dev/graphql/codegen/plugins/presets/preset-client#reducing-bundle-size
 */
type Documents = {
    "\n\t\tquery QueryItemAndCustomItem {\n\t\t\titems {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tprice\n\t\t\t\timage\n\t\t\t\tinStock\n\t\t\t\tnotes\n\t\t\t}\n\t\t\tcustomItems {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tallowOnlyOne\n\t\t\t\tprev\n\t\t\t\tvariants {\n\t\t\t\t\tid\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\t": typeof types.QueryItemAndCustomItemDocument,
    "\n\t\tmutation CreateOrder($order: CreateOrderInput!) {\n\t\t\tcreateOrder(input: $order) {\n\t\t\t\tid\n\t\t\t\tidentifier\n\t\t\t}\n\t\t}\n\t": typeof types.CreateOrderDocument,
};
const documents: Documents = {
    "\n\t\tquery QueryItemAndCustomItem {\n\t\t\titems {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tprice\n\t\t\t\timage\n\t\t\t\tinStock\n\t\t\t\tnotes\n\t\t\t}\n\t\t\tcustomItems {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tallowOnlyOne\n\t\t\t\tprev\n\t\t\t\tvariants {\n\t\t\t\t\tid\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\t": types.QueryItemAndCustomItemDocument,
    "\n\t\tmutation CreateOrder($order: CreateOrderInput!) {\n\t\t\tcreateOrder(input: $order) {\n\t\t\t\tid\n\t\t\t\tidentifier\n\t\t\t}\n\t\t}\n\t": types.CreateOrderDocument,
};

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = graphql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function graphql(source: string): unknown;

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tquery QueryItemAndCustomItem {\n\t\t\titems {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tprice\n\t\t\t\timage\n\t\t\t\tinStock\n\t\t\t\tnotes\n\t\t\t}\n\t\t\tcustomItems {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tallowOnlyOne\n\t\t\t\tprev\n\t\t\t\tvariants {\n\t\t\t\t\tid\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\t"): (typeof documents)["\n\t\tquery QueryItemAndCustomItem {\n\t\t\titems {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tprice\n\t\t\t\timage\n\t\t\t\tinStock\n\t\t\t\tnotes\n\t\t\t}\n\t\t\tcustomItems {\n\t\t\t\tid\n\t\t\t\tname\n\t\t\t\tallowOnlyOne\n\t\t\t\tprev\n\t\t\t\tvariants {\n\t\t\t\t\tid\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tmutation CreateOrder($order: CreateOrderInput!) {\n\t\t\tcreateOrder(input: $order) {\n\t\t\t\tid\n\t\t\t\tidentifier\n\t\t\t}\n\t\t}\n\t"): (typeof documents)["\n\t\tmutation CreateOrder($order: CreateOrderInput!) {\n\t\t\tcreateOrder(input: $order) {\n\t\t\t\tid\n\t\t\t\tidentifier\n\t\t\t}\n\t\t}\n\t"];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;