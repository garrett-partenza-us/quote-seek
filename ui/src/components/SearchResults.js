import React from 'react';
import './SearchResults.css';

const SearchResults = ({ result, loading, error }) => {
  if (loading) {
    return <div className="search-results-container loading">Searching...</div>;
  }

  if (error) {
    return <div className="search-results-container error">{error}</div>;
  }

  if (!result) {
    return null;
  }

  const { quote, interpretation, advice } = result; // Destructure the response object

  return (
    <div className="search-results-container">
      {quote && <p className="search-results-text section-1">{quote}</p>}
      {interpretation && <p className="search-results-text section-2">{interpretation}</p>}
      {advice && <p className="search-results-text section-3">{advice}</p>}
    </div>
  );
};

export default SearchResults;

