const template = document.createElement('template');
template.innerHTML = `
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css"
integrity="sha512-iecdLmaskl7CVkqkXNQ/ZH/XLlvWZOJyj7Yy7tcenmpD1ypASozpmT/E0iPtmFIB46ZmdtAc9eNBvH0H/ZpiBw=="
crossorigin="anonymous" referrerpolicy="no-referrer" />

<style>
    .navbar a {
        color: black;
        text-decoration: none;
    }

    .navbar a:hover {
        cursor: pointer;
        text-decoration: underline;
    }

    .icon-button i.fa {
        font-size: 1.2em;
        color: #8e8e8e;
    }

    #sidebar-close-sidebar {
        border: none;
        background: none;
        cursor: pointer;
    }

    #sidebar-close-sidebar i.fa {
        font-size: 1.2em;
        color: #333;
    }

    #sidenav {
        background-color: var(--secondary-background-color);
        height: 100vh;
    }
</style>
<aside id="sidenav" class="w-full sm:w-1/3 md:w-1/4 px-2">
    <div class="sticky top-0 pr-6 w-full">
        <button id="sidebar-close-sidebar" onclick="hamburgerClick()" class="icon-button" type="button"
            title="Toggle Table of Contents" aria-label="Toggle Table of Contents"
            aria-controls="sidebar" aria-expanded="true">
            <i class="fa fa-xmark"></i>
        </button>
        <ol class="navbar prose lg:prose-sm flex flex-col overflow-hidden">
            <slot name="nav-content"></slot>
        </ol>
    </div>
</aside>`;

class NavMenu extends HTMLElement {
    constructor() {
        super();
        this.attachShadow({ mode: "open" });
        this.shadowRoot.appendChild(template.content.cloneNode(true));
    }
}
window.customElements.define('bb-navmenu', NavMenu);
