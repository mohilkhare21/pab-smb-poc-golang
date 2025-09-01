import React, { createContext, useContext, useReducer, useEffect, ReactNode } from 'react';
import { useAuth0 } from '@auth0/auth0-react';
import apiService from '../services/api';

// Simple types for now
interface User {
  id: string;
  email: string;
  name: string;
  picture?: string;
  company_id: string;
  role: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  onboarded: boolean;
  invitation_status: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

interface LoginRequest {
  email: string;
  password: string;
}

interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

interface AuthContextType extends AuthState {
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  verifyToken: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

type AuthAction =
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_USER'; payload: User | null }
  | { type: 'SET_TOKEN'; payload: string | null }
  | { type: 'SET_AUTHENTICATED'; payload: boolean }
  | { type: 'LOGOUT' };

const authReducer = (state: AuthState, action: AuthAction): AuthState => {
  switch (action.type) {
    case 'SET_LOADING':
      return { ...state, isLoading: action.payload };
    case 'SET_USER':
      return { ...state, user: action.payload };
    case 'SET_TOKEN':
      return { ...state, token: action.payload };
    case 'SET_AUTHENTICATED':
      return { ...state, isAuthenticated: action.payload };
    case 'LOGOUT':
      return {
        user: null,
        token: null,
        isAuthenticated: false,
        isLoading: false,
      };
    default:
      return state;
  }
};

const initialState: AuthState = {
  user: null,
  token: null,
  isAuthenticated: false,
  isLoading: true,
};

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [state, dispatch] = useReducer(authReducer, initialState);
  const { 
    isAuthenticated: auth0Authenticated, 
    isLoading: auth0Loading, 
    user: auth0User, 
    loginWithRedirect, 
    logout: auth0Logout,
    getAccessTokenSilently 
  } = useAuth0();

  // Sync Auth0 state with our context
  useEffect(() => {
    if (auth0Loading) {
      dispatch({ type: 'SET_LOADING', payload: true });
    } else {
      dispatch({ type: 'SET_LOADING', payload: false });
      
      if (auth0Authenticated && auth0User) {
        // Convert Auth0 user to our User type
        const user: User = {
          id: auth0User.sub || '',
          email: auth0User.email || '',
          name: auth0User.name || '',
          picture: auth0User.picture,
          company_id: '', // Will be set when user creates/joins a company
          role: 'admin', // First user becomes admin
          is_active: true,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          onboarded: false,
          invitation_status: 'active'
        };
        
        dispatch({ type: 'SET_USER', payload: user });
        dispatch({ type: 'SET_AUTHENTICATED', payload: true });
        
        // Get access token for API calls
        getAccessTokenSilently().then(token => {
          dispatch({ type: 'SET_TOKEN', payload: token });
          localStorage.setItem('auth_token', token);
        });
      } else {
        dispatch({ type: 'LOGOUT' });
      }
    }
  }, [auth0Authenticated, auth0Loading, auth0User, getAccessTokenSilently]);

  const login = async (data: LoginRequest) => {
    // Redirect to Auth0 login
    await loginWithRedirect({
      appState: { returnTo: window.location.pathname }
    });
  };

  const register = async (data: RegisterRequest) => {
    // Redirect to Auth0 signup
    await loginWithRedirect({
      screen_hint: 'signup',
      appState: { returnTo: window.location.pathname }
    });
  };

  const logout = async () => {
    localStorage.removeItem('auth_token');
    dispatch({ type: 'LOGOUT' });
    await auth0Logout({
      logoutParams: {
        returnTo: window.location.origin
      }
    });
  };

  const verifyToken = async () => {
    // Auth0 handles token verification automatically
    return;
  };

  const value: AuthContextType = {
    ...state,
    login,
    register,
    logout,
    verifyToken,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
