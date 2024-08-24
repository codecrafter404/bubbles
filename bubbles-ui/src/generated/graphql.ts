import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Int64: { input: any; output: any; }
};

export type CustomItem = {
  __typename?: 'CustomItem';
  /** The dependency of the item for building a tree (only custom items) */
  dependsOn?: Maybe<Scalars['Int']['output']>;
  /** Wheather multiple variants can be selected at once */
  exclusive: Scalars['Boolean']['output'];
  /** Should be unique across items/customitems */
  id: Scalars['Int']['output'];
  name: Scalars['String']['output'];
  variants: Array<Item>;
};

export type CustomItemInput = {
  /** The dependency of the item for building a tree (only custom items) */
  dependsOn?: InputMaybe<Scalars['Int']['input']>;
  /** Wheather multiple variants can be selected at once */
  exclusive: Scalars['Boolean']['input'];
  /** Should be unique across items/customitems */
  id: Scalars['Int']['input'];
  name: Scalars['String']['input'];
  /** The id of the variant items of the input */
  variants: Array<Scalars['Int']['input']>;
};

export type Item = {
  __typename?: 'Item';
  /** Wheather the item is in stock; if all variants of a CustomItems are out of stock, the CustomItem will also go out of stock */
  available: Scalars['Boolean']['output'];
  /** Should be unique across items/customitems */
  id: Scalars['Int']['output'];
  /** Custom text which will be displayed to the user underneath the name */
  identifier: Scalars['String']['output'];
  /** an URL eg.: https://example.com/cheesecake */
  image: Scalars['String']['output'];
  /** False if the item is part of a customitem */
  isOneOff: Scalars['Boolean']['output'];
  /** eg.: Cheesecake */
  name: Scalars['String']['output'];
  price: Scalars['Float']['output'];
};

export type ItemInput = {
  /** Wheather the item is in stock; if all variants of a CustomItems are out of stock, the CustomItem will also go out of stock */
  available: Scalars['Boolean']['input'];
  /** Should be unique across items/customitems */
  id: Scalars['Int']['input'];
  /** Custom text which will be displayed to the user underneath the name */
  identifier: Scalars['String']['input'];
  /** an URL eg.: https://example.com/cheesecake */
  image: Scalars['String']['input'];
  /** False if the item is part of a customitem */
  isOneOff: Scalars['Boolean']['input'];
  /** eg.: Cheesecake */
  name: Scalars['String']['input'];
  price: Scalars['Float']['input'];
};

export type Mutation = {
  __typename?: 'Mutation';
  createCustomItems: Array<Scalars['Int']['output']>;
  createItems: Array<Scalars['Int']['output']>;
  createOrder: Order;
  /** Only available to admins */
  deleteOrder: Scalars['Int']['output'];
  updateCustomItem: CustomItem;
  updateItem: Item;
  updateOrder?: Maybe<Order>;
};


export type MutationCreateCustomItemsArgs = {
  items: Array<CustomItemInput>;
};


export type MutationCreateItemsArgs = {
  items: Array<ItemInput>;
};


export type MutationCreateOrderArgs = {
  order: NewOrder;
};


export type MutationDeleteOrderArgs = {
  order: Scalars['Int']['input'];
};


export type MutationUpdateCustomItemArgs = {
  id: Scalars['Int']['input'];
  item: UpdateCustomItem;
};


export type MutationUpdateItemArgs = {
  id: Scalars['Int']['input'];
  item: UpdateItem;
};


export type MutationUpdateOrderArgs = {
  order: Scalars['Int']['input'];
  state: OrderState;
};

export type NewCustomItem = {
  id: Scalars['Int']['input'];
  quantity: Scalars['Int']['input'];
  variants: Array<Scalars['Int']['input']>;
};

export type NewOrder = {
  customItems: Array<NewCustomItem>;
  items: Array<NewOrderItem>;
  /** The total amount of money earned */
  total: Scalars['Float']['input'];
};

