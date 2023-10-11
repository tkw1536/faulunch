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

(function () {
    // ensure the share api is there
    if (typeof navigator.share !== 'function') {
        console.warn('navigator.share is not a function');
        return
    }

    // find the element to add the share button to
    const element = document.getElementById('add-share-button');
    if (!element) {
        console.warn('no share to add');
        return
    }

    // create element
    const a = document.createElement('a');
    a.setAttribute('href', 'javascript:void(0);')
    a.append(document.createTextNode(document.documentElement.lang !== 'de' ? 'Share' : 'Teilen'));

    // add the link
    element.prepend(document.createTextNode(' '));
    element.prepend(a);

    a.addEventListener('click', (evt) => {
        evt.preventDefault();


        const description = document.querySelector('meta[name=description]');
        const metaDescription = (description && description.hasAttribute('content')) ? description.getAttribute('content') : undefined;

        navigator.share({
            'text': metaDescription,
            'title': document.title,
            'url': location.href,
        })
    })





})()