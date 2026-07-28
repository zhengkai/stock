import { pb } from './pb';
import Long from "long";

function n(v: number | Long | null | undefined): number {
	if (Long.isLong(v)) {
		return v.toNumber();
	}
	return v ?? 0;
}

function price(v: number | null | undefined) {
	return (n(v) / 100).toFixed(2);
}

function percent(v: number | null | undefined) {
	return (n(v) / 10000).toFixed(2) + '%';
}

function xueqiuLink(code: string | null | undefined) {
	if (!code) {
		return '#';
	}
	let prefix = 'SZ';
	if (code.startsWith('6')) {
		prefix = 'SH';
	}
	return `https://xueqiu.com/S/` + prefix + code;
}

function rate(a: number | null | undefined, b: number | null | undefined) {
	a = n(a);
	b = n(b);

	if (a == 0 || b == 0) {
		return '--rate: 0';
	}

	const f = ((b / 10) - (b - a)) / (b / 10);
	console.log('rate', a, b, f);
	if (f <= 0) {
		return '--rate: 0';
	}
	if (f > 1) {
		return '--rate: 1';
	}
	return `--rate: ${f.toFixed(2)}`;
}

// function volume(v: number | null | undefined) {
// 	return n(v).toLocaleString();
// }
//
// function money(v: number | null | undefined) {
// 	return n(v).toLocaleString() + ' 万';
// }

function timestamp(v: number | null | undefined) {
	const d = new Date(n(v) * 1000);
	const s = d.getFullYear()
		+ "-" + String(d.getMonth() + 1).padStart(2, "0")
		+ "-" + String(d.getDate()).padStart(2, "0")
		+ " " + String(d.getHours()).padStart(2, "0")
		+ ":" + String(d.getMinutes()).padStart(2, "0")
		+ ":" + String(d.getSeconds()).padStart(2, "0");
	return s;
}

function renderPool(ap: pb.AppPool) {
	const tbody = document.querySelector('#stock-table tbody')!;

	tbody.innerHTML = ap.stock.map((s) => {
		const q = s.quote;
		const a = s.alert;

		return `<tr>
	<td class="name">${s?.alert?.name}</td>
	<td class="code"><a href="${xueqiuLink(s.code)}" target="_blank">${s.code}</td>
	<td class="price">${price(q?.open)}</td>
	<td class="price bar-down" style="${rate(a?.min, q?.price)}">${price(a?.min)}</td>
	<td class="price">${price(q?.price)}</td>
	<td class="price bar-up" style="${rate(q?.price, a?.max)}">${price(a?.max)}</td>
	<td class="price">${price(q?.change)}</td>
	<td class="price">${percent(q?.changeBp)}</td>
	<td class="price">${price(q?.high)}</td>
	<td class="price">${price(q?.low)}</td>
	<td class="datetime">${timestamp(q?.ts)}</td>
</tr>`;
	}).join('');
}

(async () => {
	const res = await fetch('/api/pool');
	if (!res.ok) {
		console.error('Failed to fetch pool:', res.statusText);
	}

	const buffer = await res.arrayBuffer();
	const ap = pb.AppPool.decode(new Uint8Array(buffer));

	console.log(ap);

	renderPool(ap);
})();
