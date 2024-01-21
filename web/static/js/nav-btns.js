
const template = document.createElement('template');
template.innerHTML = `
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css"
integrity="sha512-iecdLmaskl7CVkqkXNQ/ZH/XLlvWZOJyj7Yy7tcenmpD1ypASozpmT/E0iPtmFIB46ZmdtAc9eNBvH0H/ZpiBw=="
crossorigin="anonymous" referrerpolicy="no-referrer" />

<style>
a {
    color: var(--primary-color);
}
</style>

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
window.customElements.define('bb-navbtns', NavBtns);