export type NewOrderItem = {
  item: Scalars['Int']['input'];
  quantity: Scalars['Int']['input'];
};

export type Order = {
  __typename?: 'Order';
  customItems: Array<OrderCustomItem>;
  id: Scalars['Int']['output'];
  /** A string generated sequencially to identifiy an OPEN order */
  identifier: Scalars['String']['output'];
  items: Array<OrderItem>;
  state: OrderState;
  timestamp: Scalars['Int64']['output'];
  total: Scalars['Float']['output'];
};

export type OrderCustomItem = {
  __typename?: 'OrderCustomItem';
  customItem: CustomItem;
  quantity: Scalars['Int']['output'];
};

export type OrderItem = {
  __typename?: 'OrderItem';
  item: Item;
  quantity: Scalars['Int']['output'];
};

export enum OrderState {
  Canceled = 'CANCELED',
  Compleated = 'COMPLEATED',
  Created = 'CREATED',
  Pending = 'PENDING'
}

export type Query = {
  __typename?: 'Query';
  getCustomItems: Array<CustomItem>;
  getItems: Array<Item>;
  getOrder?: Maybe<Order>;
  getPermission: User;
};


export type QueryGetOrderArgs = {
  id: Scalars['Int']['input'];
};

export type Statistics = {
  __typename?: 'Statistics';
  totalEarned: Scalars['Float']['output'];
  totalOrders: Scalars['Int']['output'];
  totalOrdersCompleated: Scalars['Int']['output'];
};

export type Subscription = {
  __typename?: 'Subscription';
  /** The next order that should be worked on by the client */
  nextOrder?: Maybe<Order>;
  orders: Array<Order>;
  stats: Statistics;
  updates?: Maybe<UpdateEvent>;
};


export type SubscriptionOrdersArgs = {
  id?: InputMaybe<Scalars['Int']['input']>;
  limit?: InputMaybe<Scalars['Int']['input']>;
  skip?: InputMaybe<Scalars['Int']['input']>;
  sortAsc?: InputMaybe<Scalars['Boolean']['input']>;
  state?: InputMaybe<OrderState>;
};

export type UpdateCustomItem = {
  /** Wheather multiple variants can be selected at once */
  exclusive?: InputMaybe<Scalars['Boolean']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};

export enum UpdateEvent {
  UpdateCustomitem = 'UPDATE_CUSTOMITEM',
  UpdateItem = 'UPDATE_ITEM'
}

export type UpdateItem = {
  /** Wheather the item is in stock; if all variants of a CustomItems are out of stock, the CustomItem will also go out of stock */
  available?: InputMaybe<Scalars['Boolean']['input']>;
  /** Custom text which will be displayed to the user underneath the name */
  identifier?: InputMaybe<Scalars['String']['input']>;
  /** an URL eg.: https://example.com/cheesecake.png */
  image?: InputMaybe<Scalars['String']['input']>;
  /** False if the item is part of a customitem */
  isOneOff?: InputMaybe<Scalars['Boolean']['input']>;
  /** eg.: Cheesecake */
  name?: InputMaybe<Scalars['String']['input']>;
  price?: InputMaybe<Scalars['Float']['input']>;
};

export enum User {
  Admin = 'ADMIN',
  User = 'USER'
}

export type CompleateOrderMutationVariables = Exact<{
  order: Scalars['Int']['input'];
}>;


export type CompleateOrderMutation = { __typename?: 'Mutation', updateOrder?: { __typename?: 'Order', id: number } | null };

export type GetItemsForStoreQueryVariables = Exact<{ [key: string]: never; }>;


export type GetItemsForStoreQuery = { __typename?: 'Query', getItems: Array<{ __typename?: 'Item', id: number, name: string, price: number, image: string, available: boolean, isOneOff: boolean, identifier: string }>, getCustomItems: Array<{ __typename?: 'CustomItem', id: number, name: string, dependsOn?: number | null, exclusive: boolean, variants: Array<{ __typename?: 'Item', id: number }> }> };

