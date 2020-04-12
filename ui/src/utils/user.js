const AccessTokenCookieName = "mpp_ac";

const GetAccessToken = () => {
	const value = document.cookie.match(
		"(^|;) ?" + AccessTokenCookieName + "=([^;]*)(;|$)"
	);

	return value ? value[2] : null;
};

const SetAccessToken = (token, expires) => {
	const d = new Date(expires);
	document.cookie =
		AccessTokenCookieName +
		"=" +
		token +
		";path=/;expires=" +
		d.toGMTString();

	LoginSubscriptions.forEach((v) => v());
};

const Logout = () => {
	SetAccessToken(null, -1);

	LogoutSubscriptions.forEach((v) => v());
};

const IsAuthenticated = () => {
	return GetAccessToken() !== null;
};

const IsInRole = (roleName) => {
	if (!IsAuthenticated()) {
		return false;
	}

	const payload = getCurrentPayload();
	if (!payload) {
		return false;
	}

	let roles = [];
	if (payload.roles) {
		roles = payload.roles.map((name) => name.toLowerCase());
	}

	return roles.indexOf(roleName.toLowerCase()) > -1;
};

const GetUsername = () => getClaim("username") ?? "User";

const GetId = () => getClaim("user_id");

const HasAccount = () => getClaim("has_account") == true;

const getClaim = (claimName) => {
	const payload = getCurrentPayload();
	if (!payload) {
		return null;
	}

	return payload[claimName];
};

const getCurrentPayload = () => {
	const token = GetAccessToken();
	if (!token) {
		return null;
	}

	const parts = token.split(".");
	if (parts.length < 2) {
		return null;
	}

	const payloadData = atob(parts[1]);

	return JSON.parse(payloadData);
};

const LoginSubscriptions = new Map();
const LogoutSubscriptions = new Map();

const SubscribeLogin = (name, callback) =>
	LoginSubscriptions.set(name, callback);
const SubscribeLogout = (name, callback) =>
	LogoutSubscriptions.set(name, callback);

const UnsubscripeLogin = (name) => LoginSubscriptions.delete(name);
const UnsubscribeLogout = (name) => LogoutSubscriptions.delete(name);

export {
	IsInRole,
	GetUsername,
	GetId,
	HasAccount,
	IsAuthenticated,
	SetAccessToken,
	SubscribeLogin,
	SubscribeLogout,
	UnsubscripeLogin,
	UnsubscribeLogout,
	Logout,
};
