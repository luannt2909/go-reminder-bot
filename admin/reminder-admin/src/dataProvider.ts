import {fetchUtils} from "react-admin"
import simpleRestProvider from "ra-data-simple-rest";
import {retrieveToken} from "./authProvider"

export const getToken = () => {
    const token = retrieveToken();
    return `Bearer ${token}`;
};

const httpClient = (url, options = {}) => {
    if (!options.headers) {
        options.headers = new Headers({ Accept: "application/json" });
    }

    options.headers.set("Authorization", getToken());
    return fetchUtils.fetchJson(url, options);
};
export const dataProvider = simpleRestProvider(import.meta.env.VITE_SIMPLE_REST_URL, httpClient);
// export const dataProvider = simpleRestProvider(
//   import.meta.env.VITE_SIMPLE_REST_URL
// );
