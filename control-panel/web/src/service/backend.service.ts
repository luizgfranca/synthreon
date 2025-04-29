const BASE_URL = `http://${window.location.hostname}:25256`
const STORAGE_SESSION_TOKEN_KEY = 'accessToken'

const DEFAULT_HEADERS = {
    'Content-Type': 'application/json'
}


function saveAccessToken(accessToken: string) {
    return window.localStorage.setItem(STORAGE_SESSION_TOKEN_KEY, accessToken)
}

function getAccessToken(): string | null {
    return window.localStorage.getItem(STORAGE_SESSION_TOKEN_KEY)
}

function clearAccessToken() {
    return window.localStorage.removeItem(STORAGE_SESSION_TOKEN_KEY)
}

function goToAuthentication() {
    return window.location.replace(`${import.meta.env.PL_PATH_PREFIX}/login`)
}

async function request(
    path: string,
    config?: RequestInit,
) {
    const maybeAccessToken = getAccessToken()
    if (!maybeAccessToken) {
        goToAuthentication()
    }
    const accessToken = maybeAccessToken;

    try {
        const response = await fetch(
            BASE_URL + path,
            {
                ...config,
                headers: {
                    ...config?.headers,
                    ...DEFAULT_HEADERS,
                    'Authorization': `Bearer ${accessToken}`
                }
            }
        )

        if (response.status === 401) {
            goToAuthentication();
        }
        return response;
    } catch (e) {
        console.error('fetch exception', e)
        throw e
    }
}

const BackendService = {
    getAccessToken,
    saveAccessToken,
    clearAccessToken,
    request
}

export default BackendService;
