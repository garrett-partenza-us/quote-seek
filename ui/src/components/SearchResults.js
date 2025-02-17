import React, { useState, useEffect } from 'react';
import { ChevronDown, ChevronUp } from 'lucide-react';
import './SearchResults.css';

const SearchResults = ({ result, loading, error, hasSearched }) => {
  // Set quote to be open by default
  const [openSections, setOpenSections] = useState({
    quote: true,
    interpretation: false,
    advice: false
  });
  
  // State for fade-in animation
  const [visible, setVisible] = useState(false);
  
  // Trigger fade-in when component mounts with results
  useEffect(() => {
    if (result && hasSearched) {
      const timer = setTimeout(() => {
        setVisible(true);
      }, 50); // Small delay to ensure DOM is ready
      return () => clearTimeout(timer);
    }
  }, [result, hasSearched]);

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
      <div className={`search-results-container ${visible ? 'visible' : 'hidden'}`}>
        {/* Quote section */}
        <div className="accordion-item">
          <div 
            className="accordion-header"
            onClick={() => toggleSection('quote')}
          >
            <h3>Quote</h3>
            {openSections.quote ? <ChevronUp /> : <ChevronDown />}
          </div>
          <div className={`accordion-content-wrapper ${openSections.quote ? 'open' : 'closed'}`}>
            <div className="accordion-content">
              <p>{quote}</p>
            </div>
          </div>
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
          <div className={`accordion-content-wrapper ${openSections.interpretation ? 'open' : 'closed'}`}>
            <div className="accordion-content">
              <p>{interpretation}</p>
            </div>
          </div>
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
          <div className={`accordion-content-wrapper ${openSections.advice ? 'open' : 'closed'}`}>
            <div className="accordion-content">
              <p>{advice}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SearchResults;
