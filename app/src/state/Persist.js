const localStorageKey = 'mokki-state-storage-json';

export const loadState = () => {
    return JSON.parse(window.localStorage.getItem(localStorageKey));
};

export const storeState = (state) => {
    window.localStorage.setItem(localStorageKey, JSON.stringify(state));
};