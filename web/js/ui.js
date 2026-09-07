(() => {
    'use strict';

    const menuButton = document.querySelector('[data-admin-toggle]');
    if (menuButton) {
        document.body.classList.add('ui-interactive');
        const setMenu = (open) => {
            document.body.toggleAttribute('data-admin-open', open);
            menuButton.setAttribute('aria-expanded', String(open));
        };
        menuButton.addEventListener('click', () => setMenu(menuButton.getAttribute('aria-expanded') !== 'true'));
        document.addEventListener('keydown', (event) => {
            if (event.key === 'Escape' && document.body.hasAttribute('data-admin-open')) {
                setMenu(false);
                menuButton.focus();
            }
        });
    }

    const navLinks = [...document.querySelectorAll('.ui-admin-nav a[href]')];
    const currentPath = window.location.pathname.replace(/\/$/, '') || '/';
    const activeLink = navLinks.filter((link) => currentPath === link.pathname || currentPath.startsWith(link.pathname + '/'))
        .sort((first, second) => second.pathname.length - first.pathname.length)[0];
    if (activeLink) {
        activeLink.classList.add('is-active');
        activeLink.setAttribute('aria-current', 'page');
        const group = activeLink.closest('details');
        if (group) group.open = true;
    }

    document.querySelectorAll('[data-password-toggle]').forEach((button) => {
        const field = document.getElementById(button.dataset.passwordToggle);
        if (!field) return;
        button.addEventListener('click', () => {
            const show = field.type === 'password';
            field.type = show ? 'text' : 'password';
            button.textContent = show ? 'Ocultar' : 'Mostrar';
            button.setAttribute('aria-label', show ? 'Ocultar senha' : 'Mostrar senha');
            button.setAttribute('aria-pressed', String(show));
        });
    });

    document.querySelectorAll('[data-password-confirmation]').forEach((form) => {
        const password = form.querySelector('#password');
        const confirmation = form.querySelector('#password2');
        const error = form.querySelector('#passwordError');
        if (!password || !confirmation || !error) return;
        const validate = () => {
            const mismatch = confirmation.value.length > 0 && password.value !== confirmation.value;
            confirmation.setCustomValidity(mismatch ? 'As senhas precisam ser iguais.' : '');
            confirmation.setAttribute('aria-invalid', String(mismatch));
            error.hidden = !mismatch;
        };
        password.addEventListener('input', validate);
        confirmation.addEventListener('input', validate);
        form.addEventListener('submit', (event) => {
            validate();
            if (!form.checkValidity()) {
                event.preventDefault();
                form.reportValidity();
            }
        });
    });

    document.querySelectorAll('[data-dismiss-message]').forEach((button) => {
        button.addEventListener('click', () => button.closest('.ui-notice')?.remove());
    });

    const queryInput = document.querySelector('form[method="get"] input[name="q"]');
    if (queryInput) queryInput.value = new URLSearchParams(window.location.search).get('q') || '';

    document.querySelectorAll('.postal-code, [data-mask="postal-code"]').forEach((input) => {
        input.addEventListener('input', () => {
            const digits = input.value.replace(/\D/g, '').slice(0, 8);
            input.value = digits.length > 5 ? `${digits.slice(0, 5)}-${digits.slice(5)}` : digits;
        });
    });

    // Keep telephone input permissive for country codes and assistive input.
    document.querySelectorAll('.phone-number, [data-mask="phone"]').forEach((input) => {
        input.type = 'tel';
        input.inputMode = 'tel';
    });

    const fileInput = document.querySelector('[data-image-input], #image');
    const imageLabel = document.querySelector('[data-image-label], #imageLabel');
    const previewImage = document.querySelector('[data-image-preview], #previewImage');
    if (fileInput?.type === 'file' && imageLabel && previewImage) {
        const originalImage = previewImage.getAttribute('src');
        let previewURL;
        imageLabel.setAttribute('role', 'status');
        fileInput.setAttribute('aria-describedby', imageLabel.id);
        const releasePreview = () => {
            if (previewURL) URL.revokeObjectURL(previewURL);
            previewURL = undefined;
        };
        fileInput.addEventListener('change', () => {
            releasePreview();
            if (originalImage) previewImage.src = originalImage;
            const file = fileInput.files[0];
            imageLabel.classList.remove('text-danger');
            if (!file) {
                imageLabel.textContent = 'Escolher imagem do produto';
                return;
            }
            if (!['image/png', 'image/jpeg', 'image/gif'].includes(file.type) || file.size > 4 * 1024 * 1024) {
                fileInput.value = '';
                imageLabel.textContent = 'Selecione uma imagem PNG, JPEG ou GIF de até 4 MB.';
                imageLabel.classList.add('text-danger');
                return;
            }
            imageLabel.textContent = 'Imagem selecionada: ' + file.name;
            previewURL = URL.createObjectURL(file);
            previewImage.src = previewURL;
        });
        window.addEventListener('pagehide', releasePreview);
    }
})();
