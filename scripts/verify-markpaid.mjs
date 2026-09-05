const base = 'http://localhost:3000';
(async () => {
  const login = await fetch(base + '/api/auth/login', {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify({ email: 'admin@reflexcms.dev', password: 'ReflexCMS@2026' })
  });
  const raw = login.headers.getSetCookie?.() ?? [login.headers.get('set-cookie') || ''];
  const cookie = (raw.find(c => c.startsWith('admin_session=')) || raw[0] || '').split(';')[0];
  if (!cookie) { console.log('NO COOKIE, login status', login.status); return; }
  const auth = { 'content-type': 'application/json', cookie };

  const list = await fetch(base + '/api/admin/orders?perPage=10', { headers: auth }).then(r => r.json());
  const o = (list.items || []).find(i => i.order_no === 'R-TEST-GUEST-001');
  if (!o) { console.log('order not found, total=', list.total); return; }
  const r = await fetch(base + '/api/admin/order-mark-paid', { method: 'POST', headers: auth, body: JSON.stringify({ order_id: o.id }) });
  console.log('mark-paid:', r.status, JSON.stringify(await r.json()));
  const view = await fetch(base + '/api/v1/orders/R-TEST-GUEST-001').then(r => r.json());
  console.log('guest view: status=' + view.status, 'code=' + view.granted_code);
})();
