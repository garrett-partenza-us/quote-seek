// App.js
import { useState, useEffect } from "react";
import PageLoader from './components/PageLoader';
import './App.css';
import SearchBar from './components/SearchBar';
import SearchResults from './components/SearchResults';
import Headline from './components/Headline';
import Header from './components/Header';
import ParticleBackground from './components/ParticleBackground';

const App = () => {
  // State for loading
  const [loading, setLoading] = useState(true);
  const [fadeIn, setFadeIn] = useState(false);
  
  // State for search
  const [searchResult, setSearchResult] = useState('');
  const [searchLoading, setSearchLoading] = useState(false);
  const [searchError, setSearchError] = useState(null);
  const [hasSearched, setHasSearched] = useState(false);

  useEffect(() => {
    // Simulate loading completion after 3 seconds
    const timer = setTimeout(() => {
      setLoading(false);
      setFadeIn(true);
    }, 3000);

    return () => clearTimeout(timer);
  }, []);

  const handleSearch = async (query) => {
    setSearchLoading(true);
    setSearchError(null);
    setHasSearched(true);
    
    try {
      // Using the current hostname with port 8080 and JSON payload format you specified
      const apiUrl = `${window.location.protocol}//${window.location.hostname}:8080/search`;
      
      console.log('Sending search request:', query);
      
      const response = await fetch(apiUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ query })
      });
      
      if (!response.ok) {
        throw new Error(`Search failed: ${response.statusText}`);
      }
      
      const data = await response.json();
      console.log('Received search result:', data);
      
      // Make sure we're handling the response correctly
      if (data && typeof data.response === 'string') {
        setSearchResult(data.response);
      } else {
        console.warn('Unexpected response format:', data);
        setSearchResult(JSON.stringify(data));
      }
    } catch (error) {
      console.error("Search error:", error);
      setSearchError(`An error occurred while searching: ${error.message}`);
      setSearchResult('');
    } finally {
      setSearchLoading(false);
    }
  };

  return (
    <div className="app-container">
      {/* Particles Background */}
      <ParticleBackground />

      {/* Conditionally render the loader */}
      {loading ? <PageLoader /> : (
        <div className={`content ${fadeIn ? 'fade-in' : ''}`}>
          {/* Render Header, Headline, and SearchBar after loading */}
          <Header />
          <Headline />
          <SearchBar onSearch={handleSearch} />
          
          {/* Render SearchResults - make sure it's visible */}
          <div className="search-results-wrapper">
            <SearchResults 
              quote={searchResult} 
              loading={searchLoading} 
              error={searchError}
              hasSearched={hasSearched}
            />
          </div>
        </div>
      )}
    </div>
  );
};

export default App;
