import config from '../config.json';

const apiUrl = `${config.backendUrl}/api`;

// returns body of the response
const doApiCall = async (token, method, path, data) => {
    try {
        if (token === '') {
            throw new Error('no authentication token provided, wont do API call');
        }
        if ((method === 'POST' || method === 'PUT') && !data) {
            throw new Error('POST or PUT method but no data given, wont do API call');
        }
        const response = await fetch(`${apiUrl}/${path}`, {
            method: method,
            mode: 'cors',
            headers: {
                'X-API-KEY': token,
            },
            body: data ? JSON.stringify(data) : null,
        });
        if (!response.ok) {
            throw new Error('network response was not ok');
        }
        return response.json();
    } catch (error) {
        throw new Error('request failed with: ' + error.message);
    }
};

// returns Promise<boolean> ok or not
const doCheckToken = async (token) => {
    console.debug(`checking validity of token ${token}`);
    
    return new Promise(resolve => {
        if (token === '') {
            resolve(false);
            return;
        }
        fetch(
            `${apiUrl}/checkToken`, {
                method: 'GET',
                mode: 'cors',
                headers: {
                    'X-API-KEY': token,
                },
            }
        ).then(response => {
            if (response.status === 200) {
                console.info('token is valid');
                resolve(true);
            } else {
                console.warn('token is invalid');
                resolve(false);
            }
        });
    });
};

// returns Promise<string> containing token if authorization ok, empty string if nok
const doAuthorize = async (password) => {
    return new Promise(resolve => {
        fetch(
            `${apiUrl}/authorize`, {
                method: 'POST',
                mode: 'cors',
                body: JSON.stringify({
                    'username': 'generic-user',
                    'password': password,
                })
            })
        .then(response => response.json())
        .then(data => {
            if (data['ok'] !== true) {
                console.error('failed to log in');
                resolve('');
            }
            console.info('log in succeeded');
            resolve(data['token']);
        });
    });
};

export { doApiCall, doCheckToken, doAuthorize };