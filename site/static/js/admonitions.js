(function () {
  function normalizeType(typeStr) {
    return typeStr
      .trim()
      .split(/\s+/)
      .map(function (part) {
        return part.charAt(0).toUpperCase() + part.slice(1).toLowerCase();
      })
      .join(" ");
  }

  function transformAdmonitions() {
    var blockquotes = document.querySelectorAll(
      ".body blockquote, .blog-post-content blockquote"
    );

    blockquotes.forEach(function (blockquote) {
      var firstP = blockquote.querySelector(":scope > p:first-child");
      if (!firstP) return;

      var headerMatch = firstP.textContent.match(/^\[!([^\]]+)\]/);
      if (!headerMatch) return;

      var rawType = headerMatch[1].trim();
      var remainder = firstP.textContent.slice(headerMatch[0].length).trim();

      var admonition = document.createElement("div");
      admonition.className = "admonition " + rawType.toLowerCase();

      var titleEl = document.createElement("div");
      titleEl.className = "admonition-title";
      titleEl.textContent = normalizeType(rawType);
      admonition.appendChild(titleEl);

      var content = document.createElement("div");
      content.className = "admonition-content";

      if (remainder) {
        var p = document.createElement("p");
        p.textContent = remainder;
        content.appendChild(p);
      }

      var siblings = Array.prototype.slice.call(blockquote.children, 1);
      siblings.forEach(function (node) {
        content.appendChild(node);
      });

      admonition.appendChild(content);
      blockquote.replaceWith(admonition);
    });
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", transformAdmonitions);
  } else {
    transformAdmonitions();
  }
})();
