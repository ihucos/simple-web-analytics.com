customElements.define(
    tagName(),
    class extends HTMLElement {
        KEY = "close-looking-for-cofounder";
        connectedCallback() {
            if (
                !window.location.href
                    .split("#")[0]
                    .endsWith("/dashboard.html?demo=1") &&
                !localStorage.getItem(this.KEY)
            ) {
                this.innerHTML += `
                   <!-- Notification -->
                   <section class="notification">
                     <div class="content">
                       <span>Help us with marketing as our Co-Founder. <a href="/cofounder.html">Learn more.</a></span>
                       <div class="btn-close"></div>
                     </div>
                   </section>`;

                $(".btn-close").click(() => {
                    this.remove();
                    localStorage.setItem(this.KEY, "yes");
                });
            }
        }
    }
);
