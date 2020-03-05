const AccessTokenCookieName = "mpp_ac";

const GetAccessToken = () => {
	var v = document.cookie.match(
		"(^|;) ?" + AccessTokenCookieName + "=([^;]*)(;|$)"
	);
	return v ? v[2] : null;
};

const IsAuthenticated = () => {
	return GetAccessToken() !== null;
};

const IsInRole = roleName => {
	if (!IsAuthenticated()) {
		return false;
	}

	const payload = getCurrentPayload();
	if (!payload) {
		return false;
	}

	return payload.roles.indexOf(roleName) > -1;
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

export { IsInRole };