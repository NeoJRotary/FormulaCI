import ws from './websocket';

async function getInfo() {
  const result = await ws.request('git/getInfo');
  return result;
}

async function setEmail(data) {
  const result = await ws.request('git/setEmail', data);
  return result;
}

async function setWebhookToken(data) {
  const result = await ws.request('git/setWebhookToken', data);
  return result;
}

async function generateSSH() {
  const result = await ws.request('git/generateSSH');
  return result;
}

export default {
  getInfo,
  setEmail,
  setWebhookToken,
  generateSSH
};
