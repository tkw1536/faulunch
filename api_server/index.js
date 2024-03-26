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
})();

(function () {
    // find the autosort list and place to put the ui
    // and make sure they exist
    const autoSortList = document.querySelector('ul#autosort-list');
    const autoSortUI = document.querySelector('#autosort-ui');
    if (!autoSortUI || !autoSortList) return;

    const increasingText = autoSortUI.getAttribute('data-increasing-text');
    const decreasingText = autoSortUI.getAttribute('data-decreasing-text');

    // known sorting criteria
    const criteriaCategories = new Map();

    // determine all the actual sections
    const items = Array.from(autoSortList.querySelectorAll('li a'))
        .map(function (a, index) {
            const li = a.parentElement;
            if (!li) return null;
            if (li.tagName !== 'LI') return;

            // get the href element
            const href = a.getAttribute('href');
            if (!href) return null;

            // get the section being linked
            if (!href.startsWith('#')) return null;
            const section = document.getElementById(href.substring(1));
            if (!section) return null;
            if (section.tagName !== 'SECTION') { return null; }

            // parse all the data values
            const values = Array.from(section.querySelectorAll('tr'))
                .map(function (tr) {
                    // get the closest summary element
                    const details = tr.closest('details');
                    if (!details) { return null; }
                    const summary = details.querySelector('summary')
                    if (!summary) { return null; }
                    const category = summary.textContent;


                    // find elements with exactly two elements
                    const tds = tr.querySelectorAll('td');
                    if (tds.length != 2) { return null; }

                    // find all the attributes of this thing
                    const value = tds[1].querySelector('math mn');
                    if (!value) { return null; }
                    const sort = parseFloat(value.textContent.replaceAll(',', '.'));

                    const attr = tds[0].textContent.trim();

                    // add it to the appropriate critera set
                    if (!criteriaCategories.has(category)) {
                        criteriaCategories.set(category, new Set());
                    }
                    criteriaCategories.get(category).add(attr);

                    // return a key-value pair
                    return [
                        attr,
                        {
                            sort: sort,
                            value: tds[1].querySelector('math'),
                        }
                    ];
                })
                .filter(function (e) { return e !== null });

            return { li: li, values: new Map(values) };
        }).filter(function (e) { return e !== null });

    const doSort = function (criterion, increasing) {
        // create a copy of the items
        const sortedItems = items.slice(0).map((item, index) => {
            const value = (criterion === null) ? { sort: index } : (item.values.get(criterion) ?? {});
            return {
                li: item.li,
                sort: value.sort ?? 0,
                value: value.value ?? null,
            };
        });

        // remove all the items from the list
        sortedItems.forEach(li => autoSortList.removeChild(li.li));

        // sort in the right order
        if (increasing) {
            sortedItems.sort((a, b) => a.sort - b.sort);
        } else {
            sortedItems.sort((a, b) => b.sort - a.sort);
        }

        // add the items back and update the value element
        sortedItems.forEach(elem => {
            // make sure there is a class for spacing
            elem.li.classList.add('sorted-list-item');

            const span = elem.li.querySelector('span.sort-value');
            if (span) {
                span.parentNode.removeChild(span);
            }

            // make a clone of the value element or create one for spacing
            let valueElem = elem.value;
            if (valueElem) {
                const value = document.createElement('span');
                elem.li.appendChild(value);
                value.setAttribute('class', 'sort-value');
                value.appendChild(document.createTextNode(' '));
                value.appendChild(valueElem.cloneNode(true));
            }

            // and append the child to it!
            autoSortList.appendChild(elem.li);
        });

        // update the ui for the sort critera
        Array.from(autoSortUI.querySelectorAll('a'))
            .forEach(a => {
                const aCriterion = a.getAttribute('data-sort-criterion') ?? '';

                a.innerHTML = '';
                a.appendChild(document.createTextNode(aCriterion));

                if (aCriterion !== criterion) {
                    // remove the class and sort stage
                    a.classList.remove('active');
                    a.removeAttribute('aria-current');
                    a.setAttribute('data-sort-stage', '0');
                    return;
                }

                a.classList.add('active');
                a.setAttribute('aria-current', 'true');
                a.appendChild(document.createTextNode(' '));

                const info = document.createElement('span');
                if (increasing) {
                    info.appendChild(document.createTextNode('+'));
                    info.setAttribute('aria-description', increasingText);
                } else {
                    info.appendChild(document.createTextNode('-'));
                    info.setAttribute('aria-description', decreasingText);
                }

                a.appendChild(info);
            })
    }

    autoSortUI.innerHTML = '';

    criteriaCategories.forEach((criteria, category) => {
        const p = document.createElement('p');
        autoSortUI.appendChild(p);

        p.appendChild(document.createTextNode(category + ': '));

        criteria.forEach(criterion => {
            // create an element that sorts increasing by default
            const a = document.createElement('a');
            a.setAttribute('role', 'menuitem');
            a.setAttribute('href', 'javascript:void(0)');
            a.setAttribute('data-sort-criterion', criterion);
            a.setAttribute('data-sort-stage', '0');

            // add the text node thing
            a.appendChild(document.createTextNode(criterion))
            a.addEventListener('click', (evt) => {
                evt.preventDefault(true);

                const stage = a.getAttribute('data-sort-stage');
                if (stage === '0') {
                    a.setAttribute('data-sort-stage', '1');
                    doSort(criterion, true);
                } else if (stage === '1') {
                    a.setAttribute('data-sort-stage', '2');
                    doSort(criterion, false);
                } else {
                    a.setAttribute('data-sort-stage', '0');
                    doSort(null, true);
                }
            });

            p.appendChild(a);
            p.appendChild(document.createTextNode(' '));
        });
    });


    doSort(null, true);
})();