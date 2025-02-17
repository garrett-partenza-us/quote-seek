import React, { useEffect, useState } from 'react';
import { helix } from 'ldrs';  // Import the helix loader from ldrs
import './PageLoader.css';  // Import custom styles for the loader

// Register the helix loader component
helix.register();

const PageLoader = ({ isLoading }) => {
  const [fadeOut, setFadeOut] = useState(false);

  useEffect(() => {
    if (!isLoading) {
      setFadeOut(true);  // Start fade-out effect when isLoading becomes false
      setTimeout(() => {
        setFadeOut(false); // Reset fade-out state after animation
      }, 1000); // Allow 1 second for the fade-out effect
    }
  }, [isLoading]);

  return (
    <div className={`page-loader-container ${fadeOut ? 'fade-out' : ''}`}>
      <l-helix
        size="45"
        speed="2.5"
        color="#FAF9F6"
      ></l-helix>
    </div>
  );
};

export default PageLoader;
