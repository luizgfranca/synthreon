export type LoginResponseDto = {
	id: number;
	name: string;
	email: string;
}

export type LoginRequestDto = {
	email: string;
	password: string;
}