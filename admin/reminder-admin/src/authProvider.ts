import { AuthProvider, HttpError } from "react-admin";
import data from "./users.json";

/**
 * This authProvider is only for test purposes. Don't use it in production.
 */
//
// const authProvider = {
//   login: ({ username, password }) =>  {
//     const request = new Request(`${import.meta.env.VITE_SIMPLE_REST_URL}/authenticate`, {
//       method: 'POST',
//       body: JSON.stringify({ username, password }),
//       headers: new Headers({ 'Content-Type': 'application/json' }),
//     });
//     return fetch(request)
//         .then(response => {
//           if (response.status < 200 || response.status >= 300) {
//             throw new Error(response.statusText);
//           }
//           return response.json();
//         })
//         .then(auth => {
//           localStorage.setItem('auth', JSON.stringify(auth));
//         })
//         .catch(() => {
//           throw new Error('Network error')
//         });
//   },
//   checkAuth: () => {
//     // Required for the authentication to work
//     return Promise.resolve();
//   },
//   getPermissions: () => {
//     // Required for the authentication to work
//     return Promise.resolve();
//   },
//   // ...
// };

export const authProvider: AuthProvider = {
  login: ({ username, password }) => {
    const request = new Request(`${import.meta.env.VITE_SIMPLE_REST_URL}/auth/authenticate`, {
      method: 'POST',
      body: JSON.stringify({ username, password }),
      headers: new Headers({ 'Content-Type': 'application/json' }),
    });
    return fetch(request)
        .then(response => {
          if (response.status < 200 || response.status >= 300) {
            throw new Error(response.statusText);
          }
          return response.json();
        })
        .then(auth => {
          localStorage.setItem('user', JSON.stringify(auth));
          localStorage.setItem('token', auth.token);
        })
        .catch(() => {
          throw new Error('Network error')
        });
  },
  logout: () => {
    localStorage.removeItem("user");
    localStorage.removeItem("token");
    return Promise.resolve();
  },
  checkError: (error) => {
    const status = error.status;
    if (status === 401 || status === 403) {
      localStorage.removeItem('user');
      localStorage.removeItem('token');
      return Promise.reject();
    }
  },
  checkAuth: () =>
    localStorage.getItem("user") ? Promise.resolve() : Promise.reject(),
  getPermissions: () => {
    const user = JSON.parse(localStorage.getItem("user"))
    return Promise.resolve(user.role);
  },
  getIdentity: () => {
    const persistedUser = localStorage.getItem("user");
    const user = persistedUser ? JSON.parse(persistedUser) : null;
    return Promise.resolve({
      id: user.id,
      fullName: user.email,
    })
    // return Promise.resolve(user);
  },
};

export function retrieveToken() {
  return localStorage.getItem("token");
}

export default authProvider;
