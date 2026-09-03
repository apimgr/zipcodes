// Extra JavaScript for Zipcodes documentation
// Dracula theme enhancements and custom functionality

document.addEventListener("DOMContentLoaded", function() {
  // Set color scheme to dracula
  if (typeof __md_scope !== 'undefined') {
    __md_scope.palette = {
      scheme: 'dracula'
    };
  }

  // Add copy button functionality enhancement
  const codeBlocks = document.querySelectorAll('pre code');
  codeBlocks.forEach(block => {
    block.addEventListener('copy', function() {
      console.log('Code copied to clipboard');
    });
  });

  // Smooth scroll for anchor links
  document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
      e.preventDefault();
      const target = document.querySelector(this.getAttribute('href'));
      if (target) {
        target.scrollIntoView({
          behavior: 'smooth',
          block: 'start'
        });
      }
    });
  });

  // Add external link indicators
  const externalLinks = document.querySelectorAll('a[href^="http"]');
  externalLinks.forEach(link => {
    if (!link.hostname.includes('zipcodes') && !link.hostname.includes('readthedocs')) {
      link.setAttribute('target', '_blank');
      link.setAttribute('rel', 'noopener noreferrer');
    }
  });

  console.log('Zipcodes documentation loaded with Dracula theme');
});
