// components/FeedbackForm.js
import React, { useState } from 'react';
import './FeedbackForm.css';

const FeedbackForm = ({ isOpen, onClose }) => {
  const [formData, setFormData] = useState({
    name: '',
    consistent: '',
    helpful: ''
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prevState => ({
      ...prevState,
      [name]: value
    }));
  };

  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitting(true);
    setSubmitError(null);
    
    try {
      const apiUrl = `${window.location.protocol}//${window.location.hostname}/api/feedback`;
      
      const response = await fetch(apiUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData)
      });
      
      if (!response.ok) {
        throw new Error(`Submission failed: ${response.statusText}`);
      }
      
      // Optional: get response data if your API returns anything
      const data = await response.json().catch(() => ({}));
      console.log('Feedback submitted successfully:', data);
      
      // Clear form and close
      setFormData({ name: '', consistent: '', helpful: '' });
      onClose();
    } catch (error) {
      console.error("Feedback submission error:", error);
      setSubmitError(`An error occurred while submitting feedback: ${error.message}`);
    } finally {
      setSubmitting(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="feedback-overlay">
      <div className="feedback-modal">
        <div className="feedback-header">
          <h2>Provide Feedback</h2>
          <button className="close-button" onClick={onClose}>×</button>
        </div>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="name">Name</label>
            <input
              type="text"
              id="name"
              name="name"
              value={formData.name}
              onChange={handleChange}
              required
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="consistent">Are the results consistent?</label>
            <textarea
              id="consistent"
              name="consistent"
              value={formData.consistent}
              onChange={handleChange}
              rows={4}
              required
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="helpful">Did it help you?</label>
            <textarea
              id="helpful"
              name="helpful"
              value={formData.helpful}
              onChange={handleChange}
              rows={4}
              required
            />
          </div>
          
          {submitError && (
            <div className="error-message">
              {submitError}
            </div>
          )}
          
          <div className="form-actions">
            <button type="button" onClick={onClose} className="cancel-button" disabled={submitting}>
              Cancel
            </button>
            <button type="submit" className="submit-button" disabled={submitting}>
              {submitting ? 'Submitting...' : 'Submit Feedback'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default FeedbackForm;
