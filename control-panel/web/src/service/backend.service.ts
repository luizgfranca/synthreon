const BASE_URL = import.meta.env.PL_BACKEND_URL
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
    return window.location.replace('/login')
}

function request(
    path: string,
    config?: RequestInit,
) {
    const maybeAccessToken = getAccessToken()
    if (!maybeAccessToken) {
        goToAuthentication()
    }
    const accessToken = maybeAccessToken;

    return fetch(
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
}

const BackendService = {
    getAccessToken,
    saveAccessToken,
    clearAccessToken,
    request
}

export default BackendService;