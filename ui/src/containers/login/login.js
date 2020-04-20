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

	const handleUpdateUsername = (e) => setUsername(e.target.value);
	const handleUpdatePassword = (e) => setPassword(e.target.value);

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

		await Fetch(
			"api/auth/token",
			{
				method: "POST",
				body: JSON.stringify({
					username,
					password,
				}),
			},
			async (res) => {
				setError(null);

				const { accessToken, expires } = await res.json();
				User.SetAccessToken(accessToken, expires * 1000);

				setRedirect("/account");
			},
			setError
		);

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
