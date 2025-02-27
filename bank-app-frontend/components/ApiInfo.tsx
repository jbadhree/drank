import React from 'react';

const ApiInfo: React.FC = () => {
  // Using the same hardcoded URL that will be replaced at container startup
  const apiUrl = 'http://localhost:8080/api/v1';

  return (
    <div className="fixed bottom-2 right-2 bg-gray-100 border border-gray-300 p-2 rounded text-xs text-gray-700">
      <div><strong>API URL:</strong> {apiUrl}</div>
      <div className="text-xs italic mt-1">If this shows localhost but you're on a remote server, the URL replacement failed.</div>
    </div>
  );
};

export default ApiInfo;
