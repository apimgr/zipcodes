// Zipcode Lookup Frontend

class ZipcodeApp {
  constructor() {
    this.searchInput = document.getElementById('search-input');
    this.searchBtn = document.getElementById('search-btn');
    this.resultsDiv = document.getElementById('results');
    this.resultsListDiv = document.getElementById('results-list');
    this.resultCountSpan = document.getElementById('result-count');
    this.loadingDiv = document.getElementById('loading');
    this.autocompleteDiv = document.getElementById('autocomplete-results');
    this.themeToggle = document.getElementById('theme-toggle');

    this.debounceTimer = null;

    this.init();
  }

  init() {
    // Event listeners
    this.searchBtn.addEventListener('click', () => this.search());
    this.searchInput.addEventListener('keypress', (e) => {
      if (e.key === 'Enter') this.search();
    });

    // Autocomplete
    this.searchInput.addEventListener('input', (e) => {
      clearTimeout(this.debounceTimer);
      this.debounceTimer = setTimeout(() => {
        this.autocomplete(e.target.value);
      }, 300);
    });

    // Close autocomplete on click outside
    document.addEventListener('click', (e) => {
      if (!this.autocompleteDiv.contains(e.target) && e.target !== this.searchInput) {
        this.autocompleteDiv.classList.remove('show');
      }
    });

    // Example links
    document.querySelectorAll('.example').forEach(link => {
      link.addEventListener('click', (e) => {
        e.preventDefault();
        const query = e.target.dataset.query;
        this.searchInput.value = query;
        this.search();
      });
    });

    // Theme toggle
    this.themeToggle.addEventListener('click', () => this.toggleTheme());

    // Load stats
    this.loadStats();

    // Initialize theme
    this.initTheme();
  }

  async search() {
    const query = this.searchInput.value.trim();
    if (!query) return;

    this.showLoading();
    this.hideResults();

    try {
      const response = await fetch(`/api/v1/zipcode/search?q=${encodeURIComponent(query)}`);
      const data = await response.json();

      this.hideLoading();

      if (data.success) {
        this.displayResults(data.data, data.count);
      } else {
        this.showError(data.error || 'Search failed');
      }
    } catch (error) {
      this.hideLoading();
      this.showError('Network error: ' + error.message);
    }
  }

  async autocomplete(query) {
    if (!query || query.length < 2) {
      this.autocompleteDiv.innerHTML = '';
      this.autocompleteDiv.classList.remove('show');
      return;
    }

    try {
      const response = await fetch(`/api/v1/zipcode/autocomplete?q=${encodeURIComponent(query)}&limit=10`);
      const data = await response.json();

      if (data.success && data.suggestions.length > 0) {
        this.displayAutocomplete(data.suggestions);
      } else {
        this.autocompleteDiv.innerHTML = '';
        this.autocompleteDiv.classList.remove('show');
      }
    } catch (error) {
      console.error('Autocomplete error:', error);
    }
  }

  displayAutocomplete(suggestions) {
    this.autocompleteDiv.innerHTML = suggestions.map(suggestion => `
      <div class="autocomplete-item" data-value="${suggestion}">
        ${suggestion}
      </div>
    `).join('');

    // Add click listeners
    this.autocompleteDiv.querySelectorAll('.autocomplete-item').forEach(item => {
      item.addEventListener('click', () => {
        this.searchInput.value = item.dataset.value;
        this.autocompleteDiv.classList.remove('show');
        this.search();
      });
    });

    this.autocompleteDiv.classList.add('show');
  }

  displayResults(data, count) {
    // Handle single result
    if (!Array.isArray(data)) {
      data = [data];
      count = 1;
    }

    if (!data || data.length === 0) {
      this.resultsListDiv.innerHTML = '<p>No results found</p>';
      this.resultCountSpan.textContent = '';
      this.resultsDiv.style.display = 'block';
      return;
    }

    this.resultCountSpan.textContent = count ? `(${count})` : '';

    this.resultsListDiv.innerHTML = data.map(zipcode => `
      <div class="result-card">
        <div class="result-zip">${zipcode.zip_code}</div>
        <div class="result-city">${zipcode.city}</div>
        <div class="result-state">${zipcode.state}${zipcode.county ? ', ' + zipcode.county : ''}</div>
        ${zipcode.latitude && zipcode.longitude ? `
          <div class="result-coords">${zipcode.latitude}, ${zipcode.longitude}</div>
        ` : ''}
      </div>
    `).join('');

    this.resultsDiv.style.display = 'block';
  }

  async loadStats() {
    try {
      const response = await fetch('/api/v1/zipcode/stats');
      const data = await response.json();

      if (data.success) {
        document.getElementById('total-zipcodes').textContent =
          this.formatNumber(data.data.total_zipcodes);
        document.getElementById('total-states').textContent =
          this.formatNumber(data.data.total_states);
        document.getElementById('total-cities').textContent =
          this.formatNumber(data.data.total_cities);
      }
    } catch (error) {
      console.error('Failed to load stats:', error);
    }
  }

  showLoading() {
    this.loadingDiv.style.display = 'block';
  }

  hideLoading() {
    this.loadingDiv.style.display = 'none';
  }

  hideResults() {
    this.resultsDiv.style.display = 'none';
  }

  showError(message) {
    this.resultsListDiv.innerHTML = `
      <div class="error-message" style="padding: 2rem; text-align: center; color: #ef4444;">
        ${message}
      </div>
    `;
    this.resultsDiv.style.display = 'block';
  }

  formatNumber(num) {
    return new Intl.NumberFormat().format(num);
  }

  initTheme() {
    const savedTheme = localStorage.getItem('theme') || 'dark';
    document.body.setAttribute('data-theme', savedTheme);
    this.updateThemeIcon(savedTheme);
  }

  toggleTheme() {
    const currentTheme = document.body.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    document.body.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
    this.updateThemeIcon(newTheme);
  }

  updateThemeIcon(theme) {
    this.themeToggle.textContent = theme === 'dark' ? '🌙' : '☀️';
  }
}

// Initialize app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  new ZipcodeApp();
});
