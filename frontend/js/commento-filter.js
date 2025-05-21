// Username filtering functionality
  global.filterByUsername = function(username) {
    // If username is empty, clear filter
    if (!username || username.trim() === "") {
      var cards = document.querySelectorAll(".commento-card");
      cards.forEach(function(card) {
        card.style.display = "";
      });
      return;
    }
    
    // Make case-insensitive comparison
    username = username.toLowerCase();
    
    var cards = document.querySelectorAll(".commento-card");
    cards.forEach(function(card) {
      var commentHex = card.id.replace(ID_CARD, "");
      var commenter = commenters[comments.find(c => c.commentHex === commentHex).commenterHex];
      
      if (commenter && commenter.name.toLowerCase().includes(username)) {
        card.style.display = "";
      } else {
        card.style.display = "none";
      }
    });
  }
  
  function usernameFilterCreate() {
    var filterContainer = create("div");
    var filterInput = create("input");
    var filterLabel = create("label");
    
    filterContainer.id = ID_FILTER + "container";
    filterInput.id = ID_FILTER + "input";
    
    classAdd(filterContainer, "filter-container");
    classAdd(filterInput, "filter-input");
    classAdd(filterLabel, "filter-label");
    
    filterLabel.innerText = "Filter by username:";
    attrSet(filterInput, "type", "text");
    attrSet(filterInput, "placeholder", "Enter username");
    
    filterInput.oninput = function() {
      global.filterByUsername(this.value);
    };
    
    append(filterContainer, filterLabel);
    append(filterContainer, filterInput);
    
    return filterContainer;
  }
