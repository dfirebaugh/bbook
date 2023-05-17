const template = document.createElement('template');
template.innerHTML = `
<aside id="sidenav" class="w-full sm:w-1/3 md:w-1/4 px-2">
		<div class="sticky top-0 pr-6 w-full">
				<slot name="nav-content"></slot>
		</div>
</aside>`;

class NavMenu extends HTMLElement {
	constructor() {
		super();
		const shadowRoot = this.attachShadow({ mode: "open" });
		shadowRoot.appendChild(template.content.cloneNode(true));
	}
}
window.customElements.define('bf-navmenu', NavMenu);