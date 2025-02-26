import React from 'react';

interface SpinnerProps {
  size?: 'small' | 'medium' | 'large';
  color?: 'primary' | 'white';
}

const Spinner: React.FC<SpinnerProps> = ({ 
  size = 'medium', 
  color = 'primary' 
}) => {
  const sizeClass = {
    small: 'w-4 h-4',
    medium: 'w-6 h-6',
    large: 'w-8 h-8',
  }[size];

  const colorClass = {
    primary: 'text-primary-600',
    white: 'text-white',
  }[color];

  return (
    <div className={`inline-block animate-spin rounded-full border-2 border-solid border-current border-r-transparent ${sizeClass} ${colorClass}`} role="status">
      <span className="sr-only">Loading...</span>
    </div>
  );
};

export default Spinner;
