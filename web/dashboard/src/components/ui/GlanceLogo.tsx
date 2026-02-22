import React from 'react';

interface GlanceLogoProps {
  size?: number;
  className?: string;
}

export const GlanceLogo: React.FC<GlanceLogoProps> = ({ size = 40, className = '' }) => {
  return (
    <svg 
      xmlns="http://www.w3.org/2000/svg" 
      width={size} 
      height={size} 
      viewBox="0 0 120 120" 
      className={className}
    >
      <defs>
        <linearGradient id="glance-logo-grad" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" style={{ stopColor: '#2563eb', stopOpacity: 1 }} />
          <stop offset="100%" style={{ stopColor: '#0f172a', stopOpacity: 1 }} />
        </linearGradient>
      </defs>
      <circle cx="60" cy="60" r="55" fill="url(#glance-logo-grad)"/>
      <circle cx="60" cy="60" r="35" fill="none" stroke="white" strokeWidth="4"/>
      <circle cx="60" cy="60" r="20" fill="none" stroke="white" strokeWidth="4"/>
      <circle cx="60" cy="60" r="5" fill="white"/>
      <path d="M 60 25 L 60 5" stroke="white" strokeWidth="3" strokeLinecap="round"/>
      <path d="M 60 95 L 60 115" stroke="white" strokeWidth="3" strokeLinecap="round"/>
      <path d="M 25 60 L 5 60" stroke="white" strokeWidth="3" strokeLinecap="round"/>
      <path d="M 95 60 L 115 60" stroke="white" strokeWidth="3" strokeLinecap="round"/>
    </svg>
  );
};
