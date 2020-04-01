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

	return await fetch(url, options);
};

const defaultFail = err => console.error(err);

const BaseUrl = "http://localhost:9789/";

const Fetch = async (
	url,
	options,
	onSuccess,
	onFail = defaultFail,
	baseUrl = BaseUrl
) => {
	try {
		const res = await Send(baseUrl + url, options);

		switch (res.status) {
			case 200:
			case 201:
				if (onSuccess) {
					await onSuccess(res);
				}
				break;
			case 401:
			case 403:
				window.location.hash = "/login";
				break;
			default:
				const { error } = await res.json();
				onFail(error);
				break;
		}
	} catch (e) {
		console.log(e);
		onFail(
			"It seems like you don't have connection to the internet. Try again later!"
		);
	}
};

export default Send;
export { Fetch, BaseUrl };
