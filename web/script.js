let phones = [];

    async function loadPhones() {
        try {
            const res = await fetch('/api/phones');
            phones = await res.json();
        } catch (error) {
            console.error("Erreur de chargement des téléphones :", error);
        }
    }

    const phoneInput = document.getElementById('phoneSearch');
    const suggestions = document.getElementById('phoneSuggestions');
    const phoneError = document.getElementById('phoneError');

    function toggleUsageFields() {
        const isPhone = document.getElementById('deviceType').value === 'telephone';
        document.getElementById('phoneModelContainer').classList.toggle('hidden', !isPhone);
        document.querySelectorAll('.usage-telephone').forEach(el => el.classList.toggle('hidden', !isPhone));
        document.querySelectorAll('.usage-ordinateur').forEach(el => el.classList.toggle('hidden', isPhone));
    }

    document.getElementById('deviceType').addEventListener('change', toggleUsageFields);
    window.addEventListener('DOMContentLoaded', () => {
        toggleUsageFields();
        loadPhones();
    });

    phoneInput.addEventListener('input', () => {
        if (phones.length === 0) return;
        const query = phoneInput.value.toLowerCase();
        const matches = phones.filter(p => p.model.toLowerCase().includes(query));
        suggestions.innerHTML = matches.map(p =>
            `<li class="px-2 py-1 cursor-pointer hover:bg-gray-100" data-model="${p.model}" data-co2="${p.co2}">${p.model}</li>`
        ).join('');
        suggestions.classList.toggle('hidden', matches.length === 0);
        phoneError.classList.add('hidden');
    });

    suggestions.addEventListener('click', (e) => {
        if (e.target.matches('li')) {
            phoneInput.value = e.target.dataset.model;
            phoneInput.dataset.co2 = e.target.dataset.co2;
            suggestions.classList.add('hidden');
        }
    });

    document.getElementById('simulator-form').addEventListener('submit', async function (e) {
        e.preventDefault();
        const form = e.target;
        const deviceType = form.deviceType.value;
        const phoneModel = form.phoneModel?.value || null;
        const phoneCO2 = parseFloat(document.getElementById('phoneSearch').dataset.co2 || "0");

        if (deviceType === 'telephone' && phoneCO2 === 0) {
            phoneError.classList.remove('hidden');
            return;
        }

        const data = {
            deviceType,
            phoneModel,
            phoneCO2,
            streaming: parseFloat(form.streaming.value) || 0,
            emails: parseInt(form.emails.value) || 0,
            videoCalls: parseFloat(form.videoCalls?.value) || 0,
            cloudStorage: parseFloat(form.cloudStorage.value) || 0,
            searchQueries: parseInt(form.searchQueries.value) || 0,
            socialMediaHours: parseFloat(form.socialMediaHours.value) || 0,
            downloads: parseFloat(form.downloads.value) || 0,
            musicStreaming: parseFloat(form.musicStreaming?.value) || 0,
            photoSharing: parseInt(form.photoSharing?.value) || 0,
            gpsUsage: parseFloat(form.gpsUsage?.value) || 0
        };

        const res = await fetch('/api/calculate', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        const result = await res.json();
        document.getElementById('co2Output').textContent = `Votre empreinte carbone estimée : ${result.co2.toFixed(2)} g de CO₂`;
        document.getElementById('tips').innerHTML = result.tips.map(t => `<li>✔️ ${t}</li>`).join('');
        document.getElementById('result').classList.remove('hidden');
    });