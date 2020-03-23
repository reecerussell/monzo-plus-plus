import React, { useState } from "react";
import * as User from "../../utils/user";
import Fetch from "../../utils/fetch";
import Login from "../../components/login/login";

const LoginContainer = () => {
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const [error, setError] = useState(null);
	const [loading, setLoading] = useState(false);
	const [redirect, setRedirect] = useState(null);

	const handleUpdateUsername = e => setUsername(e.target.value);
	const handleUpdatePassword = e => setPassword(e.target.value);

	const handleSubmit = async () => {
		if (loading) {
			return;
		}

		if (username === "") {
			setError("Enter your username!");
			return;
		}

		if (password === "") {
			setError("Enter your password!");
			return;
		}

		setLoading(true);

		try {
			const res = await Fetch("http://localhost:9789/auth/token", {
				method: "POST",
				body: JSON.stringify({
					username,
					password,
				}),
			});

			if (res.status === 200) {
				setError(null);

				const { accessToken, expires } = await res.json();
				User.SetAccessToken(accessToken, expires * 1000);

				setRedirect("/account");
				return;
			}

			const data = await res.json();
			setError(data.error);
		} catch {
			setError(
				"It seems like you don't have connection to the internet. Try again later!"
			);
		}

		setLoading(false);
	};

	return (
		<Login
			username={username}
			password={password}
			handleUpdateUsername={handleUpdateUsername}
			handleUpdatePassword={handleUpdatePassword}
			handleSubmit={handleSubmit}
			loading={loading}
			error={error}
			redirect={redirect}
		/>
	);
};

export default LoginContainer;
