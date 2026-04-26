export interface Props {
	identifier: string;
	password: string;
	orgId: string;
}

export interface Token {
	accessToken: string;
	refreshToken: string;
}
