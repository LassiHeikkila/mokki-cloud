const localStorageKey = 'mokki-state-storage-json';

export const loadState = () => {
    return JSON.parse(localStorage.getItem(localStorageKey));
};

export const storeState = (state) => {
    localStorage.setItem(localStorageKey, JSON.stringify(state));
};