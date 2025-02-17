import React from 'react';
import './Header.css';  // Make sure to add the accompanying CSS

const Header = () => {
  return (
    <header className="header">
      <div className="header-content">
        <h1 className="title">Meditations Quote Finder</h1>
        <p className="slogan">Explore the right words, no matter where you start.</p>
      </div>
    </header>
  );
};

export default Header;

