'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { isAuthenticated } from '@/lib/auth';

export default function HomePage() {
  const router = useRouter();

  useEffect(() => {
    // Redirect to dashboard if already authenticated, otherwise to login
    if (isAuthenticated()) {
      router.push('/dashboard');
    } else {
      router.push('/login');
    }
  }, [router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-primary-600">ByteBank</h1>
        <p className="mt-2 text-gray-600">A tech-driven financial institution</p>
        <p className="mt-2 text-gray-600">Redirecting you to the appropriate page...</p>
      </div>
    </div>
  );
}
