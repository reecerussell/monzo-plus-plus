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

const defaultFail = err => console.error(err);

const Fetch = async (
	url,
	options,
	onSuccess,
	onFail = defaultFail,
	successStatus = 200
) => {
	try {
		const res = await Send("http://localhost:9789/" + url, options);

		if (res.status === successStatus) {
			await onSuccess(res);
		} else {
			const { error } = await res.json();

			onFail(error);
		}
	} catch {
		onFail(
			"It seems like you don't have connection to the internet. Try again later!"
		);
	}
};

export default Send;
export { Fetch };
