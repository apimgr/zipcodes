const themeSelector = document.getElementById('theme-choice');

if (themeSelector) {
  const savedTheme = localStorage.getItem('zipcodes-theme');
  const theme = ['light', 'dark'].includes(savedTheme) ? savedTheme : 'auto';
  document.documentElement.dataset.theme = theme;
  themeSelector.value = theme;
  themeSelector.addEventListener('change', () => {
    const selectedTheme = themeSelector.value;
    document.documentElement.dataset.theme = selectedTheme;
    if (selectedTheme === 'auto') {
      localStorage.removeItem('zipcodes-theme');
    } else {
      localStorage.setItem('zipcodes-theme', selectedTheme);
    }
  });
}

if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => navigator.serviceWorker.register('/sw.js'));
}
