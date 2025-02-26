import { useAuth } from '@/lib/auth';
import Link from 'next/link';

const NavBar = () => {
  const { user, logout } = useAuth();

  return (
    <nav className="bg-primary-600 text-white shadow-md">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <Link href="/dashboard" className="text-xl font-bold">
                Drank Banking
              </Link>
            </div>
          </div>
          <div className="flex items-center">
            {user && (
              <>
                <span className="mr-4">Hello, {user.firstName}</span>
                <button 
                  onClick={logout} 
                  className="bg-primary-700 hover:bg-primary-800 text-white font-medium py-2 px-4 rounded transition-colors"
                >
                  Logout
                </button>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
};

export default NavBar;
