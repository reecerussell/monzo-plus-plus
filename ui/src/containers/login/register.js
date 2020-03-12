import React, { useState } from "react";
import Fetch from "../../utils/fetch";
import Register from "../../components/login/register";

const RegisterContainer = () => {
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const [confirm, setConfirm] = useState("");
	const [error, setError] = useState(null);
	const [loading, setLoading] = useState(false);
	const [redirect, setRedirect] = useState(null);

	const handleUpdateUsername = e => setUsername(e.target.value);
	const handleUpdatePassword = e => setPassword(e.target.value);
	const handleUpdateConfirm = e => setConfirm(e.target.value);

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

		try {
			const res = await Fetch(
				"http://localhost:9789/auth/users/register",
				{
					method: "POST",
					body: JSON.stringify({
						username,
						password,
					}),
				}
			);

			if (res.status === 201) {
				setError(null);

				setRedirect("/login");
			} else {
				const data = await res.json();
				setError(data.error);
			}
		} catch {
			setError(
				"It seems like you don't have connection to the internet. Try again later!"
			);
		}

		setLoading(false);
	};

	return (
		<Register
			username={username}
			password={password}
			confirm={confirm}
			redirect={redirect}
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
