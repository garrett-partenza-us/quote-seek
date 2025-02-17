import { useEffect, useMemo, useState, useRef } from "react";
import Particles, { initParticlesEngine } from "@tsparticles/react";
import { loadSlim } from "@tsparticles/slim";
import PageLoader from './PageLoader';  // Import the Loader component
import './App.css';
import SearchBar from './SearchBar';  // Import the SearchBar component

const App = () => {
  const [loading, setLoading] = useState(true);  // Manage loading state
  const particlesLoadedRef = useRef(false);  // Ref to track if particles have been loaded

  useEffect(() => {
    // Initialize particles only once
    if (!particlesLoadedRef.current) {
      initParticlesEngine(async (engine) => {
        await loadSlim(engine);  // Initialize particles
      }).then(() => {
        particlesLoadedRef.current = true;  // Set particles as loaded
      });
    }

    // Simulate loading completion after 3 seconds
    setTimeout(() => setLoading(false), 3000);
  }, []);

  const particlesLoaded = (container) => {
    console.log(container);
  };

  const options = useMemo(
    () => ({
      background: {
        color: {
          value: "#000000", // Background color
        },
      },
      fpsLimit: 120,
      particles: {
        color: {
          value: "#FFFFF2", // Particle color
        },
        links: {
          color: "#FFFFF2", // Link color
          distance: 200,
          enable: true,
          opacity: 0.1,
          width: 2,
        },
        move: {
          direction: "none",
          enable: true,
          speed: 1,
        },
        number: {
          value: 50, // Number of particles
        },
        opacity: {
          value: 0.1,
        },
        shape: {
          type: "circle",
        },
        size: {
          value: { min: 1, max: 5 },
        },
      },
      detectRetina: true,
    }),
    [],
  );

  return (
    <div className="App">
      {/* Particles Background (this will persist even after loading is complete) */}
      <Particles
        id="tsparticles"
        particlesLoaded={particlesLoaded}
        options={options}
      />

      {/* Conditionally render the loader */}
      {loading && <PageLoader />}

      {/* Add the custom search bar after loading */}
      {!loading && (
        <div className="App-content">
          <SearchBar />  {/* Display the search bar */}
        </div>
      )}
    </div>
  );
};

export default App;
