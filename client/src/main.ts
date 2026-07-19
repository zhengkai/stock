import { pb } from './pb';

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

	const sub = (await registration.pushManager.subscribe({
		userVisibleOnly: true,
		applicationServerKey: pubKey,
	})).toJSON();

	const t = pb.VAPIDSubscription

	const msg = t.fromObject({
		endpoint: sub.endpoint,
		p256dh: sub.keys?.p256dh,
		auth: sub.keys?.auth,
	})

	const bin = t.encode(msg).finish();

	console.log('sub', sub)

	await fetch('/api/sub', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/protobuf',
		},
		body: new Uint8Array(bin).slice().buffer,
	});
};

const checkSub = async (div: HTMLElement) => {
	const registration = await navigator.serviceWorker.ready;

	const subscription =
		await registration.pushManager.getSubscription();

	if (subscription) {
		div.innerText = 'subscribed';
	} else {
		div.innerText = 'not subscribed';
	}
};

(() => {

	const box = document.createElement('div');
	box.className = 'container';

	const btn = document.createElement('button');
	btn.className = 'btn btn-primary mb-4';
	btn.innerText = 'Subscribe';
	btn.onclick = sub;

	box.appendChild(btn);

	const div = document.createElement('div');
	div.innerText = 'checking...';
	box.appendChild(div);
	checkSub(div);

	app.appendChild(box);
})();
