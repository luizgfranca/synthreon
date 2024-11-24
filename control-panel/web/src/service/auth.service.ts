import { GenericErrorDto } from "@/dto/generic-error.dto"
import { LoginRequestDto } from "@/dto/login.dto"
import { ProjectDto } from "@/dto/project.dto"
import { LoginResponseDto } from "@/dto/login.dto"
import BackendService from "./backend.service"

export type QueryProjecstDto = ProjectDto[]

const BASE_URL = import.meta.env.PL_BACKEND_URL

function tryLogin(body: LoginRequestDto) {
    return new Promise((resolve, reject) => {
        fetch(`${BASE_URL}/api/auth/login`, {
            method: 'POST',
            body: JSON.stringify(body),
            headers: {
                'Content-Type': 'application/json'
            }
        })
            .then(response => {
                if(response.status >= 400){
                    throw response.json();
                }
                return response.json()
            })
            .then((data: LoginResponseDto) => {
                BackendService.saveAccessToken(data.access_token)
                resolve(data)
            })
            .catch(e => {
                console.log(e)
                console.log('genericerror', e)
                return reject(e as GenericErrorDto)
            })
    })
}

function isAuthenticated() {
    return BackendService.getAccessToken() != null
}

function logout() {
    BackendService.clearAccessToken();
}

const AuthService = {
    tryLogin,
    isAuthenticated,
    logout
}

export default AuthService;