import ws from './websocket';

async function getList() {
  const result = await ws.request('repo/getList');
  return JSON.parse(result);
}

async function add(name, src, branch) {
  const result = await ws.request('repo/add', { name, src, branch });
  return result;
}


async function trigger(name, branch) {
  const result = await ws.request('repo/trigger', { name, branch });
  return result;
}

async function remove(name, branch) {
  const result = await ws.request('repo/remove', { name, branch });
  return result;
}

export default {
  getList,
  add,
  trigger,
  remove
};
