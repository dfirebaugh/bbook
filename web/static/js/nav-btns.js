const template = document.createElement("template");
template.innerHTML = `
  <link 
    rel="stylesheet" 
    href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css" 
    integrity="sha512-iecdLmaskl7CVkqkXNQ/ZH/XLlvWZOJyj7Yy7tcenmpD1ypASozpmT/E0iPtmFIB46ZmdtAc9eNBvH0H/ZpiBw==" 
    crossorigin="anonymous" 
    referrerpolicy="no-referrer" 
  />

  <style>
    :host {
      display: inline-block;
    }

    .nav-wide-wrapper {
      position: fixed;
      width: 4rem;
      height: 4rem;
      display: flex;
      justify-content: center;
      align-items: center;
      /*right: 0;*/
      /*left: 0;*/
    }

    a {
      display: flex;
      justify-content: center;
      align-items: center;
      width: 100%;
      height: 100%;
      text-decoration: none;
      color: var(--primary-color);
    }

    .nav-wide-wrapper:hover {
      background-color: var(--secondary-background-color);
      cursor: pointer;
    }
  </style>

  <span class="nav-wide-wrapper">
    <a rel="prev" class="nav-chapters previous" title="Previous chapter" aria-label="Previous chapter" aria-keyshortcuts="Left">
      <i class="fa fa-angle-left"></i>
    </a>
    <a rel="next" class="nav-chapters next" title="Next chapter" aria-label="Next chapter" aria-keyshortcuts="Right">
      <i class="fa fa-angle-right"></i>
    </a>
  </span>
`;

class NavBtns extends HTMLElement {
  static get observedAttributes() {
    return ["nav-type", "href"];
  }

  constructor() {
    super();
    const shadowRoot = this.attachShadow({ mode: "open" });
    shadowRoot.appendChild(template.content.cloneNode(true));
  }

  connectedCallback() {
    this.updateLinks();
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (oldValue !== newValue) {
      this.updateLinks();
    }
  }

  updateLinks() {
    const navType = this.getAttribute("nav-type");
    const href = this.getAttribute("href");
    const prevLink = this.shadowRoot.querySelector(".previous");
    const nextLink = this.shadowRoot.querySelector(".next");

    if (navType === "prev") {
      nextLink?.remove();
      if (prevLink) prevLink.href = href || "#";
    } else if (navType === "next") {
      prevLink?.remove();
      const container = this.shadowRoot.querySelector(".nav-wide-wrapper");
      container.style.right = "0";
      if (nextLink) nextLink.href = href || "#";
    } else {
      if (prevLink) prevLink.href = "#";
      if (nextLink) nextLink.href = "#";
    }
  }
}

window.customElements.define("bb-navbtn", NavBtns);
