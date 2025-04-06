import React from 'react';
import './Headline.css';  // Make sure to add the accompanying CSS

const Headline = () => {
  return (
    <headline className="headline">
      <div className="headline-content">
        <h1 className="title">Meditations Quote Finder</h1>
        <p className="slogan">An AI search tool for exploring quotes in Marcus Aurelius's journal of Meditations.</p>
      </div>
    </headline>
  );
};

export default Headline;

