import React from 'react';
import { Card } from 'react-hephaestus';
import { Link } from 'react-router-dom';

interface CardMenuProps {
  to: string;
  iconName: string;
  text: string;
  className: string;
}

const CardMenu = ({to, iconName, text, className}: CardMenuProps) => (
  <Link to={to}>
    <Card cardType="outlined" className={className}>
      <div><i className={`fas ${iconName}`} /> {text}</div>
      <i className="fas fa-chevron-right" />
    </Card>
  </Link>
);

export default CardMenu;
