'use client';

import React, { useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import ApiInfo from './ApiInfo';

const ApiInfoWrapper: React.FC = () => {
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    return () => setMounted(false);
  }, []);

  if (!mounted) return null;

  const portalRoot = document.getElementById('api-info-root');
  if (!portalRoot) return null;

  return createPortal(<ApiInfo />, portalRoot);
};

export default ApiInfoWrapper;
