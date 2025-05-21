  // Handle clicks on timestamp links
  document.addEventListener('click', function(e) {
    if (e.target && e.target.classList.contains('timestamp-link')) {
      e.preventDefault();
      var timestamp = e.target.getAttribute('data-time');
      
      // Create a custom event with the timestamp
      var event = new CustomEvent('commentoTimestamp', {
        detail: {
          timestamp: timestamp
        }
      });
      
      // Dispatch the event for external handling
      document.dispatchEvent(event);
      
      // Optionally, you can also implement a default behavior here
      console.log('Timestamp clicked: ' + timestamp);
    }
  }, false);
