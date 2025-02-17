import React from 'react';
import './SearchResults.css';

const SearchResults = ({ quote, loading, error }) => {

	console.log(quote, loading, error)
  if (loading) {
    return <div className="search-results-container loading">Searching...</div>;
  }

  if (error) {
    return <div className="search-results-container error">{error}</div>;
  }

  if (!quote) {
    return null; // Don't show anything if there's no quote
  }

  return (
    <div className="search-results-container">
        <p className="search-results-text">"{quote}"</p>
    </div>
  );
};

export default SearchResults;
