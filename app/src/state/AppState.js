import { createSlice } from '@reduxjs/toolkit';
import { oneWeek } from '../lib/units.js';

export const appStateSlice = createSlice({
    name: 'app',
    initialState: {
        token: '',
        isAuthenticated: false,
        darkmode: false,
        selectedSensor: '',
        selectedMeasurement: '',
        selectedTimePeriod: oneWeek,
    },
    reducers: {
        setToken: (state, action) => {
            // Redux makes the passed in token an object where the original string value is in "payload" field.
            state.token = action.payload;
        },
        clearToken: (state) => {
            state.token = '';
        },
        setIsAuthenticated: (state, action) => {
            state.isAuthenticated = action.payload;
        },
        setDarkmode: (state, action) => {
            state.darkmode = action.payload;
        },
        setSelectedSensor: (state, action) => {
            state.selectedSensor = action.payload;
        },
        setSelectedMeasurement: (state, action) => {
            state.selectedMeasurement = action.payload;
        },
        setSelectedTimePeriod: (state, action) => {
            state.selectedTimePeriod = action.payload;
        }
    },
});

export const { 
    setToken, 
    clearToken, 
    setIsAuthenticated, 
    setDarkmode, 
    setSelectedSensor, 
    setSelectedMeasurement, 
    setSelectedTimePeriod 
} = appStateSlice.actions;

export const selectToken = (state) => state.app.token;
export const selectIsAuthenticated = (state) => state.app.isAuthenticated;
export const selectDarkmode = (state) => state.app.darkmode;
export const selectSelectedSensor = (state) => state.app.selectedSensor;
export const selectSelectedMeasurement = (state) => state.app.selectedMeasurement;
export const selectSelectedTimePeriod = (state) => state.app.selectedTimePeriod;

export default appStateSlice.reducer;