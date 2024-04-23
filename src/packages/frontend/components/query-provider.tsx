"use client"
import { PathInfo } from '@/types/result';
import React, { createContext, Dispatch, ReactNode, useContext, useReducer } from 'react';


export type GraphLinks = {
  source: string;
  targets: Set<string>;
};

export type LinkNodes = Record<string, GraphLinks>;

// Define the initial state and action types
export type State = {
  source: string;
  destination: string;
  selectedSource: string;
  selectedDestination: string;
  result: PathInfo[][];
  isBFS: boolean;
  bonus: boolean;
  nodes: Set<string>;
  linkNodes: LinkNodes;
};

export type Action =
  | { type: 'SET_SOURCE'; payload: string }
  | { type: 'SET_DESTINATION'; payload: string }
  | { type: 'SET_SELECTED_SOURCE'; payload: string }
  | { type: 'SET_SELECTED_DESTINATION'; payload: string }
  | { type: 'SET_RESULT'; payload: PathInfo[][] }
  | { type: 'SET_ISBFS'; payload: boolean }
  | { type: 'SET_BONUS'; payload: boolean }
  | { type: 'SET_NODES'; payload: Set<string> }
  | { type: 'SET_LINK_NODES'; payload: LinkNodes }
  | { type: 'SWAP' };

const initialState: State = {
  source: '',
  destination: '',
  selectedSource: '',
  selectedDestination: '',
  result: [],
  isBFS: true,
  bonus: false,
  nodes: new Set<string>(),
  linkNodes: {},
};

// Reducer function
const reducer = (state: State, action: Action): State => {
  switch (action.type) {
    case 'SET_SOURCE':
      return { ...state, source: action.payload };
    case 'SET_DESTINATION':
      return { ...state, destination: action.payload };
    case 'SET_SELECTED_SOURCE':
      return { ...state, selectedSource: action.payload };
    case 'SET_SELECTED_DESTINATION':
      return { ...state, selectedDestination: action.payload };
    case 'SET_RESULT':
      return { ...state, result: action.payload };
    case 'SET_ISBFS':
      return { ...state, isBFS: action.payload };
    case 'SET_BONUS':
      return { ...state, bonus: action.payload };
    case 'SET_NODES':
      return { ...state, nodes: action.payload };
    case 'SET_LINK_NODES':
      return { ...state, linkNodes: action.payload };
    case 'SWAP':
      return {
        ...state,
        source: state.destination,
        destination: state.source,
        selectedSource: state.selectedDestination,
        selectedDestination: state.selectedSource,
      };
    default:
      return state;
  }
};

// Create the context
export const QueryContext = createContext<{ state: State; dispatch: Dispatch<Action> } | undefined>(undefined);

// Provider component
type QueryProviderProps = {
  children: ReactNode;
};

export const QueryProvider: React.FC<QueryProviderProps> = ({ children }) => {
  const [state, dispatch] = useReducer(reducer, initialState);

  return <QueryContext.Provider value={{ state, dispatch }}>{children}</QueryContext.Provider>;
};

// Custom hook to use the context
export const useQueryContext = () => {
  const context = useContext(QueryContext);
  if (!context) {
    throw new Error('useQueryContext must be used within a QueryProvider');
  }
  return context;
};
