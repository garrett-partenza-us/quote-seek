import React from 'react';
import './SearchBar.css';  // We will create this custom CSS

const SearchBar = () => {
  return (
    <div className="search-bar-container">
      <input 
        type="text" 
        className="search-bar" 
        placeholder="Search..." 
      />
    </div>
  );
};

export default SearchBar;

