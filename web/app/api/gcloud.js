import ws from './websocket';

async function getInfo() {
  const result = await ws.request('gcloud/getInfo');
  return result;
}

async function setAuthKey(data) {
  const result = await ws.request('gcloud/setAuthKey', data);
  return result;
}

async function setProject(data) {
  const result = await ws.request('gcloud/setProject', data);
  return result;
}

async function setGKE(zone, name) {
  const result = await ws.request('gcloud/setGKE', { zone, name });
  return result;
}

export default {
  getInfo,
  setAuthKey,
  setProject,
  setGKE
};
