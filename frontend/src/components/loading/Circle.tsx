import React, { SVGProps, CSSProperties } from 'react';

export interface LoadingCircleProps extends SVGProps<SVGAElement> {
  color?: string;
  className?: string;
  style?: CSSProperties;
};

const LoadingCircle = (props: LoadingCircleProps) => {
  const { color, className, style } = props;
  return (
    <svg
      height="100%"
      viewBox="0 0 100 100"
      preserveAspectRatio="xMidYMid"
      className={className}
      style={style}
    >
      <circle cx="50" cy="50" fill="none" stroke={color} strokeWidth="10" r="35" strokeDasharray="164.93361431346415 56.97787143782138" transform="rotate(187.974 50 50)">
        <animateTransform attributeName="transform" type="rotate" repeatCount="indefinite" dur="1s" values="0 50 50;360 50 50" keyTimes="0;1"></animateTransform>
      </circle>
    </svg>
  );
}

LoadingCircle.defaultProps = {
  color: '#ffffff',
  className: '',
  style: {},
};

export default LoadingCircle;
