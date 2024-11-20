import { GenericErrorDto } from "@/dto/generic-error.dto"
import { LoginRequestDto } from "@/dto/login.dto"
import { ProjectDto } from "@/dto/project.dto"
import { LoginResponseDto } from "@/dto/login.dto"

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
            .then(data => resolve(data as LoginResponseDto))
            .catch(e => {
                console.log(e)
                console.log('genericerror', e)
                return reject(e as GenericErrorDto)
            })
    })
}

const AuthService = {
    tryLogin
}

export default AuthService;