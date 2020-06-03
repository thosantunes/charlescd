import mock from './mock';

const API = '/auth/realms/darwin/protocol/openid-connect/token';

const keycloak = {
  method: 'POST',
  path: `${API}`,
  handler: (req, h) => h.response(mock.auth)
};

export default {
  keycloak
};
