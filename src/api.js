import axios from 'axios';
import config from './config';

const API = `${config.BACKEND_URL}/api/v1`;

function headers() {
  return {
    Accept: 'application/json',
    'Content-Type': 'application/json',
  };
}

function queryString(params) {
  const query = Object.keys(params)
    .map((k) => `${encodeURIComponent(k)}=${encodeURIComponent(params[k])}`)
    .join('&');
  return `${query.length ? '?' : ''}${query}`;
}

export default {
  get(url, params = {}) {
    return axios.get(`${API}${url}${queryString(params)}`, { headers: headers() })
      .then((response) => response.data);
  },

  post(url, data) {
    return axios.post(`${API}${url}`, data, { headers: headers() })
      .then((response) => response.data)
      .catch((error) => error.response);
  },

  downloadFile(url, params = {}) {
    return axios.get(`${API}${url}${queryString(params)}`)
      .then((response) => {
        const disposition = response.headers['content-disposition'];
        const filename = decodeURI(disposition.match(/filename=(.*)/)[1]);

        return [response.data, filename];
      });
  },
};
