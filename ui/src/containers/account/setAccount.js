import React, { useEffect, useState } from "react";
import Fetch from "../../utils/fetch";
import * as User from "../../utils/user";
import SetAccount from "../../components/account/setAccount";

const SetAccountContainer = () => {
	const [accounts, setAccounts] = useState([]);
	const [selectedAccount, setSelectedAccount] = useState("");
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState(null);
	const [success, setSuccess] = useState(null);
	const [redirect, setRedirect] = useState(null);

	const fetchAccounts = async () => {
		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"api/auth/users/accounts/" + User.GetId(),
			null,
			async (res) => {
				const accounts = await res.json();
				setAccounts(accounts);
				setSelectedAccount(accounts[0].id);
			},
			setError
		);

		setLoading(false);
	};

	const handleSubmit = async (e) => {
		e.preventDefault();

		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"api/auth/users/account",
			{
				method: "POST",
				body: JSON.stringify({
					userId: User.GetId(),
					accountId: selectedAccount,
				}),
			},
			async () => {
				setError(null);
				setSuccess("Your changes have been saved successfully!");

				await Fetch(
					"api/auth/refresh",
					null,
					async (res) => {
						const { accessToken, expires } = await res.json();
						User.SetAccessToken(accessToken, expires * 1000);
					},
					setError
				);

				window.location.reload();
			},
			setError
		);

		setLoading(false);
	};

	const handleUpdateAccount = (e) => setSelectedAccount(e.target.value);

	useEffect(() => {
		fetchAccounts();
	}, []);

	useEffect(() => {
		if (error !== null) {
			setTimeout(() => setError(null), 5000);
		}
	}, [error]);

	return (
		<SetAccount
			redirect={redirect}
			error={error}
			success={success}
			loading={loading}
			accounts={accounts}
			selectedAccount={selectedAccount}
			handleSubmit={handleSubmit}
			handleUpdateAccount={handleUpdateAccount}
		/>
	);
};

export default SetAccountContainer;
