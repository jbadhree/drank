import React, { useEffect, useState } from 'react';

const ApiInfo: React.FC = () => {
  const [apiUrl, setApiUrl] = useState<string>('Loading...');
  
  useEffect(() => {
    // Get the API URL directly from the environment variable
    const url = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    const apiPath = url.endsWith('/api/v1') ? url : `${url}/api/v1`;
    setApiUrl(apiPath);

    // Log it to the console for debugging
    console.log('NEXT_PUBLIC_API_URL:', process.env.NEXT_PUBLIC_API_URL);
    console.log('Actual API URL used:', apiPath);
  }, []);

  return (
    <div className="fixed bottom-2 right-2 bg-gray-100 border border-gray-300 p-2 rounded text-xs text-gray-700">
      <div><strong>API URL:</strong> {apiUrl}</div>
      <div className="text-xs italic mt-1">If this is localhost but you're on a remote server, check the deployment guide.</div>
    </div>
  );
};

export default ApiInfo;
