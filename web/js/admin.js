'use strict';

// Stock search is shared by the movement form and the product history page.
for (const picker of document.querySelectorAll('[data-product-picker]')) {
    const search = picker.querySelector('[data-product-search]');
    const results = picker.querySelector('[data-product-results]');
    const status = picker.querySelector('[data-search-status]');
    const selectedID = picker.querySelector('[data-selected-product]');
    const details = picker.querySelector('[data-product-details]');
    const quantity = picker.querySelector('[data-product-quantity]');
    const submit = picker.querySelector('[data-product-submit]');
    const history = picker.querySelector('[data-product-history]');
    const currency = new Intl.NumberFormat('pt-BR', {style: 'currency', currency: 'BRL'});
    let timer;
    let request;

    function clearSelection() {
        selectedID.value = '';
        details.hidden = true;
        if (quantity) quantity.disabled = true;
        if (submit) submit.disabled = true;
        if (history) {
            history.hidden = true;
            history.removeAttribute('href');
        }
    }

    function selectProduct(product) {
        selectedID.value = String(product.ID);
        search.value = product.Name;
        for (const field of picker.querySelectorAll('[data-product-field]')) {
            const value = product[field.dataset.productField];
            field.textContent = field.dataset.productField === 'Price'
                ? currency.format(Number(value)) : String(value ?? 'Não informado');
        }
        const image = picker.querySelector('[data-product-image]');
        image.src = typeof product.Image === 'string' && product.Image.startsWith('/images/')
            ? product.Image : '/images/600x400.svg';
        image.alt = product.Name;
        details.hidden = false;
        results.hidden = true;
        results.replaceChildren();
        status.textContent = `Produto selecionado: ${product.Name}.`;
        if (quantity) quantity.disabled = false;
        if (submit) submit.disabled = false;
        if (history) {
            history.href = `/stock/${product.ID}`;
            history.hidden = false;
        }
        (quantity || history)?.focus();
    }

    async function findProducts(query) {
        request = new AbortController();
        const currentRequest = request;
        status.textContent = 'Buscando produtos…';
        picker.setAttribute('aria-busy', 'true');
        try {
            const params = new URLSearchParams({q: query, limit: '8'});
            const response = await fetch(`/products/json?${params}`, {
                signal: currentRequest.signal,
                headers: {'Accept': 'application/json'},
                credentials: 'same-origin',
                redirect: 'error'
            });
            if (!response.ok) throw new Error('Falha na busca');
            const data = await response.json();
            if (currentRequest.signal.aborted || search.value.trim() !== query) return;
            const products = Array.isArray(data.Products) ? data.Products.filter(product =>
                Number.isSafeInteger(product.ID) && product.ID > 0 && typeof product.Name === 'string'
            ) : [];
            results.replaceChildren();
            for (const product of products) {
                const item = document.createElement('li');
                const button = document.createElement('button');
                button.type = 'button';
                button.className = 'product-picker-option';
                button.textContent = `${product.Name} · #${product.ID} · ${product.CurrentStock} em estoque`;
                button.addEventListener('click', () => selectProduct(product));
                item.append(button);
                results.append(item);
            }
            results.hidden = products.length === 0;
            status.textContent = products.length
                ? `${products.length} produto(s) encontrado(s). Selecione abaixo.`
                : 'Nenhum produto encontrado. Tente outro nome.';
        } catch (error) {
            if (!currentRequest.signal.aborted && request === currentRequest && search.value.trim() === query) {
                results.hidden = true;
                status.textContent = 'Não foi possível buscar os produtos. Tente novamente.';
            }
        } finally {
            if (request === currentRequest) picker.removeAttribute('aria-busy');
        }
    }

    function searchProducts() {
        clearTimeout(timer);
        request?.abort();
        picker.removeAttribute('aria-busy');
        clearSelection();
        results.hidden = true;
        results.replaceChildren();
        const query = search.value.trim();
        if (query.length < 2) {
            status.textContent = 'Digite pelo menos 2 letras para buscar.';
            return;
        }
        timer = setTimeout(() => findProducts(query), 250);
    }

    search.addEventListener('input', searchProducts);
    search.addEventListener('keydown', event => {
        if (event.key === 'Escape') results.hidden = true;
        if (event.key === 'ArrowDown' && !results.hidden) {
            event.preventDefault();
            results.querySelector('button')?.focus();
        }
    });
    picker.querySelector('form')?.addEventListener('submit', event => {
        if (!selectedID.value) {
            event.preventDefault();
            status.textContent = 'Selecione um produto antes de registrar a movimentação.';
            search.focus();
        }
    });
    const initialQuery = new URLSearchParams(window.location.search).get('q');
    if (initialQuery) {
        search.value = initialQuery;
        searchProducts();
    }
}
