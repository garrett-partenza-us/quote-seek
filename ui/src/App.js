import { useState, useEffect } from "react";
import PageLoader from './PageLoader';  // Import the Loader component
import './App.css';
import SearchBar from './SearchBar';  // Import the SearchBar component
import Headline from './Headline';
import Header from './Header'
import ParticleBackground from './ParticleBackground';  // Import the ParticleBackground component

const App = () => {
  const [loading, setLoading] = useState(true);  // Manage loading state
  const [fadeIn, setFadeIn] = useState(false);  // Manage fade-in state

  useEffect(() => {
    // Simulate loading completion after 3 seconds
    const timer = setTimeout(() => {
      setLoading(false);  // Set loading to false after 3 seconds
      setFadeIn(true);    // Trigger fade-in after loading is done
    }, 3000);

    return () => clearTimeout(timer);  // Cleanup on unmount
  }, []);

  return (
    <div className="App">
      {/* Particles Background (this will persist even after loading is complete) */}
      <ParticleBackground />

      {/* Conditionally render the loader */}
      {loading && <PageLoader />}

      {/* Add the custom search bar and header after loading */}
      {!loading && (
        <div className={`App-content ${fadeIn ? 'fade-in' : ''}`}>
					<Header />
          <Headline />
          <SearchBar />  {/* Display the search bar */}
        </div>
      )}
    </div>
  );
};

export default App;
