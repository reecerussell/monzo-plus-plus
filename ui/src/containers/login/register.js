import React, { useState } from "react";
import { Fetch } from "../../utils/fetch";
import Register from "../../components/login/register";

const RegisterContainer = () => {
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const [confirm, setConfirm] = useState("");
	const [error, setError] = useState(null);
	const [loading, setLoading] = useState(false);

	const handleUpdateUsername = (e) => setUsername(e.target.value);
	const handleUpdatePassword = (e) => setPassword(e.target.value);
	const handleUpdateConfirm = (e) => setConfirm(e.target.value);

	const handleSubmit = async () => {
		if (loading) {
			return;
		}

		if (username === "") {
			setError("Enter a username!");
			return;
		}

		if (password === "") {
			setError("Enter a password!");
			return;
		}

		if (confirm !== password) {
			setError("Your passwords do not match!");
			return;
		}

		setLoading(true);

		await Fetch(
			"api/auth/users/register",
			{
				method: "POST",
				redirect: "manual",
				body: JSON.stringify({
					username,
					password,
				}),
			},
			async (res) => {
				const stateToken = await res.text();

				window.location.replace(
					"http://localhost:9789/api/auth/monzo/login?state=" +
						stateToken
				);
			},
			setError
		);

		setLoading(false);
	};

	return (
		<Register
			username={username}
			password={password}
			confirm={confirm}
			handleUpdateUsername={handleUpdateUsername}
			handleUpdatePassword={handleUpdatePassword}
			handleUpdateConfirm={handleUpdateConfirm}
			handleSubmit={handleSubmit}
			loading={loading}
			error={error}
		/>
	);
};

export default RegisterContainer;
