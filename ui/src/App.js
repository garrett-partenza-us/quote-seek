import { useState } from "react";
import PageLoader from './PageLoader';  // Import the Loader component
import './App.css';
import SearchBar from './SearchBar';  // Import the SearchBar component
import Header from './Header'
import ParticleBackground from './ParticleBackground';  // Import the ParticleBackground component

const App = () => {
  const [loading, setLoading] = useState(true);  // Manage loading state

  // Simulate loading completion after 3 seconds
  setTimeout(() => setLoading(false), 3000);

  return (
    <div className="App">
      {/* Particles Background (this will persist even after loading is complete) */}
      <ParticleBackground />

      {/* Conditionally render the loader */}
      {loading && <PageLoader />}

      {/* Add the custom search bar after loading */}
      {!loading && (
        <div className="App-content">
					<Header />
					<SearchBar />  {/* Display the search bar */}
        </div>
      )}
    </div>
  );
};

export default App;

