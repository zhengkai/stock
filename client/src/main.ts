if ('serviceWorker' in navigator) {
	navigator.serviceWorker.register("/worker.js")
		.then(reg => {
			console.log("SW registered", reg);
		});
}

const app = document.querySelector<HTMLDivElement>('#app')!;

const pubKey = 'BFxTRa4-p7GX1aQFWXkhqN4iEjQohbcvpyslD6ZhZPKYLiO85lx5iigt0udB4G3ONK9gzNMR-81dTMvAkv80KVw';

const sub = async () => {

	console.log('start sub')

	const registration = await navigator.serviceWorker.ready;

	const sub = await registration.pushManager.subscribe({
		userVisibleOnly: true,
		applicationServerKey: pubKey,
	});

	console.log('sub', sub)
};

(() => {

	const box = document.createElement('div');
	box.className = 'container';

	const btn = document.createElement('button');
	btn.className = 'btn btn-primary';
	btn.innerText = 'Subscribe';
	btn.onclick = sub;

	box.appendChild(btn);

	app.appendChild(box);
})();
