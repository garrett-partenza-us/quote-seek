import React, { useState } from 'react';
import { ChevronDown, ChevronUp } from 'lucide-react';
import './SearchResults.css';

const SearchResults = ({ result, loading, error, hasSearched }) => {
  // Keep the quote section open by default
  const [openSections, setOpenSections] = useState({
    quote: true,
    interpretation: false,
    advice: false
  });

  const toggleSection = (section) => {
    setOpenSections(prev => ({
      ...prev,
      [section]: !prev[section]
    }));
  };

  // Handle loading state
  if (loading) {
    return <div className="search-results-container loading">Searching...</div>;
  }

  // Handle error state
  if (error) {
    return <div className="search-results-container error">{error}</div>;
  }

  // Handle no result or not searched state
  if (!result || !hasSearched) {
    return null;
  }

  // Destructure the result after we've checked it exists
  const { quote, interpretation, advice } = result;

  return (
    <div className="results-wrapper">
      <div className="search-results-container">
        {/* Quote section */}
        <div className="accordion-item">
          <div 
            className="accordion-header"
            onClick={() => toggleSection('quote')}
          >
            <h3>Quote</h3>
            {openSections.quote ? <ChevronUp /> : <ChevronDown />}
          </div>
          {openSections.quote && (
            <div className="accordion-content">
              <p>{quote}</p>
            </div>
          )}
        </div>

        {/* Interpretation section */}
        <div className="accordion-item">
          <div 
            className="accordion-header"
            onClick={() => toggleSection('interpretation')}
          >
            <h3>Interpretation</h3>
            {openSections.interpretation ? <ChevronUp /> : <ChevronDown />}
          </div>
          {openSections.interpretation && (
            <div className="accordion-content">
              <p>{interpretation}</p>
            </div>
          )}
        </div>

        {/* Advice section */}
        <div className="accordion-item">
          <div 
            className="accordion-header"
            onClick={() => toggleSection('advice')}
          >
            <h3>Advice</h3>
            {openSections.advice ? <ChevronUp /> : <ChevronDown />}
          </div>
          {openSections.advice && (
            <div className="accordion-content">
              <p>{advice}</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default SearchResults;
