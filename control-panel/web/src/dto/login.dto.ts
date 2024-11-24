export type LoginResponseDto = {
	id: number;
	name: string;
	email: string;
	access_token: string;
}

export type LoginRequestDto = {
	email: string;
	password: string;
}