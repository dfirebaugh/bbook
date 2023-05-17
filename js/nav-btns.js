
const template = document.createElement('template');
template.innerHTML = `
<span class="nav-wide-wrapper" aria-label="Page navigation">

    <a rel="prev" class="nav-chapters previous" title="Previous chapter" aria-label="Previous chapter"
        aria-keyshortcuts="Left">
        <i class="fa fa-angle-left"></i>
    </a>
    <a rel="next" class="nav-chapters next" title="Next chapter" aria-label="Next chapter"
        aria-keyshortcuts="Right">
        <i class="fa fa-angle-right"></i>
    </a>
</span>`;

class NavBtns extends HTMLElement {
    constructor() {
        super();
        const shadowRoot = this.attachShadow({ mode: "open" });
        shadowRoot.appendChild(template.content.cloneNode(true));

        shadowRoot.querySelector(".previous").href = this.getAttribute("prev");
        shadowRoot.querySelector(".next").href = this.getAttribute("next");
    }
}
window.customElements.define('bf-navbtns', NavBtns);
