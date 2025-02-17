import { useState, useEffect } from "react";
import PageLoader from './components/PageLoader';
import './App.css';
import SearchBar from './components/SearchBar';
import SearchResults from './components/SearchResults';
import Headline from './components/Headline';
import Header from './components/Header';
import ParticleBackground from './components/ParticleBackground';

const App = () => {
  const [loading, setLoading] = useState(true);
  const [fadeIn, setFadeIn] = useState(false);
  const [searchResult, setSearchResult] = useState(null); // Store structured response
  const [searchLoading, setSearchLoading] = useState(false);
  const [searchError, setSearchError] = useState(null);
  const [hasSearched, setHasSearched] = useState(false);

  useEffect(() => {
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
      const apiUrl = `${window.location.protocol}//${window.location.hostname}:8080/search`;
      console.log('Sending search request:', query);

      const response = await fetch(apiUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
					query
				})
      });

      if (!response.ok) {
        throw new Error(`Search failed: ${response.statusText}`);
      }

      // Get the response as JSON
      const data = await response.json();

      // Check if the response contains the 'response' key with a stringified JSON
      if (data && data.response) {
				console.log(data.response)
        // Parse the 'response' field, which contains the JSON string
        const parsedResponse = JSON.parse(data.response);

        console.log('Parsed search result:', parsedResponse);

        // Set the parsed result to state
        if (parsedResponse.quote && parsedResponse.interpretation && parsedResponse.advice) {
          setSearchResult(parsedResponse); // Set the parsed result object
        } else {
          console.warn('Unexpected response format:', parsedResponse);
          setSearchResult(null);
        }
      } else {
        console.warn('Response does not contain the expected "response" key');
        setSearchResult(null);
      }
    } catch (error) {
      console.error("Search error:", error);
      setSearchError(`An error occurred while searching: ${error.message}`);
      setSearchResult(null);
    } finally {
      setSearchLoading(false);
    }
  };

  return (
    <div className="app-container">
      <ParticleBackground />
      {loading ? <PageLoader /> : (
        <div className={`content ${fadeIn ? 'fade-in' : ''}`}>
          <Header />
          <Headline />
          <SearchBar onSearch={handleSearch} />
          <div className="search-results-wrapper">
            <SearchResults
              result={searchResult}  // Pass the full parsed result object
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

