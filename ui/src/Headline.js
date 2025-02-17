import React from 'react';
import './Headline.css';  // Make sure to add the accompanying CSS

const Headline = () => {
  return (
    <headline className="headline">
      <div className="headline-content">
        <h1 className="title">Meditations Quote Finder</h1>
        <p className="slogan">Explore the right words, no matter where you start.</p>
      </div>
    </headline>
  );
};

export default Headline;

