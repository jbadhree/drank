import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { User } from './types';
import { getCurrentUser } from './api';

interface UseAuthReturn {
  user: User | null;
  loading: boolean;
  login: (token: string) => void;
  logout: () => void;
}

export const useAuth = (): UseAuthReturn => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const router = useRouter();

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const token = localStorage.getItem('token');
        if (!token) {
          setLoading(false);
          return;
        }

        const user = await getCurrentUser();
        setUser(user);
      } catch (error) {
        localStorage.removeItem('token');
        console.error('Authentication failed:', error);
      } finally {
        setLoading(false);
      }
    };

    checkAuth();
  }, []);

  const login = (token: string) => {
    localStorage.setItem('token', token);
    router.push('/dashboard');
  };

  const logout = () => {
    localStorage.removeItem('token');
    setUser(null);
    router.push('/login');
  };

  return { user, loading, login, logout };
};

export const isAuthenticated = (): boolean => {
  if (typeof window === 'undefined') return false;
  return Boolean(localStorage.getItem('token'));
};
