import React, { useEffect, useState } from 'react';
import getConfig from 'next/config';

const ApiInfo: React.FC = () => {
  const [apiUrl, setApiUrl] = useState<string>('Loading...');
  const [serverApiUrl, setServerApiUrl] = useState<string>('Checking server...');
  
  useEffect(() => {
    const { publicRuntimeConfig } = getConfig() || {};
    
    // Get API URL from runtime config or fall back to environment variable
    let baseUrl = publicRuntimeConfig?.apiUrl || 
                  process.env.NEXT_PUBLIC_API_URL || 
                  'http://localhost:8080';
    if (!baseUrl.endsWith('/api/v1')) {
      baseUrl += '/api/v1';
    }
    
    setApiUrl(baseUrl);

    // Also check what the server sees
    fetch('/api/config')
      .then(res => res.json())
      .then(data => {
        setServerApiUrl(data.apiUrl || 'Not available');
      })
      .catch(err => {
        console.error('Error fetching server config:', err);
        setServerApiUrl('Error fetching');
      });
  }, []);

  return (
    <div className="fixed bottom-2 right-2 bg-gray-100 border border-gray-300 p-2 rounded text-xs text-gray-700">
      <div>Client API URL: {apiUrl}</div>
      <div>Server ENV: {serverApiUrl}</div>
    </div>
  );
};

export default ApiInfo;
