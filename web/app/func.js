export const isNum = obj => obj !== '' && obj !== null && !Number.isNaN(Number(obj));
export const isStr = obj => typeof obj === 'string';
export const isArr = obj => Array.isArray(obj);
export const isObj = obj => (!!obj) && obj === Object(obj) && !isArr(obj);
export const isFunc = obj => typeof obj === 'function';
export const invalid = obj => obj === undefined || obj === null;
export const Undefined = obj => obj === undefined;
export const has = (obj, key) => {
  if (invalid(key)) return false;
  if (isArr(obj)) return obj.indexOf(key) !== -1;
  if (isObj(obj)) return Object.prototype.hasOwnProperty.call(obj, key);
  if (isStr(obj)) return obj.indexOf(key) !== -1;
  return undefined;
};
export const hasIndex = (obj, i) => obj[i] !== undefined;
export const keys = obj => Object.keys(obj);
export const values = obj => Object.keys(obj).map(k => obj[k]);
export const assign = (...obj) => Object.assign({}, ...obj);
export const isEmpty = (obj) => {
  if (isArr(obj)) return obj.length === 0;
  if (isObj(obj)) return keys(obj).length === 0;
  if (isStr(obj)) return obj === '';
  return false;
};
export const pick = (obj, ks) => {
  const result = {};
  ks.forEach((k) => {
    if (has(obj, k)) result[k] = obj[k];
  });
  return result;
};
export const forEach = (obj, callback) => keys(obj).forEach((key, i) => callback(obj[key], key, i));
export const map = (obj, callback) => keys(obj).map((key, i) => callback(obj[key], key, i));
export const same = (a, b) => {
  if (isObj(a) && isObj(b)) {
    const ks = keys(a);
    for (let i = 0; i < ks.length; i += 1) {
      const k = ks[i];
      if (!has(b, k)) return false;
      if (!same(a[k], b[k])) return false;
    }
  } else return a === b;
  return true;
};
export const filter = (obj, callback) => {
  const result = {};
  forEach(obj, (v, k) => {
    if (callback(v, k)) result[k] = v;
  });
  return result;
};
export const clone = obj => JSON.parse(JSON.stringify(obj));
export const remove = (obj, key) => {
  if (isObj(obj)) delete obj[key];
  if (isArr(obj)) {
    const i = obj.indexOf(key);
    if (i === -1) return;
    obj.splice(i, 1);
  }
};
export const len = (obj) => {
  if (isArr(obj)) return obj.length;
  if (isObj(obj)) return keys(obj).length;
  return null;
};
