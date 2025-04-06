// components/Header.js
import React, { useState } from 'react';
import { FaGithub, FaComments } from 'react-icons/fa';  // Added FaComments icon for feedback
import './Header.css';
import FeedbackForm from './FeedbackForm';

const Header = () => {
  const [isFeedbackOpen, setIsFeedbackOpen] = useState(false);

  const openFeedback = () => {
    setIsFeedbackOpen(true);
  };

  const closeFeedback = () => {
    setIsFeedbackOpen(false);
  };

  return (
    <header className="header">
      <div className="header-content">
        <button onClick={openFeedback} className="feedback-button">
          <FaComments className="feedback-icon" />
          <span>Feedback</span>
        </button>
        
        <div className="header-logo-container">
          {/* GitHub logo */}
          <a href="https://github.com/garrett-partenza-us/stoic-rag" target="_blank" rel="noopener noreferrer">
            <FaGithub className="github-logo" />
          </a>
        </div>
      </div>
      
      <FeedbackForm isOpen={isFeedbackOpen} onClose={closeFeedback} />
    </header>
  );
};

export default Header;
