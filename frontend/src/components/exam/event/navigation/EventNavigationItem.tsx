import React from 'react';
import { Link } from 'react-router-dom';

interface EventNavigationItemProps {
  icon: string;
  isActive: boolean;
  route: string;
  text: string;
}

const EventNavigationItem = (props: EventNavigationItemProps) => {
  const { icon, isActive, route, text } = props;
  return (
    <Link to={route}>
      <div className={`nav-item ${isActive ? 'active' : ''}`}>
          <i className={`fas ${icon}`} /> {text}
        </div>
    </Link>
  );
};

export default EventNavigationItem;
