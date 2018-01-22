import { invalid } from '../func';

let url = `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/sys-ws`;
if (window.location.hostname === 'localhost') url = 'ws://localhost:8099/sys-ws';

const ws = {
  open: false,
  i: 0,
  promise: []
};

ws.start = (ch) => {
  ws.change = ch;
  ws.connect();
  ws.interval = setInterval(() => {
    if (!ws.open) ws.connect();
    else {
      ws.request('ping', '');
    }
  }, 5000);
};

ws.stop = () => {
  clearInterval(ws.interval);
};

ws.connect = () => {
  const conn = new WebSocket(url);
  conn.onopen = () => {
    ws.open = true;
    ws.change(true);
  };
  conn.onclose = () => {
    ws.open = false;
    ws.change(false);
  };
  conn.onerror = () => {
    ws.open = false;
    ws.change(false);
  };
  conn.onmessage = (evt) => {
    evt = JSON.parse(evt.data);
    const p = ws.promise[evt.i];
    if (!invalid(p)) {
      ws.promise[evt.i] = null;
      if (evt.status !== 200) p.reject(evt.error);
      else p.resolve(evt.data);
    }
  };

  ws.conn = conn;
};

ws.request = (evt, data) => new Promise((resolve, reject) => {
  if (!ws.open) {
    setTimeout();
    reject(new Error('websocket closed'));
    return;
  }
  if (invalid(data)) data = '';
  const { i } = ws;
  ws.i += 1;
  if (ws.i === 10000) ws.i = 0;
  ws.promise[i] = { resolve, reject };
  ws.conn.send(JSON.stringify({ evt, i, data }));
});

export default ws;
