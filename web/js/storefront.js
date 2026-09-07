/* Forms remain usable without JavaScript; this file only enhances the storefront. */
(() => {
    'use strict';

    const path = window.location.pathname;
    const currentSection = path.startsWith('/cart') ? 'cart'
        : path.startsWith('/orders') ? 'orders'
        : path.startsWith('/profile') || path.startsWith('/users') ? 'profile' : 'store';
    document.querySelectorAll('[data-nav]').forEach(link => {
        if (link.dataset.nav === currentSection) link.setAttribute('aria-current', 'page');
    });

    document.querySelectorAll('[data-quantity]').forEach(control => {
        const input = control.querySelector('[data-quantity-input]');
        const minimum = Number(input.min) || 1;
        const maximum = input.max === '' ? Number.MAX_SAFE_INTEGER : Number(input.max);
        const buttons = control.querySelectorAll('[data-quantity-step]');
        const refresh = () => {
            const value = Number(input.value);
            input.setCustomValidity(input.value && !Number.isSafeInteger(value)
                ? 'Informe uma quantidade inteira válida.' : '');
            buttons.forEach(button => {
                button.disabled = maximum < minimum || (Number(button.dataset.quantityStep) < 0
                    ? value <= minimum : value >= maximum);
            });
        };
        buttons.forEach(button => button.hidden = false);
        buttons.forEach(button => button.addEventListener('click', () => {
            const value = Number(input.value);
            const current = Number.isSafeInteger(value) && value >= minimum ? value : minimum;
            input.value = String(Math.min(maximum, Math.max(minimum, current + Number(button.dataset.quantityStep))));
            input.dispatchEvent(new Event('input', {bubbles: true}));
        }));
        input.addEventListener('input', refresh);
        refresh();
    });

    const checkout = document.querySelector('[data-checkout]');
    if (checkout) {
        const delivery = checkout.querySelector('[data-delivery-select]');
        const subtotal = Number(checkout.dataset.subtotal);
        const deliveryPrice = Number(checkout.dataset.deliveryPrice);
        const formatter = new Intl.NumberFormat('pt-BR', {style: 'currency', currency: 'BRL'});
        const updateDelivery = () => {
            const isDelivery = delivery.value === '1';
            checkout.querySelectorAll('[data-delivery-address]').forEach(panel => panel.hidden = !isDelivery);
            checkout.querySelectorAll('[data-pickup-address]').forEach(panel => panel.hidden = isDelivery);
            checkout.querySelector('[data-delivery-row]').hidden = !isDelivery;
            const total = subtotal + (isDelivery ? deliveryPrice : 0);
            if (Number.isFinite(total)) checkout.querySelector('[data-checkout-total]').textContent = formatter.format(total);
        };
        if (delivery) {
            delivery.addEventListener('change', updateDelivery);
            updateDelivery();
        }
    }

    document.querySelectorAll('.store-product-image img, .store-detail-image img, .store-line-item > img').forEach(img => {
        const usePlaceholder = () => { img.src = '/images/cupcake-placeholder.svg'; };
        img.addEventListener('error', usePlaceholder, {once: true});
        if (img.complete && img.naturalWidth === 0) usePlaceholder();
    });

    const pixCopy = document.querySelector('[data-pix-copy]');
    const pixCode = document.querySelector('[data-pix-code]');
    const pixStatus = document.querySelector('[data-pix-copy-status]');
    if (pixCopy && pixCode && pixStatus && navigator.clipboard?.writeText) {
        pixCopy.hidden = false;
        pixCopy.addEventListener('click', async () => {
            try {
                await navigator.clipboard.writeText(pixCode.value);
                pixStatus.textContent = 'Código copiado. Cole no aplicativo do seu banco.';
            } catch {
                pixCode.focus();
                pixCode.select();
                pixStatus.textContent = 'Não foi possível copiar automaticamente. Copie o texto selecionado e cole no seu banco.';
            }
        });
    }

    const count = document.getElementById('cart-count');
    if (count && document.body.dataset.authenticated === 'true') {
        fetch('/cart/count', {headers: {Accept: 'application/json'}, credentials: 'same-origin', redirect: 'error'})
            .then(response => {
                if (!response.ok) throw new Error('Não foi possível consultar o carrinho.');
                return response.json();
            })
            .then(data => {
                if (Number.isSafeInteger(data.itemCount) && data.itemCount >= 0) count.textContent = String(data.itemCount);
            })
            .catch(() => { count.textContent = '–'; });
    }
})();