export type NextOrderSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type NextOrderSubscription = { __typename?: 'Subscription', nextOrder?: { __typename?: 'Order', id: number, identifier: string, items: Array<{ __typename?: 'OrderItem', quantity: number, item: { __typename?: 'Item', id: number, identifier: string, name: string } }>, customItems: Array<{ __typename?: 'OrderCustomItem', quantity: number, customItem: { __typename?: 'CustomItem', id: number, name: string, variants: Array<{ __typename?: 'Item', id: number, identifier: string, name: string }> } }> } | null };

export type SubmitOrderMutationVariables = Exact<{
  order: NewOrder;
}>;


export type SubmitOrderMutation = { __typename?: 'Mutation', createOrder: { __typename?: 'Order', id: number, identifier: string } };

export type UndoCompleateOrderMutationVariables = Exact<{
  order: Scalars['Int']['input'];
}>;


export type UndoCompleateOrderMutation = { __typename?: 'Mutation', updateOrder?: { __typename?: 'Order', id: number } | null };

export type UndoOrderMutationVariables = Exact<{
  order: Scalars['Int']['input'];
}>;


export type UndoOrderMutation = { __typename?: 'Mutation', updateOrder?: { __typename?: 'Order', id: number } | null };


export const CompleateOrderDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CompleateOrder"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"order"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateOrder"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"order"},"value":{"kind":"Variable","name":{"kind":"Name","value":"order"}}},{"kind":"Argument","name":{"kind":"Name","value":"state"},"value":{"kind":"EnumValue","value":"COMPLEATED"}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<CompleateOrderMutation, CompleateOrderMutationVariables>;
export const GetItemsForStoreDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetItemsForStore"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"getItems"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"price"}},{"kind":"Field","name":{"kind":"Name","value":"image"}},{"kind":"Field","name":{"kind":"Name","value":"available"}},{"kind":"Field","name":{"kind":"Name","value":"isOneOff"}},{"kind":"Field","name":{"kind":"Name","value":"identifier"}}]}},{"kind":"Field","name":{"kind":"Name","value":"getCustomItems"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"dependsOn"}},{"kind":"Field","name":{"kind":"Name","value":"exclusive"}},{"kind":"Field","name":{"kind":"Name","value":"variants"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]}}]} as unknown as DocumentNode<GetItemsForStoreQuery, GetItemsForStoreQueryVariables>;
export const NextOrderDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"subscription","name":{"kind":"Name","value":"NextOrder"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"nextOrder"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"identifier"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"item"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"identifier"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"quantity"}}]}},{"kind":"Field","name":{"kind":"Name","value":"customItems"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"quantity"}},{"kind":"Field","name":{"kind":"Name","value":"customItem"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"variants"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"identifier"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<NextOrderSubscription, NextOrderSubscriptionVariables>;
export const SubmitOrderDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"SubmitOrder"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"order"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"NewOrder"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createOrder"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"order"},"value":{"kind":"Variable","name":{"kind":"Name","value":"order"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"identifier"}}]}}]}}]} as unknown as DocumentNode<SubmitOrderMutation, SubmitOrderMutationVariables>;
export const UndoCompleateOrderDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UndoCompleateOrder"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"order"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateOrder"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"order"},"value":{"kind":"Variable","name":{"kind":"Name","value":"order"}}},{"kind":"Argument","name":{"kind":"Name","value":"state"},"value":{"kind":"EnumValue","value":"CREATED"}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<UndoCompleateOrderMutation, UndoCompleateOrderMutationVariables>;
export const UndoOrderDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UndoOrder"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"order"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateOrder"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"order"},"value":{"kind":"Variable","name":{"kind":"Name","value":"order"}}},{"kind":"Argument","name":{"kind":"Name","value":"state"},"value":{"kind":"EnumValue","value":"CANCELED"}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<UndoOrderMutation, UndoOrderMutationVariables>;