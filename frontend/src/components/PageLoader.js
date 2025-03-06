import React, { useEffect, useState } from 'react';
import { helix } from 'ldrs';  // Import the helix loader from ldrs
import './PageLoader.css';  // Import custom styles for the loader

// Register the helix loader component
helix.register();

const PageLoader = ({ isLoading }) => {
  const [fadeOut, setFadeOut] = useState(false);

  useEffect(() => {
    let timeoutId;
    
    if (!isLoading) {
      // Start the fade-out process after 2 seconds
      timeoutId = setTimeout(() => {
        setFadeOut(true);  // Start fading out after 2 seconds
      }, 2000);  // Delay for 2 seconds before triggering fade-out

      // Clear the timeout if the component is unmounted or if isLoading changes
      return () => clearTimeout(timeoutId);
    }

    // Reset fade-out if isLoading is true
    setFadeOut(false);
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

