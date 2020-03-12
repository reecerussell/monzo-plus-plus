const getAccessToken = () => {
	const cookieName = "mpp_ac";
	const v = document.cookie.match("(^|;) ?" + cookieName + "=([^;]*)(;|$)");
	return v ? v[2] : null;
};

const Send = async (url, options) => {
	if (!options) {
		options = {
			headers: {
				"Content-Type": "application/json",
				Authorization: "Bearer " + getAccessToken(),
			},
		};
	} else if (!options.headers) {
		options.headers = {
			"Content-Type": "application/json",
			Authorization: "Bearer " + getAccessToken(),
		};
	} else if (!options.headers["Authorization"]) {
		options.headers["Authorization"] = "Bearer " + getAccessToken();
	}

	console.log(options);

	return await fetch(url, options);
};

export default Send;
