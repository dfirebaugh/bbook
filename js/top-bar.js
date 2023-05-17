const template = document.createElement('template');
template.innerHTML = `
<span class="mt-auto flex justify-between">
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
</span>`;

class TopBar extends HTMLElement {
    constructor() {
        super();
        const shadowRoot = this.attachShadow({ mode: "open" });
        shadowRoot.appendChild(template.content.cloneNode(true));
        // shadowRoot.appendChild(document.getElementById("style-template"));

        shadowRoot.querySelector("#site-title").href = this.getAttribute("url");
        shadowRoot.querySelector("#site-title").innerText = this.getAttribute("title");
    }
}
window.customElements.define('bf-topbar', TopBar);
