import ws from './websocket';

async function getHistory() {
  const result = await ws.request('formula/getHistory');
  return JSON.parse(result).map((r) => {
    if (r.flow === '') r.flow = [];
    else r.flow = JSON.parse(r.flow);
    if (r.log === '') r.log = {};
    else r.log = JSON.parse(r.log);
    const t = new Date(r.time * 1000);
    r.time = t.toLocaleString();
    return r;
  });
}

export default {
  getHistory
};
