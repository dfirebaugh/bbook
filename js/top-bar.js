const template = document.createElement('template');
template.innerHTML = `
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css"
integrity="sha512-iecdLmaskl7CVkqkXNQ/ZH/XLlvWZOJyj7Yy7tcenmpD1ypASozpmT/E0iPtmFIB46ZmdtAc9eNBvH0H/ZpiBw=="
crossorigin="anonymous" referrerpolicy="no-referrer" />

<style>
    .top-bar {
        display: grid;
        align-items: center;
        justify-content: space-between;
        grid-template-columns: 1fr 2fr 1fr;
    }

    a {
        color: var(--primary-color);
    }

    .icon-button {
        border: none;
        background: none;
        padding: 0;
        margin: 0;
        outline: none;
        cursor: pointer;
        display: inline-flex;
        align-items: center;
        justify-content: center;
    }

    .icon-button i.fa {
        font-size: 1.2em;
        color: #333;
    }

    .icon-button:hover, .icon-button:active {
        background-color: rgba(0, 0, 0, 0.1);
    }

    #site-title:hover {
        cursor: pointer;
        text-decoration: underline;
    }
</style>
    <span class="top-bar">
        <left-buttons>
            <button id="sidebar-toggle" onclick="hamburgerClick()" class="icon-button" type="button" title="Toggle Table of Contents"
                aria-label="Toggle Table of Contents" aria-controls="sidebar" aria-expanded="true">
                <i class="fa fa-bars"></i>
            </button>
        </left-buttons>

        <h1 class="menu-title  md-6">
            <a id="site-title">
            </a>
        </h1>

        <right-buttons>
        </right-buttons>
    </span>
`;

class TopBar extends HTMLElement {
    static get observedAttributes() {
        return ['title', 'url', 'theme'];
    }

    constructor() {
        super();
        this.attachShadow({ mode: 'open' });
        this.shadowRoot.appendChild(template.content.cloneNode(true));
    }

    connectedCallback() {
        this.updateTitleAndUrl();
        this.updateTheme();
    }

    attributeChangedCallback(name, oldValue, newValue) {
        if (oldValue !== newValue) {
            if (name === 'title' || name === 'url') {
                this.updateTitleAndUrl();
            } else if (name === 'theme') {
                this.updateTheme();
            }
        }
    }

    updateTitleAndUrl() {
        const titleElement = this.shadowRoot.querySelector("#site-title");
        const title = this.getAttribute('title');
        const url = this.getAttribute('url');

        if (title) titleElement.innerText = title;
        if (url) titleElement.href = url;
    }

    updateTheme() {
        const theme = this.getAttribute('theme');
        if (theme) this.shadowRoot.host.className = theme;
    }
}

window.customElements.define('bb-topbar', TopBar);
