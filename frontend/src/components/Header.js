import React from 'react';
import { FaGithub } from 'react-icons/fa';  // GitHub icon from react-icons
import './Header.css'; // Importing the CSS for the header styling

const Header = () => {
  return (
    <header className="header">
      <div className="header-logo-container">
        {/* GitHub logo */}
        <a href="https://github.com/garrett-partenza-us/stoic-rag" target="_blank" rel="noopener noreferrer">
          <FaGithub className="github-logo" />
        </a>
      </div>
    </header>
  );
};

export default Header;

