import { configureStore } from '@reduxjs/toolkit';
import appStateReducer from './state/AppState';

export const store = configureStore({
    reducer: {
        app: appStateReducer,
    },
});
