"use strict";

(function () {
    for (let n = 1; n <= 6; n++) {
        document.querySelectorAll('h' + n).forEach((hN) => {
            // get the id of the heading
            const id = hN.getAttribute('id');
            if (!id) return;

            // create a link for it
            const a = document.createElement('a');
            a.className = 'autolink';
            a.href = '#' + id;
            a.innerHTML = '#';

            // and add the link to it
            hN.appendChild(a);
        });
    }
})();
